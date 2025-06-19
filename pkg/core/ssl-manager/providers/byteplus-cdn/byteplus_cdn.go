package bytepluscdn

import (
	"context"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	bytepluscdn "github.com/byteplus-sdk/byteplus-sdk-golang/service/cdn"

	"github.com/certimate-go/certimate/pkg/core"
	xcert "github.com/certimate-go/certimate/pkg/utils/cert"
)

type SSLManagerProviderConfig struct {
	// BytePlus AccessKey。
	AccessKey string `json:"accessKey"`
	// BytePlus SecretKey。
	SecretKey string `json:"secretKey"`
}

type SSLManagerProvider struct {
	config    *SSLManagerProviderConfig
	logger    *slog.Logger
	sdkClient *bytepluscdn.CDN
}

var _ core.SSLManager = (*SSLManagerProvider)(nil)

func NewSSLManagerProvider(config *SSLManagerProviderConfig) (*SSLManagerProvider, error) {
	if config == nil {
		return nil, errors.New("the configuration of the ssl manager provider is nil")
	}

	client := bytepluscdn.NewInstance()
	client.Client.SetAccessKey(config.AccessKey)
	client.Client.SetSecretKey(config.SecretKey)

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
	// REF: https://docs.byteplus.com/en/docs/byteplus-cdn/reference-listcertinfo
	listCertInfoPageNum := int64(1)
	listCertInfoPageSize := int64(100)
	listCertInfoTotal := 0
	listCertInfoReq := &bytepluscdn.ListCertInfoRequest{
		PageNum:  bytepluscdn.GetInt64Ptr(listCertInfoPageNum),
		PageSize: bytepluscdn.GetInt64Ptr(listCertInfoPageSize),
		Source:   bytepluscdn.GetStrPtr("cert_center"),
	}
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		listCertInfoResp, err := m.sdkClient.ListCertInfo(listCertInfoReq)
		m.logger.Debug("sdk request 'cdn.ListCertInfo'", slog.Any("request", listCertInfoReq), slog.Any("response", listCertInfoResp))
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
					m.logger.Info("ssl certificate already exists")
					return &core.SSLManageUploadResult{
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

	// 生成新证书名（需符合 BytePlus 命名规则）
	certName := fmt.Sprintf("certimate-%d", time.Now().UnixMilli())

	// 上传新证书
	// REF: https://docs.byteplus.com/en/docs/byteplus-cdn/reference-addcertificate
	addCertificateReq := &bytepluscdn.AddCertificateRequest{
		Certificate: certPEM,
		PrivateKey:  privkeyPEM,
		Source:      bytepluscdn.GetStrPtr("cert_center"),
		Desc:        bytepluscdn.GetStrPtr(certName),
	}
	addCertificateResp, err := m.sdkClient.AddCertificate(addCertificateReq)
	m.logger.Debug("sdk request 'cdn.AddCertificate'", slog.Any("request", addCertificateReq), slog.Any("response", addCertificateResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'cdn.AddCertificate': %w", err)
	}

	return &core.SSLManageUploadResult{
		CertId:   addCertificateResp.Result.CertId,
		CertName: certName,
	}, nil
}
