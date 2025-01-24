package awsacm

import (
	"context"
	"errors"
	"fmt"
	"time"

	aws "github.com/aws/aws-sdk-go-v2/aws"
	awsCfg "github.com/aws/aws-sdk-go-v2/config"
	awsCred "github.com/aws/aws-sdk-go-v2/credentials"
	awsAcm "github.com/aws/aws-sdk-go-v2/service/acm"
	xerrors "github.com/pkg/errors"

	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	"github.com/usual2970/certimate/internal/pkg/utils/certs"
)

type AWSCertificateManagerUploaderConfig struct {
	// AWS AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// AWS SecretAccessKey。
	SecretAccessKey string `json:"secretAccessKey"`
	// AWS 区域。
	Region string `json:"region"`
}

type AWSCertificateManagerUploader struct {
	config    *AWSCertificateManagerUploaderConfig
	sdkClient *awsAcm.Client
}

var _ uploader.Uploader = (*AWSCertificateManagerUploader)(nil)

func New(config *AWSCertificateManagerUploaderConfig) (*AWSCertificateManagerUploader, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}

	client, err := createSdkClient(config.AccessKeyId, config.SecretAccessKey, config.Region)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk client")
	}

	return &AWSCertificateManagerUploader{
		config:    config,
		sdkClient: client,
	}, nil
}

func (u *AWSCertificateManagerUploader) Upload(ctx context.Context, certPem string, privkeyPem string) (res *uploader.UploadResult, err error) {
	// 解析证书内容
	certX509, err := certs.ParseCertificateFromPEM(certPem)
	if err != nil {
		return nil, err
	}

	// 生成 AWS 所需的服务端证书和证书链参数
	scertPem, _ := certs.ConvertCertificateToPEM(certX509)
	bcertPem := certPem

	// 生成新证书名（需符合 AWS 命名规则）
	var certId, certName string
	certName = fmt.Sprintf("certimate_%d", time.Now().UnixMilli())

	// 导入证书
	// REF: https://docs.aws.amazon.com/en_us/acm/latest/APIReference/API_ImportCertificate.html
	importCertificateReq := &awsAcm.ImportCertificateInput{
		Certificate:      ([]byte)(scertPem),
		CertificateChain: ([]byte)(bcertPem),
		PrivateKey:       ([]byte)(privkeyPem),
	}
	importCertificateResp, err := u.sdkClient.ImportCertificate(context.TODO(), importCertificateReq)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'acm.ImportCertificate'")
	}

	certId = *importCertificateResp.CertificateArn
	return &uploader.UploadResult{
		CertId:   certId,
		CertName: certName,
	}, nil
}

func createSdkClient(accessKeyId, secretAccessKey, region string) (*awsAcm.Client, error) {
	cfg, err := awsCfg.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	client := awsAcm.NewFromConfig(cfg, func(o *awsAcm.Options) {
		o.Region = region
		o.Credentials = aws.NewCredentialsCache(awsCred.NewStaticCredentialsProvider(accessKeyId, secretAccessKey, ""))
	})
	return client, nil
}
