package deployer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	aliyunCdn "github.com/alibabacloud-go/cdn-20180510/v5/client"
	aliyunOpen "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	xerrors "github.com/pkg/errors"

	"github.com/usual2970/certimate/internal/domain"
)

type AliyunCDNDeployer struct {
	option *DeployerOption
	infos  []string

	sdkClient *aliyunCdn.Client
}

func NewAliyunCDNDeployer(option *DeployerOption) (Deployer, error) {
	access := &domain.AliyunAccess{}
	if err := json.Unmarshal([]byte(option.Access), access); err != nil {
		return nil, xerrors.Wrap(err, "failed to get access")
	}

	client, err := (&AliyunCDNDeployer{}).createSdkClient(
		access.AccessKeyId,
		access.AccessKeySecret,
	)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk client")
	}

	return &AliyunCDNDeployer{
		option:    option,
		infos:     make([]string, 0),
		sdkClient: client,
	}, nil
}

func (d *AliyunCDNDeployer) GetID() string {
	return fmt.Sprintf("%s-%s", d.option.AccessRecord.GetString("name"), d.option.AccessRecord.Id)
}

func (d *AliyunCDNDeployer) GetInfos() []string {
	return d.infos
}

func (d *AliyunCDNDeployer) Deploy(ctx context.Context) error {
	// 设置 CDN 域名域名证书
	// REF: https://help.aliyun.com/zh/cdn/developer-reference/api-cdn-2018-05-10-setcdndomainsslcertificate
	setCdnDomainSSLCertificateReq := &aliyunCdn.SetCdnDomainSSLCertificateRequest{
		DomainName:  tea.String(d.option.DeployConfig.GetConfigAsString("domain")),
		CertRegion:  tea.String(d.option.DeployConfig.GetConfigOrDefaultAsString("region", "cn-hangzhou")),
		CertName:    tea.String(fmt.Sprintf("certimate-%d", time.Now().UnixMilli())),
		CertType:    tea.String("upload"),
		SSLProtocol: tea.String("on"),
		SSLPub:      tea.String(d.option.Certificate.Certificate),
		SSLPri:      tea.String(d.option.Certificate.PrivateKey),
	}
	setCdnDomainSSLCertificateResp, err := d.sdkClient.SetCdnDomainSSLCertificate(setCdnDomainSSLCertificateReq)
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'cdn.SetCdnDomainSSLCertificate'")
	}

	d.infos = append(d.infos, toStr("已设置 CDN 域名证书", setCdnDomainSSLCertificateResp))

	return nil
}

func (d *AliyunCDNDeployer) createSdkClient(accessKeyId, accessKeySecret string) (*aliyunCdn.Client, error) {
	aConfig := &aliyunOpen.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
		Endpoint:        tea.String("cdn.aliyuncs.com"),
	}

	client, err := aliyunCdn.NewClient(aConfig)
	if err != nil {
		return nil, err
	}

	return client, nil
}
