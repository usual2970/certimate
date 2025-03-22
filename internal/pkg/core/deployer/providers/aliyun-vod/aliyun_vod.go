package aliyunvod

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	aliopen "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	alivod "github.com/alibabacloud-go/vod-20170321/v4/client"
	xerrors "github.com/pkg/errors"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
)

type DeployerConfig struct {
	// 阿里云 AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// 阿里云 AccessKeySecret。
	AccessKeySecret string `json:"accessKeySecret"`
	// 阿里云地域。
	Region string `json:"region"`
	// 点播加速域名（不支持泛域名）。
	Domain string `json:"domain"`
}

type DeployerProvider struct {
	config    *DeployerConfig
	logger    *slog.Logger
	sdkClient *alivod.Client
}

var _ deployer.Deployer = (*DeployerProvider)(nil)

func NewDeployer(config *DeployerConfig) (*DeployerProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	client, err := createSdkClient(config.AccessKeyId, config.AccessKeySecret, config.Region)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk client")
	}

	return &DeployerProvider{
		config:    config,
		logger:    slog.Default(),
		sdkClient: client,
	}, nil
}

func (d *DeployerProvider) WithLogger(logger *slog.Logger) deployer.Deployer {
	if logger == nil {
		d.logger = slog.Default()
	} else {
		d.logger = logger
	}
	return d
}

func (d *DeployerProvider) Deploy(ctx context.Context, certPem string, privkeyPem string) (*deployer.DeployResult, error) {
	// 设置域名证书
	// REF: https://help.aliyun.com/zh/vod/developer-reference/api-vod-2017-03-21-setvoddomainsslcertificate
	setVodDomainSSLCertificateReq := &alivod.SetVodDomainSSLCertificateRequest{
		DomainName:  tea.String(d.config.Domain),
		CertName:    tea.String(fmt.Sprintf("certimate-%d", time.Now().UnixMilli())),
		CertType:    tea.String("upload"),
		SSLProtocol: tea.String("on"),
		SSLPub:      tea.String(certPem),
		SSLPri:      tea.String(privkeyPem),
	}
	setVodDomainSSLCertificateResp, err := d.sdkClient.SetVodDomainSSLCertificate(setVodDomainSSLCertificateReq)
	d.logger.Debug("sdk request 'live.SetVodDomainSSLCertificate'", slog.Any("request", setVodDomainSSLCertificateReq), slog.Any("response", setVodDomainSSLCertificateResp))
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'live.SetVodDomainSSLCertificate'")
	}

	return &deployer.DeployResult{}, nil
}

func createSdkClient(accessKeyId, accessKeySecret, region string) (*alivod.Client, error) {
	// 接入点一览 https://api.aliyun.com/product/vod
	endpoint := fmt.Sprintf("vod.%s.aliyuncs.com", region)

	config := &aliopen.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
		Endpoint:        tea.String(endpoint),
	}

	client, err := alivod.NewClient(config)
	if err != nil {
		return nil, err
	}

	return client, nil
}
