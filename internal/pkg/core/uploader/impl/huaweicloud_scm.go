package impl

import (
	"context"
	"fmt"
	"time"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	scm "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/scm/v3"
	scmModel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/scm/v3/model"
	scmRegion "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/scm/v3/region"

	"certimate/internal/pkg/core/uploader"
	"certimate/internal/pkg/utils/x509"
)

type HuaweiCloudSCMUploaderConfig struct {
	Region          string `json:"region"`
	AccessKeyId     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
}

type HuaweiCloudSCMUploader struct {
	client *scm.ScmClient
}

func NewHuaweiCloudSCMUploader(config *HuaweiCloudSCMUploaderConfig) (*HuaweiCloudSCMUploader, error) {
	client, err := createClient(config.Region, config.AccessKeyId, config.SecretAccessKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	return &HuaweiCloudSCMUploader{
		client: client,
	}, nil
}

func (u *HuaweiCloudSCMUploader) Upload(ctx context.Context, certPem string, privkeyPem string) (res *uploader.UploadResult, err error) {
	// 解析证书内容
	newCert, err := x509.ParseCertificateFromPEM(certPem)
	if err != nil {
		return nil, err
	}

	// 遍历查询已有证书，避免重复上传
	// REF: https://support.huaweicloud.com/api-ccm/ListCertificates.html
	// REF: https://support.huaweicloud.com/api-ccm/ExportCertificate_0.html
	listCertificatesLimit := int32(50)
	listCertificatesOffset := int32(0)
	for {
		listCertificatesReq := &scmModel.ListCertificatesRequest{
<<<<<<< HEAD
			Limit:   int32Ptr(listCertificatesLimit),
			Offset:  int32Ptr(listCertificatesOffset),
			SortDir: stringPtr("DESC"),
			SortKey: stringPtr("certExpiredTime"),
=======
			Limit:  int32Ptr(listCertificatesLimit),
			Offset: int32Ptr(listCertificatesOffset),
>>>>>>> 1ff10bf989afaa505a3fc2fda668ffcded815d09
		}
		listCertificatesResp, err := u.client.ListCertificates(listCertificatesReq)
		if err != nil {
			return nil, fmt.Errorf("failed to execute request 'scm.ListCertificates': %w", err)
		}

		if listCertificatesResp.Certificates != nil {
			for _, certDetail := range *listCertificatesResp.Certificates {
				exportCertificateReq := &scmModel.ExportCertificateRequest{
					CertificateId: certDetail.Id,
				}
				exportCertificateResp, err := u.client.ExportCertificate(exportCertificateReq)
				if err != nil {
					if exportCertificateResp != nil && exportCertificateResp.HttpStatusCode == 404 {
						continue
					}
					return nil, fmt.Errorf("failed to execute request 'scm.ExportCertificate': %w", err)
				}

<<<<<<< HEAD
				var isSameCert bool
				if *exportCertificateResp.Certificate == certPem {
					isSameCert = true
				} else {
					cert, err := x509.ParseCertificateFromPEM(*exportCertificateResp.Certificate)
					if err != nil {
						continue
					}

					isSameCert = x509.EqualCertificate(cert, newCert)
				}

				// 如果已存在相同证书，直接返回已有的证书信息
				if isSameCert {
=======
				cert, err := x509.ParseCertificateFromPEM(*exportCertificateResp.Certificate)
				if err != nil {
					continue
				}

				if x509.EqualCertificate(cert, newCert) {
					// 如果已存在相同证书，直接返回已有的证书信息
>>>>>>> 1ff10bf989afaa505a3fc2fda668ffcded815d09
					return &uploader.UploadResult{
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
<<<<<<< HEAD
		if listCertificatesOffset >= 999 { // 避免无限获取
			break
		}
=======
>>>>>>> 1ff10bf989afaa505a3fc2fda668ffcded815d09
	}

	// 生成证书名（需符合华为云命名规则）
	var certId, certName string
	certName = fmt.Sprintf("certimate-%d", time.Now().UnixMilli())

	// 上传新证书
	// REF: https://support.huaweicloud.com/api-ccm/ImportCertificate.html
	importCertificateReq := &scmModel.ImportCertificateRequest{
		Body: &scmModel.ImportCertificateRequestBody{
			Name:        certName,
			Certificate: certPem,
			PrivateKey:  privkeyPem,
		},
	}
	importCertificateResp, err := u.client.ImportCertificate(importCertificateReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request 'scm.ImportCertificate': %w", err)
	}

	certId = *importCertificateResp.CertificateId
	return &uploader.UploadResult{
		CertId:   certId,
		CertName: certName,
	}, nil
}

func createClient(region, accessKeyId, secretAccessKey string) (*scm.ScmClient, error) {
	auth, err := basic.NewCredentialsBuilder().
		WithAk(accessKeyId).
		WithSk(secretAccessKey).
		SafeBuild()
	if err != nil {
		return nil, err
	}

	if region == "" {
		region = "cn-north-4" // SCM 服务默认区域：华北北京四
	}

	hcRegion, err := scmRegion.SafeValueOf(region)
	if err != nil {
		return nil, err
	}

	hcClient, err := scm.ScmClientBuilder().
		WithRegion(hcRegion).
		WithCredential(auth).
		SafeBuild()
	if err != nil {
		return nil, err
	}

	client := scm.NewScmClient(hcClient)
	return client, nil
}

func int32Ptr(i int32) *int32 {
	return &i
}
<<<<<<< HEAD

func stringPtr(s string) *string {
	return &s
}
=======
>>>>>>> 1ff10bf989afaa505a3fc2fda668ffcded815d09
