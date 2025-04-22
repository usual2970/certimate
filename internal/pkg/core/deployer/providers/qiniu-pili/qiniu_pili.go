package qiniupili

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/qiniu/go-sdk/v7/pili"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	uploadersp "github.com/usual2970/certimate/internal/pkg/core/uploader/providers/qiniu-sslcert"
)

type DeployerConfig struct {
	// 七牛云 AccessKey。
	AccessKey string `json:"accessKey"`
	// 七牛云 SecretKey。
	SecretKey string `json:"secretKey"`
	// 直播空间名。
	Hub string `json:"hub"`
	// 直播流域名（不支持泛域名）。
	Domain string `json:"domain"`
}

type DeployerProvider struct {
	config      *DeployerConfig
	logger      *slog.Logger
	sdkClient   *pili.Manager
	sslUploader uploader.Uploader
}

var _ deployer.Deployer = (*DeployerProvider)(nil)

func NewDeployer(config *DeployerConfig) (*DeployerProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	manager := pili.NewManager(pili.ManagerConfig{AccessKey: config.AccessKey, SecretKey: config.SecretKey})

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
		sdkClient:   manager,
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

	// 修改域名证书配置
	// REF: https://developer.qiniu.com/pili/9910/pili-service-sdk#66
	setDomainCertReq := pili.SetDomainCertRequest{
		Hub:      d.config.Hub,
		Domain:   d.config.Domain,
		CertName: upres.CertName,
	}
	err = d.sdkClient.SetDomainCert(context.TODO(), setDomainCertReq)
	d.logger.Debug("sdk request 'pili.SetDomainCert'", slog.Any("request", setDomainCertReq))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'pili.SetDomainCert': %w", err)
	}

	return &deployer.DeployResult{}, nil
}
