package baotapanelconsole

import (
	"context"
	"errors"
	"net/url"

	xerrors "github.com/pkg/errors"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/logger"
	btsdk "github.com/usual2970/certimate/internal/pkg/vendors/btpanel-sdk"
)

type BaotaPanelConsoleDeployerConfig struct {
	// 宝塔面板地址。
	ApiUrl string `json:"apiUrl"`
	// 宝塔面板接口密钥。
	ApiKey string `json:"apiKey"`
	// 是否自动重启。
	AutoRestart bool `json:"autoRestart"`
}

type BaotaPanelConsoleDeployer struct {
	config    *BaotaPanelConsoleDeployerConfig
	logger    logger.Logger
	sdkClient *btsdk.Client
}

var _ deployer.Deployer = (*BaotaPanelConsoleDeployer)(nil)

func New(config *BaotaPanelConsoleDeployerConfig) (*BaotaPanelConsoleDeployer, error) {
	return NewWithLogger(config, logger.NewNilLogger())
}

func NewWithLogger(config *BaotaPanelConsoleDeployerConfig, logger logger.Logger) (*BaotaPanelConsoleDeployer, error) {
	if config == nil {
		panic("config is nil")
	}

	if logger == nil {
		panic("logger is nil")
	}

	client, err := createSdkClient(config.ApiUrl, config.ApiKey)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk client")
	}

	return &BaotaPanelConsoleDeployer{
		logger:    logger,
		config:    config,
		sdkClient: client,
	}, nil
}

func (d *BaotaPanelConsoleDeployer) Deploy(ctx context.Context, certPem string, privkeyPem string) (*deployer.DeployResult, error) {
	// 设置面板 SSL 证书
	configSavePanelSSLReq := &btsdk.ConfigSavePanelSSLRequest{
		PrivateKey:  privkeyPem,
		Certificate: certPem,
	}
	configSavePanelSSLResp, err := d.sdkClient.ConfigSavePanelSSL(configSavePanelSSLReq)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'bt.ConfigSavePanelSSL'")
	} else {
		d.logger.Logt("已设置面板 SSL 证书", configSavePanelSSLResp)
	}

	if d.config.AutoRestart {
		// 重启面板
		systemServiceAdminReq := &btsdk.SystemServiceAdminRequest{
			Name: "nginx",
			Type: "restart",
		}
		_, err := d.sdkClient.SystemServiceAdmin(systemServiceAdminReq)
		if err != nil {
			return nil, xerrors.Wrap(err, "failed to execute sdk request 'bt.SystemServiceAdmin'")
		}
	}

	return &deployer.DeployResult{}, nil
}

func createSdkClient(apiUrl, apiKey string) (*btsdk.Client, error) {
	if _, err := url.Parse(apiUrl); err != nil {
		return nil, errors.New("invalid baota api url")
	}

	if apiKey == "" {
		return nil, errors.New("invalid baota api key")
	}

	client := btsdk.NewClient(apiUrl, apiKey)
	return client, nil
}
