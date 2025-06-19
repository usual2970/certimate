package baiducloudcdn

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	bcecdn "github.com/baidubce/bce-sdk-go/services/cdn"
	bcecdnapi "github.com/baidubce/bce-sdk-go/services/cdn/api"
	"github.com/certimate-go/certimate/pkg/core"
)

type SSLDeployerProviderConfig struct {
	// 百度智能云 AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// 百度智能云 SecretAccessKey。
	SecretAccessKey string `json:"secretAccessKey"`
	// 加速域名（支持泛域名）。
	Domain string `json:"domain"`
}

type SSLDeployerProvider struct {
	config    *SSLDeployerProviderConfig
	logger    *slog.Logger
	sdkClient *bcecdn.Client
}

var _ core.SSLDeployer = (*SSLDeployerProvider)(nil)

func NewSSLDeployerProvider(config *SSLDeployerProviderConfig) (*SSLDeployerProvider, error) {
	if config == nil {
		return nil, errors.New("the configuration of the ssl deployer provider is nil")
	}

	client, err := createSDKClient(config.AccessKeyId, config.SecretAccessKey)
	if err != nil {
		return nil, fmt.Errorf("could not create sdk client: %w", err)
	}

	return &SSLDeployerProvider{
		config:    config,
		logger:    slog.Default(),
		sdkClient: client,
	}, nil
}

func (d *SSLDeployerProvider) SetLogger(logger *slog.Logger) {
	if logger == nil {
		d.logger = slog.New(slog.DiscardHandler)
	} else {
		d.logger = logger
	}
}

func (d *SSLDeployerProvider) Deploy(ctx context.Context, certPEM string, privkeyPEM string) (*core.SSLDeployResult, error) {
	if d.config.Domain == "" {
		return nil, errors.New("config `domain` is required")
	}

	// 修改域名证书
	// REF: https://cloud.baidu.com/doc/CDN/s/qjzuz2hp8
	putCertResp, err := d.sdkClient.PutCert(
		d.config.Domain,
		&bcecdnapi.UserCertificate{
			CertName:    fmt.Sprintf("certimate-%d", time.Now().UnixMilli()),
			ServerData:  certPEM,
			PrivateData: privkeyPEM,
		},
		"ON",
	)
	d.logger.Debug("sdk request 'cdn.PutCert'", slog.String("request.domain", d.config.Domain), slog.Any("response", putCertResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'cdn.PutCert': %w", err)
	}

	return &core.SSLDeployResult{}, nil
}

func createSDKClient(accessKeyId, secretAccessKey string) (*bcecdn.Client, error) {
	client, err := bcecdn.NewClient(accessKeyId, secretAccessKey, "")
	if err != nil {
		return nil, err
	}

	return client, nil
}
