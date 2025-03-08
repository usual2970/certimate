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
	sdkClient *aliyunCas.Client
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
						oldCertX509, err := certs.ParseCertificateFromPEM(*getUserCertificateDetailResp.Body.Cert)
						if err != nil {
							continue
						}

						isSameCert = certs.EqualCertificate(certX509, oldCertX509)
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
			listUserCertificateOrderPage++
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

	// 接入点一览 https://api.aliyun.com/product/cas
	var endpoint string
	switch region {
	case "cn-hangzhou":
		endpoint = "cas.aliyuncs.com"
	default:
		endpoint = fmt.Sprintf("cas.%s.aliyuncs.com", region)
	}

	config := &aliyunOpen.Config{
		Endpoint:        tea.String(endpoint),
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
	}

	client, err := aliyunCas.NewClient(config)
	if err != nil {
		return nil, err
	}

	return client, nil
}
