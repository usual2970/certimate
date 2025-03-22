package aliyunnlb

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	aliopen "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	alinlb "github.com/alibabacloud-go/nlb-20220430/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	xerrors "github.com/pkg/errors"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	uploadersp "github.com/usual2970/certimate/internal/pkg/core/uploader/providers/aliyun-cas"
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
	// 部署资源类型为 [RESOURCE_TYPE_LOADBALANCER] 时必填。
	LoadbalancerId string `json:"loadbalancerId,omitempty"`
	// 负载均衡监听 ID。
	// 部署资源类型为 [RESOURCE_TYPE_LISTENER] 时必填。
	ListenerId string `json:"listenerId,omitempty"`
}

type DeployerProvider struct {
	config      *DeployerConfig
	logger      *slog.Logger
	sdkClient   *alinlb.Client
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

	uploader, err := createSslUploader(config.AccessKeyId, config.AccessKeySecret, config.Region)
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
	// 上传证书到 CAS
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
	// REF: https://help.aliyun.com/zh/slb/network-load-balancer/developer-reference/api-nlb-2022-04-30-getloadbalancerattribute
	getLoadBalancerAttributeReq := &alinlb.GetLoadBalancerAttributeRequest{
		LoadBalancerId: tea.String(d.config.LoadbalancerId),
	}
	getLoadBalancerAttributeResp, err := d.sdkClient.GetLoadBalancerAttribute(getLoadBalancerAttributeReq)
	d.logger.Debug("sdk request 'nlb.GetLoadBalancerAttribute'", slog.Any("request", getLoadBalancerAttributeReq), slog.Any("response", getLoadBalancerAttributeResp))
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'nlb.GetLoadBalancerAttribute'")
	}

	// 查询 TCPSSL 监听列表
	// REF: https://help.aliyun.com/zh/slb/network-load-balancer/developer-reference/api-nlb-2022-04-30-listlisteners
	listenerIds := make([]string, 0)
	listListenersLimit := int32(100)
	var listListenersToken *string = nil
	for {
		listListenersReq := &alinlb.ListListenersRequest{
			MaxResults:       tea.Int32(listListenersLimit),
			NextToken:        listListenersToken,
			LoadBalancerIds:  []*string{tea.String(d.config.LoadbalancerId)},
			ListenerProtocol: tea.String("TCPSSL"),
		}
		listListenersResp, err := d.sdkClient.ListListeners(listListenersReq)
		d.logger.Debug("sdk request 'nlb.ListListeners'", slog.Any("request", listListenersReq), slog.Any("response", listListenersResp))
		if err != nil {
			return xerrors.Wrap(err, "failed to execute sdk request 'nlb.ListListeners'")
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
		}
	}

	// 遍历更新监听证书
	if len(listenerIds) == 0 {
		d.logger.Info("no nlb listeners to deploy")
	} else {
		d.logger.Info("found tcpssl listeners to deploy", slog.Any("listenerIds", listenerIds))
		var errs []error

		for _, listenerId := range listenerIds {
			if err := d.updateListenerCertificate(ctx, listenerId, cloudCertId); err != nil {
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
	if d.config.ListenerId == "" {
		return errors.New("config `listenerId` is required")
	}

	// 更新监听
	if err := d.updateListenerCertificate(ctx, d.config.ListenerId, cloudCertId); err != nil {
		return err
	}

	return nil
}

func (d *DeployerProvider) updateListenerCertificate(ctx context.Context, cloudListenerId string, cloudCertId string) error {
	// 查询监听的属性
	// REF: https://help.aliyun.com/zh/slb/network-load-balancer/developer-reference/api-nlb-2022-04-30-getlistenerattribute
	getListenerAttributeReq := &alinlb.GetListenerAttributeRequest{
		ListenerId: tea.String(cloudListenerId),
	}
	getListenerAttributeResp, err := d.sdkClient.GetListenerAttribute(getListenerAttributeReq)
	d.logger.Debug("sdk request 'nlb.GetListenerAttribute'", slog.Any("request", getListenerAttributeReq), slog.Any("response", getListenerAttributeResp))
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'nlb.GetListenerAttribute'")
	}

	// 修改监听的属性
	// REF: https://help.aliyun.com/zh/slb/network-load-balancer/developer-reference/api-nlb-2022-04-30-updatelistenerattribute
	updateListenerAttributeReq := &alinlb.UpdateListenerAttributeRequest{
		ListenerId:     tea.String(cloudListenerId),
		CertificateIds: []*string{tea.String(cloudCertId)},
	}
	updateListenerAttributeResp, err := d.sdkClient.UpdateListenerAttribute(updateListenerAttributeReq)
	d.logger.Debug("sdk request 'nlb.UpdateListenerAttribute'", slog.Any("request", updateListenerAttributeReq), slog.Any("response", updateListenerAttributeResp))
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'nlb.UpdateListenerAttribute'")
	}

	return nil
}

func createSdkClient(accessKeyId, accessKeySecret, region string) (*alinlb.Client, error) {
	// 接入点一览 https://api.aliyun.com/product/Nlb
	var endpoint string
	switch region {
	default:
		endpoint = fmt.Sprintf("nlb.%s.aliyuncs.com", region)
	}

	config := &aliopen.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
		Endpoint:        tea.String(endpoint),
	}

	client, err := alinlb.NewClient(config)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func createSslUploader(accessKeyId, accessKeySecret, region string) (uploader.Uploader, error) {
	casRegion := region
	if casRegion != "" {
		// 阿里云 CAS 服务接入点是独立于 NLB 服务的
		// 国内版固定接入点：华东一杭州
		// 国际版固定接入点：亚太东南一新加坡
		if casRegion != "" && !strings.HasPrefix(casRegion, "cn-") {
			casRegion = "ap-southeast-1"
		} else {
			casRegion = "cn-hangzhou"
		}
	}

	uploader, err := uploadersp.NewUploader(&uploadersp.UploaderConfig{
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
		Region:          casRegion,
	})
	return uploader, err
}
