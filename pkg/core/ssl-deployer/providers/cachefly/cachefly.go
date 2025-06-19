package cachefly

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/certimate-go/certimate/pkg/core"
	cacheflysdk "github.com/certimate-go/certimate/pkg/sdk3rd/cachefly"
	xtypes "github.com/certimate-go/certimate/pkg/utils/types"
)

type SSLDeployerProviderConfig struct {
	// CacheFly API Token。
	ApiToken string `json:"apiToken"`
}

type SSLDeployerProvider struct {
	config    *SSLDeployerProviderConfig
	logger    *slog.Logger
	sdkClient *cacheflysdk.Client
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

	return &core.SSLDeployResult{}, nil
}

func createSDKClient(apiToken string) (*cacheflysdk.Client, error) {
	return cacheflysdk.NewClient(apiToken)
}
