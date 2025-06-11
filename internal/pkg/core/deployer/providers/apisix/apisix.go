package apisix

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"
	"net/url"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	apisixsdk "github.com/usual2970/certimate/internal/pkg/sdk3rd/apisix"
	certutil "github.com/usual2970/certimate/internal/pkg/utils/cert"
	typeutil "github.com/usual2970/certimate/internal/pkg/utils/type"
)

type DeployerConfig struct {
	// APISIX 服务地址。
	ServerUrl string `json:"serverUrl"`
	// APISIX Admin API Key。
	ApiKey string `json:"apiKey"`
	// 是否允许不安全的连接。
	AllowInsecureConnections bool `json:"allowInsecureConnections,omitempty"`
	// 部署资源类型。
	ResourceType ResourceType `json:"resourceType"`
	// 证书 ID。
	// 部署资源类型为 [RESOURCE_TYPE_CERTIFICATE] 时必填。
	CertificateId string `json:"certificateId,omitempty"`
}

type DeployerProvider struct {
	config    *DeployerConfig
	logger    *slog.Logger
	sdkClient *apisixsdk.Client
}

var _ deployer.Deployer = (*DeployerProvider)(nil)

func NewDeployer(config *DeployerConfig) (*DeployerProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	client, err := createSdkClient(config.ServerUrl, config.ApiKey, config.AllowInsecureConnections)
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
		d.logger = slog.New(slog.DiscardHandler)
	} else {
		d.logger = logger
	}
	return d
}

func (d *DeployerProvider) Deploy(ctx context.Context, certPEM string, privkeyPEM string) (*deployer.DeployResult, error) {
	// 根据部署资源类型决定部署方式
	switch d.config.ResourceType {
	case RESOURCE_TYPE_CERTIFICATE:
		if err := d.deployToCertificate(ctx, certPEM, privkeyPEM); err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("unsupported resource type '%s'", d.config.ResourceType)
	}

	return &deployer.DeployResult{}, nil
}

func (d *DeployerProvider) deployToCertificate(ctx context.Context, certPEM string, privkeyPEM string) error {
	if d.config.CertificateId == "" {
		return errors.New("config `certificateId` is required")
	}

	// 解析证书内容
	certX509, err := certutil.ParseCertificateFromPEM(certPEM)
	if err != nil {
		return err
	}

	// 更新 SSL 证书
	// REF: https://apisix.apache.org/zh/docs/apisix/admin-api/#ssl
	updateSSLReq := &apisixsdk.UpdateSSLRequest{
		ID:     d.config.CertificateId,
		Cert:   typeutil.ToPtr(certPEM),
		Key:    typeutil.ToPtr(privkeyPEM),
		SNIs:   typeutil.ToPtr(certX509.DNSNames),
		Type:   typeutil.ToPtr("server"),
		Status: typeutil.ToPtr(int32(1)),
	}
	updateSSLResp, err := d.sdkClient.UpdateSSL(updateSSLReq)
	d.logger.Debug("sdk request 'apisix.UpdateSSL'", slog.Any("request", updateSSLReq), slog.Any("response", updateSSLResp))
	if err != nil {
		return fmt.Errorf("failed to execute sdk request 'apisix.UpdateSSL': %w", err)
	}

	return nil
}

func createSdkClient(serverUrl, apiKey string, skipTlsVerify bool) (*apisixsdk.Client, error) {
	if _, err := url.Parse(serverUrl); err != nil {
		return nil, errors.New("invalid apisix server url")
	}

	if apiKey == "" {
		return nil, errors.New("invalid apisix api key")
	}

	client := apisixsdk.NewClient(serverUrl, apiKey)
	if skipTlsVerify {
		client.WithTLSConfig(&tls.Config{InsecureSkipVerify: true})
	}

	return client, nil
}
