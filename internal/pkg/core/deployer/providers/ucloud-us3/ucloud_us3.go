package ucloudus3

import (
	"context"
	"log/slog"

	xerrors "github.com/pkg/errors"
	usdk "github.com/ucloud/ucloud-sdk-go/ucloud"
	uAuth "github.com/ucloud/ucloud-sdk-go/ucloud/auth"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	uploadersp "github.com/usual2970/certimate/internal/pkg/core/uploader/providers/ucloud-ussl"
	usdkFile "github.com/usual2970/certimate/internal/pkg/vendors/ucloud-sdk/ufile"
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
		return nil, xerrors.Wrap(err, "failed to create sdk client")
	}

	uploader, err := uploadersp.NewUploader(&uploadersp.UploaderConfig{
		PrivateKey: config.PrivateKey,
		PublicKey:  config.PublicKey,
		ProjectId:  config.ProjectId,
	})
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create ssl uploader")
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

func (d *DeployerProvider) Deploy(ctx context.Context, certPem string, privkeyPem string) (*deployer.DeployResult, error) {
	// 上传证书到 USSL
	upres, err := d.sslUploader.Upload(ctx, certPem, privkeyPem)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to upload certificate file")
	} else {
		d.logger.Info("ssl certificate uploaded", slog.Any("result", upres))
	}

	// 添加 SSL 证书
	// REF: https://docs.ucloud.cn/api/ufile-api/add_ufile_ssl_cert
	addUFileSSLCertReq := d.sdkClient.NewAddUFileSSLCertRequest()
	addUFileSSLCertReq.BucketName = usdk.String(d.config.Bucket)
	addUFileSSLCertReq.Domain = usdk.String(d.config.Domain)
	addUFileSSLCertReq.USSLId = usdk.String(upres.CertId)
	addUFileSSLCertReq.CertificateName = usdk.String(upres.CertName)
	if d.config.ProjectId != "" {
		addUFileSSLCertReq.ProjectId = usdk.String(d.config.ProjectId)
	}
	addUFileSSLCertResp, err := d.sdkClient.AddUFileSSLCert(addUFileSSLCertReq)
	d.logger.Debug("sdk request 'us3.AddUFileSSLCert'", slog.Any("request", addUFileSSLCertReq), slog.Any("response", addUFileSSLCertResp))
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'us3.AddUFileSSLCert'")
	}

	return &deployer.DeployResult{}, nil
}

func createSdkClient(privateKey, publicKey, region string) (*usdkFile.UFileClient, error) {
	cfg := usdk.NewConfig()
	cfg.Region = region

	credential := uAuth.NewCredential()
	credential.PrivateKey = privateKey
	credential.PublicKey = publicKey

	client := usdkFile.NewClient(&cfg, &credential)
	return client, nil
}
