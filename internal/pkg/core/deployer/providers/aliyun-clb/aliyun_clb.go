package aliyunclb

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	aliyunOpen "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	aliyunSlb "github.com/alibabacloud-go/slb-20140515/v4/client"
	"github.com/alibabacloud-go/tea/tea"
	xerrors "github.com/pkg/errors"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	uploadersp "github.com/usual2970/certimate/internal/pkg/core/uploader/providers/aliyun-slb"
)

type DeployerConfig struct {
	// 阿里云 AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// 阿里云 AccessKeySecret。
	AccessKeySecret string `json:"accessKeySecret"`
	// 阿里云地域。
	Region string `json:"region"`
	// 部署资源类型。
	ResourceType ResourceType `json:"resourceType"`
	// 负载均衡实例 ID。
	// 部署资源类型为 [RESOURCE_TYPE_LOADBALANCER]、[RESOURCE_TYPE_LISTENER] 时必填。
	LoadbalancerId string `json:"loadbalancerId,omitempty"`
	// 负载均衡监听端口。
	// 部署资源类型为 [RESOURCE_TYPE_LISTENER] 时必填。
	ListenerPort int32 `json:"listenerPort,omitempty"`
	// SNI 域名（支持泛域名）。
	// 部署资源类型为 [RESOURCE_TYPE_LOADBALANCER]、[RESOURCE_TYPE_LISTENER] 时选填。
	Domain string `json:"domain,omitempty"`
}

type DeployerProvider struct {
	config      *DeployerConfig
	logger      *slog.Logger
	sdkClient   *aliyunSlb.Client
	sslUploader uploader.Uploader
}

var _ deployer.Deployer = (*DeployerProvider)(nil)

func NewDeployer(config *DeployerConfig) (*DeployerProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	client, err := createSdkClient(config.AccessKeyId, config.AccessKeySecret, config.Region)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk client")
	}

	uploader, err := uploadersp.NewUploader(&uploadersp.UploaderConfig{
		AccessKeyId:     config.AccessKeyId,
		AccessKeySecret: config.AccessKeySecret,
		Region:          config.Region,
	})
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create ssl uploader")
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
		d.logger = slog.Default()
	} else {
		d.logger = logger
	}
	d.sslUploader.WithLogger(logger)
	return d
}

func (d *DeployerProvider) Deploy(ctx context.Context, certPem string, privkeyPem string) (*deployer.DeployResult, error) {
	// 上传证书到 SLB
	upres, err := d.sslUploader.Upload(ctx, certPem, privkeyPem)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to upload certificate file")
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
		return nil, fmt.Errorf("unsupported resource type: %s", d.config.ResourceType)
	}

	return &deployer.DeployResult{}, nil
}

func (d *DeployerProvider) deployToLoadbalancer(ctx context.Context, cloudCertId string) error {
	if d.config.LoadbalancerId == "" {
		return errors.New("config `loadbalancerId` is required")
	}

	// 查询负载均衡实例的详细信息
	// REF: https://help.aliyun.com/zh/slb/classic-load-balancer/developer-reference/api-slb-2014-05-15-describeloadbalancerattribute
	describeLoadBalancerAttributeReq := &aliyunSlb.DescribeLoadBalancerAttributeRequest{
		RegionId:       tea.String(d.config.Region),
		LoadBalancerId: tea.String(d.config.LoadbalancerId),
	}
	describeLoadBalancerAttributeResp, err := d.sdkClient.DescribeLoadBalancerAttribute(describeLoadBalancerAttributeReq)
	d.logger.Debug("sdk request 'slb.DescribeLoadBalancerAttribute'", slog.Any("request", describeLoadBalancerAttributeReq), slog.Any("response", describeLoadBalancerAttributeResp))
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'slb.DescribeLoadBalancerAttribute'")
	}

	// 查询 HTTPS 监听列表
	// REF: https://help.aliyun.com/zh/slb/classic-load-balancer/developer-reference/api-slb-2014-05-15-describeloadbalancerlisteners
	listenerPorts := make([]int32, 0)
	describeLoadBalancerListenersLimit := int32(100)
	var describeLoadBalancerListenersToken *string = nil
	for {
		describeLoadBalancerListenersReq := &aliyunSlb.DescribeLoadBalancerListenersRequest{
			RegionId:         tea.String(d.config.Region),
			MaxResults:       tea.Int32(describeLoadBalancerListenersLimit),
			NextToken:        describeLoadBalancerListenersToken,
			LoadBalancerId:   []*string{tea.String(d.config.LoadbalancerId)},
			ListenerProtocol: tea.String("https"),
		}
		describeLoadBalancerListenersResp, err := d.sdkClient.DescribeLoadBalancerListeners(describeLoadBalancerListenersReq)
		d.logger.Debug("sdk request 'slb.DescribeLoadBalancerListeners'", slog.Any("request", describeLoadBalancerListenersReq), slog.Any("response", describeLoadBalancerListenersResp))
		if err != nil {
			return xerrors.Wrap(err, "failed to execute sdk request 'slb.DescribeLoadBalancerListeners'")
		}

		if describeLoadBalancerListenersResp.Body.Listeners != nil {
			for _, listener := range describeLoadBalancerListenersResp.Body.Listeners {
				listenerPorts = append(listenerPorts, *listener.ListenerPort)
			}
		}

		if len(describeLoadBalancerListenersResp.Body.Listeners) == 0 || describeLoadBalancerListenersResp.Body.NextToken == nil {
			break
		} else {
			describeLoadBalancerListenersToken = describeLoadBalancerListenersResp.Body.NextToken
		}
	}

	// 遍历更新监听证书
	if len(listenerPorts) == 0 {
		d.logger.Info("no clb listeners to deploy")
	} else {
		d.logger.Info("found https listeners to deploy", slog.Any("listenerPorts", listenerPorts))
		var errs []error

		for _, listenerPort := range listenerPorts {
			if err := d.updateListenerCertificate(ctx, d.config.LoadbalancerId, listenerPort, cloudCertId); err != nil {
				errs = append(errs, err)
			}
		}

		if len(errs) > 0 {
			return errors.Join(errs...)
		}
	}

	return nil
}

func (d *DeployerProvider) deployToListener(ctx context.Context, cloudCertId string) error {
	if d.config.LoadbalancerId == "" {
		return errors.New("config `loadbalancerId` is required")
	}
	if d.config.ListenerPort == 0 {
		return errors.New("config `listenerPort` is required")
	}

	// 更新监听
	if err := d.updateListenerCertificate(ctx, d.config.LoadbalancerId, d.config.ListenerPort, cloudCertId); err != nil {
		return err
	}

	return nil
}

func (d *DeployerProvider) updateListenerCertificate(ctx context.Context, cloudLoadbalancerId string, cloudListenerPort int32, cloudCertId string) error {
	// 查询监听配置
	// REF: https://help.aliyun.com/zh/slb/classic-load-balancer/developer-reference/api-slb-2014-05-15-describeloadbalancerhttpslistenerattribute
	describeLoadBalancerHTTPSListenerAttributeReq := &aliyunSlb.DescribeLoadBalancerHTTPSListenerAttributeRequest{
		LoadBalancerId: tea.String(cloudLoadbalancerId),
		ListenerPort:   tea.Int32(cloudListenerPort),
	}
	describeLoadBalancerHTTPSListenerAttributeResp, err := d.sdkClient.DescribeLoadBalancerHTTPSListenerAttribute(describeLoadBalancerHTTPSListenerAttributeReq)
	d.logger.Debug("sdk request 'slb.DescribeLoadBalancerHTTPSListenerAttribute'", slog.Any("request", describeLoadBalancerHTTPSListenerAttributeReq), slog.Any("response", describeLoadBalancerHTTPSListenerAttributeResp))
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'slb.DescribeLoadBalancerHTTPSListenerAttribute'")
	}

	if d.config.Domain == "" {
		// 未指定 SNI，只需部署到监听器

		// 修改监听配置
		// REF: https://help.aliyun.com/zh/slb/classic-load-balancer/developer-reference/api-slb-2014-05-15-setloadbalancerhttpslistenerattribute
		setLoadBalancerHTTPSListenerAttributeReq := &aliyunSlb.SetLoadBalancerHTTPSListenerAttributeRequest{
			RegionId:            tea.String(d.config.Region),
			LoadBalancerId:      tea.String(cloudLoadbalancerId),
			ListenerPort:        tea.Int32(cloudListenerPort),
			ServerCertificateId: tea.String(cloudCertId),
		}
		setLoadBalancerHTTPSListenerAttributeResp, err := d.sdkClient.SetLoadBalancerHTTPSListenerAttribute(setLoadBalancerHTTPSListenerAttributeReq)
		d.logger.Debug("sdk request 'slb.SetLoadBalancerHTTPSListenerAttribute'", slog.Any("request", setLoadBalancerHTTPSListenerAttributeReq), slog.Any("response", setLoadBalancerHTTPSListenerAttributeResp))
		if err != nil {
			return xerrors.Wrap(err, "failed to execute sdk request 'slb.SetLoadBalancerHTTPSListenerAttribute'")
		}
	} else {
		// 指定 SNI，需部署到扩展域名

		// 查询扩展域名
		// REF: https://help.aliyun.com/zh/slb/classic-load-balancer/developer-reference/api-slb-2014-05-15-describedomainextensions
		describeDomainExtensionsReq := &aliyunSlb.DescribeDomainExtensionsRequest{
			RegionId:       tea.String(d.config.Region),
			LoadBalancerId: tea.String(cloudLoadbalancerId),
			ListenerPort:   tea.Int32(cloudListenerPort),
		}
		describeDomainExtensionsResp, err := d.sdkClient.DescribeDomainExtensions(describeDomainExtensionsReq)
		d.logger.Debug("sdk request 'slb.DescribeDomainExtensions'", slog.Any("request", describeDomainExtensionsReq), slog.Any("response", describeDomainExtensionsResp))
		if err != nil {
			return xerrors.Wrap(err, "failed to execute sdk request 'slb.DescribeDomainExtensions'")
		}

		// 遍历修改扩展域名
		// REF: https://help.aliyun.com/zh/slb/classic-load-balancer/developer-reference/api-slb-2014-05-15-setdomainextensionattribute
		if describeDomainExtensionsResp.Body.DomainExtensions != nil && describeDomainExtensionsResp.Body.DomainExtensions.DomainExtension != nil {
			var errs []error

			for _, domainExtension := range describeDomainExtensionsResp.Body.DomainExtensions.DomainExtension {
				if *domainExtension.Domain != d.config.Domain {
					continue
				}

				setDomainExtensionAttributeReq := &aliyunSlb.SetDomainExtensionAttributeRequest{
					RegionId:            tea.String(d.config.Region),
					DomainExtensionId:   tea.String(*domainExtension.DomainExtensionId),
					ServerCertificateId: tea.String(cloudCertId),
				}
				setDomainExtensionAttributeResp, err := d.sdkClient.SetDomainExtensionAttribute(setDomainExtensionAttributeReq)
				d.logger.Debug("sdk request 'slb.SetDomainExtensionAttribute'", slog.Any("request", setDomainExtensionAttributeReq), slog.Any("response", setDomainExtensionAttributeResp))
				if err != nil {
					errs = append(errs, xerrors.Wrap(err, "failed to execute sdk request 'slb.SetDomainExtensionAttribute'"))
					continue
				}
			}

			if len(errs) > 0 {
				return errors.Join(errs...)
			}
		}
	}

	return nil
}

func createSdkClient(accessKeyId, accessKeySecret, region string) (*aliyunSlb.Client, error) {
	// 接入点一览 https://api.aliyun.com/product/Slb
	var endpoint string
	switch region {
	case
		"cn-hangzhou",
		"cn-hangzhou-finance",
		"cn-shanghai-finance-1",
		"cn-shenzhen-finance-1":
		endpoint = "slb.aliyuncs.com"
	default:
		endpoint = fmt.Sprintf("slb.%s.aliyuncs.com", region)
	}

	config := &aliyunOpen.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
		Endpoint:        tea.String(endpoint),
	}

	client, err := aliyunSlb.NewClient(config)
	if err != nil {
		return nil, err
	}

	return client, nil
}
