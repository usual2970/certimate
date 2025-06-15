package ctcccloudao

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"time"

	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	ctyunao "github.com/usual2970/certimate/internal/pkg/sdk3rd/ctyun/ao"
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
	sdkClient *ctyunao.Client
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

	// 查询用户名下证书列表，避免重复上传
	// REF: https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=113&api=13175&data=174&isNormal=1&vid=167
	listCertPage := int32(1)
	listCertPerPage := int32(1000)
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		listCertsReq := &ctyunao.ListCertsRequest{
			Page:      xtypes.ToPtr(listCertPage),
			PerPage:   xtypes.ToPtr(listCertPerPage),
			UsageMode: xtypes.ToPtr(int32(0)),
		}
		listCertsResp, err := u.sdkClient.ListCerts(listCertsReq)
		u.logger.Debug("sdk request 'ao.ListCerts'", slog.Any("request", listCertsReq), slog.Any("response", listCertsResp))
		if err != nil {
			return nil, fmt.Errorf("failed to execute sdk request 'ao.ListCerts': %w", err)
		}

		if listCertsResp.ReturnObj != nil {
			for _, certRecord := range listCertsResp.ReturnObj.Results {
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
				// REF: https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=113&api=13015&data=174&isNormal=1&vid=167
				queryCertReq := &ctyunao.QueryCertRequest{
					Id: xtypes.ToPtr(certRecord.Id),
				}
				queryCertResp, err := u.sdkClient.QueryCert(queryCertReq)
				u.logger.Debug("sdk request 'ao.QueryCert'", slog.Any("request", queryCertReq), slog.Any("response", queryCertResp))
				if err != nil {
					return nil, fmt.Errorf("failed to execute sdk request 'ao.QueryCert': %w", err)
				} else if queryCertResp.ReturnObj != nil && queryCertResp.ReturnObj.Result != nil {
					var isSameCert bool
					if queryCertResp.ReturnObj.Result.Certs == certPEM {
						isSameCert = true
					} else {
						oldCertX509, err := xcert.ParseCertificateFromPEM(queryCertResp.ReturnObj.Result.Certs)
						if err != nil {
							continue
						}

						isSameCert = xcert.EqualCertificate(certX509, oldCertX509)
					}

					// 如果已存在相同证书，直接返回
					if isSameCert {
						u.logger.Info("ssl certificate already exists")
						return &uploader.UploadResult{
							CertId:   fmt.Sprintf("%d", queryCertResp.ReturnObj.Result.Id),
							CertName: queryCertResp.ReturnObj.Result.Name,
						}, nil
					}
				}
			}
		}

		if listCertsResp.ReturnObj == nil || len(listCertsResp.ReturnObj.Results) < int(listCertPerPage) {
			break
		} else {
			listCertPage++
		}
	}

	// 生成新证书名（需符合天翼云命名规则）
	certName := fmt.Sprintf("certimate-%d", time.Now().UnixMilli())

	// 创建证书
	// REF: https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=113&api=13014&data=174&isNormal=1&vid=167
	createCertReq := &ctyunao.CreateCertRequest{
		Name:  xtypes.ToPtr(certName),
		Certs: xtypes.ToPtr(certPEM),
		Key:   xtypes.ToPtr(privkeyPEM),
	}
	createCertResp, err := u.sdkClient.CreateCert(createCertReq)
	u.logger.Debug("sdk request 'ao.CreateCert'", slog.Any("request", createCertReq), slog.Any("response", createCertResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'ao.CreateCert': %w", err)
	}

	return &uploader.UploadResult{
		CertId:   fmt.Sprintf("%d", createCertResp.ReturnObj.Id),
		CertName: certName,
	}, nil
}

func createSdkClient(accessKeyId, secretAccessKey string) (*ctyunao.Client, error) {
	return ctyunao.NewClient(accessKeyId, secretAccessKey)
}
