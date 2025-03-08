package onepanelconsole

import (
	"context"
	"crypto/tls"
	"errors"
	"net/url"

	xerrors "github.com/pkg/errors"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/logger"
	opsdk "github.com/usual2970/certimate/internal/pkg/vendors/1panel-sdk"
)

type DeployerConfig struct {
	// 1Panel 地址。
	ApiUrl string `json:"apiUrl"`
	// 1Panel 接口密钥。
	ApiKey string `json:"apiKey"`
	// 是否允许不安全的连接。
	AllowInsecureConnections bool `json:"allowInsecureConnections,omitempty"`
	// 是否自动重启。
	AutoRestart bool `json:"autoRestart"`
}

type DeployerProvider struct {
	config    *DeployerConfig
	logger    logger.Logger
	sdkClient *opsdk.Client
}

var _ deployer.Deployer = (*DeployerProvider)(nil)

func NewDeployer(config *DeployerConfig) (*DeployerProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	client, err := createSdkClient(config.ApiUrl, config.ApiKey, config.AllowInsecureConnections)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk client")
	}

	return &DeployerProvider{
		config:    config,
		logger:    logger.NewNilLogger(),
		sdkClient: client,
	}, nil
}

func (d *DeployerProvider) WithLogger(logger logger.Logger) *DeployerProvider {
	d.logger = logger
	return d
}

func (d *DeployerProvider) Deploy(ctx context.Context, certPem string, privkeyPem string) (*deployer.DeployResult, error) {
	// 设置面板 SSL 证书
	updateSystemSSLReq := &opsdk.UpdateSystemSSLRequest{
		Cert:    certPem,
		Key:     privkeyPem,
		SSL:     "enable",
		SSLType: "import-paste",
	}
	if d.config.AutoRestart {
		updateSystemSSLReq.AutoRestart = "true"
	} else {
		updateSystemSSLReq.AutoRestart = "false"
	}
	updateSystemSSLResp, err := d.sdkClient.UpdateSystemSSL(updateSystemSSLReq)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request '1panel.UpdateSystemSSL'")
	} else {
		d.logger.Logt("已设置面板 SSL 证书", updateSystemSSLResp)
	}

	return &deployer.DeployResult{}, nil
}

func createSdkClient(apiUrl, apiKey string, allowInsecure bool) (*opsdk.Client, error) {
	if _, err := url.Parse(apiUrl); err != nil {
		return nil, errors.New("invalid 1panel api url")
	}

	if apiKey == "" {
		return nil, errors.New("invalid 1panel api key")
	}

	client := opsdk.NewClient(apiUrl, apiKey)
	if allowInsecure {
		client.WithTLSConfig(&tls.Config{InsecureSkipVerify: true})
	}

	return client, nil
}
