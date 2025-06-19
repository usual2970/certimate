package ctcccloudelb

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"github.com/certimate-go/certimate/pkg/core"
	ctyunelb "github.com/certimate-go/certimate/pkg/sdk3rd/ctyun/elb"
	xcert "github.com/certimate-go/certimate/pkg/utils/cert"
	xtypes "github.com/certimate-go/certimate/pkg/utils/types"
)

type SSLManagerProviderConfig struct {
	// 天翼云 AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// 天翼云 SecretAccessKey。
	SecretAccessKey string `json:"secretAccessKey"`
	// 天翼云资源池 ID。
	RegionId string `json:"regionId"`
}

type SSLManagerProvider struct {
	config    *SSLManagerProviderConfig
	logger    *slog.Logger
	sdkClient *ctyunelb.Client
}

var _ core.SSLManager = (*SSLManagerProvider)(nil)

func NewSSLManagerProvider(config *SSLManagerProviderConfig) (*SSLManagerProvider, error) {
	if config == nil {
		return nil, errors.New("the configuration of the ssl manager provider is nil")
	}

	client, err := createSDKClient(config.AccessKeyId, config.SecretAccessKey)
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

	// 查询证书列表，避免重复上传
	// REF: https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=24&api=5692&data=88&isNormal=1&vid=82
	listCertificatesReq := &ctyunelb.ListCertificatesRequest{
		RegionID: xtypes.ToPtr(m.config.RegionId),
	}
	listCertificatesResp, err := m.sdkClient.ListCertificates(listCertificatesReq)
	m.logger.Debug("sdk request 'elb.ListCertificates'", slog.Any("request", listCertificatesReq), slog.Any("response", listCertificatesResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'elb.ListCertificates': %w", err)
	} else {
		for _, certRecord := range listCertificatesResp.ReturnObj {
			var isSameCert bool
			if certRecord.Certificate == certPEM {
				isSameCert = true
			} else {
				oldCertX509, err := xcert.ParseCertificateFromPEM(certRecord.Certificate)
				if err != nil {
					continue
				}

				isSameCert = xcert.EqualCertificate(certX509, oldCertX509)
			}

			// 如果已存在相同证书，直接返回
			if isSameCert {
				m.logger.Info("ssl certificate already exists")
				return &core.SSLManageUploadResult{
					CertId:   certRecord.ID,
					CertName: certRecord.Name,
				}, nil
			}
		}
	}

	// 生成新证书名（需符合天翼云命名规则）
	certName := fmt.Sprintf("certimate-%d", time.Now().UnixMilli())

	// 创建证书
	// REF: https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=24&api=5685&data=88&isNormal=1&vid=82
	createCertificateReq := &ctyunelb.CreateCertificateRequest{
		ClientToken: xtypes.ToPtr(generateClientToken()),
		RegionID:    xtypes.ToPtr(m.config.RegionId),
		Name:        xtypes.ToPtr(certName),
		Description: xtypes.ToPtr("upload from certimate"),
		Type:        xtypes.ToPtr("Server"),
		Certificate: xtypes.ToPtr(certPEM),
		PrivateKey:  xtypes.ToPtr(privkeyPEM),
	}
	createCertificateResp, err := m.sdkClient.CreateCertificate(createCertificateReq)
	m.logger.Debug("sdk request 'elb.CreateCertificate'", slog.Any("request", createCertificateReq), slog.Any("response", createCertificateResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'elb.CreateCertificate': %w", err)
	}

	return &core.SSLManageUploadResult{
		CertId:   createCertificateResp.ReturnObj.ID,
		CertName: certName,
	}, nil
}

func createSDKClient(accessKeyId, secretAccessKey string) (*ctyunelb.Client, error) {
	return ctyunelb.NewClient(accessKeyId, secretAccessKey)
}

func generateClientToken() string {
	return uuid.New().String()
}
