package volcenginecdn

import (
	"context"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log/slog"
	"strings"
	"time"

	vecdn "github.com/volcengine/volc-sdk-golang/service/cdn"
	ve "github.com/volcengine/volcengine-go-sdk/volcengine"

	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	certutil "github.com/usual2970/certimate/internal/pkg/utils/cert"
)

type UploaderConfig struct {
	// 火山引擎 AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// 火山引擎 AccessKeySecret。
	AccessKeySecret string `json:"accessKeySecret"`
}

type UploaderProvider struct {
	config    *UploaderConfig
	logger    *slog.Logger
	sdkClient *vecdn.CDN
}

var _ uploader.Uploader = (*UploaderProvider)(nil)

func NewUploader(config *UploaderConfig) (*UploaderProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	client := vecdn.NewInstance()
	client.Client.SetAccessKey(config.AccessKeyId)
	client.Client.SetSecretKey(config.AccessKeySecret)

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

	// 查询证书列表，避免重复上传
	// REF: https://www.volcengine.com/docs/6454/125709
	listCertInfoPageNum := int64(1)
	listCertInfoPageSize := int64(100)
	listCertInfoTotal := 0
	listCertInfoReq := &vecdn.ListCertInfoRequest{
		PageNum:  ve.Int64(listCertInfoPageNum),
		PageSize: ve.Int64(listCertInfoPageSize),
		Source:   "volc_cert_center",
	}
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		listCertInfoResp, err := u.sdkClient.ListCertInfo(listCertInfoReq)
		u.logger.Debug("sdk request 'cdn.ListCertInfo'", slog.Any("request", listCertInfoReq), slog.Any("response", listCertInfoResp))
		if err != nil {
			return nil, fmt.Errorf("failed to execute sdk request 'cdn.ListCertInfo': %w", err)
		}

		if listCertInfoResp.Result.CertInfo != nil {
			for _, certDetail := range listCertInfoResp.Result.CertInfo {
				fingerprintSha1 := sha1.Sum(certX509.Raw)
				fingerprintSha256 := sha256.Sum256(certX509.Raw)
				isSameCert := strings.EqualFold(hex.EncodeToString(fingerprintSha1[:]), certDetail.CertFingerprint.Sha1) &&
					strings.EqualFold(hex.EncodeToString(fingerprintSha256[:]), certDetail.CertFingerprint.Sha256)
				// 如果已存在相同证书，直接返回
				if isSameCert {
					u.logger.Info("ssl certificate already exists")
					return &uploader.UploadResult{
						CertId:   certDetail.CertId,
						CertName: certDetail.Desc,
					}, nil
				}
			}
		}

		listCertInfoLen := len(listCertInfoResp.Result.CertInfo)
		if listCertInfoLen < int(listCertInfoPageSize) || int(listCertInfoResp.Result.Total) <= listCertInfoTotal+listCertInfoLen {
			break
		} else {
			listCertInfoPageNum++
			listCertInfoTotal += listCertInfoLen
		}
	}

	// 生成新证书名（需符合火山引擎命名规则）
	var certId, certName string
	certName = fmt.Sprintf("certimate-%d", time.Now().UnixMilli())

	// 上传新证书
	// REF: https://www.volcengine.com/docs/6454/1245763
	addCertificateReq := &vecdn.AddCertificateRequest{
		Certificate: certPEM,
		PrivateKey:  privkeyPEM,
		Source:      ve.String("volc_cert_center"),
		Desc:        ve.String(certName),
	}
	addCertificateResp, err := u.sdkClient.AddCertificate(addCertificateReq)
	u.logger.Debug("sdk request 'cdn.AddCertificate'", slog.Any("request", addCertificateResp), slog.Any("response", addCertificateResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'cdn.AddCertificate': %w", err)
	}

	certId = addCertificateResp.Result.CertId
	return &uploader.UploadResult{
		CertId:   certId,
		CertName: certName,
	}, nil
}
