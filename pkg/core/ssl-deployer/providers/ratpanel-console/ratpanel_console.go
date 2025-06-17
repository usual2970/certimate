package ratpanelconsole

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"

	"github.com/certimate-go/certimate/pkg/core"
	rpsdk "github.com/certimate-go/certimate/pkg/sdk3rd/ratpanel"
)

type SSLDeployerProviderConfig struct {
	// 耗子面板服务地址。
	ServerUrl string `json:"serverUrl"`
	// 耗子面板访问令牌 ID。
	AccessTokenId int32 `json:"accessTokenId"`
	// 耗子面板访问令牌。
	AccessToken string `json:"accessToken"`
	// 是否允许不安全的连接。
	AllowInsecureConnections bool `json:"allowInsecureConnections,omitempty"`
}

type SSLDeployerProvider struct {
	config    *SSLDeployerProviderConfig
	logger    *slog.Logger
	sdkClient *rpsdk.Client
}

var _ core.SSLDeployer = (*SSLDeployerProvider)(nil)

func NewSSLDeployerProvider(config *SSLDeployerProviderConfig) (*SSLDeployerProvider, error) {
	if config == nil {
		return nil, errors.New("the configuration of the ssl deployer provider is nil")
	}

	client, err := createSDKClient(config.ServerUrl, config.AccessTokenId, config.AccessToken, config.AllowInsecureConnections)
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
	setSettingCertReq := &rpsdk.SetSettingCertRequest{
		Certificate: certPEM,
		PrivateKey:  privkeyPEM,
	}
	setSettingCertResp, err := d.sdkClient.SetSettingCert(setSettingCertReq)
	d.logger.Debug("sdk request 'ratpanel.SetSettingCert'", slog.Any("request", setSettingCertReq), slog.Any("response", setSettingCertResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'ratpanel.SetSettingCert': %w", err)
	}

	return &core.SSLDeployResult{}, nil
}

func createSDKClient(serverUrl string, accessTokenId int32, accessToken string, skipTlsVerify bool) (*rpsdk.Client, error) {
	client, err := rpsdk.NewClient(serverUrl, accessTokenId, accessToken)
	if err != nil {
		return nil, err
	}

	if skipTlsVerify {
		client.SetTLSConfig(&tls.Config{InsecureSkipVerify: true})
	}

	return client, nil
}
