package ratpanelconsole

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	rpsdk "github.com/usual2970/certimate/internal/pkg/sdk3rd/ratpanel"
)

type DeployerConfig struct {
	// 耗子面板服务地址。
	ServerUrl string `json:"serverUrl"`
	// 耗子面板访问令牌 ID。
	AccessTokenId int32 `json:"accessTokenId"`
	// 耗子面板访问令牌。
	AccessToken string `json:"accessToken"`
	// 是否允许不安全的连接。
	AllowInsecureConnections bool `json:"allowInsecureConnections,omitempty"`
}

type DeployerProvider struct {
	config    *DeployerConfig
	logger    *slog.Logger
	sdkClient *rpsdk.Client
}

var _ deployer.Deployer = (*DeployerProvider)(nil)

func NewDeployer(config *DeployerConfig) (*DeployerProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	client, err := createSdkClient(config.ServerUrl, config.AccessTokenId, config.AccessToken, config.AllowInsecureConnections)
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
	// 设置面板 SSL 证书
	setSettingCertReq := &rpsdk.SetSettingCertRequest{
		Certificate: certPEM,
		PrivateKey:  privkeyPEM,
	}
	setSettingCertResp, err := d.sdkClient.SetSettingCert(setSettingCertReq)
	d.logger.Debug("sdk request 'ratpanel.SetSettingCert'", slog.Any("request", setSettingCertReq), slog.Any("response", setSettingCertResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'ratpanel.SetSettingCert': %w", err)
	}

	return &deployer.DeployResult{}, nil
}

func createSdkClient(serverUrl string, accessTokenId int32, accessToken string, skipTlsVerify bool) (*rpsdk.Client, error) {
	client, err := rpsdk.NewClient(serverUrl, accessTokenId, accessToken)
	if err != nil {
		return nil, err
	}

	if skipTlsVerify {
		client.SetTLSConfig(&tls.Config{InsecureSkipVerify: true})
	}

	return client, nil
}
