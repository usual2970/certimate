package qiniucdn

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/qiniu/go-sdk/v7/auth"

	"github.com/certimate-go/certimate/pkg/core"
	sslmgrsp "github.com/certimate-go/certimate/pkg/core/ssl-manager/providers/qiniu-sslcert"
	qiniusdk "github.com/certimate-go/certimate/pkg/sdk3rd/qiniu"
)

type SSLDeployerProviderConfig struct {
	// 七牛云 AccessKey。
	AccessKey string `json:"accessKey"`
	// 七牛云 SecretKey。
	SecretKey string `json:"secretKey"`
	// 加速域名（支持泛域名）。
	Domain string `json:"domain"`
}

type SSLDeployerProvider struct {
	config     *SSLDeployerProviderConfig
	logger     *slog.Logger
	sdkClient  *qiniusdk.CdnManager
	sslManager core.SSLManager
}

var _ core.SSLDeployer = (*SSLDeployerProvider)(nil)

func NewSSLDeployerProvider(config *SSLDeployerProviderConfig) (*SSLDeployerProvider, error) {
	if config == nil {
		return nil, errors.New("the configuration of the ssl deployer provider is nil")
	}

	client := qiniusdk.NewCdnManager(auth.New(config.AccessKey, config.SecretKey))

	sslmgr, err := sslmgrsp.NewSSLManagerProvider(&sslmgrsp.SSLManagerProviderConfig{
		AccessKey: config.AccessKey,
		SecretKey: config.SecretKey,
	})
	if err != nil {
		return nil, fmt.Errorf("could not create ssl manager: %w", err)
	}

	return &SSLDeployerProvider{
		config:     config,
		logger:     slog.Default(),
		sdkClient:  client,
		sslManager: sslmgr,
	}, nil
}

func (d *SSLDeployerProvider) SetLogger(logger *slog.Logger) {
	if logger == nil {
		d.logger = slog.New(slog.DiscardHandler)
	} else {
		d.logger = logger
	}

	d.sslManager.SetLogger(logger)
}

func (d *SSLDeployerProvider) Deploy(ctx context.Context, certPEM string, privkeyPEM string) (*core.SSLDeployResult, error) {
	if d.config.Domain == "" {
		return nil, fmt.Errorf("config `domain` is required")
	}

	// 上传证书
	upres, err := d.sslManager.Upload(ctx, certPEM, privkeyPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to upload certificate file: %w", err)
	} else {
		d.logger.Info("ssl certificate uploaded", slog.Any("result", upres))
	}

	// "*.example.com" → ".example.com"，适配七牛云 CDN 要求的泛域名格式
	domain := strings.TrimPrefix(d.config.Domain, "*")

	// 获取域名信息
	// REF: https://developer.qiniu.com/fusion/4246/the-domain-name
	getDomainInfoResp, err := d.sdkClient.GetDomainInfo(context.TODO(), domain)
	d.logger.Debug("sdk request 'cdn.GetDomainInfo'", slog.String("request.domain", domain), slog.Any("response", getDomainInfoResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'cdn.GetDomainInfo': %w", err)
	}

	// 判断域名是否已启用 HTTPS。如果已启用，修改域名证书；否则，启用 HTTPS
	// REF: https://developer.qiniu.com/fusion/4246/the-domain-name
	if getDomainInfoResp.Https == nil || getDomainInfoResp.Https.CertID == "" {
		enableDomainHttpsResp, err := d.sdkClient.EnableDomainHttps(context.TODO(), domain, upres.CertId, true, true)
		d.logger.Debug("sdk request 'cdn.EnableDomainHttps'", slog.String("request.domain", domain), slog.String("request.certId", upres.CertId), slog.Any("response", enableDomainHttpsResp))
		if err != nil {
			return nil, fmt.Errorf("failed to execute sdk request 'cdn.EnableDomainHttps': %w", err)
		}
	} else if getDomainInfoResp.Https.CertID != upres.CertId {
		modifyDomainHttpsConfResp, err := d.sdkClient.ModifyDomainHttpsConf(context.TODO(), domain, upres.CertId, getDomainInfoResp.Https.ForceHttps, getDomainInfoResp.Https.Http2Enable)
		d.logger.Debug("sdk request 'cdn.ModifyDomainHttpsConf'", slog.String("request.domain", domain), slog.String("request.certId", upres.CertId), slog.Any("response", modifyDomainHttpsConfResp))
		if err != nil {
			return nil, fmt.Errorf("failed to execute sdk request 'cdn.ModifyDomainHttpsConf': %w", err)
		}
	}

	return &core.SSLDeployResult{}, nil
}
