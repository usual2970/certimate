package edgioapplications

import (
	"context"

	xerrors "github.com/pkg/errors"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/logger"
	"github.com/usual2970/certimate/internal/pkg/utils/certs"
	edgsdk "github.com/usual2970/certimate/internal/pkg/vendors/edgio-sdk/applications/v7"
	edgsdkDtos "github.com/usual2970/certimate/internal/pkg/vendors/edgio-sdk/applications/v7/dtos"
)

type EdgioApplicationsDeployerConfig struct {
	// Edgio ClientId。
	ClientId string `json:"clientId"`
	// Edgio ClientSecret。
	ClientSecret string `json:"clientSecret"`
	// Edgio 环境 ID。
	EnvironmentId string `json:"environmentId"`
}

type EdgioApplicationsDeployer struct {
	config    *EdgioApplicationsDeployerConfig
	logger    logger.Logger
	sdkClient *edgsdk.EdgioClient
}

var _ deployer.Deployer = (*EdgioApplicationsDeployer)(nil)

func New(config *EdgioApplicationsDeployerConfig) (*EdgioApplicationsDeployer, error) {
	return NewWithLogger(config, logger.NewNilLogger())
}

func NewWithLogger(config *EdgioApplicationsDeployerConfig, logger logger.Logger) (*EdgioApplicationsDeployer, error) {
	if config == nil {
		panic("config is nil")
	}

	if logger == nil {
		panic("logger is nil")
	}

	client, err := createSdkClient(config.ClientId, config.ClientSecret)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk client")
	}

	return &EdgioApplicationsDeployer{
		logger:    logger,
		config:    config,
		sdkClient: client,
	}, nil
}

func (d *EdgioApplicationsDeployer) Deploy(ctx context.Context, certPem string, privkeyPem string) (*deployer.DeployResult, error) {
	// 提取 Edgio 所需的服务端证书和中间证书内容
	privateCertPem, intermediateCertPem, err := certs.ExtractCertificatesFromPEM(certPem)
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
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'edgio.UploadTlsCert'")
	}

	d.logger.Logt("已上传 TLS 证书", uploadTlsCertResp)

	return &deployer.DeployResult{}, nil
}

func createSdkClient(clientId, clientSecret string) (*edgsdk.EdgioClient, error) {
	client := edgsdk.NewEdgioClient(clientId, clientSecret, "", "")
	return client, nil
}
