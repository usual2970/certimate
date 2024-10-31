package deployer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	xerrors "github.com/pkg/errors"
	tcClb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tcSsl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"

	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
)

type TencentCLBDeployer struct {
	option *DeployerOption
	infos  []string

	sdkClients  *tencentCLBDeployerSdkClients
	sslUploader uploader.Uploader
}

type tencentCLBDeployerSdkClients struct {
	ssl *tcSsl.Client
	clb *tcClb.Client
}

func NewTencentCLBDeployer(option *DeployerOption) (Deployer, error) {
	access := &domain.TencentAccess{}
	if err := json.Unmarshal([]byte(option.Access), access); err != nil {
		return nil, xerrors.Wrap(err, "failed to get access")
	}

	clients, err := (&TencentCLBDeployer{}).createSdkClients(
		access.SecretId,
		access.SecretKey,
		option.DeployConfig.GetConfigAsString("region"),
	)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk clients")
	}

	uploader, err := uploader.NewTencentCloudSSLUploader(&uploader.TencentCloudSSLUploaderConfig{
		SecretId:  access.SecretId,
		SecretKey: access.SecretKey,
	})
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create ssl uploader")
	}

	return &TencentCLBDeployer{
		option:      option,
		infos:       make([]string, 0),
		sdkClients:  clients,
		sslUploader: uploader,
	}, nil
}

func (d *TencentCLBDeployer) GetID() string {
	return fmt.Sprintf("%s-%s", d.option.AccessRecord.GetString("name"), d.option.AccessRecord.Id)
}

func (d *TencentCLBDeployer) GetInfos() []string {
	return d.infos
}

func (d *TencentCLBDeployer) Deploy(ctx context.Context) error {
	// TODO: 直接部署方式

	switch d.option.DeployConfig.GetConfigAsString("resourceType") {
	case "ssl-deploy":
		// 通过 SSL 服务部署到云资源实例
		err := d.deployToInstanceUseSsl(ctx)
		if err != nil {
			return err
		}
	case "loadbalancer":
		// 部署到指定负载均衡器
		if err := d.deployToLoadbalancer(ctx); err != nil {
			return err
		}
	case "listener":
		// 部署到指定监听器
		if err := d.deployToListener(ctx); err != nil {
			return err
		}
	case "ruledomain":
		// 部署到指定七层监听转发规则域名
		if err := d.deployToRuleDomain(ctx); err != nil {
			return err
		}
	default:
		return errors.New("unsupported resource type")
	}

	return nil
}

func (d *TencentCLBDeployer) createSdkClients(secretId, secretKey, region string) (*tencentCLBDeployerSdkClients, error) {
	credential := common.NewCredential(secretId, secretKey)

	sslClient, err := tcSsl.NewClient(credential, "", profile.NewClientProfile())
	if err != nil {
		return nil, err
	}

	clbClient, err := tcClb.NewClient(credential, region, profile.NewClientProfile())
	if err != nil {
		return nil, err
	}

	return &tencentCLBDeployerSdkClients{
		ssl: sslClient,
		clb: clbClient,
	}, nil
}

func (d *TencentCLBDeployer) deployToInstanceUseSsl(ctx context.Context) error {
	tcLoadbalancerId := d.option.DeployConfig.GetConfigAsString("loadbalancerId")
	tcListenerId := d.option.DeployConfig.GetConfigAsString("listenerId")
	tcDomain := d.option.DeployConfig.GetConfigAsString("domain")
	if tcLoadbalancerId == "" {
		return errors.New("`loadbalancerId` is required")
	}
	if tcListenerId == "" {
		return errors.New("`listenerId` is required")
	}

	// 上传证书到 SSL
	upres, err := d.sslUploader.Upload(ctx, d.option.Certificate.Certificate, d.option.Certificate.PrivateKey)
	if err != nil {
		return err
	}

	d.infos = append(d.infos, toStr("已上传证书", upres))

	// 证书部署到 CLB 实例
	// REF: https://cloud.tencent.com/document/product/400/91667
	deployCertificateInstanceReq := tcSsl.NewDeployCertificateInstanceRequest()
	deployCertificateInstanceReq.CertificateId = common.StringPtr(upres.CertId)
	deployCertificateInstanceReq.ResourceType = common.StringPtr("clb")
	deployCertificateInstanceReq.Status = common.Int64Ptr(1)
	if tcDomain == "" {
		// 未开启 SNI，只需指定到监听器
		deployCertificateInstanceReq.InstanceIdList = common.StringPtrs([]string{fmt.Sprintf("%s|%s", tcLoadbalancerId, tcListenerId)})
	} else {
		// 开启 SNI，需指定到域名（支持泛域名）
		deployCertificateInstanceReq.InstanceIdList = common.StringPtrs([]string{fmt.Sprintf("%s|%s|%s", tcLoadbalancerId, tcListenerId, tcDomain)})
	}
	deployCertificateInstanceResp, err := d.sdkClients.ssl.DeployCertificateInstance(deployCertificateInstanceReq)
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'ssl.DeployCertificateInstance'")
	}

	d.infos = append(d.infos, toStr("已部署证书到云资源实例", deployCertificateInstanceResp.Response))

	return nil
}

func (d *TencentCLBDeployer) deployToLoadbalancer(ctx context.Context) error {
	tcLoadbalancerId := d.option.DeployConfig.GetConfigAsString("loadbalancerId")
	tcListenerIds := make([]string, 0)
	if tcLoadbalancerId == "" {
		return errors.New("`loadbalancerId` is required")
	}

	// 查询负载均衡器详细信息
	// REF: https://cloud.tencent.com/document/api/214/46916
	describeLoadBalancersDetailReq := tcClb.NewDescribeLoadBalancersDetailRequest()
	describeLoadBalancersDetailResp, err := d.sdkClients.clb.DescribeLoadBalancersDetail(describeLoadBalancersDetailReq)
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'clb.DescribeLoadBalancersDetail'")
	}

	d.infos = append(d.infos, toStr("已查询到负载均衡详细信息", describeLoadBalancersDetailResp))

	// 查询监听器列表
	// REF: https://cloud.tencent.com/document/api/214/30686
	describeListenersReq := tcClb.NewDescribeListenersRequest()
	describeListenersReq.LoadBalancerId = common.StringPtr(tcLoadbalancerId)
	describeListenersResp, err := d.sdkClients.clb.DescribeListeners(describeListenersReq)
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'clb.DescribeListeners'")
	} else {
		if describeListenersResp.Response.Listeners != nil {
			for _, listener := range describeListenersResp.Response.Listeners {
				if listener.Protocol == nil || (*listener.Protocol != "HTTPS" && *listener.Protocol != "TCP_SSL" && *listener.Protocol != "QUIC") {
					continue
				}

				tcListenerIds = append(tcListenerIds, *listener.ListenerId)
			}
		}
	}

	d.infos = append(d.infos, toStr("已查询到负载均衡器下的监听器", tcListenerIds))

	// 上传证书到 SCM
	upres, err := d.sslUploader.Upload(ctx, d.option.Certificate.Certificate, d.option.Certificate.PrivateKey)
	if err != nil {
		return err
	}

	d.infos = append(d.infos, toStr("已上传证书", upres))

	// 批量更新监听器证书
	var errs []error
	for _, tcListenerId := range tcListenerIds {
		if err := d.modifyListenerCertificate(ctx, tcLoadbalancerId, tcListenerId, upres.CertId); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func (d *TencentCLBDeployer) deployToListener(ctx context.Context) error {
	tcLoadbalancerId := d.option.DeployConfig.GetConfigAsString("loadbalancerId")
	tcListenerId := d.option.DeployConfig.GetConfigAsString("listenerId")
	if tcLoadbalancerId == "" {
		return errors.New("`loadbalancerId` is required")
	}
	if tcListenerId == "" {
		return errors.New("`listenerId` is required")
	}

	// 上传证书到 SSL
	upres, err := d.sslUploader.Upload(ctx, d.option.Certificate.Certificate, d.option.Certificate.PrivateKey)
	if err != nil {
		return err
	}

	d.infos = append(d.infos, toStr("已上传证书", upres))

	// 更新监听器证书
	if err := d.modifyListenerCertificate(ctx, tcLoadbalancerId, tcListenerId, upres.CertId); err != nil {
		return err
	}

	return nil
}

func (d *TencentCLBDeployer) deployToRuleDomain(ctx context.Context) error {
	tcLoadbalancerId := d.option.DeployConfig.GetConfigAsString("loadbalancerId")
	tcListenerId := d.option.DeployConfig.GetConfigAsString("listenerId")
	tcDomain := d.option.DeployConfig.GetConfigAsString("domain")
	if tcLoadbalancerId == "" {
		return errors.New("`loadbalancerId` is required")
	}
	if tcListenerId == "" {
		return errors.New("`listenerId` is required")
	}
	if tcDomain == "" {
		return errors.New("`domain` is required")
	}

	// 上传证书到 SSL
	upres, err := d.sslUploader.Upload(ctx, d.option.Certificate.Certificate, d.option.Certificate.PrivateKey)
	if err != nil {
		return err
	}

	d.infos = append(d.infos, toStr("已上传证书", upres))

	// 修改负载均衡七层监听器转发规则的域名级别属性
	// REF: https://cloud.tencent.com/document/api/214/38092
	modifyDomainAttributesReq := tcClb.NewModifyDomainAttributesRequest()
	modifyDomainAttributesReq.LoadBalancerId = common.StringPtr(tcLoadbalancerId)
	modifyDomainAttributesReq.ListenerId = common.StringPtr(tcListenerId)
	modifyDomainAttributesReq.Domain = common.StringPtr(tcDomain)
	modifyDomainAttributesReq.Certificate = &tcClb.CertificateInput{
		SSLMode: common.StringPtr("UNIDIRECTIONAL"),
		CertId:  common.StringPtr(upres.CertId),
	}
	modifyDomainAttributesResp, err := d.sdkClients.clb.ModifyDomainAttributes(modifyDomainAttributesReq)
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'clb.ModifyDomainAttributes'")
	}

	d.infos = append(d.infos, toStr("已修改七层监听器转发规则的域名级别属性", modifyDomainAttributesResp.Response))

	return nil
}

func (d *TencentCLBDeployer) modifyListenerCertificate(ctx context.Context, tcLoadbalancerId, tcListenerId, tcCertId string) error {
	// 查询监听器列表
	// REF: https://cloud.tencent.com/document/api/214/30686
	describeListenersReq := tcClb.NewDescribeListenersRequest()
	describeListenersReq.LoadBalancerId = common.StringPtr(tcLoadbalancerId)
	describeListenersReq.ListenerIds = common.StringPtrs([]string{tcListenerId})
	describeListenersResp, err := d.sdkClients.clb.DescribeListeners(describeListenersReq)
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'clb.DescribeListeners'")
	}
	if len(describeListenersResp.Response.Listeners) == 0 {
		d.infos = append(d.infos, toStr("未找到监听器", nil))
		return errors.New("listener not found")
	}

	d.infos = append(d.infos, toStr("已查询到监听器属性", describeListenersResp.Response))

	// 修改监听器属性
	// REF: https://cloud.tencent.com/document/product/214/30681
	modifyListenerReq := tcClb.NewModifyListenerRequest()
	modifyListenerReq.LoadBalancerId = common.StringPtr(tcLoadbalancerId)
	modifyListenerReq.ListenerId = common.StringPtr(tcListenerId)
	modifyListenerReq.Certificate = &tcClb.CertificateInput{CertId: common.StringPtr(tcCertId)}
	if describeListenersResp.Response.Listeners[0].Certificate != nil && describeListenersResp.Response.Listeners[0].Certificate.SSLMode != nil {
		modifyListenerReq.Certificate.SSLMode = describeListenersResp.Response.Listeners[0].Certificate.SSLMode
		modifyListenerReq.Certificate.CertCaId = describeListenersResp.Response.Listeners[0].Certificate.CertCaId
	} else {
		modifyListenerReq.Certificate.SSLMode = common.StringPtr("UNIDIRECTIONAL")
	}
	modifyListenerResp, err := d.sdkClients.clb.ModifyListener(modifyListenerReq)
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'clb.ModifyListener'")
	}

	d.infos = append(d.infos, toStr("已修改监听器属性", modifyListenerResp.Response))

	return nil
}
