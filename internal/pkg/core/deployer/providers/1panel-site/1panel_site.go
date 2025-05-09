package onepanelsite

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"strconv"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	uploadersp "github.com/usual2970/certimate/internal/pkg/core/uploader/providers/1panel-ssl"
	opsdk "github.com/usual2970/certimate/internal/pkg/sdk3rd/1panel"
)

type DeployerConfig struct {
	// 1Panel 地址。
	ApiUrl string `json:"apiUrl"`
	// 1Panel 接口密钥。
	ApiKey string `json:"apiKey"`
	// 是否允许不安全的连接。
	AllowInsecureConnections bool `json:"allowInsecureConnections,omitempty"`
	// 部署资源类型。
	ResourceType ResourceType `json:"resourceType"`
	// 网站 ID。
	// 部署资源类型为 [RESOURCE_TYPE_WEBSITE] 时必填。
	WebsiteId int64 `json:"websiteId,omitempty"`
	// 证书 ID。
	// 部署资源类型为 [RESOURCE_TYPE_CERTIFICATE] 时必填。
	CertificateId int64 `json:"certificateId,omitempty"`
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
		return nil, fmt.Errorf("failed to create sdk client: %w", err)
	}

	uploader, err := uploadersp.NewUploader(&uploadersp.UploaderConfig{
		ApiUrl: config.ApiUrl,
		ApiKey: config.ApiKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create ssl uploader: %w", err)
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

func (d *DeployerProvider) Deploy(ctx context.Context, certPEM string, privkeyPEM string) (*deployer.DeployResult, error) {
	// 根据部署资源类型决定部署方式
	switch d.config.ResourceType {
	case RESOURCE_TYPE_WEBSITE:
		if err := d.deployToWebsite(ctx, certPEM, privkeyPEM); err != nil {
			return nil, err
		}

	case RESOURCE_TYPE_CERTIFICATE:
		if err := d.deployToCertificate(ctx, certPEM, privkeyPEM); err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("unsupported resource type '%s'", d.config.ResourceType)
	}

	return &deployer.DeployResult{}, nil
}

func (d *DeployerProvider) deployToWebsite(ctx context.Context, certPEM string, privkeyPEM string) error {
	if d.config.WebsiteId == 0 {
		return errors.New("config `websiteId` is required")
	}

	// 获取网站 HTTPS 配置
	getHttpsConfReq := &opsdk.GetHttpsConfRequest{
		WebsiteID: d.config.WebsiteId,
	}
	getHttpsConfResp, err := d.sdkClient.GetHttpsConf(getHttpsConfReq)
	d.logger.Debug("sdk request '1panel.GetHttpsConf'", slog.Any("request", getHttpsConfReq), slog.Any("response", getHttpsConfResp))
	if err != nil {
		return fmt.Errorf("failed to execute sdk request '1panel.GetHttpsConf': %w", err)
	}

	// 上传证书到面板
	upres, err := d.sslUploader.Upload(ctx, certPEM, privkeyPEM)
	if err != nil {
		return fmt.Errorf("failed to upload certificate file: %w", err)
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
		return fmt.Errorf("failed to execute sdk request '1panel.UpdateHttpsConf': %w", err)
	}

	return nil
}

func (d *DeployerProvider) deployToCertificate(ctx context.Context, certPEM string, privkeyPEM string) error {
	if d.config.CertificateId == 0 {
		return errors.New("config `certificateId` is required")
	}

	// 获取证书详情
	getWebsiteSSLReq := &opsdk.GetWebsiteSSLRequest{
		SSLID: d.config.CertificateId,
	}
	getWebsiteSSLResp, err := d.sdkClient.GetWebsiteSSL(getWebsiteSSLReq)
	d.logger.Debug("sdk request '1panel.GetWebsiteSSL'", slog.Any("request", getWebsiteSSLReq), slog.Any("response", getWebsiteSSLResp))
	if err != nil {
		return fmt.Errorf("failed to execute sdk request '1panel.GetWebsiteSSL': %w", err)
	}

	// 更新证书
	uploadWebsiteSSLReq := &opsdk.UploadWebsiteSSLRequest{
		Type:        "paste",
		SSLID:       d.config.CertificateId,
		Description: getWebsiteSSLResp.Data.Description,
		Certificate: certPEM,
		PrivateKey:  privkeyPEM,
	}
	uploadWebsiteSSLResp, err := d.sdkClient.UploadWebsiteSSL(uploadWebsiteSSLReq)
	d.logger.Debug("sdk request '1panel.UploadWebsiteSSL'", slog.Any("request", uploadWebsiteSSLReq), slog.Any("response", uploadWebsiteSSLResp))
	if err != nil {
		return fmt.Errorf("failed to execute sdk request '1panel.UploadWebsiteSSL': %w", err)
	}

	return nil
}

func createSdkClient(apiUrl, apiKey string, skipTlsVerify bool) (*opsdk.Client, error) {
	if _, err := url.Parse(apiUrl); err != nil {
		return nil, errors.New("invalid 1panel api url")
	}

	if apiKey == "" {
		return nil, errors.New("invalid 1panel api key")
	}

	client := opsdk.NewClient(apiUrl, apiKey)
	if skipTlsVerify {
		client.WithTLSConfig(&tls.Config{InsecureSkipVerify: true})
	}

	return client, nil
}
