package baiducloudcdn

import (
	"context"
	"errors"
	"fmt"
	"time"

	bceCdn "github.com/baidubce/bce-sdk-go/services/cdn"
	bceCdnApi "github.com/baidubce/bce-sdk-go/services/cdn/api"
	xerrors "github.com/pkg/errors"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
)

type BaiduCloudCDNDeployerConfig struct {
	// 百度智能云 AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// 百度智能云 SecretAccessKey。
	SecretAccessKey string `json:"secretAccessKey"`
	// 加速域名（不支持泛域名）。
	Domain string `json:"domain"`
}

type BaiduCloudCDNDeployer struct {
	config    *BaiduCloudCDNDeployerConfig
	logger    deployer.Logger
	sdkClient *bceCdn.Client
}

var _ deployer.Deployer = (*BaiduCloudCDNDeployer)(nil)

func New(config *BaiduCloudCDNDeployerConfig) (*BaiduCloudCDNDeployer, error) {
	return NewWithLogger(config, deployer.NewNilLogger())
}

func NewWithLogger(config *BaiduCloudCDNDeployerConfig, logger deployer.Logger) (*BaiduCloudCDNDeployer, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}

	if logger == nil {
		return nil, errors.New("logger is nil")
	}

	client, err := createSdkClient(config.AccessKeyId, config.SecretAccessKey)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk client")
	}

	return &BaiduCloudCDNDeployer{
		logger:    logger,
		config:    config,
		sdkClient: client,
	}, nil
}

func (d *BaiduCloudCDNDeployer) Deploy(ctx context.Context, certPem string, privkeyPem string) (*deployer.DeployResult, error) {
	// 修改域名证书
	// REF: https://cloud.baidu.com/doc/CDN/s/qjzuz2hp8
	putCertResp, err := d.sdkClient.PutCert(
		d.config.Domain,
		&bceCdnApi.UserCertificate{
			CertName:    fmt.Sprintf("certimate-%d", time.Now().UnixMilli()),
			ServerData:  certPem,
			PrivateData: privkeyPem,
		},
		"ON",
	)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'cdn.PutCert'")
	}

	d.logger.Appendt("已修改域名证书", putCertResp)

	return &deployer.DeployResult{}, nil
}

func createSdkClient(accessKeyId, secretAccessKey string) (*bceCdn.Client, error) {
	client, err := bceCdn.NewClient(accessKeyId, secretAccessKey, "")
	if err != nil {
		return nil, err
	}

	return client, nil
}
