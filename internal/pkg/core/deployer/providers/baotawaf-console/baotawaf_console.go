package baotapanelconsole

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	btwafsdk "github.com/usual2970/certimate/internal/pkg/sdk3rd/btwaf"
	xtypes "github.com/usual2970/certimate/internal/pkg/utils/types"
)

type DeployerConfig struct {
	// 堡塔云 WAF 服务地址。
	ServerUrl string `json:"serverUrl"`
	// 堡塔云 WAF 接口密钥。
	ApiKey string `json:"apiKey"`
	// 是否允许不安全的连接。
	AllowInsecureConnections bool `json:"allowInsecureConnections,omitempty"`
}

type DeployerProvider struct {
	config    *DeployerConfig
	logger    *slog.Logger
	sdkClient *btwafsdk.Client
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
	// 设置面板 SSL
	configSetCertReq := &btwafsdk.ConfigSetCertRequest{
		CertContent: xtypes.ToPtr(certPEM),
		KeyContent:  xtypes.ToPtr(privkeyPEM),
	}
	configSetCertResp, err := d.sdkClient.ConfigSetCert(configSetCertReq)
	d.logger.Debug("sdk request 'bt.ConfigSetCert'", slog.Any("request", configSetCertReq), slog.Any("response", configSetCertResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'bt.ConfigSetCert': %w", err)
	}

	return &deployer.DeployResult{}, nil
}

func createSdkClient(serverUrl, apiKey string, skipTlsVerify bool) (*btwafsdk.Client, error) {
	client, err := btwafsdk.NewClient(serverUrl, apiKey)
	if err != nil {
		return nil, err
	}

	if skipTlsVerify {
		client.SetTLSConfig(&tls.Config{InsecureSkipVerify: true})
	}

	return client, nil
}
