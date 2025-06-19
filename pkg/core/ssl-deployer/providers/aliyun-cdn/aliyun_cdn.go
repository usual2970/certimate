package aliyuncdn

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	alicdn "github.com/alibabacloud-go/cdn-20180510/v5/client"
	aliopen "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/certimate-go/certimate/pkg/core"
)

type SSLDeployerProviderConfig struct {
	// 阿里云 AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// 阿里云 AccessKeySecret。
	AccessKeySecret string `json:"accessKeySecret"`
	// 阿里云资源组 ID。
	ResourceGroupId string `json:"resourceGroupId,omitempty"`
	// 加速域名（支持泛域名）。
	Domain string `json:"domain"`
}

type SSLDeployerProvider struct {
	config    *SSLDeployerProviderConfig
	logger    *slog.Logger
	sdkClient *alicdn.Client
}

var _ core.SSLDeployer = (*SSLDeployerProvider)(nil)

func NewSSLDeployerProvider(config *SSLDeployerProviderConfig) (*SSLDeployerProvider, error) {
	if config == nil {
		return nil, errors.New("the configuration of the ssl deployer provider is nil")
	}

	client, err := createSDKClient(config.AccessKeyId, config.AccessKeySecret)
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
	// "*.example.com" → ".example.com"，适配阿里云 CDN 要求的泛域名格式
	domain := strings.TrimPrefix(d.config.Domain, "*")

	// 设置 CDN 域名域名证书
	// REF: https://help.aliyun.com/zh/cdn/developer-reference/api-cdn-2018-05-10-setcdndomainsslcertificate
	setCdnDomainSSLCertificateReq := &alicdn.SetCdnDomainSSLCertificateRequest{
		DomainName:  tea.String(domain),
		CertName:    tea.String(fmt.Sprintf("certimate-%d", time.Now().UnixMilli())),
		CertType:    tea.String("upload"),
		SSLProtocol: tea.String("on"),
		SSLPub:      tea.String(certPEM),
		SSLPri:      tea.String(privkeyPEM),
	}
	setCdnDomainSSLCertificateResp, err := d.sdkClient.SetCdnDomainSSLCertificate(setCdnDomainSSLCertificateReq)
	d.logger.Debug("sdk request 'cdn.SetCdnDomainSSLCertificate'", slog.Any("request", setCdnDomainSSLCertificateReq), slog.Any("response", setCdnDomainSSLCertificateResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'cdn.SetCdnDomainSSLCertificate': %w", err)
	}

	return &core.SSLDeployResult{}, nil
}

func createSDKClient(accessKeyId, accessKeySecret string) (*alicdn.Client, error) {
	config := &aliopen.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
		Endpoint:        tea.String("cdn.aliyuncs.com"),
	}

	client, err := alicdn.NewClient(config)
	if err != nil {
		return nil, err
	}

	return client, nil
}
