package upyunssl

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	upyunsdk "github.com/usual2970/certimate/internal/pkg/sdk3rd/upyun/console"
)

type UploaderConfig struct {
	// 又拍云账号用户名。
	Username string `json:"username"`
	// 又拍云账号密码。
	Password string `json:"password"`
}

type UploaderProvider struct {
	config    *UploaderConfig
	logger    *slog.Logger
	sdkClient *upyunsdk.Client
}

var _ uploader.Uploader = (*UploaderProvider)(nil)

func NewUploader(config *UploaderConfig) (*UploaderProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	client, err := createSdkClient(config.Username, config.Password)
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
	// 上传证书
	uploadHttpsCertificateReq := &upyunsdk.UploadHttpsCertificateRequest{
		Certificate: certPEM,
		PrivateKey:  privkeyPEM,
	}
	uploadHttpsCertificateResp, err := u.sdkClient.UploadHttpsCertificate(uploadHttpsCertificateReq)
	u.logger.Debug("sdk request 'console.UploadHttpsCertificate'", slog.Any("response", uploadHttpsCertificateResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'console.UploadHttpsCertificate': %w", err)
	}

	return &uploader.UploadResult{
		CertId: uploadHttpsCertificateResp.Data.Result.CertificateId,
	}, nil
}

func createSdkClient(username, password string) (*upyunsdk.Client, error) {
	if username == "" {
		return nil, errors.New("invalid upyun username")
	}

	if password == "" {
		return nil, errors.New("invalid upyun password")
	}

	client := upyunsdk.NewClient(username, password)
	return client, nil
}
