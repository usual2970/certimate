package dogecloud

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	dogesdk "github.com/usual2970/certimate/internal/pkg/sdk3rd/dogecloud"
)

type UploaderConfig struct {
	// 多吉云 AccessKey。
	AccessKey string `json:"accessKey"`
	// 多吉云 SecretKey。
	SecretKey string `json:"secretKey"`
}

type UploaderProvider struct {
	config    *UploaderConfig
	logger    *slog.Logger
	sdkClient *dogesdk.Client
}

var _ uploader.Uploader = (*UploaderProvider)(nil)

func NewUploader(config *UploaderConfig) (*UploaderProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	client, err := createSdkClient(config.AccessKey, config.SecretKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create sdk client: %w", err)
	}

	return &UploaderProvider{
		config:    config,
		logger:    slog.Default(),
		sdkClient: client,
	}, nil
}

func (u *UploaderProvider) WithLogger(logger *slog.Logger) uploader.Uploader {
	if logger == nil {
		u.logger = slog.New(slog.DiscardHandler)
	} else {
		u.logger = logger
	}
	return u
}

func (u *UploaderProvider) Upload(ctx context.Context, certPEM string, privkeyPEM string) (*uploader.UploadResult, error) {
	// 生成新证书名（需符合多吉云命名规则）
	var certId, certName string
	certName = fmt.Sprintf("certimate-%d", time.Now().UnixMilli())

	// 上传新证书
	// REF: https://docs.dogecloud.com/cdn/api-cert-upload
	uploadSslCertResp, err := u.sdkClient.UploadCdnCert(certName, certPEM, privkeyPEM)
	u.logger.Debug("sdk request 'cdn.UploadCdnCert'", slog.Any("response", uploadSslCertResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'cdn.UploadCdnCert': %w", err)
	}

	certId = fmt.Sprintf("%d", uploadSslCertResp.Data.Id)
	return &uploader.UploadResult{
		CertId:   certId,
		CertName: certName,
	}, nil
}

func createSdkClient(accessKey, secretKey string) (*dogesdk.Client, error) {
	client := dogesdk.NewClient(accessKey, secretKey)
	return client, nil
}
