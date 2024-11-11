package dogecloud

import (
	"context"
	"errors"
	"fmt"
	"time"

	xerrors "github.com/pkg/errors"

	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	doge "github.com/usual2970/certimate/internal/pkg/vendors/dogecloud-sdk"
)

type DogeCloudUploaderConfig struct {
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
}

type DogeCloudUploader struct {
	config    *DogeCloudUploaderConfig
	sdkClient *doge.Client
}

var _ uploader.Uploader = (*DogeCloudUploader)(nil)

func New(config *DogeCloudUploaderConfig) (*DogeCloudUploader, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}

	client, err := createSdkClient(
		config.AccessKey,
		config.SecretKey,
	)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk client")
	}

	return &DogeCloudUploader{
		config:    config,
		sdkClient: client,
	}, nil
}

func (u *DogeCloudUploader) Upload(ctx context.Context, certPem string, privkeyPem string) (res *uploader.UploadResult, err error) {
	// 生成新证书名（需符合多吉云命名规则）
	var certId, certName string
	certName = fmt.Sprintf("certimate-%d", time.Now().UnixMilli())

	// 上传新证书
	// REF: https://docs.dogecloud.com/cdn/api-cert-upload
	uploadSslCertResp, err := u.sdkClient.UploadCdnCert(certName, certPem, privkeyPem)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'cdn.UploadCdnCert'")
	}

	certId = fmt.Sprintf("%d", uploadSslCertResp.Data.Id)
	return &uploader.UploadResult{
		CertId:   certId,
		CertName: certName,
	}, nil
}

func createSdkClient(accessKey, secretKey string) (*doge.Client, error) {
	client := doge.NewClient(accessKey, secretKey)
	return client, nil
}
