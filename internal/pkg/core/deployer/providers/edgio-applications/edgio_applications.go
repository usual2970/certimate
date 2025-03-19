package edgioapplications

import (
	"context"
	"log/slog"

	xerrors "github.com/pkg/errors"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/utils/certutil"
	edgsdk "github.com/usual2970/certimate/internal/pkg/vendors/edgio-sdk/applications/v7"
	edgsdkDtos "github.com/usual2970/certimate/internal/pkg/vendors/edgio-sdk/applications/v7/dtos"
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
	sdkClient *edgsdk.EdgioClient
}

var _ deployer.Deployer = (*DeployerProvider)(nil)

func NewDeployer(config *DeployerConfig) (*DeployerProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	client, err := createSdkClient(config.ClientId, config.ClientSecret)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk client")
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

func (d *DeployerProvider) Deploy(ctx context.Context, certPem string, privkeyPem string) (*deployer.DeployResult, error) {
	// 提取 Edgio 所需的服务端证书和中间证书内容
	privateCertPem, intermediateCertPem, err := certutil.ExtractCertificatesFromPEM(certPem)
	if err != nil {
		return nil, err
	}

	// 上传 TLS 证书
	// REF: https://docs.edg.io/rest_api/#tag/tls-certs/operation/postConfigV01TlsCerts
	uploadTlsCertReq := edgsdkDtos.UploadTlsCertRequest{
		EnvironmentID:    d.config.EnvironmentId,
		PrimaryCert:      privateCertPem,
		IntermediateCert: intermediateCertPem,
		PrivateKey:       privkeyPem,
	}
	uploadTlsCertResp, err := d.sdkClient.UploadTlsCert(uploadTlsCertReq)
	d.logger.Debug("sdk request 'edgio.UploadTlsCert'", slog.Any("request", uploadTlsCertReq), slog.Any("response", uploadTlsCertResp))
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'edgio.UploadTlsCert'")
	}

	return &deployer.DeployResult{}, nil
}

func createSdkClient(clientId, clientSecret string) (*edgsdk.EdgioClient, error) {
	client := edgsdk.NewEdgioClient(clientId, clientSecret, "", "")
	return client, nil
}
