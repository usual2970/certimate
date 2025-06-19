package safeline

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"

	"github.com/certimate-go/certimate/pkg/core"
	safelinesdk "github.com/certimate-go/certimate/pkg/sdk3rd/safeline"
)

type SSLDeployerProviderConfig struct {
	// 雷池服务地址。
	ServerUrl string `json:"serverUrl"`
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

type SSLDeployerProvider struct {
	config    *SSLDeployerProviderConfig
	logger    *slog.Logger
	sdkClient *safelinesdk.Client
}

var _ core.SSLDeployer = (*SSLDeployerProvider)(nil)

func NewSSLDeployerProvider(config *SSLDeployerProviderConfig) (*SSLDeployerProvider, error) {
	if config == nil {
		return nil, errors.New("the configuration of the ssl deployer provider is nil")
	}

	client, err := createSDKClient(config.ServerUrl, config.ApiToken, config.AllowInsecureConnections)
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
	// 根据部署资源类型决定部署方式``
	switch d.config.ResourceType {
	case RESOURCE_TYPE_CERTIFICATE:
		if err := d.deployToCertificate(ctx, certPEM, privkeyPEM); err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("unsupported resource type '%s'", d.config.ResourceType)
	}

	return &core.SSLDeployResult{}, nil
}

func (d *SSLDeployerProvider) deployToCertificate(ctx context.Context, certPEM string, privkeyPEM string) error {
	if d.config.CertificateId == 0 {
		return errors.New("config `certificateId` is required")
	}

	// 更新证书
	updateCertificateReq := &safelinesdk.UpdateCertificateRequest{
		Id:   d.config.CertificateId,
		Type: 2,
		Manual: &safelinesdk.CertificateManul{
			Crt: certPEM,
			Key: privkeyPEM,
		},
	}
	updateCertificateResp, err := d.sdkClient.UpdateCertificate(updateCertificateReq)
	d.logger.Debug("sdk request 'safeline.UpdateCertificate'", slog.Any("request", updateCertificateReq), slog.Any("response", updateCertificateResp))
	if err != nil {
		return fmt.Errorf("failed to execute sdk request 'safeline.UpdateCertificate': %w", err)
	}

	return nil
}

func createSDKClient(serverUrl, apiToken string, skipTlsVerify bool) (*safelinesdk.Client, error) {
	client, err := safelinesdk.NewClient(serverUrl, apiToken)
	if err != nil {
		return nil, err
	}

	if skipTlsVerify {
		client.SetTLSConfig(&tls.Config{InsecureSkipVerify: true})
	}

	return client, nil
}
