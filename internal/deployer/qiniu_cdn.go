package deployer

import (
	"context"
	"encoding/json"
	"fmt"

	xerrors "github.com/pkg/errors"
	"github.com/qiniu/go-sdk/v7/auth"

	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	qiniuEx "github.com/usual2970/certimate/internal/pkg/vendors/qiniu-sdk"
)

type QiniuCDNDeployer struct {
	option *DeployerOption
	infos  []string

	sdkClient   *qiniuEx.Client
	sslUploader uploader.Uploader
}

func NewQiniuCDNDeployer(option *DeployerOption) (Deployer, error) {
	access := &domain.QiniuAccess{}
	if err := json.Unmarshal([]byte(option.Access), access); err != nil {
		return nil, xerrors.Wrap(err, "failed to get access")
	}

	client, err := (&QiniuCDNDeployer{}).createSdkClient(
		access.AccessKey,
		access.SecretKey,
	)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk client")
	}

	uploader, err := uploader.NewQiniuSSLCertUploader(&uploader.QiniuSSLCertUploaderConfig{
		AccessKey: access.AccessKey,
		SecretKey: access.SecretKey,
	})
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create ssl uploader")
	}

	return &QiniuCDNDeployer{
		option:      option,
		infos:       make([]string, 0),
		sdkClient:   client,
		sslUploader: uploader,
	}, nil
}

func (d *QiniuCDNDeployer) GetID() string {
	return fmt.Sprintf("%s-%s", d.option.AccessRecord.GetString("name"), d.option.AccessRecord.Id)
}

func (d *QiniuCDNDeployer) GetInfo() []string {
	return d.infos
}

func (d *QiniuCDNDeployer) Deploy(ctx context.Context) error {
	// 上传证书
	upres, err := d.sslUploader.Upload(ctx, d.option.Certificate.Certificate, d.option.Certificate.PrivateKey)
	if err != nil {
		return err
	}

	d.infos = append(d.infos, toStr("已上传证书", upres))

	// 获取域名信息
	// REF: https://developer.qiniu.com/fusion/4246/the-domain-name
	domain := d.option.DeployConfig.GetConfigAsString("domain")
	getDomainInfoResp, err := d.sdkClient.GetDomainInfo(domain)
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'cdn.GetDomainInfo'")
	}

	d.infos = append(d.infos, toStr("已获取域名信息", getDomainInfoResp))

	// 判断域名是否已启用 HTTPS。如果已启用，修改域名证书；否则，启用 HTTPS
	// REF: https://developer.qiniu.com/fusion/4246/the-domain-name
	if getDomainInfoResp.Https != nil && getDomainInfoResp.Https.CertID != "" {
		modifyDomainHttpsConfResp, err := d.sdkClient.ModifyDomainHttpsConf(domain, upres.CertId, getDomainInfoResp.Https.ForceHttps, getDomainInfoResp.Https.Http2Enable)
		if err != nil {
			return xerrors.Wrap(err, "failed to execute sdk request 'cdn.ModifyDomainHttpsConf'")
		}

		d.infos = append(d.infos, toStr("已修改域名证书", modifyDomainHttpsConfResp))
	} else {
		enableDomainHttpsResp, err := d.sdkClient.EnableDomainHttps(domain, upres.CertId, true, true)
		if err != nil {
			return xerrors.Wrap(err, "failed to execute sdk request 'cdn.EnableDomainHttps'")
		}

		d.infos = append(d.infos, toStr("已将域名升级为 HTTPS", enableDomainHttpsResp))
	}

	return nil
}

func (u *QiniuCDNDeployer) createSdkClient(accessKey, secretKey string) (*qiniuEx.Client, error) {
	credential := auth.New(accessKey, secretKey)
	client := qiniuEx.NewClient(credential)
	return client, nil
}
