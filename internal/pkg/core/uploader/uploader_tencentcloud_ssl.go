package uploader

import (
	"context"

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
	// 上传新证书
	// REF: https://cloud.tencent.com/document/product/400/41665
	uploadCertificateReq := tcSsl.NewUploadCertificateRequest()
	uploadCertificateReq.CertificatePublicKey = common.StringPtr(certPem)
	uploadCertificateReq.CertificatePrivateKey = common.StringPtr(privkeyPem)
	uploadCertificateReq.Repeatable = common.BoolPtr(false)
	uploadCertificateResp, err := u.sdkClient.UploadCertificate(uploadCertificateReq)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'ssl.UploadCertificate'")
	}

	certId := *describeCertificateDetailResp.Response.CertificateId
	return &UploadResult{
		CertId:   certId,
		CertName: "",
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
