package bunnycdn

import (
	"context"
	"encoding/base64"
	"log/slog"

	xerrors "github.com/pkg/errors"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	bunnysdk "github.com/usual2970/certimate/internal/pkg/vendors/bunny-sdk"
)

type DeployerConfig struct {
	// Bunny API Key
	ApiKey string `json:"apiKey"`
	// Bunny Pull Zone ID
	PullZoneId string `json:"pullZoneId"`
	// Bunny CDN Hostname（支持泛域名）
	HostName string `json:"hostName"`
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
		d.logger = slog.Default()
	} else {
		d.logger = logger
	}
	return d
}

func (d *DeployerProvider) Deploy(ctx context.Context, certPem string, privkeyPem string) (*deployer.DeployResult, error) {
	// Prepare
	certPemBase64 := base64.StdEncoding.EncodeToString([]byte(certPem))
	privkeyPemBase64 := base64.StdEncoding.EncodeToString([]byte(privkeyPem))
	// 上传证书
	createCertificateReq := &bunnysdk.AddCustomCertificateRequest{
		Hostname:       d.config.HostName,
		PullZoneId:     d.config.PullZoneId,
		Certificate:    certPemBase64,
		CertificateKey: privkeyPemBase64,
	}
	createCertificateResp, err := d.sdkClient.AddCustomCertificate(createCertificateReq)
	d.logger.Debug("sdk request 'bunny-cdn.AddCustomCertificate'", slog.Any("request", createCertificateReq), slog.Any("response", createCertificateResp))
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'bunny-cdn.AddCustomCertificate'")
	}

	return &deployer.DeployResult{}, nil
}
