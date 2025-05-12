package gcorecdn

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/G-Core/gcorelabscdn-go/gcore/provider"
	"github.com/G-Core/gcorelabscdn-go/sslcerts"

	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	gcoresdk "github.com/usual2970/certimate/internal/pkg/sdk3rd/gcore/common"
)

type UploaderConfig struct {
	// Gcore API Token。
	ApiToken string `json:"apiToken"`
}

type UploaderProvider struct {
	config    *UploaderConfig
	logger    *slog.Logger
	sdkClient *sslcerts.Service
}

var _ uploader.Uploader = (*UploaderProvider)(nil)

func NewUploader(config *UploaderConfig) (*UploaderProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	client, err := createSdkClient(config.ApiToken)
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
	// 生成新证书名（需符合 Gcore 命名规则）
	var certId, certName string
	certName = fmt.Sprintf("certimate_%d", time.Now().UnixMilli())

	// 新增证书
	// REF: https://api.gcore.com/docs/cdn#tag/SSL-certificates/operation/add_ssl_certificates
	createCertificateReq := &sslcerts.CreateRequest{
		Name:           certName,
		Cert:           certPEM,
		PrivateKey:     privkeyPEM,
		Automated:      false,
		ValidateRootCA: false,
	}
	createCertificateResp, err := u.sdkClient.Create(context.TODO(), createCertificateReq)
	u.logger.Debug("sdk request 'sslcerts.Create'", slog.Any("request", createCertificateReq), slog.Any("response", createCertificateResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'sslcerts.Create': %w", err)
	}

	certId = fmt.Sprintf("%d", createCertificateResp.ID)
	certName = createCertificateResp.Name
	return &uploader.UploadResult{
		CertId:   certId,
		CertName: certName,
	}, nil
}

func createSdkClient(apiToken string) (*sslcerts.Service, error) {
	if apiToken == "" {
		return nil, errors.New("invalid gcore api token")
	}

	requester := provider.NewClient(
		gcoresdk.BASE_URL,
		provider.WithSigner(gcoresdk.NewAuthRequestSigner(apiToken)),
	)
	service := sslcerts.NewService(requester)
	return service, nil
}
