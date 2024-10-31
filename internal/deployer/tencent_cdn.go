package deployer

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	xerrors "github.com/pkg/errors"
	tcCdn "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn/v20180606"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tcSsl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"
	"golang.org/x/exp/slices"

	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
)

type TencentCDNDeployer struct {
	option *DeployerOption
	infos  []string

	sdkClients  *tencentCDNDeployerSdkClients
	sslUploader uploader.Uploader
}

type tencentCDNDeployerSdkClients struct {
	ssl *tcSsl.Client
	cdn *tcCdn.Client
}

func NewTencentCDNDeployer(option *DeployerOption) (Deployer, error) {
	access := &domain.TencentAccess{}
	if err := json.Unmarshal([]byte(option.Access), access); err != nil {
		return nil, xerrors.Wrap(err, "failed to get access")
	}

	clients, err := (&TencentCDNDeployer{}).createSdkClients(
		access.SecretId,
		access.SecretKey,
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

	return &TencentCDNDeployer{
		option:      option,
		infos:       make([]string, 0),
		sdkClients:  clients,
		sslUploader: uploader,
	}, nil
}

func (d *TencentCDNDeployer) GetID() string {
	return fmt.Sprintf("%s-%s", d.option.AccessRecord.GetString("name"), d.option.AccessRecord.Id)
}

func (d *TencentCDNDeployer) GetInfos() []string {
	return d.infos
}

func (d *TencentCDNDeployer) Deploy(ctx context.Context) error {
	// 上传证书到 SSL
	upres, err := d.sslUploader.Upload(ctx, d.option.Certificate.Certificate, d.option.Certificate.PrivateKey)
	if err != nil {
		return err
	}

	d.infos = append(d.infos, toStr("已上传证书", upres))

	// 获取待部署的 CDN 实例
	// 如果是泛域名，根据证书匹配 CDN 实例
	tcInstanceIds := make([]string, 0)
	domain := d.option.DeployConfig.GetConfigAsString("domain")
	if strings.HasPrefix(domain, "*") {
		domains, err := d.getDomainsByCertificateId(upres.CertId)
		if err != nil {
			return err
		}

		tcInstanceIds = domains
	} else {
		tcInstanceIds = append(tcInstanceIds, domain)
	}

	// 跳过已部署的 CDN 实例
	if len(tcInstanceIds) > 0 {
		deployedDomains, err := d.getDeployedDomainsByCertificateId(upres.CertId)
		if err != nil {
			return err
		}

		temp := make([]string, 0)
		for _, aliInstanceId := range tcInstanceIds {
			if !slices.Contains(deployedDomains, aliInstanceId) {
				temp = append(temp, aliInstanceId)
			}
		}
		tcInstanceIds = temp
	}
	if len(tcInstanceIds) == 0 {
		d.infos = append(d.infos, "已部署过或没有要部署的 CDN 实例")
		return nil
	}

	// 证书部署到 CDN 实例
	// REF: https://cloud.tencent.com/document/product/400/91667
	deployCertificateInstanceReq := tcSsl.NewDeployCertificateInstanceRequest()
	deployCertificateInstanceReq.CertificateId = common.StringPtr(upres.CertId)
	deployCertificateInstanceReq.ResourceType = common.StringPtr("cdn")
	deployCertificateInstanceReq.Status = common.Int64Ptr(1)
	deployCertificateInstanceReq.InstanceIdList = common.StringPtrs(tcInstanceIds)
	deployCertificateInstanceResp, err := d.sdkClients.ssl.DeployCertificateInstance(deployCertificateInstanceReq)
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'ssl.DeployCertificateInstance'")
	}

	d.infos = append(d.infos, toStr("已部署证书到云资源实例", deployCertificateInstanceResp.Response))

	return nil
}

func (d *TencentCDNDeployer) createSdkClients(secretId, secretKey string) (*tencentCDNDeployerSdkClients, error) {
	credential := common.NewCredential(secretId, secretKey)

	sslClient, err := tcSsl.NewClient(credential, "", profile.NewClientProfile())
	if err != nil {
		return nil, err
	}

	cdnClient, err := tcCdn.NewClient(credential, "", profile.NewClientProfile())
	if err != nil {
		return nil, err
	}

	return &tencentCDNDeployerSdkClients{
		ssl: sslClient,
		cdn: cdnClient,
	}, nil
}

func (d *TencentCDNDeployer) getDomainsByCertificateId(tcCertId string) ([]string, error) {
	// 获取证书中的可用域名
	// REF: https://cloud.tencent.com/document/product/228/42491
	describeCertDomainsReq := tcCdn.NewDescribeCertDomainsRequest()
	describeCertDomainsReq.CertId = common.StringPtr(tcCertId)
	describeCertDomainsReq.Product = common.StringPtr("cdn")
	describeCertDomainsResp, err := d.sdkClients.cdn.DescribeCertDomains(describeCertDomainsReq)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'cdn.DescribeCertDomains'")
	}

	domains := make([]string, 0)
	if describeCertDomainsResp.Response.Domains == nil {
		for _, domain := range describeCertDomainsResp.Response.Domains {
			domains = append(domains, *domain)
		}
	}

	return domains, nil
}

func (d *TencentCDNDeployer) getDeployedDomainsByCertificateId(tcCertId string) ([]string, error) {
	// 根据证书查询关联 CDN 域名
	// REF: https://cloud.tencent.com/document/product/400/62674
	describeDeployedResourcesReq := tcSsl.NewDescribeDeployedResourcesRequest()
	describeDeployedResourcesReq.CertificateIds = common.StringPtrs([]string{tcCertId})
	describeDeployedResourcesReq.ResourceType = common.StringPtr("cdn")
	describeDeployedResourcesResp, err := d.sdkClients.ssl.DescribeDeployedResources(describeDeployedResourcesReq)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'cdn.DescribeDeployedResources'")
	}

	domains := make([]string, 0)
	if describeDeployedResourcesResp.Response.DeployedResources != nil {
		for _, deployedResource := range describeDeployedResourcesResp.Response.DeployedResources {
			for _, resource := range deployedResource.Resources {
				domains = append(domains, *resource)
			}
		}
	}

	return domains, nil
}
