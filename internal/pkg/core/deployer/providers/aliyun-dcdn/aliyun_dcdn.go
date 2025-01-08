package aliyundcdn

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	aliyunOpen "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	aliyunDcdn "github.com/alibabacloud-go/dcdn-20180115/v3/client"
	"github.com/alibabacloud-go/tea/tea"
	xerrors "github.com/pkg/errors"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/logger"
)

type AliyunDCDNDeployerConfig struct {
	// 阿里云 AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// 阿里云 AccessKeySecret。
	AccessKeySecret string `json:"accessKeySecret"`
	// 加速域名（支持泛域名）。
	Domain string `json:"domain"`
}

type AliyunDCDNDeployer struct {
	config    *AliyunDCDNDeployerConfig
	logger    logger.Logger
	sdkClient *aliyunDcdn.Client
}

var _ deployer.Deployer = (*AliyunDCDNDeployer)(nil)

func New(config *AliyunDCDNDeployerConfig) (*AliyunDCDNDeployer, error) {
	return NewWithLogger(config, logger.NewNilLogger())
}

func NewWithLogger(config *AliyunDCDNDeployerConfig, logger logger.Logger) (*AliyunDCDNDeployer, error) {
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

	return &AliyunDCDNDeployer{
		logger:    logger,
		config:    config,
		sdkClient: client,
	}, nil
}

func (d *AliyunDCDNDeployer) Deploy(ctx context.Context, certPem string, privkeyPem string) (*deployer.DeployResult, error) {
	// "*.example.com" → ".example.com"，适配阿里云 DCDN 要求的泛域名格式
	domain := strings.TrimPrefix(d.config.Domain, "*")

	// 配置域名证书
	// REF: https://help.aliyun.com/zh/edge-security-acceleration/dcdn/developer-reference/api-dcdn-2018-01-15-setdcdndomainsslcertificate
	setDcdnDomainSSLCertificateReq := &aliyunDcdn.SetDcdnDomainSSLCertificateRequest{
		DomainName:  tea.String(domain),
		CertName:    tea.String(fmt.Sprintf("certimate-%d", time.Now().UnixMilli())),
		CertType:    tea.String("upload"),
		SSLProtocol: tea.String("on"),
		SSLPub:      tea.String(certPem),
		SSLPri:      tea.String(privkeyPem),
	}
	setDcdnDomainSSLCertificateResp, err := d.sdkClient.SetDcdnDomainSSLCertificate(setDcdnDomainSSLCertificateReq)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'dcdn.SetDcdnDomainSSLCertificate'")
	}

	d.logger.Logt("已配置 DCDN 域名证书", setDcdnDomainSSLCertificateResp)

	return &deployer.DeployResult{}, nil
}

func createSdkClient(accessKeyId, accessKeySecret string) (*aliyunDcdn.Client, error) {
	config := &aliyunOpen.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
		Endpoint:        tea.String("dcdn.aliyuncs.com"),
	}

	client, err := aliyunDcdn.NewClient(config)
	if err != nil {
		return nil, err
	}

	return client, nil
}
