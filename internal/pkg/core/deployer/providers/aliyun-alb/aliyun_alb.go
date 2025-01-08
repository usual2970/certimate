package aliyunalb

import (
	"context"
	"errors"
	"fmt"
	"strings"

	aliyunAlb "github.com/alibabacloud-go/alb-20200616/v2/client"
	aliyunOpen "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	xerrors "github.com/pkg/errors"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/logger"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	providerCas "github.com/usual2970/certimate/internal/pkg/core/uploader/providers/aliyun-cas"
)

type AliyunALBDeployerConfig struct {
	// 阿里云 AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// 阿里云 AccessKeySecret。
	AccessKeySecret string `json:"accessKeySecret"`
	// 阿里云地域。
	Region string `json:"region"`
	// 部署资源类型。
	ResourceType DeployResourceType `json:"resourceType"`
	// 负载均衡实例 ID。
	// 部署资源类型为 [DEPLOY_RESOURCE_LOADBALANCER] 时必填。
	LoadbalancerId string `json:"loadbalancerId,omitempty"`
	// 负载均衡监听 ID。
	// 部署资源类型为 [DEPLOY_RESOURCE_LISTENER] 时必填。
	ListenerId string `json:"listenerId,omitempty"`
}

type AliyunALBDeployer struct {
	config      *AliyunALBDeployerConfig
	logger      logger.Logger
	sdkClient   *aliyunAlb.Client
	sslUploader uploader.Uploader
}

var _ deployer.Deployer = (*AliyunALBDeployer)(nil)

func New(config *AliyunALBDeployerConfig) (*AliyunALBDeployer, error) {
	return NewWithLogger(config, logger.NewNilLogger())
}

func NewWithLogger(config *AliyunALBDeployerConfig, logger logger.Logger) (*AliyunALBDeployer, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}

	if logger == nil {
		return nil, errors.New("logger is nil")
	}

	client, err := createSdkClient(config.AccessKeyId, config.AccessKeySecret, config.Region)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk client")
	}

	aliyunCasRegion := config.Region
	if aliyunCasRegion != "" {
		// 阿里云 CAS 服务接入点是独立于 ALB 服务的
		// 国内版固定接入点：华东一杭州
		// 国际版固定接入点：亚太东南一新加坡
		if !strings.HasPrefix(aliyunCasRegion, "cn-") {
			aliyunCasRegion = "ap-southeast-1"
		} else {
			aliyunCasRegion = "cn-hangzhou"
		}
	}
	uploader, err := providerCas.New(&providerCas.AliyunCASUploaderConfig{
		AccessKeyId:     config.AccessKeyId,
		AccessKeySecret: config.AccessKeySecret,
		Region:          aliyunCasRegion,
	})
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create ssl uploader")
	}

	return &AliyunALBDeployer{
		logger:      logger,
		config:      config,
		sdkClient:   client,
		sslUploader: uploader,
	}, nil
}

func (d *AliyunALBDeployer) Deploy(ctx context.Context, certPem string, privkeyPem string) (*deployer.DeployResult, error) {
	// 上传证书到 CAS
	upres, err := d.sslUploader.Upload(ctx, certPem, privkeyPem)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to upload certificate file")
	}

	d.logger.Logt("certificate file uploaded", upres)

	// 根据部署资源类型决定部署方式
	switch d.config.ResourceType {
	case DEPLOY_RESOURCE_LOADBALANCER:
		if err := d.deployToLoadbalancer(ctx, upres.CertId); err != nil {
			return nil, err
		}

	case DEPLOY_RESOURCE_LISTENER:
		if err := d.deployToListener(ctx, upres.CertId); err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("unsupported resource type: %s", d.config.ResourceType)
	}

	return &deployer.DeployResult{}, nil
}

func (d *AliyunALBDeployer) deployToLoadbalancer(ctx context.Context, cloudCertId string) error {
	if d.config.LoadbalancerId == "" {
		return errors.New("config `loadbalancerId` is required")
	}

	listenerIds := make([]string, 0)

	// 查询负载均衡实例的详细信息
	// REF: https://help.aliyun.com/zh/slb/application-load-balancer/developer-reference/api-alb-2020-06-16-getloadbalancerattribute
	getLoadBalancerAttributeReq := &aliyunAlb.GetLoadBalancerAttributeRequest{
		LoadBalancerId: tea.String(d.config.LoadbalancerId),
	}
	getLoadBalancerAttributeResp, err := d.sdkClient.GetLoadBalancerAttribute(getLoadBalancerAttributeReq)
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'alb.GetLoadBalancerAttribute'")
	}

	d.logger.Logt("已查询到 ALB 负载均衡实例", getLoadBalancerAttributeResp)

	// 查询 HTTPS 监听列表
	// REF: https://help.aliyun.com/zh/slb/application-load-balancer/developer-reference/api-alb-2020-06-16-listlisteners
	listListenersPage := 1
	listListenersLimit := int32(100)
	var listListenersToken *string = nil
	for {
		listListenersReq := &aliyunAlb.ListListenersRequest{
			MaxResults:       tea.Int32(listListenersLimit),
			NextToken:        listListenersToken,
			LoadBalancerIds:  []*string{tea.String(d.config.LoadbalancerId)},
			ListenerProtocol: tea.String("HTTPS"),
		}
		listListenersResp, err := d.sdkClient.ListListeners(listListenersReq)
		if err != nil {
			return xerrors.Wrap(err, "failed to execute sdk request 'alb.ListListeners'")
		}

		if listListenersResp.Body.Listeners != nil {
			for _, listener := range listListenersResp.Body.Listeners {
				listenerIds = append(listenerIds, *listener.ListenerId)
			}
		}

		if len(listListenersResp.Body.Listeners) == 0 || listListenersResp.Body.NextToken == nil {
			break
		} else {
			listListenersToken = listListenersResp.Body.NextToken
			listListenersPage += 1
		}
	}

	d.logger.Logt("已查询到 ALB 负载均衡实例下的全部 HTTPS 监听", listenerIds)

	// 查询 QUIC 监听列表
	// REF: https://help.aliyun.com/zh/slb/application-load-balancer/developer-reference/api-alb-2020-06-16-listlisteners
	listListenersPage = 1
	listListenersToken = nil
	for {
		listListenersReq := &aliyunAlb.ListListenersRequest{
			MaxResults:       tea.Int32(listListenersLimit),
			NextToken:        listListenersToken,
			LoadBalancerIds:  []*string{tea.String(d.config.LoadbalancerId)},
			ListenerProtocol: tea.String("QUIC"),
		}
		listListenersResp, err := d.sdkClient.ListListeners(listListenersReq)
		if err != nil {
			return xerrors.Wrap(err, "failed to execute sdk request 'alb.ListListeners'")
		}

		if listListenersResp.Body.Listeners != nil {
			for _, listener := range listListenersResp.Body.Listeners {
				listenerIds = append(listenerIds, *listener.ListenerId)
			}
		}

		if len(listListenersResp.Body.Listeners) == 0 || listListenersResp.Body.NextToken == nil {
			break
		} else {
			listListenersToken = listListenersResp.Body.NextToken
			listListenersPage += 1
		}
	}

	d.logger.Logt("已查询到 ALB 负载均衡实例下的全部 QUIC 监听", listenerIds)

	// 批量更新监听证书
	var errs []error
	for _, listenerId := range listenerIds {
		if err := d.updateListenerCertificate(ctx, listenerId, cloudCertId); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func (d *AliyunALBDeployer) deployToListener(ctx context.Context, cloudCertId string) error {
	if d.config.ListenerId == "" {
		return errors.New("config `listenerId` is required")
	}

	// 更新监听
	if err := d.updateListenerCertificate(ctx, d.config.ListenerId, cloudCertId); err != nil {
		return err
	}

	return nil
}

func (d *AliyunALBDeployer) updateListenerCertificate(ctx context.Context, cloudListenerId string, cloudCertId string) error {
	// 查询监听的属性
	// REF: https://help.aliyun.com/zh/slb/application-load-balancer/developer-reference/api-alb-2020-06-16-getlistenerattribute
	getListenerAttributeReq := &aliyunAlb.GetListenerAttributeRequest{
		ListenerId: tea.String(cloudListenerId),
	}
	getListenerAttributeResp, err := d.sdkClient.GetListenerAttribute(getListenerAttributeReq)
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'alb.GetListenerAttribute'")
	}

	d.logger.Logt("已查询到 ALB 监听配置", getListenerAttributeResp)

	// 修改监听的属性
	// REF: https://help.aliyun.com/zh/slb/application-load-balancer/developer-reference/api-alb-2020-06-16-updatelistenerattribute
	updateListenerAttributeReq := &aliyunAlb.UpdateListenerAttributeRequest{
		ListenerId: tea.String(cloudListenerId),
		Certificates: []*aliyunAlb.UpdateListenerAttributeRequestCertificates{{
			CertificateId: tea.String(cloudCertId),
		}},
	}
	updateListenerAttributeResp, err := d.sdkClient.UpdateListenerAttribute(updateListenerAttributeReq)
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'alb.UpdateListenerAttribute'")
	}

	d.logger.Logt("已更新 ALB 监听配置", updateListenerAttributeResp)

	// TODO: #347

	return nil
}

func createSdkClient(accessKeyId, accessKeySecret, region string) (*aliyunAlb.Client, error) {
	if region == "" {
		region = "cn-hangzhou"
	}

	// 接入点一览 https://www.alibabacloud.com/help/zh/slb/application-load-balancer/developer-reference/api-alb-2020-06-16-endpoint
	var endpoint string
	switch region {
	case "cn-hangzhou-finance":
		endpoint = "alb.cn-hangzhou.aliyuncs.com"
	default:
		endpoint = fmt.Sprintf("alb.%s.aliyuncs.com", region)
	}

	config := &aliyunOpen.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
		Endpoint:        tea.String(endpoint),
	}

	client, err := aliyunAlb.NewClient(config)
	if err != nil {
		return nil, err
	}

	return client, nil
}
