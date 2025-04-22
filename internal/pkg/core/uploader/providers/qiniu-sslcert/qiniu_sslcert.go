package qiniusslcert

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/qiniu/go-sdk/v7/auth"

	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	qiniusdk "github.com/usual2970/certimate/internal/pkg/sdk3rd/qiniu"
	certutil "github.com/usual2970/certimate/internal/pkg/utils/cert"
)

type UploaderConfig struct {
	// 七牛云 AccessKey。
	AccessKey string `json:"accessKey"`
	// 七牛云 SecretKey。
	SecretKey string `json:"secretKey"`
}

type UploaderProvider struct {
	config    *UploaderConfig
	logger    *slog.Logger
	sdkClient *qiniusdk.Client
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
		u.logger = slog.Default()
	} else {
		u.logger = logger
	}
	return u
}

func (u *UploaderProvider) Upload(ctx context.Context, certPEM string, privkeyPEM string) (res *uploader.UploadResult, err error) {
	// 解析证书内容
	certX509, err := certutil.ParseCertificateFromPEM(certPEM)
	if err != nil {
		return nil, err
	}

	// 生成新证书名（需符合七牛云命名规则）
	var certId, certName string
	certName = fmt.Sprintf("certimate-%d", time.Now().UnixMilli())

	// 上传新证书
	// REF: https://developer.qiniu.com/fusion/8593/interface-related-certificate
	uploadSslCertResp, err := u.sdkClient.UploadSslCert(context.TODO(), certName, certX509.Subject.CommonName, certPEM, privkeyPEM)
	u.logger.Debug("sdk request 'cdn.UploadSslCert'", slog.Any("response", uploadSslCertResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'cdn.UploadSslCert': %w", err)
	}

	certId = uploadSslCertResp.CertID
	return &uploader.UploadResult{
		CertId:   certId,
		CertName: certName,
	}, nil
}

func createSdkClient(accessKey, secretKey string) (*qiniusdk.Client, error) {
	if secretKey == "" {
		return nil, errors.New("invalid qiniu access key")
	}

	if secretKey == "" {
		return nil, errors.New("invalid qiniu secret key")
	}

	credential := auth.New(accessKey, secretKey)
	client := qiniusdk.NewClient(credential)
	return client, nil
}
