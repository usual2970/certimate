package ctcccloudcms

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	ctyuncms "github.com/usual2970/certimate/internal/pkg/sdk3rd/ctyun/cms"
	certutil "github.com/usual2970/certimate/internal/pkg/utils/cert"
	typeutil "github.com/usual2970/certimate/internal/pkg/utils/type"
)

type UploaderConfig struct {
	// 天翼云 AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// 天翼云 SecretAccessKey。
	SecretAccessKey string `json:"secretAccessKey"`
}

type UploaderProvider struct {
	config    *UploaderConfig
	logger    *slog.Logger
	sdkClient *ctyuncms.Client
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
	// 遍历证书列表，避免重复上传
	if res, _ := u.findCertIfExists(ctx, certPEM); res != nil {
		return res, nil
	}

	// 提取服务器证书
	serverCertPEM, intermediaCertPEM, err := certutil.ExtractCertificatesFromPEM(certPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to extract certs: %w", err)
	}

	// 生成新证书名（需符合天翼云命名规则）
	certName := fmt.Sprintf("cm%d", time.Now().Unix())

	// 上传证书
	// REF: https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=152&api=17243&data=204&isNormal=1&vid=283
	uploadCertificateReq := &ctyuncms.UploadCertificateRequest{
		Name:               typeutil.ToPtr(certName),
		Certificate:        typeutil.ToPtr(serverCertPEM),
		CertificateChain:   typeutil.ToPtr(intermediaCertPEM),
		PrivateKey:         typeutil.ToPtr(privkeyPEM),
		EncryptionStandard: typeutil.ToPtr("INTERNATIONAL"),
	}
	uploadCertificateResp, err := u.sdkClient.UploadCertificate(uploadCertificateReq)
	u.logger.Debug("sdk request 'cms.UploadCertificate'", slog.Any("request", uploadCertificateReq), slog.Any("response", uploadCertificateResp))
	if err != nil {
		if uploadCertificateResp != nil && uploadCertificateResp.GetError() == "CCMS_100000067" {
			if res, err := u.findCertIfExists(ctx, certPEM); err != nil {
				return nil, err
			} else if res == nil {
				return nil, errors.New("ctyun cms: no certificate found")
			} else {
				u.logger.Info("ssl certificate already exists")
				return res, nil
			}
		}

		return nil, fmt.Errorf("failed to execute sdk request 'cms.UploadCertificate': %w", err)
	}

	// 遍历证书列表，获取刚刚上传证书 ID
	if res, err := u.findCertIfExists(ctx, certPEM); err != nil {
		return nil, err
	} else if res == nil {
		return nil, fmt.Errorf("no ssl certificate found, may be upload failed")
	} else {
		return res, nil
	}
}

func (u *UploaderProvider) findCertIfExists(ctx context.Context, certPEM string) (*uploader.UploadResult, error) {
	// 解析证书内容
	certX509, err := certutil.ParseCertificateFromPEM(certPEM)
	if err != nil {
		return nil, err
	}

	// 查询用户证书列表
	// REF: https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=152&api=17233&data=204&isNormal=1&vid=283
	getCertificateListPageNum := int32(1)
	getCertificateListPageSize := int32(10)
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		getCertificateListReq := &ctyuncms.GetCertificateListRequest{
			PageNum:  typeutil.ToPtr(getCertificateListPageNum),
			PageSize: typeutil.ToPtr(getCertificateListPageSize),
			Keyword:  typeutil.ToPtr(certX509.Subject.CommonName),
			Origin:   typeutil.ToPtr("UPLOAD"),
		}
		getCertificateListResp, err := u.sdkClient.GetCertificateList(getCertificateListReq)
		u.logger.Debug("sdk request 'cms.GetCertificateList'", slog.Any("request", getCertificateListReq), slog.Any("response", getCertificateListResp))
		if err != nil {
			return nil, fmt.Errorf("failed to execute sdk request 'cms.GetCertificateList': %w", err)
		}

		if getCertificateListResp.ReturnObj != nil {
			fingerprint := sha1.Sum(certX509.Raw)
			fingerprintHex := hex.EncodeToString(fingerprint[:])

			for _, certRecord := range getCertificateListResp.ReturnObj.List {
				// 对比证书名称
				if !strings.EqualFold(strings.Join(certX509.DNSNames, ","), certRecord.DomainName) {
					continue
				}

				// 对比证书有效期
				oldCertNotBefore, _ := time.Parse("2006-01-02T15:04:05Z", certRecord.IssueTime)
				oldCertNotAfter, _ := time.Parse("2006-01-02T15:04:05Z", certRecord.ExpireTime)
				if !certX509.NotBefore.Equal(oldCertNotBefore) {
					continue
				} else if !certX509.NotAfter.Equal(oldCertNotAfter) {
					continue
				}

				// 对比证书指纹
				if !strings.EqualFold(fingerprintHex, certRecord.Fingerprint) {
					continue
				}

				// 如果以上信息都一致，则视为已存在相同证书，直接返回
				u.logger.Info("ssl certificate already exists")
				return &uploader.UploadResult{
					CertId:   string(*&certRecord.Id),
					CertName: certRecord.Name,
				}, nil
			}
		}

		if getCertificateListResp.ReturnObj == nil || len(getCertificateListResp.ReturnObj.List) < int(getCertificateListPageSize) {
			break
		} else {
			getCertificateListPageNum++
		}
	}

	return nil, nil
}

func createSdkClient(accessKeyId, secretAccessKey string) (*ctyuncms.Client, error) {
	return ctyuncms.NewClient(accessKeyId, secretAccessKey)
}
