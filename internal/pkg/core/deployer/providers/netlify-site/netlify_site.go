package netlifysite

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	netlifysdk "github.com/usual2970/certimate/internal/pkg/sdk3rd/netlify"
	certutil "github.com/usual2970/certimate/internal/pkg/utils/cert"
)

type DeployerConfig struct {
	// netlify API Token。
	ApiToken string `json:"apiToken"`
	// netlify 网站 ID。
	SiteId string `json:"siteId"`
}

type DeployerProvider struct {
	config    *DeployerConfig
	logger    *slog.Logger
	sdkClient *netlifysdk.Client
}

var _ deployer.Deployer = (*DeployerProvider)(nil)

func NewDeployer(config *DeployerConfig) (*DeployerProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	client, err := createSdkClient(config.ApiToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create sdk client: %w", err)
	}

	return &DeployerProvider{
		config:    config,
		logger:    slog.Default(),
		sdkClient: client,
	}, nil
}

func (d *DeployerProvider) WithLogger(logger *slog.Logger) deployer.Deployer {
	if logger == nil {
		d.logger = slog.New(slog.DiscardHandler)
	} else {
		d.logger = logger
	}
	return d
}

func (d *DeployerProvider) Deploy(ctx context.Context, certPEM string, privkeyPEM string) (*deployer.DeployResult, error) {
	if d.config.SiteId == "" {
		return nil, errors.New("config `siteId` is required")
	}

	// 提取服务器证书和中间证书
	serverCertPEM, intermediaCertPEM, err := certutil.ExtractCertificatesFromPEM(certPEM)
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

	return &deployer.DeployResult{}, nil
}

func createSdkClient(apiToken string) (*netlifysdk.Client, error) {
	if apiToken == "" {
		return nil, errors.New("invalid netlify api token")
	}

	client := netlifysdk.NewClient(apiToken)
	return client, nil
}
