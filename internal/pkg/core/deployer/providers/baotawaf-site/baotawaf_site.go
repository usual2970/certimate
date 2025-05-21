package baotapanelwaf

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"
	"net/url"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	btsdk "github.com/usual2970/certimate/internal/pkg/sdk3rd/btwaf"
	typeutil "github.com/usual2970/certimate/internal/pkg/utils/type"
)

type DeployerConfig struct {
	// 堡塔云 WAF 地址。
	ApiUrl string `json:"apiUrl"`
	// 堡塔云 WAF 接口密钥。
	ApiKey string `json:"apiKey"`
	// 是否允许不安全的连接。
	AllowInsecureConnections bool `json:"allowInsecureConnections,omitempty"`
	// 网站名称。
	SiteName string `json:"siteName"`
	// 网站 SSL 端口。
	// 零值时默认为 443。
	SitePort int32 `json:"sitePort,omitempty"`
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
	if d.config.SiteName == "" {
		return nil, errors.New("config `siteName` is required")
	}
	if d.config.SitePort == 0 {
		d.config.SitePort = 443
	}

	// 遍历获取网站列表，获取网站 ID
	// REF: https://support.huaweicloud.com/api-waf/ListHost.html
	siteId := ""
	getSitListPage := int32(1)
	getSitListPageSize := int32(100)
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		getSiteListReq := &btsdk.GetSiteListRequest{
			SiteName: typeutil.ToPtr(d.config.SiteName),
			Page:     typeutil.ToPtr(getSitListPage),
			PageSize: typeutil.ToPtr(getSitListPageSize),
		}
		getSiteListResp, err := d.sdkClient.GetSiteList(getSiteListReq)
		d.logger.Debug("sdk request 'bt.GetSiteList'", slog.Any("request", getSiteListReq), slog.Any("response", getSiteListResp))
		if err != nil {
			return nil, fmt.Errorf("failed to execute sdk request 'bt.GetSiteList': %w", err)
		}

		if getSiteListResp.Result != nil && getSiteListResp.Result.List != nil {
			for _, siteItem := range getSiteListResp.Result.List {
				if siteItem.SiteName == d.config.SiteName {
					siteId = siteItem.SiteId
					break
				}
			}
		}

		if getSiteListResp.Result == nil || len(getSiteListResp.Result.List) < int(getSitListPageSize) {
			break
		} else {
			getSitListPage++
		}
	}
	if siteId == "" {
		return nil, errors.New("site not found")
	}

	// 修改站点配置
	modifySiteReq := &btsdk.ModifySiteRequest{
		SiteId: siteId,
		Type:   typeutil.ToPtr("openCert"),
		Server: &btsdk.SiteServerInfo{
			ListenSSLPort: typeutil.ToPtr(d.config.SitePort),
			SSL: &btsdk.SiteServerSSLInfo{
				IsSSL:      typeutil.ToPtr(int32(1)),
				FullChain:  typeutil.ToPtr(certPEM),
				PrivateKey: typeutil.ToPtr(privkeyPEM),
			},
		},
	}
	modifySiteResp, err := d.sdkClient.ModifySite(modifySiteReq)
	d.logger.Debug("sdk request 'bt.ModifySite'", slog.Any("request", modifySiteReq), slog.Any("response", modifySiteResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'bt.ModifySite': %w", err)
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
