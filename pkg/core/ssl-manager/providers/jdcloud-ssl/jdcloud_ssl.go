package jdcloudssl

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	jdcore "github.com/jdcloud-api/jdcloud-sdk-go/core"
	jdsslapi "github.com/jdcloud-api/jdcloud-sdk-go/services/ssl/apis"
	jdsslclient "github.com/jdcloud-api/jdcloud-sdk-go/services/ssl/client"
	"golang.org/x/exp/slices"

	"github.com/certimate-go/certimate/pkg/core"
	xcert "github.com/certimate-go/certimate/pkg/utils/cert"
)

type SSLManagerProviderConfig struct {
	// 京东云 AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// 京东云 AccessKeySecret。
	AccessKeySecret string `json:"accessKeySecret"`
}

type SSLManagerProvider struct {
	config    *SSLManagerProviderConfig
	logger    *slog.Logger
	sdkClient *jdsslclient.SslClient
}

var _ core.SSLManager = (*SSLManagerProvider)(nil)

func NewSSLManagerProvider(config *SSLManagerProviderConfig) (*SSLManagerProvider, error) {
	if config == nil {
		return nil, errors.New("the configuration of the ssl manager provider is nil")
	}

	client, err := createSDKClient(config.AccessKeyId, config.AccessKeySecret)
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

	// 格式化私钥内容，以便后续计算私钥摘要
	privkeyPEM = strings.TrimSpace(privkeyPEM)
	privkeyPEM = strings.ReplaceAll(privkeyPEM, "\r", "")
	privkeyPEM = strings.ReplaceAll(privkeyPEM, "\n", "\r\n")
	privkeyPEM = privkeyPEM + "\r\n"

	// 遍历查看证书列表，避免重复上传
	// REF: https://docs.jdcloud.com/cn/ssl-certificate/api/describecerts
	describeCertsPageNumber := 1
	describeCertsPageSize := 10
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		describeCertsReq := jdsslapi.NewDescribeCertsRequest()
		describeCertsReq.SetDomainName(certX509.Subject.CommonName)
		describeCertsReq.SetPageNumber(describeCertsPageNumber)
		describeCertsReq.SetPageSize(describeCertsPageSize)
		describeCertsResp, err := m.sdkClient.DescribeCerts(describeCertsReq)
		m.logger.Debug("sdk request 'ssl.DescribeCerts'", slog.Any("request", describeCertsReq), slog.Any("response", describeCertsResp))
		if err != nil {
			return nil, fmt.Errorf("failed to execute sdk request 'ssl.DescribeCerts': %w", err)
		}

		for _, certDetail := range describeCertsResp.Result.CertListDetails {
			// 先对比证书通用名称
			if !strings.EqualFold(certX509.Subject.CommonName, certDetail.CommonName) {
				continue
			}

			// 再对比证书多域名
			if !slices.Equal(certX509.DNSNames, certDetail.DnsNames) {
				continue
			}

			// 再对比证书有效期
			oldCertNotBefore, _ := time.Parse(time.RFC3339, certDetail.StartTime)
			oldCertNotAfter, _ := time.Parse(time.RFC3339, certDetail.EndTime)
			if !certX509.NotBefore.Equal(oldCertNotBefore) || !certX509.NotAfter.Equal(oldCertNotAfter) {
				continue
			}

			// 最后对比私钥摘要
			newKeyDigest := sha256.Sum256([]byte(privkeyPEM))
			newKeyDigestHex := hex.EncodeToString(newKeyDigest[:])
			if !strings.EqualFold(newKeyDigestHex, certDetail.Digest) {
				continue
			}

			// 如果以上信息都一致，则视为已存在相同证书，直接返回
			m.logger.Info("ssl certificate already exists")
			return &core.SSLManageUploadResult{
				CertId:   certDetail.CertId,
				CertName: certDetail.CertName,
			}, nil
		}

		if len(describeCertsResp.Result.CertListDetails) < int(describeCertsPageSize) {
			break
		} else {
			describeCertsPageNumber++
		}
	}

	// 生成新证书名（需符合京东云命名规则）
	certName := fmt.Sprintf("certimate-%d", time.Now().UnixMilli())

	// 上传证书
	// REF: https://docs.jdcloud.com/cn/ssl-certificate/api/uploadcert
	uploadCertReq := jdsslapi.NewUploadCertRequest(certName, privkeyPEM, certPEM)
	uploadCertResp, err := m.sdkClient.UploadCert(uploadCertReq)
	m.logger.Debug("sdk request 'ssl.UploadCertificate'", slog.Any("request", uploadCertReq), slog.Any("response", uploadCertResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'ssl.UploadCertificate': %w", err)
	}

	return &core.SSLManageUploadResult{
		CertId:   uploadCertResp.Result.CertId,
		CertName: certName,
	}, nil
}

func createSDKClient(accessKeyId, accessKeySecret string) (*jdsslclient.SslClient, error) {
	clientCredentials := jdcore.NewCredentials(accessKeyId, accessKeySecret)
	client := jdsslclient.NewSslClient(clientCredentials)
	client.SetLogger(jdcore.NewDefaultLogger(jdcore.LogWarn))
	return client, nil
}
