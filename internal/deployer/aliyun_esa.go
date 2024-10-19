/*
 * @Author: Bin
 * @Date: 2024-09-17
 * @FilePath: /certimate/internal/deployer/aliyun_esa.go
 */
package deployer

import (
	"context"
	"encoding/json"
	"fmt"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dcdn20180115 "github.com/alibabacloud-go/dcdn-20180115/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"

	"certimate/internal/domain"
	"certimate/internal/utils/rand"
)

type AliyunESADeployer struct {
	client *dcdn20180115.Client
	option *DeployerOption
	infos  []string
}

func NewAliyunESADeployer(option *DeployerOption) (*AliyunESADeployer, error) {
	access := &domain.AliyunAccess{}
	json.Unmarshal([]byte(option.Access), access)

	d := &AliyunESADeployer{
		option: option,
	}

	client, err := d.createClient(access.AccessKeyId, access.AccessKeySecret)
	if err != nil {
		return nil, err
	}

	return &AliyunESADeployer{
		client: client,
		option: option,
		infos:  make([]string, 0),
	}, nil
}

func (d *AliyunESADeployer) GetID() string {
	return fmt.Sprintf("%s-%s", d.option.AceessRecord.GetString("name"), d.option.AceessRecord.Id)
}

func (d *AliyunESADeployer) GetInfo() []string {
	return d.infos
}

func (d *AliyunESADeployer) Deploy(ctx context.Context) error {
	certName := fmt.Sprintf("%s-%s-%s", d.option.Domain, d.option.DomainId, rand.RandStr(6))
	setDcdnDomainSSLCertificateRequest := &dcdn20180115.SetDcdnDomainSSLCertificateRequest{
		DomainName:  tea.String(getDeployString(d.option.DeployConfig, "domain")),
		CertName:    tea.String(certName),
		CertType:    tea.String("upload"),
		SSLProtocol: tea.String("on"),
		SSLPub:      tea.String(d.option.Certificate.Certificate),
		SSLPri:      tea.String(d.option.Certificate.PrivateKey),
		CertRegion:  tea.String("cn-hangzhou"),
	}

	runtime := &util.RuntimeOptions{}

	resp, err := d.client.SetDcdnDomainSSLCertificateWithOptions(setDcdnDomainSSLCertificateRequest, runtime)
	if err != nil {
		return err
	}

	d.infos = append(d.infos, toStr("dcdn设置证书", resp))

	return nil
}

func (d *AliyunESADeployer) createClient(accessKeyId, accessKeySecret string) (_result *dcdn20180115.Client, _err error) {
	config := &openapi.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
	}
	config.Endpoint = tea.String("dcdn.aliyuncs.com")
	_result = &dcdn20180115.Client{}
	_result, _err = dcdn20180115.NewClient(config)
	return _result, _err
}
