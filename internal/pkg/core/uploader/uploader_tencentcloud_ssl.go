package uploader

import (
	"context"
	"fmt"
	"time"

	xerrors "github.com/pkg/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tcSsl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"
)

type TencentCloudSSLUploaderConfig struct {
	SecretId  string `json:"secretId"`
	SecretKey string `json:"secretKey"`
}

type TencentCloudSSLUploader struct {
	config    *TencentCloudSSLUploaderConfig
	sdkClient *tcSsl.Client
}

func NewTencentCloudSSLUploader(config *TencentCloudSSLUploaderConfig) (Uploader, error) {
	client, err := (&TencentCloudSSLUploader{}).createSdkClient(
		config.SecretId,
		config.SecretKey,
	)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk client")
	}

	return &TencentCloudSSLUploader{
		config:    config,
		sdkClient: client,
	}, nil
}

func (u *TencentCloudSSLUploader) Upload(ctx context.Context, certPem string, privkeyPem string) (res *UploadResult, err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic: %+v", r)
			fmt.Println()
		}
	}()

	// 生成新证书名（需符合腾讯云命名规则）
	var certId, certName string
	certName = fmt.Sprintf("certimate-%d", time.Now().UnixMilli())

	// 上传新证书
	// REF: https://cloud.tencent.com/document/product/400/41665
	uploadCertificateReq := tcSsl.NewUploadCertificateRequest()
	uploadCertificateReq.Alias = common.StringPtr(certName)
	uploadCertificateReq.CertificatePublicKey = common.StringPtr(certPem)
	uploadCertificateReq.CertificatePrivateKey = common.StringPtr(privkeyPem)
	uploadCertificateReq.Repeatable = common.BoolPtr(false)
	uploadCertificateResp, err := u.sdkClient.UploadCertificate(uploadCertificateReq)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'ssl.UploadCertificate'")
	}

	// 获取证书详情
	// REF: https://cloud.tencent.com/document/api/400/41673
	//
	// P.S. 上传重复证书会返回上一次的证书 ID，这里需要重新获取一遍证书名（https://github.com/usual2970/certimate/pull/227）
	describeCertificateDetailReq := tcSsl.NewDescribeCertificateDetailRequest()
	describeCertificateDetailReq.CertificateId = uploadCertificateResp.Response.CertificateId
	describeCertificateDetailResp, err := u.sdkClient.DescribeCertificateDetail(describeCertificateDetailReq)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'ssl.DescribeCertificateDetail'")
	}

	certId = *describeCertificateDetailResp.Response.CertificateId
	certName = *describeCertificateDetailResp.Response.Alias
	return &UploadResult{
		CertId:   certId,
		CertName: certName,
	}, nil
}

func (u *TencentCloudSSLUploader) createSdkClient(secretId, secretKey string) (*tcSsl.Client, error) {
	credential := common.NewCredential(secretId, secretKey)
	client, err := tcSsl.NewClient(credential, "", profile.NewClientProfile())
	if err != nil {
		return nil, err
	}

	return client, nil
}
