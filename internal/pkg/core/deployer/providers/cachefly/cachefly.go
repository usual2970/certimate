package cachefly

import (
	"context"
	"errors"

	xerrors "github.com/pkg/errors"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/logger"
	cfsdk "github.com/usual2970/certimate/internal/pkg/vendors/cachefly-sdk"
)

type CacheFlyDeployerConfig struct {
	// CacheFly API Token。
	ApiToken string `json:"apiToken"`
}

type CacheFlyDeployer struct {
	config    *CacheFlyDeployerConfig
	logger    logger.Logger
	sdkClient *cfsdk.Client
}

var _ deployer.Deployer = (*CacheFlyDeployer)(nil)

func New(config *CacheFlyDeployerConfig) (*CacheFlyDeployer, error) {
	return NewWithLogger(config, logger.NewNilLogger())
}

func NewWithLogger(config *CacheFlyDeployerConfig, logger logger.Logger) (*CacheFlyDeployer, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}

	if logger == nil {
		return nil, errors.New("logger is nil")
	}

	client, err := createSdkClient(config.ApiToken)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk client")
	}

	return &CacheFlyDeployer{
		logger:    logger,
		config:    config,
		sdkClient: client,
	}, nil
}

func (d *CacheFlyDeployer) Deploy(ctx context.Context, certPem string, privkeyPem string) (*deployer.DeployResult, error) {
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
	client := cfsdk.NewClient(apiToken)
	return client, nil
}
