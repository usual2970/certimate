package ratpanelconsole

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"
	"net/url"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	rpsdk "github.com/usual2970/certimate/internal/pkg/sdk3rd/ratpanel"
)

type DeployerConfig struct {
	// 耗子面板地址。
	ApiUrl string `json:"apiUrl"`
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

	client, err := createSdkClient(config.ApiUrl, config.AccessTokenId, config.AccessToken, config.AllowInsecureConnections)
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
	settingCertReq := &rpsdk.SettingCertRequest{
		Certificate: certPEM,
		PrivateKey:  privkeyPEM,
	}
	settingCertResp, err := d.sdkClient.SettingCert(settingCertReq)
	d.logger.Debug("sdk request 'ratpanel.SettingCert'", slog.Any("request", settingCertReq), slog.Any("response", settingCertResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'ratpanel.SettingCert': %w", err)
	}

	return &deployer.DeployResult{}, nil
}

func createSdkClient(apiUrl string, accessTokenId int32, accessToken string, skipTlsVerify bool) (*rpsdk.Client, error) {
	if _, err := url.Parse(apiUrl); err != nil {
		return nil, errors.New("invalid ratpanel api url")
	}

	if accessTokenId == 0 {
		return nil, errors.New("invalid ratpanel access token id")
	}

	if accessToken == "" {
		return nil, errors.New("invalid ratpanel access token")
	}

	client := rpsdk.NewClient(apiUrl, accessTokenId, accessToken)
	if skipTlsVerify {
		client.WithTLSConfig(&tls.Config{InsecureSkipVerify: true})
	}

	return client, nil
}
