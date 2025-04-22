package dogecloudcdn

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	uploadersp "github.com/usual2970/certimate/internal/pkg/core/uploader/providers/dogecloud"
	dogesdk "github.com/usual2970/certimate/internal/pkg/sdk3rd/dogecloud"
)

type DeployerConfig struct {
	// 多吉云 AccessKey。
	AccessKey string `json:"accessKey"`
	// 多吉云 SecretKey。
	SecretKey string `json:"secretKey"`
	// 加速域名（不支持泛域名）。
	Domain string `json:"domain"`
}

type DeployerProvider struct {
	config      *DeployerConfig
	logger      *slog.Logger
	sdkClient   *dogesdk.Client
	sslUploader uploader.Uploader
}

var _ deployer.Deployer = (*DeployerProvider)(nil)

func NewDeployer(config *DeployerConfig) (*DeployerProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	client := dogesdk.NewClient(config.AccessKey, config.SecretKey)

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
		d.logger = slog.Default()
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

	// 绑定证书
	// REF: https://docs.dogecloud.com/cdn/api-cert-bind
	bindCdnCertId, _ := strconv.ParseInt(upres.CertId, 10, 64)
	bindCdnCertResp, err := d.sdkClient.BindCdnCertWithDomain(bindCdnCertId, d.config.Domain)
	d.logger.Debug("sdk request 'cdn.BindCdnCert'", slog.Int64("request.certId", bindCdnCertId), slog.String("request.domain", d.config.Domain), slog.Any("response", bindCdnCertResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'cdn.BindCdnCert': %w", err)
	}

	return &deployer.DeployResult{}, nil
}
