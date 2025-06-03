package awsiam

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	aws "github.com/aws/aws-sdk-go-v2/aws"
	awscfg "github.com/aws/aws-sdk-go-v2/config"
	awscred "github.com/aws/aws-sdk-go-v2/credentials"
	awsiam "github.com/aws/aws-sdk-go-v2/service/iam"

	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	certutil "github.com/usual2970/certimate/internal/pkg/utils/cert"
)

type UploaderConfig struct {
	// AWS AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// AWS SecretAccessKey。
	SecretAccessKey string `json:"secretAccessKey"`
	// AWS 区域。
	Region string `json:"region"`
	// IAM 证书路径。
	// 选填。
	CertificatePath string `json:"certificatePath,omitempty"`
}

type UploaderProvider struct {
	config    *UploaderConfig
	logger    *slog.Logger
	sdkClient *awsiam.Client
}

var _ uploader.Uploader = (*UploaderProvider)(nil)

func NewUploader(config *UploaderConfig) (*UploaderProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	client, err := createSdkClient(config.AccessKeyId, config.SecretAccessKey, config.Region)
	if err != nil {
		return nil, fmt.Errorf("failed to create sdk client: %w", err)
	}

	return &UploaderProvider{
		config:    config,
		logger:    slog.Default(),
		sdkClient: client,
	}, nil
}

func (u *UploaderProvider) WithLogger(logger *slog.Logger) uploader.Uploader {
	if logger == nil {
		u.logger = slog.New(slog.DiscardHandler)
	} else {
		u.logger = logger
	}
	return u
}

func (u *UploaderProvider) Upload(ctx context.Context, certPEM string, privkeyPEM string) (*uploader.UploadResult, error) {
	// 解析证书内容
	certX509, err := certutil.ParseCertificateFromPEM(certPEM)
	if err != nil {
		return nil, err
	}

	// 提取服务器证书
	serverCertPEM, intermediaCertPEM, err := certutil.ExtractCertificatesFromPEM(certPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to extract certs: %w", err)
	}

	// 获取证书列表，避免重复上传
	// REF: https://docs.aws.amazon.com/en_us/IAM/latest/APIReference/API_ListServerCertificates.html
	var listServerCertificatesMarker *string = nil
	var listServerCertificatesMaxItems int32 = 1000
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		listServerCertificatesReq := &awsiam.ListServerCertificatesInput{
			Marker:   listServerCertificatesMarker,
			MaxItems: aws.Int32(listServerCertificatesMaxItems),
		}
		if u.config.CertificatePath != "" {
			listServerCertificatesReq.PathPrefix = aws.String(u.config.CertificatePath)
		}
		listServerCertificatesResp, err := u.sdkClient.ListServerCertificates(context.TODO(), listServerCertificatesReq)
		u.logger.Debug("sdk request 'iam.ListServerCertificates'", slog.Any("request", listServerCertificatesReq), slog.Any("response", listServerCertificatesResp))
		if err != nil {
			return nil, fmt.Errorf("failed to execute sdk request 'iam.ListServerCertificates': %w", err)
		}

		for _, certMeta := range listServerCertificatesResp.ServerCertificateMetadataList {
			// 先对比证书路径
			if u.config.CertificatePath != "" && aws.ToString(certMeta.Path) != u.config.CertificatePath {
				continue
			}

			// 先对比证书有效期
			if certMeta.Expiration == nil || !certMeta.Expiration.Equal(certX509.NotAfter) {
				continue
			}

			// 最后对比证书内容
			// REF: https://docs.aws.amazon.com/en_us/IAM/latest/APIReference/API_GetServerCertificate.html
			getServerCertificateReq := &awsiam.GetServerCertificateInput{
				ServerCertificateName: certMeta.ServerCertificateName,
			}
			getServerCertificateResp, err := u.sdkClient.GetServerCertificate(context.TODO(), getServerCertificateReq)
			if err != nil {
				return nil, fmt.Errorf("failed to execute sdk request 'iam.GetServerCertificate': %w", err)
			} else {
				oldCertPEM := aws.ToString(getServerCertificateResp.ServerCertificate.CertificateBody)
				oldCertX509, err := certutil.ParseCertificateFromPEM(oldCertPEM)
				if err != nil {
					continue
				}

				if !certutil.EqualCertificate(certX509, oldCertX509) {
					continue
				}
			}

			// 如果以上信息都一致，则视为已存在相同证书，直接返回
			u.logger.Info("ssl certificate already exists")
			return &uploader.UploadResult{
				CertId:   aws.ToString(certMeta.ServerCertificateId),
				CertName: aws.ToString(certMeta.ServerCertificateName),
			}, nil
		}

		if listServerCertificatesResp.Marker == nil || len(listServerCertificatesResp.ServerCertificateMetadataList) < int(listServerCertificatesMaxItems) {
			break
		} else {
			listServerCertificatesMarker = listServerCertificatesResp.Marker
		}
	}

	// 生成新证书名（需符合 AWS IAM 命名规则）
	certName := fmt.Sprintf("certimate-%d", time.Now().UnixMilli())

	// 导入证书
	// REF: https://docs.aws.amazon.com/en_us/IAM/latest/APIReference/API_UploadServerCertificate.html
	uploadServerCertificateReq := &awsiam.UploadServerCertificateInput{
		ServerCertificateName: aws.String(certName),
		Path:                  aws.String(u.config.CertificatePath),
		CertificateBody:       aws.String(serverCertPEM),
		CertificateChain:      aws.String(intermediaCertPEM),
		PrivateKey:            aws.String(privkeyPEM),
	}
	if u.config.CertificatePath == "" {
		uploadServerCertificateReq.Path = aws.String("/")
	}
	uploadServerCertificateResp, err := u.sdkClient.UploadServerCertificate(context.TODO(), uploadServerCertificateReq)
	u.logger.Debug("sdk request 'iam.UploadServerCertificate'", slog.Any("request", uploadServerCertificateReq), slog.Any("response", uploadServerCertificateResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'iam.UploadServerCertificate': %w", err)
	}

	return &uploader.UploadResult{
		CertId:   aws.ToString(uploadServerCertificateResp.ServerCertificateMetadata.ServerCertificateId),
		CertName: certName,
	}, nil
}

func createSdkClient(accessKeyId, secretAccessKey, region string) (*awsiam.Client, error) {
	cfg, err := awscfg.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	client := awsiam.NewFromConfig(cfg, func(o *awsiam.Options) {
		o.Region = region
		o.Credentials = aws.NewCredentialsCache(awscred.NewStaticCredentialsProvider(accessKeyId, secretAccessKey, ""))
	})
	return client, nil
}
