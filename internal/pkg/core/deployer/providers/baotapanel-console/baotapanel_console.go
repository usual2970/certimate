package baotapanelconsole

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"
	"net/url"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	btsdk "github.com/usual2970/certimate/internal/pkg/sdk3rd/btpanel"
)

type DeployerConfig struct {
	// 宝塔面板服务地址。
	ServerUrl string `json:"serverUrl"`
	// 宝塔面板接口密钥。
	ApiKey string `json:"apiKey"`
	// 是否允许不安全的连接。
	AllowInsecureConnections bool `json:"allowInsecureConnections,omitempty"`
	// 是否自动重启。
	AutoRestart bool `json:"autoRestart"`
}

type DeployerProvider struct {
	config    *DeployerConfig
	logger    *slog.Logger
	sdkClient *btsdk.Client
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
		d.logger = slog.Default()
	} else {
		d.logger = logger
	}
	return d
}

func (d *DeployerProvider) Deploy(ctx context.Context, certPEM string, privkeyPEM string) (*deployer.DeployResult, error) {
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

	return &deployer.DeployResult{}, nil
}

func createSdkClient(serverUrl, apiKey string, skipTlsVerify bool) (*btsdk.Client, error) {
	if _, err := url.Parse(serverUrl); err != nil {
		return nil, errors.New("invalid baota server url")
	}

	if apiKey == "" {
		return nil, errors.New("invalid baota api key")
	}

	client := btsdk.NewClient(serverUrl, apiKey)
	if skipTlsVerify {
		client.WithTLSConfig(&tls.Config{InsecureSkipVerify: true})
	}

	return client, nil
}
