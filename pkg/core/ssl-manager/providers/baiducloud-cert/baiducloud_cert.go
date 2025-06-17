package baiducloudcert

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/certimate-go/certimate/pkg/core"
	bdsdk "github.com/certimate-go/certimate/pkg/sdk3rd/baiducloud/cert"
	xcert "github.com/certimate-go/certimate/pkg/utils/cert"
)

type SSLManagerProviderConfig struct {
	// 百度智能云 AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// 百度智能云 SecretAccessKey。
	SecretAccessKey string `json:"secretAccessKey"`
}

type SSLManagerProvider struct {
	config    *SSLManagerProviderConfig
	logger    *slog.Logger
	sdkClient *bdsdk.Client
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

	// 遍历证书列表，避免重复上传
	// REF: https://cloud.baidu.com/doc/Reference/s/Gjwvz27xu#35-%E6%9F%A5%E7%9C%8B%E8%AF%81%E4%B9%A6%E5%88%97%E8%A1%A8%E8%AF%A6%E6%83%85
	listCertDetail, err := m.sdkClient.ListCertDetail()
	m.logger.Debug("sdk request 'cert.ListCertDetail'", slog.Any("response", listCertDetail))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'cert.ListCertDetail': %w", err)
	} else {
		for _, certDetail := range listCertDetail.Certs {
			// 先对比证书通用名称
			if !strings.EqualFold(certX509.Subject.CommonName, certDetail.CertCommonName) {
				continue
			}

			// 再对比证书有效期
			oldCertNotBefore, _ := time.Parse("2006-01-02T15:04:05Z", certDetail.CertStartTime)
			oldCertNotAfter, _ := time.Parse("2006-01-02T15:04:05Z", certDetail.CertStopTime)
			if !certX509.NotBefore.Equal(oldCertNotBefore) || !certX509.NotAfter.Equal(oldCertNotAfter) {
				continue
			}

			// 再对比证书多域名
			if certDetail.CertDNSNames != strings.Join(certX509.DNSNames, ",") {
				continue
			}

			// 最后对比证书内容
			getCertDetailResp, err := m.sdkClient.GetCertRawData(certDetail.CertId)
			m.logger.Debug("sdk request 'cert.GetCertRawData'", slog.Any("certId", certDetail.CertId), slog.Any("response", getCertDetailResp))
			if err != nil {
				return nil, fmt.Errorf("failed to execute sdk request 'cert.GetCertRawData': %w", err)
			} else {
				oldCertX509, err := xcert.ParseCertificateFromPEM(getCertDetailResp.CertServerData)
				if err != nil {
					continue
				}
				if !xcert.EqualCertificate(certX509, oldCertX509) {
					continue
				}
			}

			// 如果以上信息都一致，则视为已存在相同证书，直接返回
			m.logger.Info("ssl certificate already exists")
			return &core.SSLManageUploadResult{
				CertId:   certDetail.CertId,
				CertName: certDetail.CertName,
			}, nil
		}
	}

	// 创建证书
	// REF: https://cloud.baidu.com/doc/Reference/s/Gjwvz27xu#31-%E5%88%9B%E5%BB%BA%E8%AF%81%E4%B9%A6
	createCertReq := &bdsdk.CreateCertArgs{}
	createCertReq.CertName = fmt.Sprintf("certimate-%d", time.Now().UnixMilli())
	createCertReq.CertServerData = certPEM
	createCertReq.CertPrivateData = privkeyPEM
	createCertResp, err := m.sdkClient.CreateCert(createCertReq)
	m.logger.Debug("sdk request 'cert.CreateCert'", slog.Any("request", createCertReq), slog.Any("response", createCertResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'cert.CreateCert': %w", err)
	}

	return &core.SSLManageUploadResult{
		CertId:   createCertResp.CertId,
		CertName: createCertResp.CertName,
	}, nil
}

func createSDKClient(accessKeyId, secretAccessKey string) (*bdsdk.Client, error) {
	client, err := bdsdk.NewClient(accessKeyId, secretAccessKey, "")
	if err != nil {
		return nil, err
	}

	return client, nil
}
