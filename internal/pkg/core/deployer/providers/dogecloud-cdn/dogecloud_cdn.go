package dogecloudcdn

import (
	"context"
	"errors"
	"strconv"

	xerrors "github.com/pkg/errors"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/logger"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	uploaderp "github.com/usual2970/certimate/internal/pkg/core/uploader/providers/dogecloud"
	dogesdk "github.com/usual2970/certimate/internal/pkg/vendors/dogecloud-sdk"
)

type DogeCloudCDNDeployerConfig struct {
	// 多吉云 AccessKey。
	AccessKey string `json:"accessKey"`
	// 多吉云 SecretKey。
	SecretKey string `json:"secretKey"`
	// 加速域名（不支持泛域名）。
	Domain string `json:"domain"`
}

type DogeCloudCDNDeployer struct {
	config      *DogeCloudCDNDeployerConfig
	logger      logger.Logger
	sdkClient   *dogesdk.Client
	sslUploader uploader.Uploader
}

var _ deployer.Deployer = (*DogeCloudCDNDeployer)(nil)

func New(config *DogeCloudCDNDeployerConfig) (*DogeCloudCDNDeployer, error) {
	return NewWithLogger(config, logger.NewNilLogger())
}

func NewWithLogger(config *DogeCloudCDNDeployerConfig, logger logger.Logger) (*DogeCloudCDNDeployer, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}

	if logger == nil {
		return nil, errors.New("logger is nil")
	}

	client := dogesdk.NewClient(config.AccessKey, config.SecretKey)

	uploader, err := uploaderp.New(&uploaderp.DogeCloudUploaderConfig{
		AccessKey: config.AccessKey,
		SecretKey: config.SecretKey,
	})
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create ssl uploader")
	}

	return &DogeCloudCDNDeployer{
		logger:      logger,
		config:      config,
		sdkClient:   client,
		sslUploader: uploader,
	}, nil
}

func (d *DogeCloudCDNDeployer) Deploy(ctx context.Context, certPem string, privkeyPem string) (*deployer.DeployResult, error) {
	// 上传证书到 CDN
	upres, err := d.sslUploader.Upload(ctx, certPem, privkeyPem)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to upload certificate file")
	}

	d.logger.Logt("certificate file uploaded", upres)

	// 绑定证书
	// REF: https://docs.dogecloud.com/cdn/api-cert-bind
	bindCdnCertId, _ := strconv.ParseInt(upres.CertId, 10, 64)
	bindCdnCertResp, err := d.sdkClient.BindCdnCertWithDomain(bindCdnCertId, d.config.Domain)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'cdn.BindCdnCert'")
	}

	d.logger.Logt("已绑定证书", bindCdnCertResp)

	return &deployer.DeployResult{}, nil
}
