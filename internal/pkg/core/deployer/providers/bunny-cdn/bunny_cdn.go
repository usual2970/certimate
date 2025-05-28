package bunnycdn

import (
	"context"
	"encoding/base64"
	"fmt"
	"log/slog"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	bunnysdk "github.com/usual2970/certimate/internal/pkg/sdk3rd/bunny"
)

type DeployerConfig struct {
	// Bunny API Key。
	ApiKey string `json:"apiKey"`
	// Bunny Pull Zone ID。
	PullZoneId string `json:"pullZoneId"`
	// Bunny CDN Hostname（支持泛域名）。
	Hostname string `json:"hostname"`
}

type DeployerProvider struct {
	config    *DeployerConfig
	logger    *slog.Logger
	sdkClient *bunnysdk.Client
}

var _ deployer.Deployer = (*DeployerProvider)(nil)

func NewDeployer(config *DeployerConfig) (*DeployerProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	return &DeployerProvider{
		config:    config,
		logger:    slog.Default(),
		sdkClient: bunnysdk.NewClient(config.ApiKey),
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
	createCertificateReq := &bunnysdk.AddCustomCertificateRequest{
		Hostname:       d.config.Hostname,
		PullZoneId:     d.config.PullZoneId,
		Certificate:    base64.StdEncoding.EncodeToString([]byte(certPEM)),
		CertificateKey: base64.StdEncoding.EncodeToString([]byte(privkeyPEM)),
	}
	createCertificateResp, err := d.sdkClient.AddCustomCertificate(createCertificateReq)
	d.logger.Debug("sdk request 'bunny.AddCustomCertificate'", slog.Any("request", createCertificateReq), slog.Any("response", createCertificateResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'bunny.AddCustomCertificate': %w", err)
	}

	return &deployer.DeployResult{}, nil
}
