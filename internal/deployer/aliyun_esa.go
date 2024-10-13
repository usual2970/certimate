/*
 * @Author: Bin
 * @Date: 2024-09-17
 * @FilePath: /certimate/internal/deployer/aliyun_esa.go
 */
package deployer

import (
	"certimate/internal/domain"
	"certimate/internal/utils/rand"
	"context"
	"encoding/json"
	"fmt"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dcdn20180115 "github.com/alibabacloud-go/dcdn-20180115/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

type AliyunEsa struct {
	client *dcdn20180115.Client
	option *DeployerOption
	infos  []string
}

func NewAliyunEsa(option *DeployerOption) (*AliyunEsa, error) {
	access := &domain.AliyunAccess{}
	json.Unmarshal([]byte(option.Access), access)
	a := &AliyunEsa{
		option: option,
	}
	client, err := a.createClient(access.AccessKeyId, access.AccessKeySecret)
	if err != nil {
		return nil, err
	}

	return &AliyunEsa{
		client: client,
		option: option,
		infos:  make([]string, 0),
	}, nil
}

func (a *AliyunEsa) GetID() string {
	return fmt.Sprintf("%s-%s", a.option.AceessRecord.GetString("name"), a.option.AceessRecord.Id)
}

func (a *AliyunEsa) GetInfo() []string {
	return a.infos
}

func (a *AliyunEsa) Deploy(ctx context.Context) error {

	certName := fmt.Sprintf("%s-%s-%s", a.option.Domain, a.option.DomainId, rand.RandStr(6))
	setDcdnDomainSSLCertificateRequest := &dcdn20180115.SetDcdnDomainSSLCertificateRequest{
		DomainName:  tea.String(getDeployString(a.option.DeployConfig, "domain")),
		CertName:    tea.String(certName),
		CertType:    tea.String("upload"),
		SSLProtocol: tea.String("on"),
		SSLPub:      tea.String(a.option.Certificate.Certificate),
		SSLPri:      tea.String(a.option.Certificate.PrivateKey),
		CertRegion:  tea.String("cn-hangzhou"),
	}

	runtime := &util.RuntimeOptions{}

	resp, err := a.client.SetDcdnDomainSSLCertificateWithOptions(setDcdnDomainSSLCertificateRequest, runtime)
	if err != nil {
		return err
	}

	a.infos = append(a.infos, toStr("dcdn设置证书", resp))

	return nil
}

func (a *AliyunEsa) createClient(accessKeyId, accessKeySecret string) (_result *dcdn20180115.Client, _err error) {
	config := &openapi.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
	}
	config.Endpoint = tea.String("dcdn.aliyuncs.com")
	_result = &dcdn20180115.Client{}
	_result, _err = dcdn20180115.NewClient(config)
	return _result, _err
}
