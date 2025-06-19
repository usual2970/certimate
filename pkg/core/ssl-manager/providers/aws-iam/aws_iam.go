package awsiam

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	aws "github.com/aws/aws-sdk-go-v2/aws"
	awscfg "github.com/aws/aws-sdk-go-v2/config"
	awscred "github.com/aws/aws-sdk-go-v2/credentials"
	awsiam "github.com/aws/aws-sdk-go-v2/service/iam"

	"github.com/certimate-go/certimate/pkg/core"
	xcert "github.com/certimate-go/certimate/pkg/utils/cert"
)

type SSLManagerProviderConfig struct {
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

type SSLManagerProvider struct {
	config    *SSLManagerProviderConfig
	logger    *slog.Logger
	sdkClient *awsiam.Client
}

var _ core.SSLManager = (*SSLManagerProvider)(nil)

func NewSSLManagerProvider(config *SSLManagerProviderConfig) (*SSLManagerProvider, error) {
	if config == nil {
		return nil, errors.New("the configuration of the ssl manager provider is nil")
	}

	client, err := createSDKClient(config.AccessKeyId, config.SecretAccessKey, config.Region)
	if err != nil {
		return nil, fmt.Errorf("could not create sdk client: %w", err)
	}

	return &SSLManagerProvider{
		config:    config,
		logger:    slog.Default(),
		sdkClient: client,
	}, nil
}

func (m *SSLManagerProvider) SetLogger(logger *slog.Logger) {
	if logger == nil {
		m.logger = slog.New(slog.DiscardHandler)
	} else {
		m.logger = logger
	}
}

func (m *SSLManagerProvider) Upload(ctx context.Context, certPEM string, privkeyPEM string) (*core.SSLManageUploadResult, error) {
	// 解析证书内容
	certX509, err := xcert.ParseCertificateFromPEM(certPEM)
	if err != nil {
		return nil, err
	}

	// 提取服务器证书
	serverCertPEM, intermediaCertPEM, err := xcert.ExtractCertificatesFromPEM(certPEM)
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
		if m.config.CertificatePath != "" {
			listServerCertificatesReq.PathPrefix = aws.String(m.config.CertificatePath)
		}
		listServerCertificatesResp, err := m.sdkClient.ListServerCertificates(context.TODO(), listServerCertificatesReq)
		m.logger.Debug("sdk request 'iam.ListServerCertificates'", slog.Any("request", listServerCertificatesReq), slog.Any("response", listServerCertificatesResp))
		if err != nil {
			return nil, fmt.Errorf("failed to execute sdk request 'iam.ListServerCertificates': %w", err)
		}

		for _, certMeta := range listServerCertificatesResp.ServerCertificateMetadataList {
			// 先对比证书路径
			if m.config.CertificatePath != "" && aws.ToString(certMeta.Path) != m.config.CertificatePath {
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
			getServerCertificateResp, err := m.sdkClient.GetServerCertificate(context.TODO(), getServerCertificateReq)
			if err != nil {
				return nil, fmt.Errorf("failed to execute sdk request 'iam.GetServerCertificate': %w", err)
			} else {
				oldCertPEM := aws.ToString(getServerCertificateResp.ServerCertificate.CertificateBody)
				oldCertX509, err := xcert.ParseCertificateFromPEM(oldCertPEM)
				if err != nil {
					continue
				}

				if !xcert.EqualCertificate(certX509, oldCertX509) {
					continue
				}
			}

			// 如果以上信息都一致，则视为已存在相同证书，直接返回
			m.logger.Info("ssl certificate already exists")
			return &core.SSLManageUploadResult{
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
		Path:                  aws.String(m.config.CertificatePath),
		CertificateBody:       aws.String(serverCertPEM),
		CertificateChain:      aws.String(intermediaCertPEM),
		PrivateKey:            aws.String(privkeyPEM),
	}
	if m.config.CertificatePath == "" {
		uploadServerCertificateReq.Path = aws.String("/")
	}
	uploadServerCertificateResp, err := m.sdkClient.UploadServerCertificate(context.TODO(), uploadServerCertificateReq)
	m.logger.Debug("sdk request 'iam.UploadServerCertificate'", slog.Any("request", uploadServerCertificateReq), slog.Any("response", uploadServerCertificateResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'iam.UploadServerCertificate': %w", err)
	}

	return &core.SSLManageUploadResult{
		CertId:   aws.ToString(uploadServerCertificateResp.ServerCertificateMetadata.ServerCertificateId),
		CertName: certName,
	}, nil
}

func createSDKClient(accessKeyId, secretAccessKey, region string) (*awsiam.Client, error) {
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
