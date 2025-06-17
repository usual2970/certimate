package cachefly

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	cacheflysdk "github.com/usual2970/certimate/internal/pkg/sdk3rd/cachefly"
	xtypes "github.com/usual2970/certimate/internal/pkg/utils/types"
)

type DeployerConfig struct {
	// CacheFly API Token。
	ApiToken string `json:"apiToken"`
}

type DeployerProvider struct {
	config    *DeployerConfig
	logger    *slog.Logger
	sdkClient *cacheflysdk.Client
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
	// 上传证书
	// REF: https://api.cachefly.com/api/2.5/docs#tag/Certificates/paths/~1certificates/post
	createCertificateReq := &cacheflysdk.CreateCertificateRequest{
		Certificate:    xtypes.ToPtr(certPEM),
		CertificateKey: xtypes.ToPtr(privkeyPEM),
	}
	createCertificateResp, err := d.sdkClient.CreateCertificate(createCertificateReq)
	d.logger.Debug("sdk request 'cachefly.CreateCertificate'", slog.Any("request", createCertificateReq), slog.Any("response", createCertificateResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'cachefly.CreateCertificate': %w", err)
	}

	return &deployer.DeployResult{}, nil
}

func createSdkClient(apiToken string) (*cacheflysdk.Client, error) {
	return cacheflysdk.NewClient(apiToken)
}
