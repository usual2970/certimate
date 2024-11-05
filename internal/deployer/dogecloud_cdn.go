package deployer

import (
	"context"
	"encoding/json"
	"fmt"

	xerrors "github.com/pkg/errors"

	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	uploaderDoge "github.com/usual2970/certimate/internal/pkg/core/uploader/providers/dogecloud"
	doge "github.com/usual2970/certimate/internal/pkg/vendors/dogecloud-sdk"
)

type DogeCloudCDNDeployer struct {
	option *DeployerOption
	infos  []string

	sdkClient   *doge.Client
	sslUploader uploader.Uploader
}

func NewDogeCloudCDNDeployer(option *DeployerOption) (Deployer, error) {
	access := &domain.DogeCloudAccess{}
	if err := json.Unmarshal([]byte(option.Access), access); err != nil {
		return nil, xerrors.Wrap(err, "failed to get access")
	}

	client, err := (&DogeCloudCDNDeployer{}).createSdkClient(
		access.AccessKey,
		access.SecretKey,
	)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk client")
	}

	uploader, err := uploaderDoge.New(&uploaderDoge.DogeCloudUploaderConfig{
		AccessKey: access.AccessKey,
		SecretKey: access.SecretKey,
	})
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create ssl uploader")
	}

	return &DogeCloudCDNDeployer{
		option:      option,
		infos:       make([]string, 0),
		sdkClient:   client,
		sslUploader: uploader,
	}, nil
}

func (d *DogeCloudCDNDeployer) GetID() string {
	return fmt.Sprintf("%s-%s", d.option.AccessRecord.GetString("name"), d.option.AccessRecord.Id)
}

func (d *DogeCloudCDNDeployer) GetInfos() []string {
	return d.infos
}

func (d *DogeCloudCDNDeployer) Deploy(ctx context.Context) error {
	// 上传证书到 CDN
	upres, err := d.sslUploader.Upload(ctx, d.option.Certificate.Certificate, d.option.Certificate.PrivateKey)
	if err != nil {
		return err
	}

	d.infos = append(d.infos, toStr("已上传证书", upres))

	// 绑定证书
	// REF: https://docs.dogecloud.com/cdn/api-cert-bind
	bindCdnCertResp, err := d.sdkClient.BindCdnCertWithDomain(upres.CertId, d.option.DeployConfig.GetConfigAsString("domain"))
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'cdn.BindCdnCert'")
	}

	d.infos = append(d.infos, toStr("已绑定证书", bindCdnCertResp))

	return nil
}

func (d *DogeCloudCDNDeployer) createSdkClient(accessKey, secretKey string) (*doge.Client, error) {
	client := doge.NewClient(accessKey, secretKey)
	return client, nil
}
