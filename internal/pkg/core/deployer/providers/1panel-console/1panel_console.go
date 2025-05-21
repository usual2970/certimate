package onepanelconsole

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"
	"net/url"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	onepanelsdk "github.com/usual2970/certimate/internal/pkg/sdk3rd/1panel"
)

type DeployerConfig struct {
	// 1Panel 地址。
	ApiUrl string `json:"apiUrl"`
	// 1Panel 版本。
	// 可取值 "v1"、"v2"。
	ApiVersion string `json:"apiVersion"`
	// 1Panel 接口密钥。
	ApiKey string `json:"apiKey"`
	// 是否允许不安全的连接。
	AllowInsecureConnections bool `json:"allowInsecureConnections,omitempty"`
	// 是否自动重启。
	AutoRestart bool `json:"autoRestart"`
}

type DeployerProvider struct {
	config    *DeployerConfig
	logger    *slog.Logger
	sdkClient *onepanelsdk.Client
}

var _ deployer.Deployer = (*DeployerProvider)(nil)

func NewDeployer(config *DeployerConfig) (*DeployerProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	client, err := createSdkClient(config.ApiUrl, config.ApiVersion, config.ApiKey, config.AllowInsecureConnections)
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
	updateSystemSSLReq := &onepanelsdk.UpdateSystemSSLRequest{
		Cert:    certPEM,
		Key:     privkeyPEM,
		SSL:     "enable",
		SSLType: "import-paste",
	}
	if d.config.AutoRestart {
		updateSystemSSLReq.AutoRestart = "true"
	} else {
		updateSystemSSLReq.AutoRestart = "false"
	}
	updateSystemSSLResp, err := d.sdkClient.UpdateSystemSSL(updateSystemSSLReq)
	d.logger.Debug("sdk request '1panel.UpdateSystemSSL'", slog.Any("request", updateSystemSSLReq), slog.Any("response", updateSystemSSLResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request '1panel.UpdateSystemSSL': %w", err)
	}

	return &deployer.DeployResult{}, nil
}

func createSdkClient(apiUrl, apiVersion, apiKey string, skipTlsVerify bool) (*onepanelsdk.Client, error) {
	if _, err := url.Parse(apiUrl); err != nil {
		return nil, errors.New("invalid 1panel api url")
	}

	if apiVersion == "" {
		return nil, errors.New("invalid 1panel api version")
	}

	if apiKey == "" {
		return nil, errors.New("invalid 1panel api key")
	}

	client := onepanelsdk.NewClient(apiUrl, apiVersion, apiKey)
	if skipTlsVerify {
		client.WithTLSConfig(&tls.Config{InsecureSkipVerify: true})
	}

	return client, nil
}
