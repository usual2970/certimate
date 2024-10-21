package uploader

import (
	"context"
	"fmt"
	"time"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	hcScm "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/scm/v3"
	hcScmModel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/scm/v3/model"
	hcScmRegion "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/scm/v3/region"

	"github.com/usual2970/certimate/internal/pkg/utils/cast"
	"github.com/usual2970/certimate/internal/pkg/utils/x509"
)

type HuaweiCloudSCMUploaderConfig struct {
	Region          string `json:"region"`
	AccessKeyId     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
}

type HuaweiCloudSCMUploader struct {
	config    *HuaweiCloudSCMUploaderConfig
	sdkClient *hcScm.ScmClient
}

func NewHuaweiCloudSCMUploader(config *HuaweiCloudSCMUploaderConfig) (*HuaweiCloudSCMUploader, error) {
	client, err := (&HuaweiCloudSCMUploader{config: config}).createSdkClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create sdk client: %w", err)
	}

	return &HuaweiCloudSCMUploader{
		config:    config,
		sdkClient: client,
	}, nil
}

func (u *HuaweiCloudSCMUploader) Upload(ctx context.Context, certPem string, privkeyPem string) (res *UploadResult, err error) {
	// 解析证书内容
	certX509, err := x509.ParseCertificateFromPEM(certPem)
	if err != nil {
		return nil, err
	}

	// 遍历查询已有证书，避免重复上传
	// REF: https://support.huaweicloud.com/api-ccm/ListCertificates.html
	// REF: https://support.huaweicloud.com/api-ccm/ExportCertificate_0.html
	listCertificatesPage := 1
	listCertificatesLimit := int32(50)
	listCertificatesOffset := int32(0)
	for {
		listCertificatesReq := &hcScmModel.ListCertificatesRequest{
			Limit:   cast.Int32Ptr(listCertificatesLimit),
			Offset:  cast.Int32Ptr(listCertificatesOffset),
			SortDir: cast.StringPtr("DESC"),
			SortKey: cast.StringPtr("certExpiredTime"),
		}
		listCertificatesResp, err := u.sdkClient.ListCertificates(listCertificatesReq)
		if err != nil {
			return nil, fmt.Errorf("failed to execute sdk request 'scm.ListCertificates': %w", err)
		}

		if listCertificatesResp.Certificates != nil {
			for _, certDetail := range *listCertificatesResp.Certificates {
				exportCertificateReq := &hcScmModel.ExportCertificateRequest{
					CertificateId: certDetail.Id,
				}
				exportCertificateResp, err := u.sdkClient.ExportCertificate(exportCertificateReq)
				if err != nil {
					if exportCertificateResp != nil && exportCertificateResp.HttpStatusCode == 404 {
						continue
					}
					return nil, fmt.Errorf("failed to execute sdk request 'scm.ExportCertificate': %w", err)
				}

				var isSameCert bool
				if *exportCertificateResp.Certificate == certPem {
					isSameCert = true
				} else {
					cert, err := x509.ParseCertificateFromPEM(*exportCertificateResp.Certificate)
					if err != nil {
						continue
					}

					isSameCert = x509.EqualCertificate(certX509, cert)
				}

				// 如果已存在相同证书，直接返回已有的证书信息
				if isSameCert {
					return &UploadResult{
						CertId:   certDetail.Id,
						CertName: certDetail.Name,
					}, nil
				}
			}
		}

		if listCertificatesResp.Certificates == nil || len(*listCertificatesResp.Certificates) < int(listCertificatesLimit) {
			break
		}

		listCertificatesOffset += listCertificatesLimit
		listCertificatesPage += 1
		if listCertificatesPage > 99 { // 避免死循环
			break
		}
	}

	// 生成新证书名（需符合华为云命名规则）
	var certId, certName string
	certName = fmt.Sprintf("certimate-%d", time.Now().UnixMilli())

	// 上传新证书
	// REF: https://support.huaweicloud.com/api-ccm/ImportCertificate.html
	importCertificateReq := &hcScmModel.ImportCertificateRequest{
		Body: &hcScmModel.ImportCertificateRequestBody{
			Name:        certName,
			Certificate: certPem,
			PrivateKey:  privkeyPem,
		},
	}
	importCertificateResp, err := u.sdkClient.ImportCertificate(importCertificateReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'scm.ImportCertificate': %w", err)
	}

	certId = *importCertificateResp.CertificateId
	return &UploadResult{
		CertId:   certId,
		CertName: certName,
	}, nil
}

func (u *HuaweiCloudSCMUploader) createSdkClient() (*hcScm.ScmClient, error) {
	region := u.config.Region
	accessKeyId := u.config.AccessKeyId
	secretAccessKey := u.config.SecretAccessKey
	if region == "" {
		region = "cn-north-4" // SCM 服务默认区域：华北四北京
	}

	auth, err := basic.NewCredentialsBuilder().
		WithAk(accessKeyId).
		WithSk(secretAccessKey).
		SafeBuild()
	if err != nil {
		return nil, err
	}

	hcRegion, err := hcScmRegion.SafeValueOf(region)
	if err != nil {
		return nil, err
	}

	hcClient, err := hcScm.ScmClientBuilder().
		WithRegion(hcRegion).
		WithCredential(auth).
		SafeBuild()
	if err != nil {
		return nil, err
	}

	client := hcScm.NewScmClient(hcClient)
	return client, nil
}
