package awsacm

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	aws "github.com/aws/aws-sdk-go-v2/aws"
	awscfg "github.com/aws/aws-sdk-go-v2/config"
	awscred "github.com/aws/aws-sdk-go-v2/credentials"
	awsacm "github.com/aws/aws-sdk-go-v2/service/acm"
	"golang.org/x/exp/slices"

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
}

type SSLManagerProvider struct {
	config    *SSLManagerProviderConfig
	logger    *slog.Logger
	sdkClient *awsacm.Client
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
	// REF: https://docs.aws.amazon.com/en_us/acm/latest/APIReference/API_ListCertificates.html
	var listCertificatesNextToken *string = nil
	var listCertificatesMaxItems int32 = 1000
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		listCertificatesReq := &awsacm.ListCertificatesInput{
			NextToken: listCertificatesNextToken,
			MaxItems:  aws.Int32(listCertificatesMaxItems),
		}
		listCertificatesResp, err := m.sdkClient.ListCertificates(context.TODO(), listCertificatesReq)
		m.logger.Debug("sdk request 'acm.ListCertificates'", slog.Any("request", listCertificatesReq), slog.Any("response", listCertificatesResp))
		if err != nil {
			return nil, fmt.Errorf("failed to execute sdk request 'acm.ListCertificates': %w", err)
		}

		for _, certSummary := range listCertificatesResp.CertificateSummaryList {
			// 先对比证书有效期
			if certSummary.NotBefore == nil || !certSummary.NotBefore.Equal(certX509.NotBefore) {
				continue
			}
			if certSummary.NotAfter == nil || !certSummary.NotAfter.Equal(certX509.NotAfter) {
				continue
			}

			// 再对比证书多域名
			if !slices.Equal(certX509.DNSNames, certSummary.SubjectAlternativeNameSummaries) {
				continue
			}

			// 最后对比证书内容
			// REF: https://docs.aws.amazon.com/en_us/acm/latest/APIReference/API_GetCertificate.html
			getCertificateReq := &awsacm.GetCertificateInput{
				CertificateArn: certSummary.CertificateArn,
			}
			getCertificateResp, err := m.sdkClient.GetCertificate(context.TODO(), getCertificateReq)
			if err != nil {
				return nil, fmt.Errorf("failed to execute sdk request 'acm.GetCertificate': %w", err)
			} else {
				oldCertPEM := aws.ToString(getCertificateResp.Certificate)
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
				CertId: *certSummary.CertificateArn,
			}, nil
		}

		if listCertificatesResp.NextToken == nil || len(listCertificatesResp.CertificateSummaryList) < int(listCertificatesMaxItems) {
			break
		} else {
			listCertificatesNextToken = listCertificatesResp.NextToken
		}
	}

	// 导入证书
	// REF: https://docs.aws.amazon.com/en_us/acm/latest/APIReference/API_ImportCertificate.html
	importCertificateReq := &awsacm.ImportCertificateInput{
		Certificate:      ([]byte)(serverCertPEM),
		CertificateChain: ([]byte)(intermediaCertPEM),
		PrivateKey:       ([]byte)(privkeyPEM),
	}
	importCertificateResp, err := m.sdkClient.ImportCertificate(context.TODO(), importCertificateReq)
	m.logger.Debug("sdk request 'acm.ImportCertificate'", slog.Any("request", importCertificateReq), slog.Any("response", importCertificateResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'acm.ImportCertificate': %w", err)
	}

	return &core.SSLManageUploadResult{
		CertId: aws.ToString(importCertificateResp.CertificateArn),
	}, nil
}

func createSDKClient(accessKeyId, secretAccessKey, region string) (*awsacm.Client, error) {
	cfg, err := awscfg.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	client := awsacm.NewFromConfig(cfg, func(o *awsacm.Options) {
		o.Region = region
		o.Credentials = aws.NewCredentialsCache(awscred.NewStaticCredentialsProvider(accessKeyId, secretAccessKey, ""))
	})
	return client, nil
}
