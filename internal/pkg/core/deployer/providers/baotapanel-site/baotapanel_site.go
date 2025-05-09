package baotapanelsite

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"
	"net/url"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	btsdk "github.com/usual2970/certimate/internal/pkg/sdk3rd/btpanel"
	sliceutil "github.com/usual2970/certimate/internal/pkg/utils/slice"
)

type DeployerConfig struct {
	// 宝塔面板地址。
	ApiUrl string `json:"apiUrl"`
	// 宝塔面板接口密钥。
	ApiKey string `json:"apiKey"`
	// 是否允许不安全的连接。
	AllowInsecureConnections bool `json:"allowInsecureConnections,omitempty"`
	// 站点类型。
	SiteType string `json:"siteType"`
	// 站点名称（单个）。
	SiteName string `json:"siteName,omitempty"`
	// 站点名称（多个）。
	SiteNames []string `json:"siteNames,omitempty"`
}

type DeployerProvider struct {
	config    *DeployerConfig
	logger    *slog.Logger
	sdkClient *btsdk.Client
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
				BatchInfo: sliceutil.Map(d.config.SiteNames, func(siteName string) *btsdk.SSLSetBatchCertToSiteRequestBatchInfo {
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

	return &deployer.DeployResult{}, nil
}

func createSdkClient(apiUrl, apiKey string, skipTlsVerify bool) (*btsdk.Client, error) {
	if _, err := url.Parse(apiUrl); err != nil {
		return nil, errors.New("invalid baota api url")
	}

	if apiKey == "" {
		return nil, errors.New("invalid baota api key")
	}

	client := btsdk.NewClient(apiUrl, apiKey)
	if skipTlsVerify {
		client.WithTLSConfig(&tls.Config{InsecureSkipVerify: true})
	}

	return client, nil
}
