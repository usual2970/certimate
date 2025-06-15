package ctcccloudelb

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	ctyunelb "github.com/usual2970/certimate/internal/pkg/sdk3rd/ctyun/elb"
	certutil "github.com/usual2970/certimate/internal/pkg/utils/cert"
	typeutil "github.com/usual2970/certimate/internal/pkg/utils/type"
)

type UploaderConfig struct {
	// 天翼云 AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// 天翼云 SecretAccessKey。
	SecretAccessKey string `json:"secretAccessKey"`
	// 天翼云资源池 ID。
	RegionId string `json:"regionId"`
}

type UploaderProvider struct {
	config    *UploaderConfig
	logger    *slog.Logger
	sdkClient *ctyunelb.Client
}

var _ uploader.Uploader = (*UploaderProvider)(nil)

func NewUploader(config *UploaderConfig) (*UploaderProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	client, err := createSdkClient(config.AccessKeyId, config.SecretAccessKey)
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
	// 解析证书内容
	certX509, err := certutil.ParseCertificateFromPEM(certPEM)
	if err != nil {
		return nil, err
	}

	// 查询证书列表，避免重复上传
	// REF: https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=24&api=5692&data=88&isNormal=1&vid=82
	listCertificatesReq := &ctyunelb.ListCertificatesRequest{
		RegionID: typeutil.ToPtr(u.config.RegionId),
	}
	listCertificatesResp, err := u.sdkClient.ListCertificates(listCertificatesReq)
	u.logger.Debug("sdk request 'elb.ListCertificates'", slog.Any("request", listCertificatesReq), slog.Any("response", listCertificatesResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'elb.ListCertificates': %w", err)
	} else {
		for _, certRecord := range listCertificatesResp.ReturnObj {
			var isSameCert bool
			if certRecord.Certificate == certPEM {
				isSameCert = true
			} else {
				oldCertX509, err := certutil.ParseCertificateFromPEM(certRecord.Certificate)
				if err != nil {
					continue
				}

				isSameCert = certutil.EqualCertificate(certX509, oldCertX509)
			}

			// 如果已存在相同证书，直接返回
			if isSameCert {
				u.logger.Info("ssl certificate already exists")
				return &uploader.UploadResult{
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
		ClientToken: typeutil.ToPtr(generateClientToken()),
		RegionID:    typeutil.ToPtr(u.config.RegionId),
		Name:        typeutil.ToPtr(certName),
		Description: typeutil.ToPtr("upload from certimate"),
		Type:        typeutil.ToPtr("Server"),
		Certificate: typeutil.ToPtr(certPEM),
		PrivateKey:  typeutil.ToPtr(privkeyPEM),
	}
	createCertificateResp, err := u.sdkClient.CreateCertificate(createCertificateReq)
	u.logger.Debug("sdk request 'elb.CreateCertificate'", slog.Any("request", createCertificateReq), slog.Any("response", createCertificateResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'elb.CreateCertificate': %w", err)
	}

	return &uploader.UploadResult{
		CertId:   createCertificateResp.ReturnObj.ID,
		CertName: certName,
	}, nil
}

func createSdkClient(accessKeyId, secretAccessKey string) (*ctyunelb.Client, error) {
	return ctyunelb.NewClient(accessKeyId, secretAccessKey)
}

func generateClientToken() string {
	return uuid.New().String()
}
