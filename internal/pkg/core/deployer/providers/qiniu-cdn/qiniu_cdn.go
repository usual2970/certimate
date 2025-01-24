package qiniucdn

import (
	"context"
	"errors"
	"strings"

	xerrors "github.com/pkg/errors"
	"github.com/qiniu/go-sdk/v7/auth"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/logger"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	uploaderp "github.com/usual2970/certimate/internal/pkg/core/uploader/providers/qiniu-sslcert"
	qiniusdk "github.com/usual2970/certimate/internal/pkg/vendors/qiniu-sdk"
)

type QiniuCDNDeployerConfig struct {
	// 七牛云 AccessKey。
	AccessKey string `json:"accessKey"`
	// 七牛云 SecretKey。
	SecretKey string `json:"secretKey"`
	// 加速域名（支持泛域名）。
	Domain string `json:"domain"`
}

type QiniuCDNDeployer struct {
	config      *QiniuCDNDeployerConfig
	logger      logger.Logger
	sdkClient   *qiniusdk.Client
	sslUploader uploader.Uploader
}

var _ deployer.Deployer = (*QiniuCDNDeployer)(nil)

func New(config *QiniuCDNDeployerConfig) (*QiniuCDNDeployer, error) {
	return NewWithLogger(config, logger.NewNilLogger())
}

func NewWithLogger(config *QiniuCDNDeployerConfig, logger logger.Logger) (*QiniuCDNDeployer, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}

	if logger == nil {
		return nil, errors.New("logger is nil")
	}

	client := qiniusdk.NewClient(auth.New(config.AccessKey, config.SecretKey))

	uploader, err := uploaderp.New(&uploaderp.QiniuSSLCertUploaderConfig{
		AccessKey: config.AccessKey,
		SecretKey: config.SecretKey,
	})
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create ssl uploader")
	}

	return &QiniuCDNDeployer{
		logger:      logger,
		config:      config,
		sdkClient:   client,
		sslUploader: uploader,
	}, nil
}

func (d *QiniuCDNDeployer) Deploy(ctx context.Context, certPem string, privkeyPem string) (*deployer.DeployResult, error) {
	// 上传证书到 CDN
	upres, err := d.sslUploader.Upload(ctx, certPem, privkeyPem)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to upload certificate file")
	}

	d.logger.Logt("certificate file uploaded", upres)

	// "*.example.com" → ".example.com"，适配七牛云 CDN 要求的泛域名格式
	domain := strings.TrimPrefix(d.config.Domain, "*")

	// 获取域名信息
	// REF: https://developer.qiniu.com/fusion/4246/the-domain-name
	getDomainInfoResp, err := d.sdkClient.GetDomainInfo(domain)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'cdn.GetDomainInfo'")
	}

	d.logger.Logt("已获取域名信息", getDomainInfoResp)

	// 判断域名是否已启用 HTTPS。如果已启用，修改域名证书；否则，启用 HTTPS
	// REF: https://developer.qiniu.com/fusion/4246/the-domain-name
	if getDomainInfoResp.Https != nil && getDomainInfoResp.Https.CertID != "" {
		modifyDomainHttpsConfResp, err := d.sdkClient.ModifyDomainHttpsConf(domain, upres.CertId, getDomainInfoResp.Https.ForceHttps, getDomainInfoResp.Https.Http2Enable)
		if err != nil {
			return nil, xerrors.Wrap(err, "failed to execute sdk request 'cdn.ModifyDomainHttpsConf'")
		}

		d.logger.Logt("已修改域名证书", modifyDomainHttpsConfResp)
	} else {
		enableDomainHttpsResp, err := d.sdkClient.EnableDomainHttps(domain, upres.CertId, true, true)
		if err != nil {
			return nil, xerrors.Wrap(err, "failed to execute sdk request 'cdn.EnableDomainHttps'")
		}

		d.logger.Logt("已将域名升级为 HTTPS", enableDomainHttpsResp)
	}

	return &deployer.DeployResult{}, nil
}
