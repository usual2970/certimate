package safeline

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"
	"net/url"

	xerrors "github.com/pkg/errors"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	safelinesdk "github.com/usual2970/certimate/internal/pkg/vendors/safeline-sdk"
)

type DeployerConfig struct {
	// 雷池 URL。
	ApiUrl string `json:"apiUrl"`
	// 雷池 API Token。
	ApiToken string `json:"apiToken"`
	// 是否允许不安全的连接。
	AllowInsecureConnections bool `json:"allowInsecureConnections,omitempty"`
	// 部署资源类型。
	ResourceType ResourceType `json:"resourceType"`
	// 证书 ID。
	// 部署资源类型为 [RESOURCE_TYPE_CERTIFICATE] 时必填。
	CertificateId int32 `json:"certificateId,omitempty"`
}

type DeployerProvider struct {
	config    *DeployerConfig
	logger    *slog.Logger
	sdkClient *safelinesdk.Client
}

var _ deployer.Deployer = (*DeployerProvider)(nil)

func NewDeployer(config *DeployerConfig) (*DeployerProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	client, err := createSdkClient(config.ApiUrl, config.ApiToken, config.AllowInsecureConnections)
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

func (d *DeployerProvider) deployToCertificate(ctx context.Context, certPem string, privkeyPem string) error {
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
	d.logger.Debug("sdk request 'safeline.UpdateCertificate'", slog.Any("request", updateCertificateReq), slog.Any("response", updateCertificateResp))
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'safeline.UpdateCertificate'")
	}

	return nil
}

func createSdkClient(apiUrl, apiToken string, allowInsecure bool) (*safelinesdk.Client, error) {
	if _, err := url.Parse(apiUrl); err != nil {
		return nil, errors.New("invalid safeline api url")
	}

	if apiToken == "" {
		return nil, errors.New("invalid safeline api token")
	}

	client := safelinesdk.NewClient(apiUrl, apiToken)
	if allowInsecure {
		client.WithTLSConfig(&tls.Config{InsecureSkipVerify: true})
	}

	return client, nil
}
