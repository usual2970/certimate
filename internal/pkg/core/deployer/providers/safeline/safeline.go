package safeline

import (
	"context"
	"errors"
	"fmt"

	xerrors "github.com/pkg/errors"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/logger"
	safelinesdk "github.com/usual2970/certimate/internal/pkg/vendors/safeline-sdk"
)

type SafeLineDeployerConfig struct {
	// 雷池 URL。
	ApiUrl string `json:"apiUrl"`
	// 雷池 API Token。
	ApiToken string `json:"apiToken"`
	// 部署资源类型。
	ResourceType ResourceType `json:"resourceType"`
	// 证书 ID。
	// 部署资源类型为 [RESOURCE_TYPE_CERTIFICATE] 时必填。
	CertificateId int32 `json:"certificateId,omitempty"`
}

type SafeLineDeployer struct {
	config    *SafeLineDeployerConfig
	logger    logger.Logger
	sdkClient *safelinesdk.Client
}

var _ deployer.Deployer = (*SafeLineDeployer)(nil)

func New(config *SafeLineDeployerConfig) (*SafeLineDeployer, error) {
	return NewWithLogger(config, logger.NewNilLogger())
}

func NewWithLogger(config *SafeLineDeployerConfig, logger logger.Logger) (*SafeLineDeployer, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}

	if logger == nil {
		return nil, errors.New("logger is nil")
	}

	client, err := createSdkClient(config.ApiUrl, config.ApiToken)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk clients")
	}

	return &SafeLineDeployer{
		logger:    logger,
		config:    config,
		sdkClient: client,
	}, nil
}

func (d *SafeLineDeployer) Deploy(ctx context.Context, certPem string, privkeyPem string) (*deployer.DeployResult, error) {
	// 根据部署资源类型决定部署方式
	switch d.config.ResourceType {
	case RESOURCE_TYPE_CERTIFICATE:
		if err := d.deployToCertificate(ctx, certPem, privkeyPem); err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("unsupported resource type: %s", d.config.ResourceType)
	}

	return &deployer.DeployResult{}, nil
}

func (d *SafeLineDeployer) deployToCertificate(ctx context.Context, certPem string, privkeyPem string) error {
	if d.config.CertificateId == 0 {
		return errors.New("config `certificateId` is required")
	}

	// 更新证书
	updateCertificateReq := &safelinesdk.UpdateCertificateRequest{
		Id:   d.config.CertificateId,
		Type: 2,
		Manual: &safelinesdk.UpdateCertificateRequestBodyManul{
			Crt: certPem,
			Key: privkeyPem,
		},
	}
	updateCertificateResp, err := d.sdkClient.UpdateCertificate(updateCertificateReq)
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'safeline.UpdateCertificate'")
	} else {
		d.logger.Logt("已更新证书", updateCertificateResp)
	}

	return nil
}

func createSdkClient(apiUrl, apiToken string) (*safelinesdk.Client, error) {
	client := safelinesdk.NewClient(apiUrl, apiToken)
	return client, nil
}
