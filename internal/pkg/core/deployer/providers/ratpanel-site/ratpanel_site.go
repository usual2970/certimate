package ratpanelsite

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
	// 耗子面板服务地址。
	ServerUrl string `json:"serverUrl"`
	// 耗子面板访问令牌 ID。
	AccessTokenId int32 `json:"accessTokenId"`
	// 耗子面板访问令牌。
	AccessToken string `json:"accessToken"`
	// 是否允许不安全的连接。
	AllowInsecureConnections bool `json:"allowInsecureConnections,omitempty"`
	// 网站名称。
	SiteName string `json:"siteName"`
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
	if d.config.SiteName == "" {
		return nil, errors.New("config `siteName` is required")
	}

	// 设置站点 SSL 证书
	websiteCertReq := &rpsdk.WebsiteCertRequest{
		SiteName:    d.config.SiteName,
		Certificate: certPEM,
		PrivateKey:  privkeyPEM,
	}
	websiteCertResp, err := d.sdkClient.WebsiteCert(websiteCertReq)
	d.logger.Debug("sdk request 'ratpanel.WebsiteCert'", slog.Any("request", websiteCertReq), slog.Any("response", websiteCertResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'ratpanel.WebsiteCert': %w", err)
	}

	return &deployer.DeployResult{}, nil
}

func createSdkClient(serverUrl string, accessTokenId int32, accessToken string, skipTlsVerify bool) (*rpsdk.Client, error) {
	if _, err := url.Parse(serverUrl); err != nil {
		return nil, errors.New("invalid ratpanel server url")
	}

	if accessTokenId == 0 {
		return nil, errors.New("invalid ratpanel access token id")
	}

	if accessToken == "" {
		return nil, errors.New("invalid ratpanel access token")
	}

	client := rpsdk.NewClient(serverUrl, accessTokenId, accessToken)
	if skipTlsVerify {
		client.WithTLSConfig(&tls.Config{InsecureSkipVerify: true})
	}

	return client, nil
}
