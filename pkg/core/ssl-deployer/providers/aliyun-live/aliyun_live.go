package aliyunlive

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	aliopen "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	alilive "github.com/alibabacloud-go/live-20161101/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/certimate-go/certimate/pkg/core"
)

type SSLDeployerProviderConfig struct {
	// 阿里云 AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// 阿里云 AccessKeySecret。
	AccessKeySecret string `json:"accessKeySecret"`
	// 阿里云资源组 ID。
	ResourceGroupId string `json:"resourceGroupId,omitempty"`
	// 阿里云地域。
	Region string `json:"region"`
	// 直播流域名（支持泛域名）。
	Domain string `json:"domain"`
}

type SSLDeployerProvider struct {
	config    *SSLDeployerProviderConfig
	logger    *slog.Logger
	sdkClient *alilive.Client
}

var _ core.SSLDeployer = (*SSLDeployerProvider)(nil)

func NewSSLDeployerProvider(config *SSLDeployerProviderConfig) (*SSLDeployerProvider, error) {
	if config == nil {
		return nil, errors.New("the configuration of the ssl deployer provider is nil")
	}

	client, err := createSDKClient(config.AccessKeyId, config.AccessKeySecret, config.Region)
	if err != nil {
		return nil, fmt.Errorf("could not create sdk client: %w", err)
	}

	return &SSLDeployerProvider{
		config:    config,
		logger:    slog.Default(),
		sdkClient: client,
	}, nil
}

func (d *SSLDeployerProvider) SetLogger(logger *slog.Logger) {
	if logger == nil {
		d.logger = slog.New(slog.DiscardHandler)
	} else {
		d.logger = logger
	}
}

func (d *SSLDeployerProvider) Deploy(ctx context.Context, certPEM string, privkeyPEM string) (*core.SSLDeployResult, error) {
	if d.config.Domain == "" {
		return nil, errors.New("config `domain` is required")
	}

	// "*.example.com" → ".example.com"，适配阿里云 Live 要求的泛域名格式
	domain := strings.TrimPrefix(d.config.Domain, "*")

	// 设置域名证书
	// REF: https://help.aliyun.com/zh/live/developer-reference/api-live-2016-11-01-setlivedomaincertificate
	setLiveDomainSSLCertificateReq := &alilive.SetLiveDomainCertificateRequest{
		DomainName:  tea.String(domain),
		CertName:    tea.String(fmt.Sprintf("certimate-%d", time.Now().UnixMilli())),
		CertType:    tea.String("upload"),
		SSLProtocol: tea.String("on"),
		SSLPub:      tea.String(certPEM),
		SSLPri:      tea.String(privkeyPEM),
	}
	setLiveDomainSSLCertificateResp, err := d.sdkClient.SetLiveDomainCertificate(setLiveDomainSSLCertificateReq)
	d.logger.Debug("sdk request 'live.SetLiveDomainCertificate'", slog.Any("request", setLiveDomainSSLCertificateReq), slog.Any("response", setLiveDomainSSLCertificateResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'live.SetLiveDomainCertificate': %w", err)
	}

	return &core.SSLDeployResult{}, nil
}

func createSDKClient(accessKeyId, accessKeySecret, region string) (*alilive.Client, error) {
	// 接入点一览 https://api.aliyun.com/product/live
	var endpoint string
	switch region {
	case "",
		"cn-qingdao",
		"cn-beijing",
		"cn-shanghai",
		"cn-shenzhen",
		"ap-northeast-1",
		"ap-southeast-5",
		"me-central-1":
		endpoint = "live.aliyuncs.com"
	default:
		endpoint = fmt.Sprintf("live.%s.aliyuncs.com", region)
	}

	config := &aliopen.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
		Endpoint:        tea.String(endpoint),
	}

	client, err := alilive.NewClient(config)
	if err != nil {
		return nil, err
	}

	return client, nil
}
