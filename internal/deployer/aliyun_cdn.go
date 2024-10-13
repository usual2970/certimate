package deployer

import (
	"certimate/internal/domain"
	"certimate/internal/utils/rand"
	"context"
	"encoding/json"
	"fmt"

	cdn20180510 "github.com/alibabacloud-go/cdn-20180510/v5/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

type AliyunCdn struct {
	client *cdn20180510.Client
	option *DeployerOption
	infos  []string
}

func NewAliyunCdn(option *DeployerOption) (*AliyunCdn, error) {
	access := &domain.AliyunAccess{}
	json.Unmarshal([]byte(option.Access), access)
	a := &AliyunCdn{
		option: option,
	}
	client, err := a.createClient(access.AccessKeyId, access.AccessKeySecret)
	if err != nil {
		return nil, err
	}

	return &AliyunCdn{
		client: client,
		option: option,
		infos:  make([]string, 0),
	}, nil
}

func (a *AliyunCdn) GetID() string {
	return fmt.Sprintf("%s-%s", a.option.AceessRecord.GetString("name"), a.option.AceessRecord.Id)
}

func (a *AliyunCdn) GetInfo() []string {
	return a.infos
}

func (a *AliyunCdn) Deploy(ctx context.Context) error {

	certName := fmt.Sprintf("%s-%s-%s", a.option.Domain, a.option.DomainId, rand.RandStr(6))
	setCdnDomainSSLCertificateRequest := &cdn20180510.SetCdnDomainSSLCertificateRequest{
		DomainName:  tea.String(getDeployString(a.option.DeployConfig, "domain")),
		CertName:    tea.String(certName),
		CertType:    tea.String("upload"),
		SSLProtocol: tea.String("on"),
		SSLPub:      tea.String(a.option.Certificate.Certificate),
		SSLPri:      tea.String(a.option.Certificate.PrivateKey),
		CertRegion:  tea.String("cn-hangzhou"),
	}

	runtime := &util.RuntimeOptions{}

	resp, err := a.client.SetCdnDomainSSLCertificateWithOptions(setCdnDomainSSLCertificateRequest, runtime)
	if err != nil {
		return err
	}

	a.infos = append(a.infos, toStr("cdn设置证书", resp))

	return nil
}

func (a *AliyunCdn) createClient(accessKeyId, accessKeySecret string) (_result *cdn20180510.Client, _err error) {
	config := &openapi.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
	}
	config.Endpoint = tea.String("cdn.aliyuncs.com")
	_result = &cdn20180510.Client{}
	_result, _err = cdn20180510.NewClient(config)
	return _result, _err
}
