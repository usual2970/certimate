package ctcccloudcdn

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"time"

	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	ctyuncdn "github.com/usual2970/certimate/internal/pkg/sdk3rd/ctyun/cdn"
	xcert "github.com/usual2970/certimate/internal/pkg/utils/cert"
	xtypes "github.com/usual2970/certimate/internal/pkg/utils/types"
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
	sdkClient *ctyuncdn.Client
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
	certX509, err := xcert.ParseCertificateFromPEM(certPEM)
	if err != nil {
		return nil, err
	}

	// 查询证书列表，避免重复上传
	// REF: https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=108&api=10901&data=161&isNormal=1&vid=154
	queryCertListPage := int32(1)
	queryCertListPerPage := int32(1000)
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		queryCertListReq := &ctyuncdn.QueryCertListRequest{
			Page:      xtypes.ToPtr(queryCertListPage),
			PerPage:   xtypes.ToPtr(queryCertListPerPage),
			UsageMode: xtypes.ToPtr(int32(0)),
		}
		queryCertListResp, err := u.sdkClient.QueryCertList(queryCertListReq)
		u.logger.Debug("sdk request 'cdn.QueryCertList'", slog.Any("request", queryCertListReq), slog.Any("response", queryCertListResp))
		if err != nil {
			return nil, fmt.Errorf("failed to execute sdk request 'cdn.QueryCertList': %w", err)
		}

		if queryCertListResp.ReturnObj != nil {
			for _, certRecord := range queryCertListResp.ReturnObj.Results {
				// 对比证书通用名称
				if !strings.EqualFold(certX509.Subject.CommonName, certRecord.CN) {
					continue
				}

				// 对比证书扩展名称
				if !slices.Equal(certX509.DNSNames, certRecord.SANs) {
					continue
				}

				// 对比证书有效期
				if !certX509.NotBefore.Equal(time.Unix(certRecord.IssueTime, 0).UTC()) {
					continue
				} else if !certX509.NotAfter.Equal(time.Unix(certRecord.ExpiresTime, 0).UTC()) {
					continue
				}

				// 查询证书详情
				// REF: https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=108&api=10899&data=161&isNormal=1&vid=154
				queryCertDetailReq := &ctyuncdn.QueryCertDetailRequest{
					Id: xtypes.ToPtr(certRecord.Id),
				}
				queryCertDetailResp, err := u.sdkClient.QueryCertDetail(queryCertDetailReq)
				u.logger.Debug("sdk request 'cdn.QueryCertDetail'", slog.Any("request", queryCertDetailReq), slog.Any("response", queryCertDetailResp))
				if err != nil {
					return nil, fmt.Errorf("failed to execute sdk request 'cdn.QueryCertDetail': %w", err)
				} else if queryCertDetailResp.ReturnObj != nil && queryCertDetailResp.ReturnObj.Result != nil {
					var isSameCert bool
					if queryCertDetailResp.ReturnObj.Result.Certs == certPEM {
						isSameCert = true
					} else {
						oldCertX509, err := xcert.ParseCertificateFromPEM(queryCertDetailResp.ReturnObj.Result.Certs)
						if err != nil {
							continue
						}

						isSameCert = xcert.EqualCertificate(certX509, oldCertX509)
					}

					// 如果已存在相同证书，直接返回
					if isSameCert {
						u.logger.Info("ssl certificate already exists")
						return &uploader.UploadResult{
							CertId:   fmt.Sprintf("%d", queryCertDetailResp.ReturnObj.Result.Id),
							CertName: queryCertDetailResp.ReturnObj.Result.Name,
						}, nil
					}
				}
			}
		}

		if queryCertListResp.ReturnObj == nil || len(queryCertListResp.ReturnObj.Results) < int(queryCertListPerPage) {
			break
		} else {
			queryCertListPage++
		}
	}

	// 生成新证书名（需符合天翼云命名规则）
	certName := fmt.Sprintf("certimate-%d", time.Now().UnixMilli())

	// 创建证书
	// REF: https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=108&api=10893&data=161&isNormal=1&vid=154
	createCertReq := &ctyuncdn.CreateCertRequest{
		Name:  xtypes.ToPtr(certName),
		Certs: xtypes.ToPtr(certPEM),
		Key:   xtypes.ToPtr(privkeyPEM),
	}
	createCertResp, err := u.sdkClient.CreateCert(createCertReq)
	u.logger.Debug("sdk request 'cdn.CreateCert'", slog.Any("request", createCertReq), slog.Any("response", createCertResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'cdn.CreateCert': %w", err)
	}

	return &uploader.UploadResult{
		CertId:   fmt.Sprintf("%d", createCertResp.ReturnObj.Id),
		CertName: certName,
	}, nil
}

func createSdkClient(accessKeyId, secretAccessKey string) (*ctyuncdn.Client, error) {
	return ctyuncdn.NewClient(accessKeyId, secretAccessKey)
}
