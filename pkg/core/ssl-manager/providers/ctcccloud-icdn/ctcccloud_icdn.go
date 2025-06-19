package ctcccloudicdn

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"time"

	"github.com/certimate-go/certimate/pkg/core"
	ctyunicdn "github.com/certimate-go/certimate/pkg/sdk3rd/ctyun/icdn"
	xcert "github.com/certimate-go/certimate/pkg/utils/cert"
	xtypes "github.com/certimate-go/certimate/pkg/utils/types"
)

type SSLManagerProviderConfig struct {
	// 天翼云 AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// 天翼云 SecretAccessKey。
	SecretAccessKey string `json:"secretAccessKey"`
}

type SSLManagerProvider struct {
	config    *SSLManagerProviderConfig
	logger    *slog.Logger
	sdkClient *ctyunicdn.Client
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
	// REF: https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=112&api=10838&data=173&isNormal=1&vid=166
	queryCertListPage := int32(1)
	queryCertListPerPage := int32(1000)
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		queryCertListReq := &ctyunicdn.QueryCertListRequest{
			Page:      xtypes.ToPtr(queryCertListPage),
			PerPage:   xtypes.ToPtr(queryCertListPerPage),
			UsageMode: xtypes.ToPtr(int32(0)),
		}
		queryCertListResp, err := m.sdkClient.QueryCertList(queryCertListReq)
		m.logger.Debug("sdk request 'icdn.QueryCertList'", slog.Any("request", queryCertListReq), slog.Any("response", queryCertListResp))
		if err != nil {
			return nil, fmt.Errorf("failed to execute sdk request 'icdn.QueryCertList': %w", err)
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
				// REF: https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=112&api=10837&data=173&isNormal=1&vid=166
				queryCertDetailReq := &ctyunicdn.QueryCertDetailRequest{
					Id: xtypes.ToPtr(certRecord.Id),
				}
				queryCertDetailResp, err := m.sdkClient.QueryCertDetail(queryCertDetailReq)
				m.logger.Debug("sdk request 'icdn.QueryCertDetail'", slog.Any("request", queryCertDetailReq), slog.Any("response", queryCertDetailResp))
				if err != nil {
					return nil, fmt.Errorf("failed to execute sdk request 'icdn.QueryCertDetail': %w", err)
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
						m.logger.Info("ssl certificate already exists")
						return &core.SSLManageUploadResult{
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
	// REF: https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=112&api=10835&data=173&isNormal=1&vid=166
	createCertReq := &ctyunicdn.CreateCertRequest{
		Name:  xtypes.ToPtr(certName),
		Certs: xtypes.ToPtr(certPEM),
		Key:   xtypes.ToPtr(privkeyPEM),
	}
	createCertResp, err := m.sdkClient.CreateCert(createCertReq)
	m.logger.Debug("sdk request 'icdn.CreateCert'", slog.Any("request", createCertReq), slog.Any("response", createCertResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'icdn.CreateCert': %w", err)
	}

	return &core.SSLManageUploadResult{
		CertId:   fmt.Sprintf("%d", createCertResp.ReturnObj.Id),
		CertName: certName,
	}, nil
}

func createSDKClient(accessKeyId, secretAccessKey string) (*ctyunicdn.Client, error) {
	return ctyunicdn.NewClient(accessKeyId, secretAccessKey)
}
