package aliyuncdn

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	aliyunCdn "github.com/alibabacloud-go/cdn-20180510/v5/client"
	aliyunOpen "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	xerrors "github.com/pkg/errors"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
)

type DeployerConfig struct {
	// 阿里云 AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// 阿里云 AccessKeySecret。
	AccessKeySecret string `json:"accessKeySecret"`
	// 加速域名（支持泛域名）。
	Domain string `json:"domain"`
}

type DeployerProvider struct {
	config    *DeployerConfig
	logger    *slog.Logger
	sdkClient *aliyunCdn.Client
}

var _ deployer.Deployer = (*DeployerProvider)(nil)

func NewDeployer(config *DeployerConfig) (*DeployerProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	client, err := createSdkClient(config.AccessKeyId, config.AccessKeySecret)
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
	// "*.example.com" → ".example.com"，适配阿里云 CDN 要求的泛域名格式
	domain := strings.TrimPrefix(d.config.Domain, "*")

	// 设置 CDN 域名域名证书
	// REF: https://help.aliyun.com/zh/cdn/developer-reference/api-cdn-2018-05-10-setcdndomainsslcertificate
	setCdnDomainSSLCertificateReq := &aliyunCdn.SetCdnDomainSSLCertificateRequest{
		DomainName:  tea.String(domain),
		CertName:    tea.String(fmt.Sprintf("certimate-%d", time.Now().UnixMilli())),
		CertType:    tea.String("upload"),
		SSLProtocol: tea.String("on"),
		SSLPub:      tea.String(certPem),
		SSLPri:      tea.String(privkeyPem),
	}
	setCdnDomainSSLCertificateResp, err := d.sdkClient.SetCdnDomainSSLCertificate(setCdnDomainSSLCertificateReq)
	d.logger.Debug("sdk request 'cdn.SetCdnDomainSSLCertificate'", slog.Any("request", setCdnDomainSSLCertificateReq), slog.Any("response", setCdnDomainSSLCertificateResp))
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'cdn.SetCdnDomainSSLCertificate'")
	}

	return &deployer.DeployResult{}, nil
}

func createSdkClient(accessKeyId, accessKeySecret string) (*aliyunCdn.Client, error) {
	config := &aliyunOpen.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
		Endpoint:        tea.String("cdn.aliyuncs.com"),
	}

	client, err := aliyunCdn.NewClient(config)
	if err != nil {
		return nil, err
	}

	return client, nil
}
