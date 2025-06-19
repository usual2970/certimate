package baotapanelconsole

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"

	"github.com/certimate-go/certimate/pkg/core"
	btsdk "github.com/certimate-go/certimate/pkg/sdk3rd/btpanel"
)

type SSLDeployerProviderConfig struct {
	// 宝塔面板服务地址。
	ServerUrl string `json:"serverUrl"`
	// 宝塔面板接口密钥。
	ApiKey string `json:"apiKey"`
	// 是否允许不安全的连接。
	AllowInsecureConnections bool `json:"allowInsecureConnections,omitempty"`
	// 是否自动重启。
	AutoRestart bool `json:"autoRestart"`
}

type SSLDeployerProvider struct {
	config    *SSLDeployerProviderConfig
	logger    *slog.Logger
	sdkClient *btsdk.Client
}

var _ core.SSLDeployer = (*SSLDeployerProvider)(nil)

func NewSSLDeployerProvider(config *SSLDeployerProviderConfig) (*SSLDeployerProvider, error) {
	if config == nil {
		return nil, errors.New("the configuration of the ssl deployer provider is nil")
	}

	client, err := createSDKClient(config.ServerUrl, config.ApiKey, config.AllowInsecureConnections)
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
	// 设置面板 SSL 证书
	configSavePanelSSLReq := &btsdk.ConfigSavePanelSSLRequest{
		PrivateKey:  privkeyPEM,
		Certificate: certPEM,
	}
	configSavePanelSSLResp, err := d.sdkClient.ConfigSavePanelSSL(configSavePanelSSLReq)
	d.logger.Debug("sdk request 'bt.ConfigSavePanelSSL'", slog.Any("request", configSavePanelSSLReq), slog.Any("response", configSavePanelSSLResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'bt.ConfigSavePanelSSL': %w", err)
	}

	if d.config.AutoRestart {
		// 重启面板（无需关心响应，因为宝塔重启时会断开连接产生 error）
		systemServiceAdminReq := &btsdk.SystemServiceAdminRequest{
			Name: "nginx",
			Type: "restart",
		}
		systemServiceAdminResp, _ := d.sdkClient.SystemServiceAdmin(systemServiceAdminReq)
		d.logger.Debug("sdk request 'bt.SystemServiceAdmin'", slog.Any("request", systemServiceAdminReq), slog.Any("response", systemServiceAdminResp))
	}

	return &core.SSLDeployResult{}, nil
}

func createSDKClient(serverUrl, apiKey string, skipTlsVerify bool) (*btsdk.Client, error) {
	client, err := btsdk.NewClient(serverUrl, apiKey)
	if err != nil {
		return nil, err
	}

	if skipTlsVerify {
		client.SetTLSConfig(&tls.Config{InsecureSkipVerify: true})
	}

	return client, nil
}
