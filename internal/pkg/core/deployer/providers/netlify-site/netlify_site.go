package netlifysite

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/netlify/open-api/v2/go/porcelain"
	porcelainctx "github.com/netlify/open-api/v2/go/porcelain/context"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	certutil "github.com/usual2970/certimate/internal/pkg/utils/cert"
)

type DeployerConfig struct {
	// netlify API Token。
	ApiToken string `json:"apiToken"`
	// netlify 网站 ID。
	SiteId string `json:"siteId"`
}

type DeployerProvider struct {
	config          *DeployerConfig
	logger          *slog.Logger
	sdkClient       *porcelain.Netlify
	sdkClientAuther runtime.ClientAuthInfoWriter
}

var _ deployer.Deployer = (*DeployerProvider)(nil)

func NewDeployer(config *DeployerConfig) (*DeployerProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	client, clientAuther, err := createSdkClient(config.ApiToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create sdk client: %w", err)
	}

	return &DeployerProvider{
		config:          config,
		logger:          slog.Default(),
		sdkClient:       client,
		sdkClientAuther: clientAuther,
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
	configureSiteTLSCertificateCtx := porcelainctx.WithAuthInfo(context.TODO(), d.sdkClientAuther)
	configureSiteTLSCertificateReq := &porcelain.CustomTLSCertificate{
		Certificate:    serverCertPEM,
		CACertificates: intermediaCertPEM,
		Key:            privkeyPEM,
	}
	configureSiteTLSCertificateResp, err := d.sdkClient.ConfigureSiteTLSCertificate(configureSiteTLSCertificateCtx, d.config.SiteId, configureSiteTLSCertificateReq)
	d.logger.Debug("sdk request 'netlify.provisionSiteTLSCertificate'", slog.String("siteId", d.config.SiteId), slog.Any("request", configureSiteTLSCertificateReq), slog.Any("response", configureSiteTLSCertificateResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'netlify.provisionSiteTLSCertificate': %w", err)
	}

	return &deployer.DeployResult{}, nil
}

func createSdkClient(apiToken string) (*porcelain.Netlify, runtime.ClientAuthInfoWriter, error) {
	if apiToken == "" {
		return nil, nil, errors.New("invalid netlify api token")
	}

	creds := runtime.ClientAuthInfoWriterFunc(func(r runtime.ClientRequest, _ strfmt.Registry) error {
		r.SetHeaderParam("User-Agent", "Certimate")
		r.SetHeaderParam("Authorization", "Bearer "+apiToken)
		return nil
	})

	return porcelain.Default, creds, nil
}
