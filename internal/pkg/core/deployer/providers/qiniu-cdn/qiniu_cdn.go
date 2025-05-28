package qiniucdn

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/qiniu/go-sdk/v7/auth"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	uploadersp "github.com/usual2970/certimate/internal/pkg/core/uploader/providers/qiniu-sslcert"
	qiniusdk "github.com/usual2970/certimate/internal/pkg/sdk3rd/qiniu"
)

type DeployerConfig struct {
	// 七牛云 AccessKey。
	AccessKey string `json:"accessKey"`
	// 七牛云 SecretKey。
	SecretKey string `json:"secretKey"`
	// 加速域名（支持泛域名）。
	Domain string `json:"domain"`
}

type DeployerProvider struct {
	config      *DeployerConfig
	logger      *slog.Logger
	sdkClient   *qiniusdk.Client
	sslUploader uploader.Uploader
}

var _ deployer.Deployer = (*DeployerProvider)(nil)

func NewDeployer(config *DeployerConfig) (*DeployerProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	client := qiniusdk.NewClient(auth.New(config.AccessKey, config.SecretKey))

	uploader, err := uploadersp.NewUploader(&uploadersp.UploaderConfig{
		AccessKey: config.AccessKey,
		SecretKey: config.SecretKey,
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
		d.logger = slog.New(slog.DiscardHandler)
	} else {
		d.logger = logger
	}
	d.sslUploader.WithLogger(logger)
	return d
}

func (d *DeployerProvider) Deploy(ctx context.Context, certPEM string, privkeyPEM string) (*deployer.DeployResult, error) {
	// 上传证书到 CDN
	upres, err := d.sslUploader.Upload(ctx, certPEM, privkeyPEM)
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

	return &deployer.DeployResult{}, nil
}
