package jdcloudssl

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	jdCore "github.com/jdcloud-api/jdcloud-sdk-go/core"
	jdSslApi "github.com/jdcloud-api/jdcloud-sdk-go/services/ssl/apis"
	jdSslClient "github.com/jdcloud-api/jdcloud-sdk-go/services/ssl/client"
	xerrors "github.com/pkg/errors"
	"golang.org/x/exp/slices"

	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	"github.com/usual2970/certimate/internal/pkg/utils/certs"
)

type UploaderConfig struct {
	// 京东云 AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// 京东云 AccessKeySecret。
	AccessKeySecret string `json:"accessKeySecret"`
}

type UploaderProvider struct {
	config    *UploaderConfig
	sdkClient *jdSslClient.SslClient
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
		sdkClient: client,
	}, nil
}

func (u *UploaderProvider) Upload(ctx context.Context, certPem string, privkeyPem string) (res *uploader.UploadResult, err error) {
	// 解析证书内容
	certX509, err := certs.ParseCertificateFromPEM(certPem)
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
	describeCertsPageSize := 100
	for {
		describeCertsReq := jdSslApi.NewDescribeCertsRequest()
		describeCertsReq.SetDomainName(certX509.Subject.CommonName)
		describeCertsReq.SetPageNumber(describeCertsPageNumber)
		describeCertsReq.SetPageSize(describeCertsPageSize)
		describeCertsResp, err := u.sdkClient.DescribeCerts(describeCertsReq)
		if err != nil {
			return nil, xerrors.Wrap(err, "failed to execute sdk request 'ssl.DescribeCerts'")
		}

		for _, certDetail := range describeCertsResp.Result.CertListDetails {
			// 先尝试匹配 CN
			if !strings.EqualFold(certX509.Subject.CommonName, certDetail.CommonName) {
				continue
			}

			// 再尝试匹配 SAN
			if !slices.Equal(certX509.DNSNames, certDetail.DnsNames) {
				continue
			}

			// 再尝试匹配证书有效期
			oldCertNotBefore, _ := time.Parse(time.RFC3339, certDetail.StartTime)
			oldCertNotAfter, _ := time.Parse(time.RFC3339, certDetail.EndTime)
			if !certX509.NotBefore.Equal(oldCertNotBefore) || !certX509.NotAfter.Equal(oldCertNotAfter) {
				continue
			}

			// 最后尝试匹配私钥摘要
			newKeyDigest := sha256.Sum256([]byte(privkeyPem))
			newKeyDigestHex := hex.EncodeToString(newKeyDigest[:])
			if !strings.EqualFold(newKeyDigestHex, certDetail.Digest) {
				continue
			}

			// 如果以上都匹配，则视为已存在相同证书，直接返回已有的证书信息
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
	uploadCertReq := jdSslApi.NewUploadCertRequest(certName, privkeyPem, certPem)
	uploadCertResp, err := u.sdkClient.UploadCert(uploadCertReq)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'ssl.UploadCertificate'")
	}

	return &uploader.UploadResult{
		CertId:   uploadCertResp.Result.CertId,
		CertName: certName,
	}, nil
}

func createSdkClient(accessKeyId, accessKeySecret string) (*jdSslClient.SslClient, error) {
	clientCredentials := jdCore.NewCredentials(accessKeyId, accessKeySecret)
	client := jdSslClient.NewSslClient(clientCredentials)
	client.SetLogger(jdCore.NewDefaultLogger(jdCore.LogWarn))
	return client, nil
}
