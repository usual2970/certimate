package baotapanelsite

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"

	"github.com/certimate-go/certimate/pkg/core"
	btsdk "github.com/certimate-go/certimate/pkg/sdk3rd/btpanel"
	xslices "github.com/certimate-go/certimate/pkg/utils/slices"
)

type SSLDeployerProviderConfig struct {
	// 宝塔面板服务地址。
	ServerUrl string `json:"serverUrl"`
	// 宝塔面板接口密钥。
	ApiKey string `json:"apiKey"`
	// 是否允许不安全的连接。
	AllowInsecureConnections bool `json:"allowInsecureConnections,omitempty"`
	// 网站类型。
	SiteType string `json:"siteType"`
	// 网站名称（单个）。
	SiteName string `json:"siteName,omitempty"`
	// 网站名称（多个）。
	SiteNames []string `json:"siteNames,omitempty"`
}

type SSLDeployerProvider struct {
	config    *SSLDeployerProviderConfig
	logger    *slog.Logger
	sdkClient *btsdk.Client
}

var _ core.SSLDeployer = (*SSLDeployerProvider)(nil)

func NewSSLDeployerProvider(config *SSLDeployerProviderConfig) (*SSLDeployerProvider, error) {
	if config == nil {
		return nil, errors.New("the configuration of the ssl deployer provider is nil")
	}

	client, err := createSDKClient(config.ServerUrl, config.ApiKey, config.AllowInsecureConnections)
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
	switch d.config.SiteType {
	case "php":
		{
			if d.config.SiteName == "" {
				return nil, errors.New("config `siteName` is required")
			}

			// 设置站点 SSL 证书
			siteSetSSLReq := &btsdk.SiteSetSSLRequest{
				SiteName:    d.config.SiteName,
				Type:        "0",
				Certificate: certPEM,
				PrivateKey:  privkeyPEM,
			}
			siteSetSSLResp, err := d.sdkClient.SiteSetSSL(siteSetSSLReq)
			d.logger.Debug("sdk request 'bt.SiteSetSSL'", slog.Any("request", siteSetSSLReq), slog.Any("response", siteSetSSLResp))
			if err != nil {
				return nil, fmt.Errorf("failed to execute sdk request 'bt.SiteSetSSL': %w", err)
			}
		}

	case "other":
		{
			if len(d.config.SiteNames) == 0 {
				return nil, errors.New("config `siteNames` is required")
			}

			// 上传证书
			sslCertSaveCertReq := &btsdk.SSLCertSaveCertRequest{
				Certificate: certPEM,
				PrivateKey:  privkeyPEM,
			}
			sslCertSaveCertResp, err := d.sdkClient.SSLCertSaveCert(sslCertSaveCertReq)
			d.logger.Debug("sdk request 'bt.SSLCertSaveCert'", slog.Any("request", sslCertSaveCertReq), slog.Any("response", sslCertSaveCertResp))
			if err != nil {
				return nil, fmt.Errorf("failed to execute sdk request 'bt.SSLCertSaveCert': %w", err)
			}

			// 设置站点证书
			sslSetBatchCertToSiteReq := &btsdk.SSLSetBatchCertToSiteRequest{
				BatchInfo: xslices.Map(d.config.SiteNames, func(siteName string) *btsdk.SSLSetBatchCertToSiteRequestBatchInfo {
					return &btsdk.SSLSetBatchCertToSiteRequestBatchInfo{
						SiteName: siteName,
						SSLHash:  sslCertSaveCertResp.SSLHash,
					}
				}),
			}
			sslSetBatchCertToSiteResp, err := d.sdkClient.SSLSetBatchCertToSite(sslSetBatchCertToSiteReq)
			d.logger.Debug("sdk request 'bt.SSLSetBatchCertToSite'", slog.Any("request", sslSetBatchCertToSiteReq), slog.Any("response", sslSetBatchCertToSiteResp))
			if err != nil {
				return nil, fmt.Errorf("failed to execute sdk request 'bt.SSLSetBatchCertToSite': %w", err)
			}
		}

	default:
		return nil, fmt.Errorf("unsupported site type '%s'", d.config.SiteType)
	}

	return &core.SSLDeployResult{}, nil
}

func createSDKClient(serverUrl, apiKey string, skipTlsVerify bool) (*btsdk.Client, error) {
	client, err := btsdk.NewClient(serverUrl, apiKey)
	if err != nil {
		return nil, err
	}

	if skipTlsVerify {
		client.SetTLSConfig(&tls.Config{InsecureSkipVerify: true})
	}

	return client, nil
}
