package ucloudus3

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ucloud/ucloud-sdk-go/ucloud"
	"github.com/ucloud/ucloud-sdk-go/ucloud/auth"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	uploadersp "github.com/usual2970/certimate/internal/pkg/core/uploader/providers/ucloud-ussl"
	usdkFile "github.com/usual2970/certimate/internal/pkg/sdk3rd/ucloud/ufile"
)

type DeployerConfig struct {
	// 优刻得 API 私钥。
	PrivateKey string `json:"privateKey"`
	// 优刻得 API 公钥。
	PublicKey string `json:"publicKey"`
	// 优刻得项目 ID。
	ProjectId string `json:"projectId,omitempty"`
	// 优刻得地域。
	Region string `json:"region"`
	// 存储桶名。
	Bucket string `json:"bucket"`
	// 自定义域名（不支持泛域名）。
	Domain string `json:"domain"`
}

type DeployerProvider struct {
	config      *DeployerConfig
	logger      *slog.Logger
	sdkClient   *usdkFile.UFileClient
	sslUploader uploader.Uploader
}

var _ deployer.Deployer = (*DeployerProvider)(nil)

func NewDeployer(config *DeployerConfig) (*DeployerProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	client, err := createSdkClient(config.PrivateKey, config.PublicKey, config.Region)
	if err != nil {
		return nil, fmt.Errorf("failed to create sdk client: %w", err)
	}

	uploader, err := uploadersp.NewUploader(&uploadersp.UploaderConfig{
		PrivateKey: config.PrivateKey,
		PublicKey:  config.PublicKey,
		ProjectId:  config.ProjectId,
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
	// 上传证书到 USSL
	upres, err := d.sslUploader.Upload(ctx, certPEM, privkeyPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to upload certificate file: %w", err)
	} else {
		d.logger.Info("ssl certificate uploaded", slog.Any("result", upres))
	}

	// 添加 SSL 证书
	// REF: https://docs.ucloud.cn/api/ufile-api/add_ufile_ssl_cert
	addUFileSSLCertReq := d.sdkClient.NewAddUFileSSLCertRequest()
	addUFileSSLCertReq.BucketName = ucloud.String(d.config.Bucket)
	addUFileSSLCertReq.Domain = ucloud.String(d.config.Domain)
	addUFileSSLCertReq.USSLId = ucloud.String(upres.CertId)
	addUFileSSLCertReq.CertificateName = ucloud.String(upres.CertName)
	if d.config.ProjectId != "" {
		addUFileSSLCertReq.ProjectId = ucloud.String(d.config.ProjectId)
	}
	addUFileSSLCertResp, err := d.sdkClient.AddUFileSSLCert(addUFileSSLCertReq)
	d.logger.Debug("sdk request 'us3.AddUFileSSLCert'", slog.Any("request", addUFileSSLCertReq), slog.Any("response", addUFileSSLCertResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'us3.AddUFileSSLCert': %w", err)
	}

	return &deployer.DeployResult{}, nil
}

func createSdkClient(privateKey, publicKey, region string) (*usdkFile.UFileClient, error) {
	cfg := ucloud.NewConfig()
	cfg.Region = region

	credential := auth.NewCredential()
	credential.PrivateKey = privateKey
	credential.PublicKey = publicKey

	client := usdkFile.NewClient(&cfg, &credential)
	return client, nil
}
