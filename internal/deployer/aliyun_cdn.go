package deployer

import (
	"context"
	"encoding/json"
	"fmt"

	cdn20180510 "github.com/alibabacloud-go/cdn-20180510/v5/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"

	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/utils/rand"
)

type AliyunCDNDeployer struct {
	client *cdn20180510.Client
	option *DeployerOption
	infos  []string
}

func NewAliyunCDNDeployer(option *DeployerOption) (*AliyunCDNDeployer, error) {
	access := &domain.AliyunAccess{}
	json.Unmarshal([]byte(option.Access), access)

	d := &AliyunCDNDeployer{
		option: option,
	}

	client, err := d.createClient(access.AccessKeyId, access.AccessKeySecret)
	if err != nil {
		return nil, err
	}

	return &AliyunCDNDeployer{
		client: client,
		option: option,
		infos:  make([]string, 0),
	}, nil
}

func (d *AliyunCDNDeployer) GetID() string {
	return fmt.Sprintf("%s-%s", d.option.AccessRecord.GetString("name"), d.option.AccessRecord.Id)
}

func (d *AliyunCDNDeployer) GetInfo() []string {
	return d.infos
}

func (d *AliyunCDNDeployer) Deploy(ctx context.Context) error {
	certName := fmt.Sprintf("%s-%s-%s", d.option.Domain, d.option.DomainId, rand.RandStr(6))
	setCdnDomainSSLCertificateRequest := &cdn20180510.SetCdnDomainSSLCertificateRequest{
		DomainName:  tea.String(getDeployString(d.option.DeployConfig, "domain")),
		CertName:    tea.String(certName),
		CertType:    tea.String("upload"),
		SSLProtocol: tea.String("on"),
		SSLPub:      tea.String(d.option.Certificate.Certificate),
		SSLPri:      tea.String(d.option.Certificate.PrivateKey),
		CertRegion:  tea.String("cn-hangzhou"),
	}

	runtime := &util.RuntimeOptions{}

	resp, err := d.client.SetCdnDomainSSLCertificateWithOptions(setCdnDomainSSLCertificateRequest, runtime)
	if err != nil {
		return err
	}

	d.infos = append(d.infos, toStr("cdn设置证书", resp))

	return nil
}

func (d *AliyunCDNDeployer) createClient(accessKeyId, accessKeySecret string) (_result *cdn20180510.Client, _err error) {
	config := &openapi.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
	}
	config.Endpoint = tea.String("cdn.aliyuncs.com")
	_result = &cdn20180510.Client{}
	_result, _err = cdn20180510.NewClient(config)
	return _result, _err
}
