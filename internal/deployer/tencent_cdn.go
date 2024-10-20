package deployer

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	cdn "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn/v20180606"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"

	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/utils/rand"
)

type TencentCDNDeployer struct {
	option     *DeployerOption
	credential *common.Credential
	infos      []string
}

func NewTencentCDNDeployer(option *DeployerOption) (Deployer, error) {
	access := &domain.TencentAccess{}
	if err := json.Unmarshal([]byte(option.Access), access); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tencent access: %w", err)
	}

	credential := common.NewCredential(
		access.SecretId,
		access.SecretKey,
	)

	return &TencentCDNDeployer{
		option:     option,
		credential: credential,
		infos:      make([]string, 0),
	}, nil
}

func (d *TencentCDNDeployer) GetID() string {
	return fmt.Sprintf("%s-%s", d.option.AccessRecord.GetString("name"), d.option.AccessRecord.Id)
}

func (d *TencentCDNDeployer) GetInfo() []string {
	return d.infos
}

func (d *TencentCDNDeployer) Deploy(ctx context.Context) error {
	// 上传证书
	certId, err := d.uploadCert()
	if err != nil {
		return fmt.Errorf("failed to upload certificate: %w", err)
	}
	d.infos = append(d.infos, toStr("上传证书", certId))

	if err := d.deploy(certId); err != nil {
		return fmt.Errorf("failed to deploy: %w", err)
	}

	return nil
}

func (d *TencentCDNDeployer) uploadCert() (string, error) {
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "ssl.tencentcloudapi.com"

	client, _ := ssl.NewClient(d.credential, "", cpf)

	request := ssl.NewUploadCertificateRequest()

	request.CertificatePublicKey = common.StringPtr(d.option.Certificate.Certificate)
	request.CertificatePrivateKey = common.StringPtr(d.option.Certificate.PrivateKey)
	request.Alias = common.StringPtr(d.option.Domain + "_" + rand.RandStr(6))
	request.Repeatable = common.BoolPtr(false)

	response, err := client.UploadCertificate(request)
	if err != nil {
		return "", fmt.Errorf("failed to upload certificate: %w", err)
	}

	return *response.Response.CertificateId, nil
}

func (d *TencentCDNDeployer) deploy(certId string) error {
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "ssl.tencentcloudapi.com"
	// 实例化要请求产品的client对象,clientProfile是可选的
	client, _ := ssl.NewClient(d.credential, "", cpf)

	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := ssl.NewDeployCertificateInstanceRequest()

	request.CertificateId = common.StringPtr(certId)
	request.ResourceType = common.StringPtr("cdn")
	request.Status = common.Int64Ptr(1)

	// 如果是泛域名就从cdn列表下获取SSL证书中的可用域名
	domain := getDeployString(d.option.DeployConfig, "domain")
	if strings.Contains(domain, "*") {
		list, errGetList := d.getDomainList()
		if errGetList != nil {
			return fmt.Errorf("failed to get certificate domain list: %w", errGetList)
		}
		if list == nil || len(list) == 0 {
			return fmt.Errorf("failed to get certificate domain list: empty list.")
		}
		request.InstanceIdList = common.StringPtrs(list)
	} else { // 否则直接使用传入的域名
		request.InstanceIdList = common.StringPtrs([]string{domain})
	}

	// 返回的resp是一个DeployCertificateInstanceResponse的实例，与请求对象对应
	resp, err := client.DeployCertificateInstance(request)
	if err != nil {
		return fmt.Errorf("failed to deploy certificate: %w", err)
	}
	d.infos = append(d.infos, toStr("部署证书", resp.Response))
	return nil
}

func (d *TencentCDNDeployer) getDomainList() ([]string, error) {
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "cdn.tencentcloudapi.com"
	client, _ := cdn.NewClient(d.credential, "", cpf)

	request := cdn.NewDescribeCertDomainsRequest()

	cert := base64.StdEncoding.EncodeToString([]byte(d.option.Certificate.Certificate))
	request.Cert = &cert

	response, err := client.DescribeCertDomains(request)
	if err != nil {
		return nil, fmt.Errorf("failed to get domain list: %w", err)
	}

	domains := make([]string, 0)
	for _, domain := range response.Response.Domains {
		domains = append(domains, *domain)
	}

	return domains, nil
}
