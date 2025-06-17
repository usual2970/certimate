package qiniusslcert

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/qiniu/go-sdk/v7/auth"

	"github.com/certimate-go/certimate/pkg/core"
	qiniusdk "github.com/certimate-go/certimate/pkg/sdk3rd/qiniu"
	xcert "github.com/certimate-go/certimate/pkg/utils/cert"
)

type SSLManagerProviderConfig struct {
	// 七牛云 AccessKey。
	AccessKey string `json:"accessKey"`
	// 七牛云 SecretKey。
	SecretKey string `json:"secretKey"`
}

type SSLManagerProvider struct {
	config    *SSLManagerProviderConfig
	logger    *slog.Logger
	sdkClient *qiniusdk.CdnManager
}

var _ core.SSLManager = (*SSLManagerProvider)(nil)

func NewSSLManagerProvider(config *SSLManagerProviderConfig) (*SSLManagerProvider, error) {
	if config == nil {
		return nil, errors.New("the configuration of the ssl manager provider is nil")
	}

	client, err := createSDKClient(config.AccessKey, config.SecretKey)
	if err != nil {
		return nil, fmt.Errorf("could not create sdk client: %w", err)
	}

	return &SSLManagerProvider{
		config:    config,
		logger:    slog.Default(),
		sdkClient: client,
	}, nil
}

func (m *SSLManagerProvider) SetLogger(logger *slog.Logger) {
	if logger == nil {
		m.logger = slog.New(slog.DiscardHandler)
	} else {
		m.logger = logger
	}
}

func (m *SSLManagerProvider) Upload(ctx context.Context, certPEM string, privkeyPEM string) (*core.SSLManageUploadResult, error) {
	// 解析证书内容
	certX509, err := xcert.ParseCertificateFromPEM(certPEM)
	if err != nil {
		return nil, err
	}

	// 生成新证书名（需符合七牛云命名规则）
	certName := fmt.Sprintf("certimate-%d", time.Now().UnixMilli())

	// 上传新证书
	// REF: https://developer.qiniu.com/fusion/8593/interface-related-certificate
	uploadSslCertResp, err := m.sdkClient.UploadSslCert(context.TODO(), certName, certX509.Subject.CommonName, certPEM, privkeyPEM)
	m.logger.Debug("sdk request 'cdn.UploadSslCert'", slog.Any("response", uploadSslCertResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'cdn.UploadSslCert': %w", err)
	}

	return &core.SSLManageUploadResult{
		CertId:   uploadSslCertResp.CertID,
		CertName: certName,
	}, nil
}

func createSDKClient(accessKey, secretKey string) (*qiniusdk.CdnManager, error) {
	if secretKey == "" {
		return nil, errors.New("invalid qiniu access key")
	}

	if secretKey == "" {
		return nil, errors.New("invalid qiniu secret key")
	}

	credential := auth.New(accessKey, secretKey)
	client := qiniusdk.NewCdnManager(credential)
	return client, nil
}
