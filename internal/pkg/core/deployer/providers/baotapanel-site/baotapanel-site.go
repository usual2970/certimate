package baotapanelsite

import (
	"context"
	"errors"

	xerrors "github.com/pkg/errors"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/logger"
	btsdk "github.com/usual2970/certimate/internal/pkg/vendors/btpanel-sdk"
)

type BaotaPanelSiteDeployerConfig struct {
	// 宝塔面板地址。
	ApiUrl string `json:"apiUrl"`
	// 宝塔面板接口密钥。
	ApiKey string `json:"apiKey"`
	// 站点名称
	SiteName string `json:"siteName"`
}

type BaotaPanelSiteDeployer struct {
	config    *BaotaPanelSiteDeployerConfig
	logger    logger.Logger
	sdkClient *btsdk.BaoTaPanelClient
}

var _ deployer.Deployer = (*BaotaPanelSiteDeployer)(nil)

func New(config *BaotaPanelSiteDeployerConfig) (*BaotaPanelSiteDeployer, error) {
	return NewWithLogger(config, logger.NewNilLogger())
}

func NewWithLogger(config *BaotaPanelSiteDeployerConfig, logger logger.Logger) (*BaotaPanelSiteDeployer, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}

	if logger == nil {
		return nil, errors.New("logger is nil")
	}

	client, err := createSdkClient(config.ApiUrl, config.ApiKey)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk client")
	}

	return &BaotaPanelSiteDeployer{
		logger:    logger,
		config:    config,
		sdkClient: client,
	}, nil
}

func (d *BaotaPanelSiteDeployer) Deploy(ctx context.Context, certPem string, privkeyPem string) (*deployer.DeployResult, error) {
	if d.config.SiteName == "" {
		return nil, errors.New("config `siteName` is required")
	}

	// 设置站点 SSL 证书
	setSiteSSLReq := &btsdk.SetSiteSSLRequest{
		SiteName: d.config.SiteName,
		Type:     "1",
		Key:      privkeyPem,
		Csr:      certPem,
	}
	setSiteSSLResp, err := d.sdkClient.SetSiteSSL(setSiteSSLReq)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'bt.SetSiteSSL'")
	}

	d.logger.Logt("已设置站点 SSL 证书", setSiteSSLResp)

	return &deployer.DeployResult{}, nil
}

func createSdkClient(apiUrl, apiKey string) (*btsdk.BaoTaPanelClient, error) {
	client := btsdk.NewBaoTaPanelClient(apiUrl, apiKey)
	return client, nil
}
