package deployer

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	aliyunOpen "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	aliyunDcdn "github.com/alibabacloud-go/dcdn-20180115/v3/client"
	"github.com/alibabacloud-go/tea/tea"
	xerrors "github.com/pkg/errors"

	"github.com/usual2970/certimate/internal/domain"
)

type AliyunDCDNDeployer struct {
	option *DeployerOption
	infos  []string

	sdkClient *aliyunDcdn.Client
}

func NewAliyunDCDNDeployer(option *DeployerOption) (Deployer, error) {
	access := &domain.AliyunAccess{}
	if err := json.Unmarshal([]byte(option.Access), access); err != nil {
		return nil, xerrors.Wrap(err, "failed to get access")
	}

	client, err := (&AliyunDCDNDeployer{}).createSdkClient(
		access.AccessKeyId,
		access.AccessKeySecret,
	)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk client")
	}

	return &AliyunDCDNDeployer{
		option:    option,
		infos:     make([]string, 0),
		sdkClient: client,
	}, nil
}

func (d *AliyunDCDNDeployer) GetID() string {
	return fmt.Sprintf("%s-%s", d.option.AccessRecord.GetString("name"), d.option.AccessRecord.Id)
}

func (d *AliyunDCDNDeployer) GetInfos() []string {
	return d.infos
}

func (d *AliyunDCDNDeployer) Deploy(ctx context.Context) error {
	// 支持泛解析域名，在 Aliyun DCDN 中泛解析域名表示为 .example.com
	domain := d.option.DeployConfig.GetConfigAsString("domain")
	if strings.HasPrefix(domain, "*") {
		domain = strings.TrimPrefix(domain, "*")
	}

	// 配置域名证书
	// REF: https://help.aliyun.com/zh/edge-security-acceleration/dcdn/developer-reference/api-dcdn-2018-01-15-setdcdndomainsslcertificate
	setDcdnDomainSSLCertificateReq := &aliyunDcdn.SetDcdnDomainSSLCertificateRequest{
		DomainName:  tea.String(domain),
		CertRegion:  tea.String(d.option.DeployConfig.GetConfigOrDefaultAsString("region", "cn-hangzhou")),
		CertName:    tea.String(fmt.Sprintf("certimate-%d", time.Now().UnixMilli())),
		CertType:    tea.String("upload"),
		SSLProtocol: tea.String("on"),
		SSLPub:      tea.String(d.option.Certificate.Certificate),
		SSLPri:      tea.String(d.option.Certificate.PrivateKey),
	}
	setDcdnDomainSSLCertificateResp, err := d.sdkClient.SetDcdnDomainSSLCertificate(setDcdnDomainSSLCertificateReq)
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'dcdn.SetDcdnDomainSSLCertificate'")
	}

	d.infos = append(d.infos, toStr("已配置 DCDN 域名证书", setDcdnDomainSSLCertificateResp))

	return nil
}

func (d *AliyunDCDNDeployer) createSdkClient(accessKeyId, accessKeySecret string) (*aliyunDcdn.Client, error) {
	aConfig := &aliyunOpen.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
		Endpoint:        tea.String("dcdn.aliyuncs.com"),
	}

	client, err := aliyunDcdn.NewClient(aConfig)
	if err != nil {
		return nil, err
	}

	return client, nil
}
