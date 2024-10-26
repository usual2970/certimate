package deployer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	nlb20220430 "github.com/alibabacloud-go/nlb-20220430/v2/client"
	"github.com/alibabacloud-go/tea/tea"

	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
)

type AliyunNLBDeployer struct {
	option *DeployerOption
	infos  []string

	sdkClient   *nlb20220430.Client
	sslUploader uploader.Uploader
}

func NewAliyunNLBDeployer(option *DeployerOption) (Deployer, error) {
	access := &domain.AliyunAccess{}
	json.Unmarshal([]byte(option.Access), access)

	client, err := (&AliyunNLBDeployer{}).createSdkClient(
		access.AccessKeyId,
		access.AccessKeySecret,
		option.DeployConfig.GetConfigAsString("region"),
	)
	if err != nil {
		return nil, err
	}

	uploader, err := uploader.NewAliyunCASUploader(&uploader.AliyunCASUploaderConfig{
		AccessKeyId:     access.AccessKeyId,
		AccessKeySecret: access.AccessKeySecret,
		Region:          option.DeployConfig.GetConfigAsString("region"),
	})
	if err != nil {
		return nil, err
	}

	return &AliyunNLBDeployer{
		option:      option,
		infos:       make([]string, 0),
		sdkClient:   client,
		sslUploader: uploader,
	}, nil
}

func (d *AliyunNLBDeployer) GetID() string {
	return fmt.Sprintf("%s-%s", d.option.AccessRecord.GetString("name"), d.option.AccessRecord.Id)
}

func (d *AliyunNLBDeployer) GetInfo() []string {
	return d.infos
}

func (d *AliyunNLBDeployer) Deploy(ctx context.Context) error {
	switch d.option.DeployConfig.GetConfigAsString("resourceType") {
	case "loadbalancer":
		if err := d.deployToLoadbalancer(ctx); err != nil {
			return err
		}
	case "listener":
		if err := d.deployToListener(ctx); err != nil {
			return err
		}
	default:
		return errors.New("unsupported resource type")
	}

	return nil
}

func (d *AliyunNLBDeployer) createSdkClient(accessKeyId, accessKeySecret, region string) (*nlb20220430.Client, error) {
	if region == "" {
		region = "cn-hangzhou" // NLB 服务默认区域：华东一杭州
	}

	aConfig := &openapi.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
	}

	var endpoint string
	switch region {
	case "cn-hangzhou-finance":
		endpoint = "nlb.cn-hangzhou.aliyuncs.com"
	default:
		endpoint = fmt.Sprintf("nlb.%s.aliyuncs.com", region)
	}
	aConfig.Endpoint = tea.String(endpoint)

	client, err := nlb20220430.NewClient(aConfig)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (d *AliyunNLBDeployer) deployToLoadbalancer(ctx context.Context) error {
	aliLoadbalancerId := d.option.DeployConfig.GetConfigAsString("loadbalancerId")
	if aliLoadbalancerId == "" {
		return errors.New("`loadbalancerId` is required")
	}

	// 查询负载均衡实例的详细信息
	// REF: https://help.aliyun.com/zh/slb/network-load-balancer/developer-reference/api-nlb-2022-04-30-getloadbalancerattribute
	getLoadBalancerAttributeReq := &nlb20220430.GetLoadBalancerAttributeRequest{
		LoadBalancerId: tea.String(aliLoadbalancerId),
	}
	getLoadBalancerAttributeResp, err := d.sdkClient.GetLoadBalancerAttribute(getLoadBalancerAttributeReq)
	if err != nil {
		return fmt.Errorf("failed to execute sdk request 'nlb.GetLoadBalancerAttribute': %w", err)
	}

	d.infos = append(d.infos, toStr("已查询到 NLB 负载均衡实例", getLoadBalancerAttributeResp))

	// 查询监听列表
	// REF: https://help.aliyun.com/zh/slb/network-load-balancer/developer-reference/api-nlb-2022-04-30-listlisteners
	aliListenerIds := make([]string, 0)
	listListenersPage := 1
	listListenersLimit := int32(100)
	var listListenersToken *string = nil
	for {
		listListenersReq := &nlb20220430.ListListenersRequest{
			MaxResults:       tea.Int32(listListenersLimit),
			NextToken:        listListenersToken,
			LoadBalancerIds:  []*string{tea.String(aliLoadbalancerId)},
			ListenerProtocol: tea.String("TCPSSL"),
		}
		listListenersResp, err := d.sdkClient.ListListeners(listListenersReq)
		if err != nil {
			return fmt.Errorf("failed to execute sdk request 'nlb.ListListeners': %w", err)
		}

		if listListenersResp.Body.Listeners != nil {
			for _, listener := range listListenersResp.Body.Listeners {
				aliListenerIds = append(aliListenerIds, *listener.ListenerId)
			}
		}

		if listListenersResp.Body.NextToken == nil {
			break
		} else {
			listListenersToken = listListenersResp.Body.NextToken
			listListenersPage += 1
		}
	}

	d.infos = append(d.infos, toStr("已查询到 NLB 负载均衡实例下的全部 TCPSSL 监听", aliListenerIds))

	// 上传证书到 SSL
	uploadResult, err := d.sslUploader.Upload(ctx, d.option.Certificate.Certificate, d.option.Certificate.PrivateKey)
	if err != nil {
		return err
	}

	d.infos = append(d.infos, toStr("已上传证书", uploadResult))

	// 批量更新监听证书
	var errs []error
	for _, aliListenerId := range aliListenerIds {
		if err := d.updateListenerCertificate(ctx, aliListenerId, uploadResult.CertId); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func (d *AliyunNLBDeployer) deployToListener(ctx context.Context) error {
	aliListenerId := d.option.DeployConfig.GetConfigAsString("listenerId")
	if aliListenerId == "" {
		return errors.New("`listenerId` is required")
	}

	// 上传证书到 SSL
	uploadResult, err := d.sslUploader.Upload(ctx, d.option.Certificate.Certificate, d.option.Certificate.PrivateKey)
	if err != nil {
		return err
	}

	d.infos = append(d.infos, toStr("已上传证书", uploadResult))

	// 更新监听
	if err := d.updateListenerCertificate(ctx, aliListenerId, uploadResult.CertId); err != nil {
		return err
	}

	return nil
}

func (d *AliyunNLBDeployer) updateListenerCertificate(ctx context.Context, aliListenerId string, aliCertId string) error {
	// 查询监听的属性
	// REF: https://help.aliyun.com/zh/slb/network-load-balancer/developer-reference/api-nlb-2022-04-30-getlistenerattribute
	getListenerAttributeReq := &nlb20220430.GetListenerAttributeRequest{
		ListenerId: tea.String(aliListenerId),
	}
	getListenerAttributeResp, err := d.sdkClient.GetListenerAttribute(getListenerAttributeReq)
	if err != nil {
		return fmt.Errorf("failed to execute sdk request 'nlb.GetListenerAttribute': %w", err)
	}

	d.infos = append(d.infos, toStr("已查询到 NLB 监听配置", getListenerAttributeResp))

	// 修改监听的属性
	// REF: https://help.aliyun.com/zh/slb/network-load-balancer/developer-reference/api-nlb-2022-04-30-updatelistenerattribute
	updateListenerAttributeReq := &nlb20220430.UpdateListenerAttributeRequest{
		ListenerId:     tea.String(aliListenerId),
		CertificateIds: []*string{tea.String(aliCertId)},
	}
	updateListenerAttributeResp, err := d.sdkClient.UpdateListenerAttribute(updateListenerAttributeReq)
	if err != nil {
		return fmt.Errorf("failed to execute sdk request 'nlb.UpdateListenerAttribute': %w", err)
	}

	d.infos = append(d.infos, toStr("已更新 NLB 监听配置", updateListenerAttributeResp))

	return nil
}
