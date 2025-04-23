package rainyunsslcenter

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	rainyunsdk "github.com/usual2970/certimate/internal/pkg/sdk3rd/rainyun"
	certutil "github.com/usual2970/certimate/internal/pkg/utils/cert"
)

type UploaderConfig struct {
	// 雨云 API 密钥。
	ApiKey string `json:"ApiKey"`
}

type UploaderProvider struct {
	config    *UploaderConfig
	logger    *slog.Logger
	sdkClient *rainyunsdk.Client
}

var _ uploader.Uploader = (*UploaderProvider)(nil)

func NewUploader(config *UploaderConfig) (*UploaderProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	client, err := createSdkClient(config.ApiKey)
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
	if res, err := u.getCertIfExists(ctx, certPEM); err != nil {
		return nil, err
	} else if res != nil {
		u.logger.Info("ssl certificate already exists")
		return res, nil
	}

	// SSL 证书上传
	// REF: https://apifox.com/apidoc/shared/a4595cc8-44c5-4678-a2a3-eed7738dab03/api-69943046
	sslCenterCreateReq := &rainyunsdk.SslCenterCreateRequest{
		Cert: certPEM,
		Key:  privkeyPEM,
	}
	sslCenterCreateResp, err := u.sdkClient.SslCenterCreate(sslCenterCreateReq)
	u.logger.Debug("sdk request 'sslcenter.Create'", slog.Any("request", sslCenterCreateReq), slog.Any("response", sslCenterCreateResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'sslcenter.Create': %w", err)
	}

	if res, err := u.getCertIfExists(ctx, certPEM); err != nil {
		return nil, err
	} else if res == nil {
		return nil, errors.New("rainyun sslcenter: no certificate found")
	} else {
		return res, nil
	}
}

func (u *UploaderProvider) getCertIfExists(ctx context.Context, certPEM string) (res *uploader.UploadResult, err error) {
	// 解析证书内容
	certX509, err := certutil.ParseCertificateFromPEM(certPEM)
	if err != nil {
		return nil, err
	}

	// 遍历 SSL 证书列表，避免重复上传
	// REF: https://apifox.com/apidoc/shared/a4595cc8-44c5-4678-a2a3-eed7738dab03/api-69943046
	// REF: https://apifox.com/apidoc/shared/a4595cc8-44c5-4678-a2a3-eed7738dab03/api-69943048
	sslCenterListPage := int32(1)
	sslCenterListPerPage := int32(100)
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		sslCenterListReq := &rainyunsdk.SslCenterListRequest{
			Filters: &rainyunsdk.SslCenterListFilters{
				Domain: &certX509.Subject.CommonName,
			},
			Page:    &sslCenterListPage,
			PerPage: &sslCenterListPerPage,
		}
		sslCenterListResp, err := u.sdkClient.SslCenterList(sslCenterListReq)
		u.logger.Debug("sdk request 'sslcenter.List'", slog.Any("request", sslCenterListReq), slog.Any("response", sslCenterListResp))
		if err != nil {
			return nil, fmt.Errorf("failed to execute sdk request 'sslcenter.List': %w", err)
		}

		if sslCenterListResp.Data != nil && sslCenterListResp.Data.Records != nil {
			for _, sslItem := range sslCenterListResp.Data.Records {
				// 先对比证书的多域名
				if sslItem.Domain != strings.Join(certX509.DNSNames, ", ") {
					continue
				}

				// 再对比证书的有效期
				if sslItem.StartDate != certX509.NotBefore.Unix() || sslItem.ExpireDate != certX509.NotAfter.Unix() {
					continue
				}

				// 最后对比证书内容
				sslCenterGetResp, err := u.sdkClient.SslCenterGet(sslItem.ID)
				if err != nil {
					return nil, fmt.Errorf("failed to execute sdk request 'sslcenter.Get': %w", err)
				}

				var isSameCert bool
				if sslCenterGetResp.Data != nil {
					if sslCenterGetResp.Data.Cert == certPEM {
						isSameCert = true
					} else {
						oldCertX509, err := certutil.ParseCertificateFromPEM(sslCenterGetResp.Data.Cert)
						if err != nil {
							continue
						}

						isSameCert = certutil.EqualCertificate(certX509, oldCertX509)
					}
				}

				// 如果已存在相同证书，直接返回
				if isSameCert {
					return &uploader.UploadResult{
						CertId: fmt.Sprintf("%d", sslItem.ID),
					}, nil
				}
			}
		}

		if sslCenterListResp.Data == nil || len(sslCenterListResp.Data.Records) < int(sslCenterListPerPage) {
			break
		} else {
			sslCenterListPage++
		}
	}

	return nil, nil
}

func createSdkClient(apiKey string) (*rainyunsdk.Client, error) {
	if apiKey == "" {
		return nil, errors.New("invalid rainyun api key")
	}

	client := rainyunsdk.NewClient(apiKey)
	return client, nil
}
