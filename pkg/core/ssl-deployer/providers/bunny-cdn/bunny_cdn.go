package bunnycdn

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"

	"github.com/certimate-go/certimate/pkg/core"
	bunnysdk "github.com/certimate-go/certimate/pkg/sdk3rd/bunny"
)

type SSLDeployerProviderConfig struct {
	// Bunny API Key。
	ApiKey string `json:"apiKey"`
	// Bunny Pull Zone ID。
	PullZoneId string `json:"pullZoneId"`
	// Bunny CDN Hostname（支持泛域名）。
	Hostname string `json:"hostname"`
}

type SSLDeployerProvider struct {
	config    *SSLDeployerProviderConfig
	logger    *slog.Logger
	sdkClient *bunnysdk.Client
}

var _ core.SSLDeployer = (*SSLDeployerProvider)(nil)

func NewSSLDeployerProvider(config *SSLDeployerProviderConfig) (*SSLDeployerProvider, error) {
	if config == nil {
		return nil, errors.New("the configuration of the ssl deployer provider is nil")
	}

	client, err := createSDKClient(config.ApiKey)
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
	if d.config.PullZoneId == "" {
		return nil, fmt.Errorf("config `pullZoneId` is required")
	}
	if d.config.Hostname == "" {
		return nil, fmt.Errorf("config `hostname` is required")
	}

	// 上传证书
	createCertificateReq := &bunnysdk.AddCustomCertificateRequest{
		Hostname:       d.config.Hostname,
		Certificate:    base64.StdEncoding.EncodeToString([]byte(certPEM)),
		CertificateKey: base64.StdEncoding.EncodeToString([]byte(privkeyPEM)),
	}
	err := d.sdkClient.AddCustomCertificate(d.config.PullZoneId, createCertificateReq)
	d.logger.Debug("sdk request 'bunny.AddCustomCertificate'", slog.Any("request", createCertificateReq))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'bunny.AddCustomCertificate': %w", err)
	}

	return &core.SSLDeployResult{}, nil
}

func createSDKClient(apiKey string) (*bunnysdk.Client, error) {
	return bunnysdk.NewClient(apiKey)
}
