package deployer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	bceCdn "github.com/baidubce/bce-sdk-go/services/cdn"
	bceCdnApi "github.com/baidubce/bce-sdk-go/services/cdn/api"
	xerrors "github.com/pkg/errors"

	"github.com/usual2970/certimate/internal/domain"
)

type BaiduCloudCDNDeployer struct {
	option *DeployerOption
	infos  []string

	sdkClient *bceCdn.Client
}

func NewBaiduCloudCDNDeployer(option *DeployerOption) (Deployer, error) {
	access := &domain.BaiduCloudAccess{}
	if err := json.Unmarshal([]byte(option.Access), access); err != nil {
		return nil, xerrors.Wrap(err, "failed to get access")
	}

	client, err := (&BaiduCloudCDNDeployer{}).createSdkClient(
		access.AccessKeyId,
		access.SecretAccessKey,
	)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk client")
	}

	return &BaiduCloudCDNDeployer{
		option:    option,
		infos:     make([]string, 0),
		sdkClient: client,
	}, nil
}

func (d *BaiduCloudCDNDeployer) GetID() string {
	return fmt.Sprintf("%s-%s", d.option.AccessRecord.GetString("name"), d.option.AccessRecord.Id)
}

func (d *BaiduCloudCDNDeployer) GetInfos() []string {
	return d.infos
}

func (d *BaiduCloudCDNDeployer) Deploy(ctx context.Context) error {
	// 修改域名证书
	// REF: https://cloud.baidu.com/doc/CDN/s/qjzuz2hp8
	putCertResp, err := d.sdkClient.PutCert(
		d.option.DeployConfig.GetConfigAsString("domain"),
		&bceCdnApi.UserCertificate{
			CertName:    fmt.Sprintf("certimate-%d", time.Now().UnixMilli()),
			ServerData:  d.option.Certificate.Certificate,
			PrivateData: d.option.Certificate.PrivateKey,
		},
		"ON",
	)
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'cdn.PutCert'")
	}

	d.infos = append(d.infos, toStr("已修改域名证书", putCertResp))

	return nil
}

func (d *BaiduCloudCDNDeployer) createSdkClient(accessKeyId, secretAccessKey string) (*bceCdn.Client, error) {
	client, err := bceCdn.NewClient(accessKeyId, secretAccessKey, "")
	if err != nil {
		return nil, err
	}

	return client, nil
}
