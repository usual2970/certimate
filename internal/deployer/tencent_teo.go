package deployer

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"

	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/utils/rand"
)

type TencentTEODeployer struct {
	option     *DeployerOption
	credential *common.Credential
	infos      []string
}

func NewTencentTEODeployer(option *DeployerOption) (Deployer, error) {
	access := &domain.TencentAccess{}
	if err := json.Unmarshal([]byte(option.Access), access); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tencent access: %w", err)
	}

	credential := common.NewCredential(
		access.SecretId,
		access.SecretKey,
	)

	return &TencentTEODeployer{
		option:     option,
		credential: credential,
		infos:      make([]string, 0),
	}, nil
}

func (d *TencentTEODeployer) GetID() string {
	return fmt.Sprintf("%s-%s", d.option.AccessRecord.GetString("name"), d.option.AccessRecord.Id)
}

func (d *TencentTEODeployer) GetInfo() []string {
	return d.infos
}

func (d *TencentTEODeployer) Deploy(ctx context.Context) error {
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

func (d *TencentTEODeployer) uploadCert() (string, error) {
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

func (d *TencentTEODeployer) deploy(certId string) error {
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "teo.tencentcloudapi.com"
	// 实例化要请求产品的client对象,clientProfile是可选的
	client, _ := teo.NewClient(d.credential, "", cpf)

	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := teo.NewModifyHostsCertificateRequest()

	request.ZoneId = common.StringPtr(getDeployString(d.option.DeployConfig, "zoneId"))
	request.Mode = common.StringPtr("sslcert")
	request.ServerCertInfo = []*teo.ServerCertInfo{{
		CertId: common.StringPtr(certId),
	}}

	domains := strings.Split(strings.ReplaceAll(d.option.Domain, "\r\n", "\n"),"\n")
	request.Hosts = common.StringPtrs(domains)

	// 返回的resp是一个DeployCertificateInstanceResponse的实例，与请求对象对应
	resp, err := client.ModifyHostsCertificate(request)
	if err != nil {
		return fmt.Errorf("failed to deploy certificate: %w", err)
	}
	d.infos = append(d.infos, toStr("部署证书", resp.Response))
	return nil
}
