package cachefly

import (
	"context"
	"errors"

	xerrors "github.com/pkg/errors"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/logger"
	cfsdk "github.com/usual2970/certimate/internal/pkg/vendors/cachefly-sdk"
)

type DeployerConfig struct {
	// CacheFly API Token。
	ApiToken string `json:"apiToken"`
}

type DeployerProvider struct {
	config    *DeployerConfig
	logger    logger.Logger
	sdkClient *cfsdk.Client
}

var _ deployer.Deployer = (*DeployerProvider)(nil)

func NewDeployer(config *DeployerConfig) (*DeployerProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	client, err := createSdkClient(config.ApiToken)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk client")
	}

	return &DeployerProvider{
		config:    config,
		logger:    logger.NewNilLogger(),
		sdkClient: client,
	}, nil
}

func (d *DeployerProvider) WithLogger(logger logger.Logger) *DeployerProvider {
	d.logger = logger
	return d
}

func (d *DeployerProvider) Deploy(ctx context.Context, certPem string, privkeyPem string) (*deployer.DeployResult, error) {
	// 上传证书
	createCertificateReq := &cfsdk.CreateCertificateRequest{
		Certificate:    certPem,
		CertificateKey: privkeyPem,
	}
	createCertificateResp, err := d.sdkClient.CreateCertificate(createCertificateReq)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'cachefly.CreateCertificate'")
	} else {
		d.logger.Logt("已上传证书", createCertificateResp)
	}

	return &deployer.DeployResult{}, nil
}

func createSdkClient(apiToken string) (*cfsdk.Client, error) {
	if apiToken == "" {
		return nil, errors.New("invalid cachefly api token")
	}

	client := cfsdk.NewClient(apiToken)
	return client, nil
}
