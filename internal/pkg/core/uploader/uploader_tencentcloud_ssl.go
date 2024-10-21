package uploader

import (
	"context"
	"fmt"
	"time"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tcSsl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"

	"github.com/usual2970/certimate/internal/pkg/utils/cast"
)

type TencentCloudSSLUploaderConfig struct {
	Region    string `json:"region"`
	SecretId  string `json:"secretId"`
	SecretKey string `json:"secretKey"`
}

type TencentCloudSSLUploader struct {
	config    *TencentCloudSSLUploaderConfig
	sdkClient *tcSsl.Client
}

func NewTencentCloudSSLUploader(config *TencentCloudSSLUploaderConfig) (*TencentCloudSSLUploader, error) {
	client, err := (&TencentCloudSSLUploader{config: config}).createSdkClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create sdk client: %w", err)
	}

	return &TencentCloudSSLUploader{
		config:    config,
		sdkClient: client,
	}, nil
}

func (u *TencentCloudSSLUploader) Upload(ctx context.Context, certPem string, privkeyPem string) (res *UploadResult, err error) {
	// 生成新证书名（需符合腾讯云命名规则）
	var certId, certName string
	certName = fmt.Sprintf("certimate-%d", time.Now().UnixMilli())

	// 上传新证书
	// REF: https://cloud.tencent.com/document/product/400/41665
	uploadCertificateReq := &tcSsl.UploadCertificateRequest{
		Alias:                 cast.StringPtr(certName),
		CertificatePublicKey:  cast.StringPtr(certPem),
		CertificatePrivateKey: cast.StringPtr(privkeyPem),
		Repeatable:            cast.BoolPtr(false),
	}
	uploadCertificateResp, err := u.sdkClient.UploadCertificate(uploadCertificateReq)
	if uploadCertificateResp != nil && uploadCertificateResp.Response != nil && uploadCertificateResp.Response.RepeatCertId != nil {
		// 获取重复证书
		// REF: https://cloud.tencent.com/document/api/400/41673
		describeCertificateDetailReq := &tcSsl.DescribeCertificateDetailRequest{
			CertificateId: uploadCertificateResp.Response.RepeatCertId,
		}
		describeCertificateDetailResp, err := u.sdkClient.DescribeCertificateDetail(describeCertificateDetailReq)
		if err != nil {
			return nil, fmt.Errorf("failed to execute sdk request 'ssl.DescribeCertificateDetail': %w", err)
		}

		certId = *uploadCertificateResp.Response.RepeatCertId
		certName = *describeCertificateDetailResp.Response.Alias
		return &UploadResult{
			CertId:   certId,
			CertName: certName,
		}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'ssl.UploadCertificate': %w", err)
	}

	certId = *uploadCertificateResp.Response.CertificateId
	return &UploadResult{
		CertId:   certId,
		CertName: certName,
	}, nil
}

func (u *TencentCloudSSLUploader) createSdkClient() (*tcSsl.Client, error) {
	region := u.config.Region
	secretId := u.config.SecretId
	secretKey := u.config.SecretKey
	if region == "" {
		region = "ap-guangzhou" // SSL 服务默认区域：广州
	}

	credential := common.NewCredential(secretId, secretKey)
	client, err := tcSsl.NewClient(credential, region, profile.NewClientProfile())
	if err != nil {
		return nil, err
	}

	return client, nil
}
