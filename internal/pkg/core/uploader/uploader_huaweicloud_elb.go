package uploader

import (
	"context"
	"fmt"
	"time"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	hcElb "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/elb/v3"
	hcElbModel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/elb/v3/model"
	hcElbRegion "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/elb/v3/region"

	"github.com/usual2970/certimate/internal/pkg/utils/cast"
	"github.com/usual2970/certimate/internal/pkg/utils/x509"
)

type HuaweiCloudELBUploaderConfig struct {
	Region          string `json:"region"`
	ProjectId       string `json:"projectId"`
	AccessKeyId     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
}

type HuaweiCloudELBUploader struct {
	config    *HuaweiCloudELBUploaderConfig
	sdkClient *hcElb.ElbClient
}

func NewHuaweiCloudELBUploader(config *HuaweiCloudELBUploaderConfig) (Uploader, error) {
	client, err := (&HuaweiCloudELBUploader{}).createSdkClient(
		config.Region,
		config.AccessKeyId,
		config.SecretAccessKey,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create sdk client: %w", err)
	}

	return &HuaweiCloudELBUploader{
		config:    config,
		sdkClient: client,
	}, nil
}

func (u *HuaweiCloudELBUploader) Upload(ctx context.Context, certPem string, privkeyPem string) (res *UploadResult, err error) {
	// 解析证书内容
	newCert, err := x509.ParseCertificateFromPEM(certPem)
	if err != nil {
		return nil, err
	}

	// 遍历查询已有证书，避免重复上传
	// REF: https://support.huaweicloud.com/api-elb/ListCertificates.html
	listCertificatesPage := 1
	listCertificatesLimit := int32(2000)
	var listCertificatesMarker *string = nil
	for {
		listCertificatesReq := &hcElbModel.ListCertificatesRequest{
			Limit:  cast.Int32Ptr(listCertificatesLimit),
			Marker: listCertificatesMarker,
			Type:   &[]string{"server"},
		}
		listCertificatesResp, err := u.sdkClient.ListCertificates(listCertificatesReq)
		if err != nil {
			return nil, fmt.Errorf("failed to execute sdk request 'elb.ListCertificates': %w", err)
		}

		if listCertificatesResp.Certificates != nil {
			for _, certDetail := range *listCertificatesResp.Certificates {
				var isSameCert bool
				if certDetail.Certificate == certPem {
					isSameCert = true
				} else {
					cert, err := x509.ParseCertificateFromPEM(certDetail.Certificate)
					if err != nil {
						continue
					}

					isSameCert = x509.EqualCertificate(cert, newCert)
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
		} else {
			listCertificatesMarker = listCertificatesResp.PageInfo.NextMarker
			listCertificatesPage++
			if listCertificatesPage >= 9 { // 避免死循环
				break
			}
		}
	}

	// 生成新证书名（需符合华为云命名规则）
	var certId, certName string
	certName = fmt.Sprintf("certimate-%d", time.Now().UnixMilli())

	// 创建新证书
	// REF: https://support.huaweicloud.com/api-elb/CreateCertificate.html
	createCertificateReq := &hcElbModel.CreateCertificateRequest{
		Body: &hcElbModel.CreateCertificateRequestBody{
			Certificate: &hcElbModel.CreateCertificateOption{
				ProjectId:   cast.StringPtr(u.config.ProjectId),
				Name:        cast.StringPtr(certName),
				Certificate: cast.StringPtr(certPem),
				PrivateKey:  cast.StringPtr(privkeyPem),
			},
		},
	}
	createCertificateResp, err := u.sdkClient.CreateCertificate(createCertificateReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'elb.CreateCertificate': %w", err)
	}

	certId = createCertificateResp.Certificate.Id
	certName = createCertificateResp.Certificate.Name
	return &UploadResult{
		CertId:   certId,
		CertName: certName,
	}, nil
}

func (u *HuaweiCloudELBUploader) createSdkClient(region, accessKeyId, secretAccessKey string) (*hcElb.ElbClient, error) {
	if region == "" {
		region = "cn-north-4" // ELB 服务默认区域：华北四北京
	}

	auth, err := basic.NewCredentialsBuilder().
		WithAk(accessKeyId).
		WithSk(secretAccessKey).
		SafeBuild()
	if err != nil {
		return nil, err
	}

	hcRegion, err := hcElbRegion.SafeValueOf(region)
	if err != nil {
		return nil, err
	}

	hcClient, err := hcElb.ElbClientBuilder().
		WithRegion(hcRegion).
		WithCredential(auth).
		SafeBuild()
	if err != nil {
		return nil, err
	}

	client := hcElb.NewElbClient(hcClient)
	return client, nil
}
