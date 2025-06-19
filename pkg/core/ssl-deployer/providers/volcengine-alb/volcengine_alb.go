package volcenginealb

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	vealb "github.com/volcengine/volcengine-go-sdk/service/alb"
	ve "github.com/volcengine/volcengine-go-sdk/volcengine"
	vesession "github.com/volcengine/volcengine-go-sdk/volcengine/session"

	"github.com/certimate-go/certimate/pkg/core"
	sslmgrsp "github.com/certimate-go/certimate/pkg/core/ssl-manager/providers/volcengine-certcenter"
	xslices "github.com/certimate-go/certimate/pkg/utils/slices"
)

type SSLDeployerProviderConfig struct {
	// 火山引擎 AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// 火山引擎 AccessKeySecret。
	AccessKeySecret string `json:"accessKeySecret"`
	// 火山引擎地域。
	Region string `json:"region"`
	// 部署资源类型。
	ResourceType ResourceType `json:"resourceType"`
	// 负载均衡实例 ID。
	// 部署资源类型为 [RESOURCE_TYPE_LOADBALANCER] 时必填。
	LoadbalancerId string `json:"loadbalancerId,omitempty"`
	// 负载均衡监听器 ID。
	// 部署资源类型为 [RESOURCE_TYPE_LISTENER] 时必填。
	ListenerId string `json:"listenerId,omitempty"`
	// SNI 域名（支持泛域名）。
	// 部署资源类型为 [RESOURCE_TYPE_LOADBALANCER]、[RESOURCE_TYPE_LISTENER] 时选填。
	Domain string `json:"domain,omitempty"`
}

type SSLDeployerProvider struct {
	config     *SSLDeployerProviderConfig
	logger     *slog.Logger
	sdkClient  *vealb.ALB
	sslManager core.SSLManager
}

var _ core.SSLDeployer = (*SSLDeployerProvider)(nil)

func NewSSLDeployerProvider(config *SSLDeployerProviderConfig) (*SSLDeployerProvider, error) {
	if config == nil {
		return nil, errors.New("the configuration of the ssl deployer provider is nil")
	}

	client, err := createSDKClient(config.AccessKeyId, config.AccessKeySecret, config.Region)
	if err != nil {
		return nil, fmt.Errorf("could not create sdk client: %w", err)
	}

	sslmgr, err := sslmgrsp.NewSSLManagerProvider(&sslmgrsp.SSLManagerProviderConfig{
		AccessKeyId:     config.AccessKeyId,
		AccessKeySecret: config.AccessKeySecret,
		Region:          config.Region,
	})
	if err != nil {
		return nil, fmt.Errorf("could not create ssl manager: %w", err)
	}

	return &SSLDeployerProvider{
		config:     config,
		logger:     slog.Default(),
		sdkClient:  client,
		sslManager: sslmgr,
	}, nil
}

func (d *SSLDeployerProvider) SetLogger(logger *slog.Logger) {
	if logger == nil {
		d.logger = slog.New(slog.DiscardHandler)
	} else {
		d.logger = logger
	}

	d.sslManager.SetLogger(logger)
}

func (d *SSLDeployerProvider) Deploy(ctx context.Context, certPEM string, privkeyPEM string) (*core.SSLDeployResult, error) {
	// 上传证书
	upres, err := d.sslManager.Upload(ctx, certPEM, privkeyPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to upload certificate file: %w", err)
	} else {
		d.logger.Info("ssl certificate uploaded", slog.Any("result", upres))
	}

	// 根据部署资源类型决定部署方式
	switch d.config.ResourceType {
	case RESOURCE_TYPE_LOADBALANCER:
		if err := d.deployToLoadbalancer(ctx, upres.CertId); err != nil {
			return nil, err
		}

	case RESOURCE_TYPE_LISTENER:
		if err := d.deployToListener(ctx, upres.CertId); err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("unsupported resource type '%s'", d.config.ResourceType)
	}

	return &core.SSLDeployResult{}, nil
}

func (d *SSLDeployerProvider) deployToLoadbalancer(ctx context.Context, cloudCertId string) error {
	if d.config.LoadbalancerId == "" {
		return errors.New("config `loadbalancerId` is required")
	}

	// 查询 ALB 实例的详细信息
	// REF: https://www.volcengine.com/docs/6767/113596
	describeLoadBalancerAttributesReq := &vealb.DescribeLoadBalancerAttributesInput{
		LoadBalancerId: ve.String(d.config.LoadbalancerId),
	}
	describeLoadBalancerAttributesResp, err := d.sdkClient.DescribeLoadBalancerAttributes(describeLoadBalancerAttributesReq)
	d.logger.Debug("sdk request 'alb.DescribeLoadBalancerAttributes'", slog.Any("request", describeLoadBalancerAttributesReq), slog.Any("response", describeLoadBalancerAttributesResp))
	if err != nil {
		return fmt.Errorf("failed to execute sdk request 'alb.DescribeLoadBalancerAttributes': %w", err)
	}

	// 查询 HTTPS 监听器列表
	// REF: https://www.volcengine.com/docs/6767/113684
	listenerIds := make([]string, 0)
	describeListenersPageSize := int64(100)
	describeListenersPageNumber := int64(1)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		describeListenersReq := &vealb.DescribeListenersInput{
			LoadBalancerId: ve.String(d.config.LoadbalancerId),
			Protocol:       ve.String("HTTPS"),
			PageNumber:     ve.Int64(describeListenersPageNumber),
			PageSize:       ve.Int64(describeListenersPageSize),
		}
		describeListenersResp, err := d.sdkClient.DescribeListeners(describeListenersReq)
		d.logger.Debug("sdk request 'alb.DescribeListeners'", slog.Any("request", describeListenersReq), slog.Any("response", describeListenersResp))
		if err != nil {
			return fmt.Errorf("failed to execute sdk request 'alb.DescribeListeners': %w", err)
		}

		for _, listener := range describeListenersResp.Listeners {
			listenerIds = append(listenerIds, *listener.ListenerId)
		}

		if len(describeListenersResp.Listeners) < int(describeListenersPageSize) {
			break
		} else {
			describeListenersPageNumber++
		}
	}

	// 遍历更新监听证书
	if len(listenerIds) == 0 {
		d.logger.Info("no alb listeners to deploy")
	} else {
		d.logger.Info("found https listeners to deploy", slog.Any("listenerIds", listenerIds))
		var errs []error

		for _, listenerId := range listenerIds {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				if err := d.updateListenerCertificate(ctx, listenerId, cloudCertId); err != nil {
					errs = append(errs, err)
				}
			}
		}

		if len(errs) > 0 {
			return errors.Join(errs...)
		}
	}

	return nil
}

func (d *SSLDeployerProvider) deployToListener(ctx context.Context, cloudCertId string) error {
	if d.config.ListenerId == "" {
		return errors.New("config `listenerId` is required")
	}

	if err := d.updateListenerCertificate(ctx, d.config.ListenerId, cloudCertId); err != nil {
		return err
	}

	return nil
}

func (d *SSLDeployerProvider) updateListenerCertificate(ctx context.Context, cloudListenerId string, cloudCertId string) error {
	// 查询指定监听器的详细信息
	// REF: https://www.volcengine.com/docs/6767/113686
	describeListenerAttributesReq := &vealb.DescribeListenerAttributesInput{
		ListenerId: ve.String(cloudListenerId),
	}
	describeListenerAttributesResp, err := d.sdkClient.DescribeListenerAttributes(describeListenerAttributesReq)
	d.logger.Debug("sdk request 'alb.DescribeListenerAttributes'", slog.Any("request", describeListenerAttributesReq), slog.Any("response", describeListenerAttributesResp))
	if err != nil {
		return fmt.Errorf("failed to execute sdk request 'alb.DescribeListenerAttributes': %w", err)
	}

	if d.config.Domain == "" {
		// 未指定 SNI，只需部署到监听器

		// 修改指定监听器
		// REF: https://www.volcengine.com/docs/6767/113683
		modifyListenerAttributesReq := &vealb.ModifyListenerAttributesInput{
			ListenerId:              ve.String(cloudListenerId),
			CertificateSource:       ve.String("cert_center"),
			CertCenterCertificateId: ve.String(cloudCertId),
		}
		modifyListenerAttributesResp, err := d.sdkClient.ModifyListenerAttributes(modifyListenerAttributesReq)
		d.logger.Debug("sdk request 'alb.ModifyListenerAttributes'", slog.Any("request", modifyListenerAttributesReq), slog.Any("response", modifyListenerAttributesResp))
		if err != nil {
			return fmt.Errorf("failed to execute sdk request 'alb.ModifyListenerAttributes': %w", err)
		}
	} else {
		// 指定 SNI，需部署到扩展域名

		// 修改指定监听器
		// REF: https://www.volcengine.com/docs/6767/113683
		modifyListenerAttributesReq := &vealb.ModifyListenerAttributesInput{
			ListenerId: ve.String(cloudListenerId),
			DomainExtensions: xslices.Map(
				xslices.Filter(
					describeListenerAttributesResp.DomainExtensions,
					func(domain *vealb.DomainExtensionForDescribeListenerAttributesOutput) bool {
						return *domain.Domain == d.config.Domain
					},
				),
				func(domain *vealb.DomainExtensionForDescribeListenerAttributesOutput) *vealb.DomainExtensionForModifyListenerAttributesInput {
					return &vealb.DomainExtensionForModifyListenerAttributesInput{
						DomainExtensionId:       domain.DomainExtensionId,
						Domain:                  domain.Domain,
						CertificateSource:       ve.String("cert_center"),
						CertCenterCertificateId: ve.String(cloudCertId),
						Action:                  ve.String("modify"),
					}
				}),
		}
		modifyListenerAttributesResp, err := d.sdkClient.ModifyListenerAttributes(modifyListenerAttributesReq)
		d.logger.Debug("sdk request 'alb.ModifyListenerAttributes'", slog.Any("request", modifyListenerAttributesReq), slog.Any("response", modifyListenerAttributesResp))
		if err != nil {
			return fmt.Errorf("failed to execute sdk request 'alb.ModifyListenerAttributes': %w", err)
		}
	}

	return nil
}

func createSDKClient(accessKeyId, accessKeySecret, region string) (*vealb.ALB, error) {
	config := ve.NewConfig().WithRegion(region).WithAkSk(accessKeyId, accessKeySecret)

	session, err := vesession.NewSession(config)
	if err != nil {
		return nil, err
	}

	client := vealb.New(session)
	return client, nil
}
