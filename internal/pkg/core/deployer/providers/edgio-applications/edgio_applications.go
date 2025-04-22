package edgioapplications

import (
	"context"
	"fmt"
	"log/slog"

	edgio "github.com/Edgio/edgio-api/applications/v7"
	edgiodtos "github.com/Edgio/edgio-api/applications/v7/dtos"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	certutil "github.com/usual2970/certimate/internal/pkg/utils/cert"
)

type DeployerConfig struct {
	// Edgio ClientId。
	ClientId string `json:"clientId"`
	// Edgio ClientSecret。
	ClientSecret string `json:"clientSecret"`
	// Edgio 环境 ID。
	EnvironmentId string `json:"environmentId"`
}

type DeployerProvider struct {
	config    *DeployerConfig
	logger    *slog.Logger
	sdkClient *edgio.EdgioClient
}

var _ deployer.Deployer = (*DeployerProvider)(nil)

func NewDeployer(config *DeployerConfig) (*DeployerProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	client, err := createSdkClient(config.ClientId, config.ClientSecret)
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
		d.logger = slog.Default()
	} else {
		d.logger = logger
	}
	return d
}

func (d *DeployerProvider) Deploy(ctx context.Context, certPEM string, privkeyPEM string) (*deployer.DeployResult, error) {
	// 提取 Edgio 所需的服务端证书和中间证书内容
	privateCertPEM, intermediateCertPEM, err := certutil.ExtractCertificatesFromPEM(certPEM)
	if err != nil {
		return nil, err
	}

	// 上传 TLS 证书
	// REF: https://docs.edg.io/rest_api/#tag/tls-certs/operation/postConfigV01TlsCerts
	uploadTlsCertReq := edgiodtos.UploadTlsCertRequest{
		EnvironmentID:    d.config.EnvironmentId,
		PrimaryCert:      privateCertPEM,
		IntermediateCert: intermediateCertPEM,
		PrivateKey:       privkeyPEM,
	}
	uploadTlsCertResp, err := d.sdkClient.UploadTlsCert(uploadTlsCertReq)
	d.logger.Debug("sdk request 'edgio.UploadTlsCert'", slog.Any("request", uploadTlsCertReq), slog.Any("response", uploadTlsCertResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'edgio.UploadTlsCert': %w", err)
	}

	return &deployer.DeployResult{}, nil
}

func createSdkClient(clientId, clientSecret string) (*edgio.EdgioClient, error) {
	client := edgio.NewEdgioClient(clientId, clientSecret, "", "")
	return client, nil
}
