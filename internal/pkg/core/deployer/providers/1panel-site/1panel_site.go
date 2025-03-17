package onepanelsite

import (
	"context"
	"crypto/tls"
	"errors"
	"log/slog"
	"net/url"
	"strconv"

	xerrors "github.com/pkg/errors"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	uploadersp "github.com/usual2970/certimate/internal/pkg/core/uploader/providers/1panel-ssl"
	opsdk "github.com/usual2970/certimate/internal/pkg/vendors/1panel-sdk"
)

type DeployerConfig struct {
	// 1Panel 地址。
	ApiUrl string `json:"apiUrl"`
	// 1Panel 接口密钥。
	ApiKey string `json:"apiKey"`
	// 是否允许不安全的连接。
	AllowInsecureConnections bool `json:"allowInsecureConnections,omitempty"`
	// 网站 ID。
	WebsiteId int64 `json:"websiteId"`
}

type DeployerProvider struct {
	config      *DeployerConfig
	logger      *slog.Logger
	sdkClient   *opsdk.Client
	sslUploader uploader.Uploader
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

	uploader, err := uploadersp.NewUploader(&uploadersp.UploaderConfig{
		ApiUrl: config.ApiUrl,
		ApiKey: config.ApiKey,
	})
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create ssl uploader")
	}

	return &DeployerProvider{
		config:      config,
		logger:      slog.Default(),
		sdkClient:   client,
		sslUploader: uploader,
	}, nil
}

func (d *DeployerProvider) WithLogger(logger *slog.Logger) deployer.Deployer {
	if logger == nil {
		d.logger = slog.Default()
	} else {
		d.logger = logger
	}
	d.sslUploader.WithLogger(logger)
	return d
}

func (d *DeployerProvider) Deploy(ctx context.Context, certPem string, privkeyPem string) (*deployer.DeployResult, error) {
	// 获取网站 HTTPS 配置
	getHttpsConfReq := &opsdk.GetHttpsConfRequest{
		WebsiteID: d.config.WebsiteId,
	}
	getHttpsConfResp, err := d.sdkClient.GetHttpsConf(getHttpsConfReq)
	d.logger.Debug("sdk request '1panel.GetHttpsConf'", slog.Any("request", getHttpsConfReq), slog.Any("response", getHttpsConfResp))
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request '1panel.GetHttpsConf'")
	}

	// 上传证书到面板
	upres, err := d.sslUploader.Upload(ctx, certPem, privkeyPem)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to upload certificate file")
	} else {
		d.logger.Info("ssl certificate uploaded", slog.Any("result", upres))
	}

	// 修改网站 HTTPS 配置
	certId, _ := strconv.ParseInt(upres.CertId, 10, 64)
	updateHttpsConfReq := &opsdk.UpdateHttpsConfRequest{
		WebsiteID:    d.config.WebsiteId,
		Type:         "existed",
		WebsiteSSLID: certId,
		Enable:       getHttpsConfResp.Data.Enable,
		HttpConfig:   getHttpsConfResp.Data.HttpConfig,
		SSLProtocol:  getHttpsConfResp.Data.SSLProtocol,
		Algorithm:    getHttpsConfResp.Data.Algorithm,
		Hsts:         getHttpsConfResp.Data.Hsts,
	}
	updateHttpsConfResp, err := d.sdkClient.UpdateHttpsConf(updateHttpsConfReq)
	d.logger.Debug("sdk request '1panel.UpdateHttpsConf'", slog.Any("request", updateHttpsConfReq), slog.Any("response", updateHttpsConfResp))
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request '1panel.UpdateHttpsConf'")
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
