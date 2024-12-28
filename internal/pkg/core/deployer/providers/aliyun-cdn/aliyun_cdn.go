package aliyuncdn

import (
	"context"
	"errors"
	"fmt"
	"time"

	aliyunCdn "github.com/alibabacloud-go/cdn-20180510/v5/client"
	aliyunOpen "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	xerrors "github.com/pkg/errors"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/logger"
)

type AliyunCDNDeployerConfig struct {
	// 阿里云 AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// 阿里云 AccessKeySecret。
	AccessKeySecret string `json:"accessKeySecret"`
	// 加速域名（不支持泛域名）。
	Domain string `json:"domain"`
}

type AliyunCDNDeployer struct {
	config    *AliyunCDNDeployerConfig
	logger    logger.Logger
	sdkClient *aliyunCdn.Client
}

var _ deployer.Deployer = (*AliyunCDNDeployer)(nil)

func New(config *AliyunCDNDeployerConfig) (*AliyunCDNDeployer, error) {
	return NewWithLogger(config, logger.NewNilLogger())
}

func NewWithLogger(config *AliyunCDNDeployerConfig, logger logger.Logger) (*AliyunCDNDeployer, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}

	if logger == nil {
		return nil, errors.New("logger is nil")
	}

	client, err := createSdkClient(config.AccessKeyId, config.AccessKeySecret)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk client")
	}

	return &AliyunCDNDeployer{
		logger:    logger,
		config:    config,
		sdkClient: client,
	}, nil
}

func (d *AliyunCDNDeployer) Deploy(ctx context.Context, certPem string, privkeyPem string) (*deployer.DeployResult, error) {
	// 设置 CDN 域名域名证书
	// REF: https://help.aliyun.com/zh/cdn/developer-reference/api-cdn-2018-05-10-setcdndomainsslcertificate
	setCdnDomainSSLCertificateReq := &aliyunCdn.SetCdnDomainSSLCertificateRequest{
		DomainName:  tea.String(d.config.Domain),
		CertName:    tea.String(fmt.Sprintf("certimate-%d", time.Now().UnixMilli())),
		CertType:    tea.String("upload"),
		SSLProtocol: tea.String("on"),
		SSLPub:      tea.String(certPem),
		SSLPri:      tea.String(privkeyPem),
	}
	setCdnDomainSSLCertificateResp, err := d.sdkClient.SetCdnDomainSSLCertificate(setCdnDomainSSLCertificateReq)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'cdn.SetCdnDomainSSLCertificate'")
	}

	d.logger.Logt("已设置 CDN 域名证书", setCdnDomainSSLCertificateResp)

	return &deployer.DeployResult{}, nil
}

func createSdkClient(accessKeyId, accessKeySecret string) (*aliyunCdn.Client, error) {
	config := &aliyunOpen.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
		Endpoint:        tea.String("cdn.aliyuncs.com"),
	}

	client, err := aliyunCdn.NewClient(config)
	if err != nil {
		return nil, err
	}

	return client, nil
}
