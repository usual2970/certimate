package netlifysite

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/certimate-go/certimate/pkg/core"
	netlifysdk "github.com/certimate-go/certimate/pkg/sdk3rd/netlify"
	xcert "github.com/certimate-go/certimate/pkg/utils/cert"
)

type SSLDeployerProviderConfig struct {
	// netlify API Token。
	ApiToken string `json:"apiToken"`
	// netlify 网站 ID。
	SiteId string `json:"siteId"`
}

type SSLDeployerProvider struct {
	config    *SSLDeployerProviderConfig
	logger    *slog.Logger
	sdkClient *netlifysdk.Client
}

var _ core.SSLDeployer = (*SSLDeployerProvider)(nil)

func NewSSLDeployerProvider(config *SSLDeployerProviderConfig) (*SSLDeployerProvider, error) {
	if config == nil {
		return nil, errors.New("the configuration of the ssl deployer provider is nil")
	}

	client, err := createSDKClient(config.ApiToken)
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
	if d.config.SiteId == "" {
		return nil, errors.New("config `siteId` is required")
	}

	// 提取服务器证书和中间证书
	serverCertPEM, intermediaCertPEM, err := xcert.ExtractCertificatesFromPEM(certPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to extract certs: %w", err)
	}

	// 上传网站证书
	// REF: https://open-api.netlify.com/#tag/sniCertificate/operation/provisionSiteTLSCertificate
	provisionSiteTLSCertificateReq := &netlifysdk.ProvisionSiteTLSCertificateParams{
		Certificate:    serverCertPEM,
		CACertificates: intermediaCertPEM,
		Key:            privkeyPEM,
	}
	provisionSiteTLSCertificateResp, err := d.sdkClient.ProvisionSiteTLSCertificate(d.config.SiteId, provisionSiteTLSCertificateReq)
	d.logger.Debug("sdk request 'netlify.provisionSiteTLSCertificate'", slog.String("siteId", d.config.SiteId), slog.Any("request", provisionSiteTLSCertificateReq), slog.Any("response", provisionSiteTLSCertificateResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'netlify.provisionSiteTLSCertificate': %w", err)
	}

	return &core.SSLDeployResult{}, nil
}

func createSDKClient(apiToken string) (*netlifysdk.Client, error) {
	return netlifysdk.NewClient(apiToken)
}
