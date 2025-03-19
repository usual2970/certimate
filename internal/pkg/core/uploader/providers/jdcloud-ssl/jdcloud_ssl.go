package jdcloudssl

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log/slog"
	"strings"
	"time"

	jdcore "github.com/jdcloud-api/jdcloud-sdk-go/core"
	jdsslapi "github.com/jdcloud-api/jdcloud-sdk-go/services/ssl/apis"
	jdsslclient "github.com/jdcloud-api/jdcloud-sdk-go/services/ssl/client"
	xerrors "github.com/pkg/errors"
	"golang.org/x/exp/slices"

	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	"github.com/usual2970/certimate/internal/pkg/utils/certutil"
)

type UploaderConfig struct {
	// 京东云 AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// 京东云 AccessKeySecret。
	AccessKeySecret string `json:"accessKeySecret"`
}

type UploaderProvider struct {
	config    *UploaderConfig
	logger    *slog.Logger
	sdkClient *jdsslclient.SslClient
}

var _ uploader.Uploader = (*UploaderProvider)(nil)

func NewUploader(config *UploaderConfig) (*UploaderProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	client, err := createSdkClient(config.AccessKeyId, config.AccessKeySecret)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk client")
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

func (u *UploaderProvider) Upload(ctx context.Context, certPem string, privkeyPem string) (res *uploader.UploadResult, err error) {
	// 解析证书内容
	certX509, err := certutil.ParseCertificateFromPEM(certPem)
	if err != nil {
		return nil, err
	}

	// 格式化私钥内容，以便后续计算私钥摘要
	privkeyPem = strings.TrimSpace(privkeyPem)
	privkeyPem = strings.ReplaceAll(privkeyPem, "\r", "")
	privkeyPem = strings.ReplaceAll(privkeyPem, "\n", "\r\n")
	privkeyPem = privkeyPem + "\r\n"

	// 遍历查看证书列表，避免重复上传
	// REF: https://docs.jdcloud.com/cn/ssl-certificate/api/describecerts
	describeCertsPageNumber := 1
	describeCertsPageSize := 10
	for {
		describeCertsReq := jdsslapi.NewDescribeCertsRequest()
		describeCertsReq.SetDomainName(certX509.Subject.CommonName)
		describeCertsReq.SetPageNumber(describeCertsPageNumber)
		describeCertsReq.SetPageSize(describeCertsPageSize)
		describeCertsResp, err := u.sdkClient.DescribeCerts(describeCertsReq)
		u.logger.Debug("sdk request 'ssl.DescribeCerts'", slog.Any("request", describeCertsReq), slog.Any("response", describeCertsResp))
		if err != nil {
			return nil, xerrors.Wrap(err, "failed to execute sdk request 'ssl.DescribeCerts'")
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
			newKeyDigest := sha256.Sum256([]byte(privkeyPem))
			newKeyDigestHex := hex.EncodeToString(newKeyDigest[:])
			if !strings.EqualFold(newKeyDigestHex, certDetail.Digest) {
				continue
			}

			// 如果以上信息都一致，则视为已存在相同证书，直接返回
			u.logger.Info("ssl certificate already exists")
			return &uploader.UploadResult{
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
	uploadCertReq := jdsslapi.NewUploadCertRequest(certName, privkeyPem, certPem)
	uploadCertResp, err := u.sdkClient.UploadCert(uploadCertReq)
	u.logger.Debug("sdk request 'ssl.UploadCertificate'", slog.Any("request", uploadCertReq), slog.Any("response", uploadCertResp))
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'ssl.UploadCertificate'")
	}

	return &uploader.UploadResult{
		CertId:   uploadCertResp.Result.CertId,
		CertName: certName,
	}, nil
}

func createSdkClient(accessKeyId, accessKeySecret string) (*jdsslclient.SslClient, error) {
	clientCredentials := jdcore.NewCredentials(accessKeyId, accessKeySecret)
	client := jdsslclient.NewSslClient(clientCredentials)
	client.SetLogger(jdcore.NewDefaultLogger(jdcore.LogWarn))
	return client, nil
}
