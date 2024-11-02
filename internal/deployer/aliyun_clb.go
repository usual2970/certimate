package deployer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	aliyunOpen "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	aliyunSlb "github.com/alibabacloud-go/slb-20140515/v4/client"
	"github.com/alibabacloud-go/tea/tea"
	xerrors "github.com/pkg/errors"

	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	uploaderAliyunSlb "github.com/usual2970/certimate/internal/pkg/core/uploader/providers/aliyun-slb"
)

type AliyunCLBDeployer struct {
	option *DeployerOption
	infos  []string

	sdkClient   *aliyunSlb.Client
	sslUploader uploader.Uploader
}

func NewAliyunCLBDeployer(option *DeployerOption) (Deployer, error) {
	access := &domain.AliyunAccess{}
	if err := json.Unmarshal([]byte(option.Access), access); err != nil {
		return nil, xerrors.Wrap(err, "failed to get access")
	}

	client, err := (&AliyunCLBDeployer{}).createSdkClient(
		access.AccessKeyId,
		access.AccessKeySecret,
		option.DeployConfig.GetConfigAsString("region"),
	)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk client")
	}

	uploader, err := uploaderAliyunSlb.New(&uploaderAliyunSlb.AliyunSLBUploaderConfig{
		AccessKeyId:     access.AccessKeyId,
		AccessKeySecret: access.AccessKeySecret,
		Region:          option.DeployConfig.GetConfigAsString("region"),
	})
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create ssl uploader")
	}

	return &AliyunCLBDeployer{
		option:      option,
		infos:       make([]string, 0),
		sdkClient:   client,
		sslUploader: uploader,
	}, nil
}

func (d *AliyunCLBDeployer) GetID() string {
	return fmt.Sprintf("%s-%s", d.option.AccessRecord.GetString("name"), d.option.AccessRecord.Id)
}

func (d *AliyunCLBDeployer) GetInfos() []string {
	return d.infos
}

func (d *AliyunCLBDeployer) Deploy(ctx context.Context) error {
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

func (d *AliyunCLBDeployer) createSdkClient(accessKeyId, accessKeySecret, region string) (*aliyunSlb.Client, error) {
	if region == "" {
		region = "cn-hangzhou" // CLB(SLB) 服务默认区域：华东一杭州
	}

	aConfig := &aliyunOpen.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
	}

	var endpoint string
	switch region {
	case "cn-hangzhou":
	case "cn-hangzhou-finance":
	case "cn-shanghai-finance-1":
	case "cn-shenzhen-finance-1":
		endpoint = "slb.aliyuncs.com"
	default:
		endpoint = fmt.Sprintf("slb.%s.aliyuncs.com", region)
	}
	aConfig.Endpoint = tea.String(endpoint)

	client, err := aliyunSlb.NewClient(aConfig)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (d *AliyunCLBDeployer) deployToLoadbalancer(ctx context.Context) error {
	aliLoadbalancerId := d.option.DeployConfig.GetConfigAsString("loadbalancerId")
	aliListenerPorts := make([]int32, 0)
	if aliLoadbalancerId == "" {
		return errors.New("`loadbalancerId` is required")
	}

	// 查询负载均衡实例的详细信息
	// REF: https://help.aliyun.com/zh/slb/classic-load-balancer/developer-reference/api-slb-2014-05-15-describeloadbalancerattribute
	describeLoadBalancerAttributeReq := &aliyunSlb.DescribeLoadBalancerAttributeRequest{
		RegionId:       tea.String(d.option.DeployConfig.GetConfigAsString("region")),
		LoadBalancerId: tea.String(aliLoadbalancerId),
	}
	describeLoadBalancerAttributeResp, err := d.sdkClient.DescribeLoadBalancerAttribute(describeLoadBalancerAttributeReq)
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'slb.DescribeLoadBalancerAttribute'")
	}

	d.infos = append(d.infos, toStr("已查询到 CLB 负载均衡实例", describeLoadBalancerAttributeResp))

	// 查询 HTTPS 监听列表
	// REF: https://help.aliyun.com/zh/slb/classic-load-balancer/developer-reference/api-slb-2014-05-15-describeloadbalancerlisteners
	listListenersPage := 1
	listListenersLimit := int32(100)
	var listListenersToken *string = nil
	for {
		describeLoadBalancerListenersReq := &aliyunSlb.DescribeLoadBalancerListenersRequest{
			RegionId:         tea.String(d.option.DeployConfig.GetConfigAsString("region")),
			MaxResults:       tea.Int32(listListenersLimit),
			NextToken:        listListenersToken,
			LoadBalancerId:   []*string{tea.String(aliLoadbalancerId)},
			ListenerProtocol: tea.String("https"),
		}
		describeLoadBalancerListenersResp, err := d.sdkClient.DescribeLoadBalancerListeners(describeLoadBalancerListenersReq)
		if err != nil {
			return xerrors.Wrap(err, "failed to execute sdk request 'slb.DescribeLoadBalancerListeners'")
		}

		if describeLoadBalancerListenersResp.Body.Listeners != nil {
			for _, listener := range describeLoadBalancerListenersResp.Body.Listeners {
				aliListenerPorts = append(aliListenerPorts, *listener.ListenerPort)
			}
		}

		if describeLoadBalancerListenersResp.Body.NextToken == nil {
			break
		} else {
			listListenersToken = describeLoadBalancerListenersResp.Body.NextToken
			listListenersPage += 1
		}
	}

	d.infos = append(d.infos, toStr("已查询到 CLB 负载均衡实例下的全部 HTTPS 监听", aliListenerPorts))

	// 上传证书到 SLB
	upres, err := d.sslUploader.Upload(ctx, d.option.Certificate.Certificate, d.option.Certificate.PrivateKey)
	if err != nil {
		return err
	}

	d.infos = append(d.infos, toStr("已上传证书", upres))

	// 批量更新监听证书
	var errs []error
	for _, aliListenerPort := range aliListenerPorts {
		if err := d.updateListenerCertificate(ctx, aliLoadbalancerId, aliListenerPort, upres.CertId); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func (d *AliyunCLBDeployer) deployToListener(ctx context.Context) error {
	aliLoadbalancerId := d.option.DeployConfig.GetConfigAsString("loadbalancerId")
	if aliLoadbalancerId == "" {
		return errors.New("`loadbalancerId` is required")
	}

	aliListenerPort := d.option.DeployConfig.GetConfigAsInt32("listenerPort")
	if aliListenerPort == 0 {
		return errors.New("`listenerPort` is required")
	}

	// 上传证书到 SLB
	upres, err := d.sslUploader.Upload(ctx, d.option.Certificate.Certificate, d.option.Certificate.PrivateKey)
	if err != nil {
		return err
	}

	d.infos = append(d.infos, toStr("已上传证书", upres))

	// 更新监听
	if err := d.updateListenerCertificate(ctx, aliLoadbalancerId, aliListenerPort, upres.CertId); err != nil {
		return err
	}

	return nil
}

func (d *AliyunCLBDeployer) updateListenerCertificate(ctx context.Context, aliLoadbalancerId string, aliListenerPort int32, aliCertId string) error {
	// 查询监听配置
	// REF: https://help.aliyun.com/zh/slb/classic-load-balancer/developer-reference/api-slb-2014-05-15-describeloadbalancerhttpslistenerattribute
	describeLoadBalancerHTTPSListenerAttributeReq := &aliyunSlb.DescribeLoadBalancerHTTPSListenerAttributeRequest{
		LoadBalancerId: tea.String(aliLoadbalancerId),
		ListenerPort:   tea.Int32(aliListenerPort),
	}
	describeLoadBalancerHTTPSListenerAttributeResp, err := d.sdkClient.DescribeLoadBalancerHTTPSListenerAttribute(describeLoadBalancerHTTPSListenerAttributeReq)
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'slb.DescribeLoadBalancerHTTPSListenerAttribute'")
	}

	d.infos = append(d.infos, toStr("已查询到 CLB HTTPS 监听配置", describeLoadBalancerHTTPSListenerAttributeResp))

	// 查询扩展域名
	// REF: https://help.aliyun.com/zh/slb/classic-load-balancer/developer-reference/api-slb-2014-05-15-describedomainextensions
	describeDomainExtensionsReq := &aliyunSlb.DescribeDomainExtensionsRequest{
		RegionId:       tea.String(d.option.DeployConfig.GetConfigAsString("region")),
		LoadBalancerId: tea.String(aliLoadbalancerId),
		ListenerPort:   tea.Int32(aliListenerPort),
	}
	describeDomainExtensionsResp, err := d.sdkClient.DescribeDomainExtensions(describeDomainExtensionsReq)
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'slb.DescribeDomainExtensions'")
	}

	d.infos = append(d.infos, toStr("已查询到 CLB 扩展域名", describeDomainExtensionsResp))

	// 遍历修改扩展域名
	// REF: https://help.aliyun.com/zh/slb/classic-load-balancer/developer-reference/api-slb-2014-05-15-setdomainextensionattribute
	//
	// 这里仅修改跟被替换证书一致的扩展域名
	if describeDomainExtensionsResp.Body.DomainExtensions == nil && describeDomainExtensionsResp.Body.DomainExtensions.DomainExtension == nil {
		for _, domainExtension := range describeDomainExtensionsResp.Body.DomainExtensions.DomainExtension {
			if *domainExtension.ServerCertificateId == *describeLoadBalancerHTTPSListenerAttributeResp.Body.ServerCertificateId {
				break
			}

			setDomainExtensionAttributeReq := &aliyunSlb.SetDomainExtensionAttributeRequest{
				RegionId:            tea.String(d.option.DeployConfig.GetConfigAsString("region")),
				DomainExtensionId:   tea.String(*domainExtension.DomainExtensionId),
				ServerCertificateId: tea.String(aliCertId),
			}
			_, err := d.sdkClient.SetDomainExtensionAttribute(setDomainExtensionAttributeReq)
			if err != nil {
				return xerrors.Wrap(err, "failed to execute sdk request 'slb.SetDomainExtensionAttribute'")
			}
		}
	}

	// 修改监听配置
	// REF: https://help.aliyun.com/zh/slb/classic-load-balancer/developer-reference/api-slb-2014-05-15-setloadbalancerhttpslistenerattribute
	//
	// 注意修改监听配置要放在修改扩展域名之后
	setLoadBalancerHTTPSListenerAttributeReq := &aliyunSlb.SetLoadBalancerHTTPSListenerAttributeRequest{
		RegionId:            tea.String(d.option.DeployConfig.GetConfigAsString("region")),
		LoadBalancerId:      tea.String(aliLoadbalancerId),
		ListenerPort:        tea.Int32(aliListenerPort),
		ServerCertificateId: tea.String(aliCertId),
	}
	setLoadBalancerHTTPSListenerAttributeResp, err := d.sdkClient.SetLoadBalancerHTTPSListenerAttribute(setLoadBalancerHTTPSListenerAttributeReq)
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'slb.SetLoadBalancerHTTPSListenerAttribute'")
	}

	d.infos = append(d.infos, toStr("已更新 CLB HTTPS 监听配置", setLoadBalancerHTTPSListenerAttributeResp))

	return nil
}
