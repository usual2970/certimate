package aliyunslb

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
	"time"

	aliyunOpen "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	aliyunSlb "github.com/alibabacloud-go/slb-20140515/v4/client"
	"github.com/alibabacloud-go/tea/tea"
	xerrors "github.com/pkg/errors"

	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	"github.com/usual2970/certimate/internal/pkg/utils/certs"
)

type UploaderConfig struct {
	// 阿里云 AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// 阿里云 AccessKeySecret。
	AccessKeySecret string `json:"accessKeySecret"`
	// 阿里云地域。
	Region string `json:"region"`
}

type UploaderProvider struct {
	config    *UploaderConfig
	sdkClient *aliyunSlb.Client
}

var _ uploader.Uploader = (*UploaderProvider)(nil)

func NewUploader(config *UploaderConfig) (*UploaderProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	client, err := createSdkClient(
		config.AccessKeyId,
		config.AccessKeySecret,
		config.Region,
	)
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

	// 查询证书列表，避免重复上传
	// REF: https://help.aliyun.com/zh/slb/classic-load-balancer/developer-reference/api-slb-2014-05-15-describeservercertificates
	describeServerCertificatesReq := &aliyunSlb.DescribeServerCertificatesRequest{
		RegionId: tea.String(u.config.Region),
	}
	describeServerCertificatesResp, err := u.sdkClient.DescribeServerCertificates(describeServerCertificatesReq)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'slb.DescribeServerCertificates'")
	}

	if describeServerCertificatesResp.Body.ServerCertificates != nil && describeServerCertificatesResp.Body.ServerCertificates.ServerCertificate != nil {
		fingerprint := sha256.Sum256(certX509.Raw)
		fingerprintHex := hex.EncodeToString(fingerprint[:])
		for _, certDetail := range describeServerCertificatesResp.Body.ServerCertificates.ServerCertificate {
			isSameCert := *certDetail.IsAliCloudCertificate == 0 &&
				strings.EqualFold(fingerprintHex, strings.ReplaceAll(*certDetail.Fingerprint, ":", "")) &&
				strings.EqualFold(certX509.Subject.CommonName, *certDetail.CommonName)
			// 如果已存在相同证书，直接返回已有的证书信息
			if isSameCert {
				return &uploader.UploadResult{
					CertId:   *certDetail.ServerCertificateId,
					CertName: *certDetail.ServerCertificateName,
				}, nil
			}
		}
	}

	// 生成新证书名（需符合阿里云命名规则）
	var certId, certName string
	certName = fmt.Sprintf("certimate_%d", time.Now().UnixMilli())

	// 去除证书和私钥内容中的空白行，以符合阿里云 API 要求
	// REF: https://github.com/usual2970/certimate/issues/326
	re := regexp.MustCompile(`(?m)^\s*$\n?`)
	certPem = strings.TrimSpace(re.ReplaceAllString(certPem, ""))
	privkeyPem = strings.TrimSpace(re.ReplaceAllString(privkeyPem, ""))

	// 上传新证书
	// REF: https://help.aliyun.com/zh/slb/classic-load-balancer/developer-reference/api-slb-2014-05-15-uploadservercertificate
	uploadServerCertificateReq := &aliyunSlb.UploadServerCertificateRequest{
		RegionId:              tea.String(u.config.Region),
		ServerCertificateName: tea.String(certName),
		ServerCertificate:     tea.String(certPem),
		PrivateKey:            tea.String(privkeyPem),
	}
	uploadServerCertificateResp, err := u.sdkClient.UploadServerCertificate(uploadServerCertificateReq)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'slb.UploadServerCertificate'")
	}

	certId = *uploadServerCertificateResp.Body.ServerCertificateId
	return &uploader.UploadResult{
		CertId:   certId,
		CertName: certName,
	}, nil
}

func createSdkClient(accessKeyId, accessKeySecret, region string) (*aliyunSlb.Client, error) {
	// 接入点一览 https://api.aliyun.com/product/Slb
	var endpoint string
	switch region {
	case
		"cn-hangzhou",
		"cn-hangzhou-finance",
		"cn-shanghai-finance-1",
		"cn-shenzhen-finance-1":
		endpoint = "slb.aliyuncs.com"
	default:
		endpoint = fmt.Sprintf("slb.%s.aliyuncs.com", region)
	}

	config := &aliyunOpen.Config{
		Endpoint:        tea.String(endpoint),
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
	}

	client, err := aliyunSlb.NewClient(config)
	if err != nil {
		return nil, err
	}

	return client, nil
}
