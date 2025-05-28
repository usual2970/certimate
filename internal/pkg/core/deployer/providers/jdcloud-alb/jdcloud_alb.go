package jdcloudalb

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	jdcore "github.com/jdcloud-api/jdcloud-sdk-go/core"
	jdcommon "github.com/jdcloud-api/jdcloud-sdk-go/services/common/models"
	jdlbapi "github.com/jdcloud-api/jdcloud-sdk-go/services/lb/apis"
	jdlbclient "github.com/jdcloud-api/jdcloud-sdk-go/services/lb/client"
	jdlbmodel "github.com/jdcloud-api/jdcloud-sdk-go/services/lb/models"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	uploadersp "github.com/usual2970/certimate/internal/pkg/core/uploader/providers/jdcloud-ssl"
	sliceutil "github.com/usual2970/certimate/internal/pkg/utils/slice"
)

type DeployerConfig struct {
	// 京东云 AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// 京东云 AccessKeySecret。
	AccessKeySecret string `json:"accessKeySecret"`
	// 京东云地域 ID。
	RegionId string `json:"regionId"`
	// 部署资源类型。
	ResourceType ResourceType `json:"resourceType"`
	// 负载均衡器 ID。
	// 部署资源类型为 [RESOURCE_TYPE_LOADBALANCER] 时必填。
	LoadbalancerId string `json:"loadbalancerId,omitempty"`
	// 监听器 ID。
	// 部署资源类型为 [RESOURCE_TYPE_LISTENER] 时必填。
	ListenerId string `json:"listenerId,omitempty"`
	// SNI 域名（支持泛域名）。
	// 部署资源类型为 [RESOURCE_TYPE_LOADBALANCER]、[RESOURCE_TYPE_LISTENER] 时选填。
	Domain string `json:"domain,omitempty"`
}

type DeployerProvider struct {
	config      *DeployerConfig
	logger      *slog.Logger
	sdkClient   *jdlbclient.LbClient
	sslUploader uploader.Uploader
}

var _ deployer.Deployer = (*DeployerProvider)(nil)

func NewDeployer(config *DeployerConfig) (*DeployerProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	client, err := createSdkClient(config.AccessKeyId, config.AccessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("failed to create sdk client: %w", err)
	}

	uploader, err := uploadersp.NewUploader(&uploadersp.UploaderConfig{
		AccessKeyId:     config.AccessKeyId,
		AccessKeySecret: config.AccessKeySecret,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create ssl uploader: %w", err)
	}

	return &DeployerProvider{
		config:      config,
		logger:      slog.Default(),
		sdkClient:   client,
		sslUploader: uploader,
	}, nil
}

func (d *DeployerProvider) WithLogger(logger *slog.Logger) deployer.Deployer {
	if logger == nil {
		d.logger = slog.New(slog.DiscardHandler)
	} else {
		d.logger = logger
	}
	d.sslUploader.WithLogger(logger)
	return d
}

func (d *DeployerProvider) Deploy(ctx context.Context, certPEM string, privkeyPEM string) (*deployer.DeployResult, error) {
	// 上传证书到 SSL
	upres, err := d.sslUploader.Upload(ctx, certPEM, privkeyPEM)
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

	return &deployer.DeployResult{}, nil
}

func (d *DeployerProvider) deployToLoadbalancer(ctx context.Context, cloudCertId string) error {
	if d.config.LoadbalancerId == "" {
		return errors.New("config `loadbalancerId` is required")
	}

	// 查询负载均衡器详情
	// REF: https://docs.jdcloud.com/cn/load-balancer/api/describeloadbalancer
	describeLoadBalancerReq := jdlbapi.NewDescribeLoadBalancerRequest(d.config.RegionId, d.config.LoadbalancerId)
	describeLoadBalancerResp, err := d.sdkClient.DescribeLoadBalancer(describeLoadBalancerReq)
	d.logger.Debug("sdk request 'lb.DescribeLoadBalancer'", slog.Any("request", describeLoadBalancerReq), slog.Any("response", describeLoadBalancerResp))
	if err != nil {
		return fmt.Errorf("failed to execute sdk request 'lb.DescribeLoadBalancer': %w", err)
	}

	// 查询监听器列表
	// REF: https://docs.jdcloud.com/cn/load-balancer/api/describelisteners
	listenerIds := make([]string, 0)
	describeListenersPageNumber := 1
	describeListenersPageSize := 100
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		describeListenersReq := jdlbapi.NewDescribeListenersRequest(d.config.RegionId)
		describeListenersReq.SetFilters([]jdcommon.Filter{{Name: "loadBalancerId", Values: []string{d.config.LoadbalancerId}}})
		describeListenersReq.SetPageSize(describeListenersPageNumber)
		describeListenersReq.SetPageSize(describeListenersPageSize)
		describeListenersResp, err := d.sdkClient.DescribeListeners(describeListenersReq)
		d.logger.Debug("sdk request 'lb.DescribeListeners'", slog.Any("request", describeListenersReq), slog.Any("response", describeListenersResp))
		if err != nil {
			return fmt.Errorf("failed to execute sdk request 'lb.DescribeListeners': %w", err)
		}

		for _, listener := range describeListenersResp.Result.Listeners {
			if strings.EqualFold(listener.Protocol, "https") || strings.EqualFold(listener.Protocol, "tls") {
				listenerIds = append(listenerIds, listener.ListenerId)
			}
		}

		if len(describeListenersResp.Result.Listeners) < int(describeListenersPageSize) {
			break
		} else {
			describeListenersPageNumber++
		}
	}

	// 遍历更新监听器证书
	if len(listenerIds) == 0 {
		d.logger.Info("no listeners to deploy")
	} else {
		d.logger.Info("found https/tls listeners to deploy", slog.Any("listenerIds", listenerIds))

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

func (d *DeployerProvider) deployToListener(ctx context.Context, cloudCertId string) error {
	if d.config.ListenerId == "" {
		return errors.New("config `listenerId` is required")
	}

	// 更新监听器证书
	if err := d.updateListenerCertificate(ctx, d.config.ListenerId, cloudCertId); err != nil {
		return err
	}

	return nil
}

func (d *DeployerProvider) updateListenerCertificate(ctx context.Context, cloudListenerId string, cloudCertId string) error {
	// 查询监听器详情
	// REF: https://docs.jdcloud.com/cn/load-balancer/api/describelistener
	describeListenerReq := jdlbapi.NewDescribeListenerRequest(d.config.RegionId, cloudListenerId)
	describeListenerResp, err := d.sdkClient.DescribeListener(describeListenerReq)
	d.logger.Debug("sdk request 'lb.DescribeListener'", slog.Any("request", describeListenerReq), slog.Any("response", describeListenerResp))
	if err != nil {
		return fmt.Errorf("failed to execute sdk request 'lb.DescribeListener': %w", err)
	}

	if d.config.Domain == "" {
		// 未指定 SNI，只需部署到监听器

		// 修改监听器信息
		// REF: https://docs.jdcloud.com/cn/load-balancer/api/updatelistener
		updateListenerReq := jdlbapi.NewUpdateListenerRequest(d.config.RegionId, cloudListenerId)
		updateListenerReq.SetCertificateSpecs([]jdlbmodel.CertificateSpec{{CertificateId: cloudCertId}})
		updateListenerResp, err := d.sdkClient.UpdateListener(updateListenerReq)
		d.logger.Debug("sdk request 'lb.UpdateListener'", slog.Any("request", updateListenerReq), slog.Any("response", updateListenerResp))
		if err != nil {
			return fmt.Errorf("failed to execute sdk request 'lb.UpdateListener': %w", err)
		}
	} else {
		// 指定 SNI，需部署到扩展证书

		extCertSpecs := sliceutil.Filter(describeListenerResp.Result.Listener.ExtensionCertificateSpecs, func(extCertSpec jdlbmodel.ExtensionCertificateSpec) bool {
			return extCertSpec.Domain == d.config.Domain
		})
		if len(extCertSpecs) == 0 {
			return errors.New("extension certificate spec not found")
		}

		// 批量修改扩展证书
		// REF: https://docs.jdcloud.com/cn/load-balancer/api/updatelistenercertificates
		updateListenerCertificatesReq := jdlbapi.NewUpdateListenerCertificatesRequest(
			d.config.RegionId,
			cloudListenerId,
			sliceutil.Map(extCertSpecs, func(extCertSpec jdlbmodel.ExtensionCertificateSpec) jdlbmodel.ExtCertificateUpdateSpec {
				return jdlbmodel.ExtCertificateUpdateSpec{
					CertificateBindId: extCertSpec.CertificateBindId,
					CertificateId:     &cloudCertId,
					Domain:            &extCertSpec.Domain,
				}
			}),
		)
		updateListenerCertificatesResp, err := d.sdkClient.UpdateListenerCertificates(updateListenerCertificatesReq)
		d.logger.Debug("sdk request 'lb.UpdateListenerCertificates'", slog.Any("request", updateListenerCertificatesReq), slog.Any("response", updateListenerCertificatesResp))
		if err != nil {
			return fmt.Errorf("failed to execute sdk request 'lb.UpdateListenerCertificates': %w", err)
		}
	}

	return nil
}

func createSdkClient(accessKeyId, accessKeySecret string) (*jdlbclient.LbClient, error) {
	clientCredentials := jdcore.NewCredentials(accessKeyId, accessKeySecret)
	client := jdlbclient.NewLbClient(clientCredentials)
	client.SetLogger(jdcore.NewDefaultLogger(jdcore.LogWarn))
	return client, nil
}
