package aliyuncas

import (
	"context"
	"fmt"
	"strings"
	"time"

	aliyunCas "github.com/alibabacloud-go/cas-20200407/v3/client"
	aliyunOpen "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	xerrors "github.com/pkg/errors"

	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	"github.com/usual2970/certimate/internal/pkg/utils/x509"
)

type AliyunCASUploaderConfig struct {
	AccessKeyId     string `json:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret"`
	Region          string `json:"region"`
}

type AliyunCASUploader struct {
	config    *AliyunCASUploaderConfig
	sdkClient *aliyunCas.Client
}

var _ uploader.Uploader = (*AliyunCASUploader)(nil)

func New(config *AliyunCASUploaderConfig) (*AliyunCASUploader, error) {
	client, err := createSdkClient(
		config.AccessKeyId,
		config.AccessKeySecret,
		config.Region,
	)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk client")
	}

	return &AliyunCASUploader{
		config:    config,
		sdkClient: client,
	}, nil
}

func (u *AliyunCASUploader) Upload(ctx context.Context, certPem string, privkeyPem string) (res *uploader.UploadResult, err error) {
	// 解析证书内容
	certX509, err := x509.ParseCertificateFromPEM(certPem)
	if err != nil {
		return nil, err
	}

	// 查询证书列表，避免重复上传
	// REF: https://help.aliyun.com/zh/ssl-certificate/developer-reference/api-cas-2020-04-07-listusercertificateorder
	// REF: https://help.aliyun.com/zh/ssl-certificate/developer-reference/api-cas-2020-04-07-getusercertificatedetail
	listUserCertificateOrderPage := int64(1)
	listUserCertificateOrderLimit := int64(50)
	for {
		listUserCertificateOrderReq := &aliyunCas.ListUserCertificateOrderRequest{
			CurrentPage: tea.Int64(listUserCertificateOrderPage),
			ShowSize:    tea.Int64(listUserCertificateOrderLimit),
			OrderType:   tea.String("CERT"),
		}
		listUserCertificateOrderResp, err := u.sdkClient.ListUserCertificateOrder(listUserCertificateOrderReq)
		if err != nil {
			return nil, xerrors.Wrap(err, "failed to execute sdk request 'cas.ListUserCertificateOrder'")
		}

		if listUserCertificateOrderResp.Body.CertificateOrderList != nil {
			for _, certDetail := range listUserCertificateOrderResp.Body.CertificateOrderList {
				if strings.EqualFold(certX509.SerialNumber.Text(16), *certDetail.SerialNo) {
					getUserCertificateDetailReq := &aliyunCas.GetUserCertificateDetailRequest{
						CertId: certDetail.CertificateId,
					}
					getUserCertificateDetailResp, err := u.sdkClient.GetUserCertificateDetail(getUserCertificateDetailReq)
					if err != nil {
						return nil, xerrors.Wrap(err, "failed to execute sdk request 'cas.GetUserCertificateDetail'")
					}

					var isSameCert bool
					if *getUserCertificateDetailResp.Body.Cert == certPem {
						isSameCert = true
					} else {
						oldCertX509, err := x509.ParseCertificateFromPEM(*getUserCertificateDetailResp.Body.Cert)
						if err != nil {
							continue
						}

						isSameCert = x509.EqualCertificate(certX509, oldCertX509)
					}

					// 如果已存在相同证书，直接返回已有的证书信息
					if isSameCert {
						return &uploader.UploadResult{
							CertId:   fmt.Sprintf("%d", tea.Int64Value(certDetail.CertificateId)),
							CertName: *certDetail.Name,
						}, nil
					}
				}
			}
		}

		if listUserCertificateOrderResp.Body.CertificateOrderList == nil || len(listUserCertificateOrderResp.Body.CertificateOrderList) < int(listUserCertificateOrderLimit) {
			break
		} else {
			listUserCertificateOrderPage += 1
			if listUserCertificateOrderPage > 99 { // 避免死循环
				break
			}
		}
	}

	// 生成新证书名（需符合阿里云命名规则）
	var certId, certName string
	certName = fmt.Sprintf("certimate_%d", time.Now().UnixMilli())

	// 上传新证书
	// REF: https://help.aliyun.com/zh/ssl-certificate/developer-reference/api-cas-2020-04-07-uploadusercertificate
	uploadUserCertificateReq := &aliyunCas.UploadUserCertificateRequest{
		Name: tea.String(certName),
		Cert: tea.String(certPem),
		Key:  tea.String(privkeyPem),
	}
	uploadUserCertificateResp, err := u.sdkClient.UploadUserCertificate(uploadUserCertificateReq)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'cas.UploadUserCertificate'")
	}

	certId = fmt.Sprintf("%d", tea.Int64Value(uploadUserCertificateResp.Body.CertId))
	return &uploader.UploadResult{
		CertId:   certId,
		CertName: certName,
	}, nil
}

func createSdkClient(accessKeyId, accessKeySecret, region string) (*aliyunCas.Client, error) {
	if region == "" {
		region = "cn-hangzhou" // CAS 服务默认区域：华东一杭州
	}

	aConfig := &aliyunOpen.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
	}

	var endpoint string
	switch region {
	case "cn-hangzhou":
		endpoint = "cas.aliyuncs.com"
	default:
		endpoint = fmt.Sprintf("cas.%s.aliyuncs.com", region)
	}
	aConfig.Endpoint = tea.String(endpoint)

	client, err := aliyunCas.NewClient(aConfig)
	if err != nil {
		return nil, err
	}

	return client, nil
}
