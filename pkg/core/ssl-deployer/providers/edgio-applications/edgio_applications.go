package edgioapplications

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	edgio "github.com/Edgio/edgio-api/applications/v7"
	edgiodtos "github.com/Edgio/edgio-api/applications/v7/dtos"

	"github.com/certimate-go/certimate/pkg/core"
	xcert "github.com/certimate-go/certimate/pkg/utils/cert"
)

type SSLDeployerProviderConfig struct {
	// Edgio ClientId。
	ClientId string `json:"clientId"`
	// Edgio ClientSecret。
	ClientSecret string `json:"clientSecret"`
	// Edgio 环境 ID。
	EnvironmentId string `json:"environmentId"`
}

type SSLDeployerProvider struct {
	config    *SSLDeployerProviderConfig
	logger    *slog.Logger
	sdkClient *edgio.EdgioClient
}

var _ core.SSLDeployer = (*SSLDeployerProvider)(nil)

func NewSSLDeployerProvider(config *SSLDeployerProviderConfig) (*SSLDeployerProvider, error) {
	if config == nil {
		return nil, errors.New("the configuration of the ssl deployer provider is nil")
	}

	client, err := createSDKClient(config.ClientId, config.ClientSecret)
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
	// 提取服务器证书和中间证书
	serverCertPEM, intermediaCertPEM, err := xcert.ExtractCertificatesFromPEM(certPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to extract certs: %w", err)
	}

	// 上传 TLS 证书
	// REF: https://docs.edg.io/rest_api/#tag/tls-certs/operation/postConfigV01TlsCerts
	uploadTlsCertReq := edgiodtos.UploadTlsCertRequest{
		EnvironmentID:    d.config.EnvironmentId,
		PrimaryCert:      serverCertPEM,
		IntermediateCert: intermediaCertPEM,
		PrivateKey:       privkeyPEM,
	}
	uploadTlsCertResp, err := d.sdkClient.UploadTlsCert(uploadTlsCertReq)
	d.logger.Debug("sdk request 'edgio.UploadTlsCert'", slog.Any("request", uploadTlsCertReq), slog.Any("response", uploadTlsCertResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'edgio.UploadTlsCert': %w", err)
	}

	return &core.SSLDeployResult{}, nil
}

func createSDKClient(clientId, clientSecret string) (*edgio.EdgioClient, error) {
	client := edgio.NewEdgioClient(clientId, clientSecret, "", "")
	return client, nil
}
