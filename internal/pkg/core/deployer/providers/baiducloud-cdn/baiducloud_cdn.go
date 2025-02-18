package baiducloudcdn

import (
	"context"
	"fmt"
	"time"

	bceCdn "github.com/baidubce/bce-sdk-go/services/cdn"
	bceCdnApi "github.com/baidubce/bce-sdk-go/services/cdn/api"
	xerrors "github.com/pkg/errors"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/logger"
)

type DeployerConfig struct {
	// 百度智能云 AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// 百度智能云 SecretAccessKey。
	SecretAccessKey string `json:"secretAccessKey"`
	// 加速域名（支持泛域名）。
	Domain string `json:"domain"`
}

type DeployerProvider struct {
	config    *DeployerConfig
	logger    logger.Logger
	sdkClient *bceCdn.Client
}

var _ deployer.Deployer = (*DeployerProvider)(nil)

func NewDeployer(config *DeployerConfig) (*DeployerProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	client, err := createSdkClient(config.AccessKeyId, config.SecretAccessKey)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk client")
	}

	return &DeployerProvider{
		config:    config,
		logger:    logger.NewNilLogger(),
		sdkClient: client,
	}, nil
}

func (d *DeployerProvider) WithLogger(logger logger.Logger) *DeployerProvider {
	d.logger = logger
	return d
}

func (d *DeployerProvider) Deploy(ctx context.Context, certPem string, privkeyPem string) (*deployer.DeployResult, error) {
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

	d.logger.Logt("已修改域名证书", putCertResp)

	return &deployer.DeployResult{}, nil
}

func createSdkClient(accessKeyId, secretAccessKey string) (*bceCdn.Client, error) {
	client, err := bceCdn.NewClient(accessKeyId, secretAccessKey, "")
	if err != nil {
		return nil, err
	}

	return client, nil
}
