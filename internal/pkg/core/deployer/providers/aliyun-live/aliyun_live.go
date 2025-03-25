package aliyunlive

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	aliopen "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	alilive "github.com/alibabacloud-go/live-20161101/client"
	"github.com/alibabacloud-go/tea/tea"
	xerrors "github.com/pkg/errors"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
)

type DeployerConfig struct {
	// 阿里云 AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// 阿里云 AccessKeySecret。
	AccessKeySecret string `json:"accessKeySecret"`
	// 阿里云地域。
	Region string `json:"region"`
	// 直播流域名（支持泛域名）。
	Domain string `json:"domain"`
}

type DeployerProvider struct {
	config    *DeployerConfig
	logger    *slog.Logger
	sdkClient *alilive.Client
}

var _ deployer.Deployer = (*DeployerProvider)(nil)

func NewDeployer(config *DeployerConfig) (*DeployerProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	client, err := createSdkClient(config.AccessKeyId, config.AccessKeySecret, config.Region)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk client")
	}

	return &DeployerProvider{
		config:    config,
		logger:    slog.Default(),
		sdkClient: client,
	}, nil
}

func (d *DeployerProvider) WithLogger(logger *slog.Logger) deployer.Deployer {
	if logger == nil {
		d.logger = slog.Default()
	} else {
		d.logger = logger
	}
	return d
}

func (d *DeployerProvider) Deploy(ctx context.Context, certPem string, privkeyPem string) (*deployer.DeployResult, error) {
	// "*.example.com" → ".example.com"，适配阿里云 Live 要求的泛域名格式
	domain := strings.TrimPrefix(d.config.Domain, "*")

	// 设置域名证书
	// REF: https://help.aliyun.com/zh/live/developer-reference/api-live-2016-11-01-setlivedomaincertificate
	setLiveDomainSSLCertificateReq := &alilive.SetLiveDomainCertificateRequest{
		DomainName:  tea.String(domain),
		CertName:    tea.String(fmt.Sprintf("certimate-%d", time.Now().UnixMilli())),
		CertType:    tea.String("upload"),
		SSLProtocol: tea.String("on"),
		SSLPub:      tea.String(certPem),
		SSLPri:      tea.String(privkeyPem),
	}
	setLiveDomainSSLCertificateResp, err := d.sdkClient.SetLiveDomainCertificate(setLiveDomainSSLCertificateReq)
	d.logger.Debug("sdk request 'live.SetLiveDomainCertificate'", slog.Any("request", setLiveDomainSSLCertificateReq), slog.Any("response", setLiveDomainSSLCertificateResp))
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'live.SetLiveDomainCertificate'")
	}

	return &deployer.DeployResult{}, nil
}

func createSdkClient(accessKeyId, accessKeySecret, region string) (*alilive.Client, error) {
	// 接入点一览 https://api.aliyun.com/product/live
	var endpoint string
	switch region {
	case
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
