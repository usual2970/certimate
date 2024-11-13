package aliyunslb

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	aliyunOpen "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	aliyunSlb "github.com/alibabacloud-go/slb-20140515/v4/client"
	"github.com/alibabacloud-go/tea/tea"
	xerrors "github.com/pkg/errors"

	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	"github.com/usual2970/certimate/internal/pkg/utils/x509"
)

type AliyunSLBUploaderConfig struct {
	AccessKeyId     string `json:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret"`
	Region          string `json:"region"`
}

type AliyunSLBUploader struct {
	config    *AliyunSLBUploaderConfig
	sdkClient *aliyunSlb.Client
}

var _ uploader.Uploader = (*AliyunSLBUploader)(nil)

func New(config *AliyunSLBUploaderConfig) (*AliyunSLBUploader, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}

	client, err := createSdkClient(
		config.AccessKeyId,
		config.AccessKeySecret,
		config.Region,
	)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk client")
	}

	return &AliyunSLBUploader{
		config:    config,
		sdkClient: client,
	}, nil
}

func (u *AliyunSLBUploader) Upload(ctx context.Context, certPem string, privkeyPem string) (res *uploader.UploadResult, err error) {
	// 解析证书内容
	certX509, err := x509.ParseCertificateFromPEM(certPem)
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
	if region == "" {
		region = "cn-hangzhou" // SLB 服务默认区域：华东一杭州
	}

	aConfig := &aliyunOpen.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
	}

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
	aConfig.Endpoint = tea.String(endpoint)

	client, err := aliyunSlb.NewClient(aConfig)
	if err != nil {
		return nil, err
	}

	return client, nil
}
