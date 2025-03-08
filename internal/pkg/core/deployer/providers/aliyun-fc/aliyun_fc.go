package aliyunfc

import (
	"context"
	"fmt"
	"time"

	aliyunOpen "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	aliyunFc3 "github.com/alibabacloud-go/fc-20230330/v4/client"
	aliyunFc2 "github.com/alibabacloud-go/fc-open-20210406/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	xerrors "github.com/pkg/errors"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/logger"
)

type DeployerConfig struct {
	// 阿里云 AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// 阿里云 AccessKeySecret。
	AccessKeySecret string `json:"accessKeySecret"`
	// 阿里云地域。
	Region string `json:"region"`
	// 服务版本。
	// 零值时默认为 "3.0"。
	ServiceVersion string `json:"serviceVersion"`
	// 自定义域名（不支持泛域名）。
	Domain string `json:"domain"`
}

type DeployerProvider struct {
	config     *DeployerConfig
	logger     logger.Logger
	sdkClients *wSdkClients
}

var _ deployer.Deployer = (*DeployerProvider)(nil)

type wSdkClients struct {
	fc2 *aliyunFc2.Client
	fc3 *aliyunFc3.Client
}

func NewDeployer(config *DeployerConfig) (*DeployerProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	clients, err := createSdkClients(config.AccessKeyId, config.AccessKeySecret, config.Region)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk clients")
	}

	return &DeployerProvider{
		config:     config,
		logger:     logger.NewNilLogger(),
		sdkClients: clients,
	}, nil
}

func (d *DeployerProvider) WithLogger(logger logger.Logger) *DeployerProvider {
	d.logger = logger
	return d
}

func (d *DeployerProvider) Deploy(ctx context.Context, certPem string, privkeyPem string) (*deployer.DeployResult, error) {
	switch d.config.ServiceVersion {
	case "", "3.0":
		if err := d.deployToFC3(ctx, certPem, privkeyPem); err != nil {
			return nil, err
		}

	case "2.0":
		if err := d.deployToFC2(ctx, certPem, privkeyPem); err != nil {
			return nil, err
		}

	default:
		return nil, xerrors.Errorf("unsupported service version: %s", d.config.ServiceVersion)
	}

	return &deployer.DeployResult{}, nil
}

func (d *DeployerProvider) deployToFC3(ctx context.Context, certPem string, privkeyPem string) error {
	// 获取自定义域名
	// REF: https://help.aliyun.com/zh/functioncompute/fc-3-0/developer-reference/api-fc-2023-03-30-getcustomdomain
	getCustomDomainResp, err := d.sdkClients.fc3.GetCustomDomain(tea.String(d.config.Domain))
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'fc.GetCustomDomain'")
	} else {
		d.logger.Logt("已获取自定义域名", getCustomDomainResp)
	}

	// 更新自定义域名
	// REF: https://help.aliyun.com/zh/functioncompute/fc-3-0/developer-reference/api-fc-2023-03-30-updatecustomdomain
	updateCustomDomainReq := &aliyunFc3.UpdateCustomDomainRequest{
		Body: &aliyunFc3.UpdateCustomDomainInput{
			CertConfig: &aliyunFc3.CertConfig{
				CertName:    tea.String(fmt.Sprintf("certimate-%d", time.Now().UnixMilli())),
				Certificate: tea.String(certPem),
				PrivateKey:  tea.String(privkeyPem),
			},
			Protocol:  getCustomDomainResp.Body.Protocol,
			TlsConfig: getCustomDomainResp.Body.TlsConfig,
		},
	}
	updateCustomDomainResp, err := d.sdkClients.fc3.UpdateCustomDomain(tea.String(d.config.Domain), updateCustomDomainReq)
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'fc.UpdateCustomDomain'")
	} else {
		d.logger.Logt("已更新自定义域名", updateCustomDomainResp)
	}

	return nil
}

func (d *DeployerProvider) deployToFC2(ctx context.Context, certPem string, privkeyPem string) error {
	// 获取自定义域名
	// REF: https://help.aliyun.com/zh/functioncompute/fc-2-0/developer-reference/api-fc-open-2021-04-06-getcustomdomain
	getCustomDomainResp, err := d.sdkClients.fc2.GetCustomDomain(tea.String(d.config.Domain))
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'fc.GetCustomDomain'")
	} else {
		d.logger.Logt("已获取自定义域名", getCustomDomainResp)
	}

	// 更新自定义域名
	// REF: https://help.aliyun.com/zh/functioncompute/fc-2-0/developer-reference/api-fc-open-2021-04-06-updatecustomdomain
	updateCustomDomainReq := &aliyunFc2.UpdateCustomDomainRequest{
		CertConfig: &aliyunFc2.CertConfig{
			CertName:    tea.String(fmt.Sprintf("certimate-%d", time.Now().UnixMilli())),
			Certificate: tea.String(certPem),
			PrivateKey:  tea.String(privkeyPem),
		},
		Protocol:  getCustomDomainResp.Body.Protocol,
		TlsConfig: getCustomDomainResp.Body.TlsConfig,
	}
	updateCustomDomainResp, err := d.sdkClients.fc2.UpdateCustomDomain(tea.String(d.config.Domain), updateCustomDomainReq)
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'fc.UpdateCustomDomain'")
	} else {
		d.logger.Logt("已更新自定义域名", updateCustomDomainResp)
	}

	return nil
}

func createSdkClients(accessKeyId, accessKeySecret, region string) (*wSdkClients, error) {
	// 接入点一览 https://api.aliyun.com/product/FC-Open
	var fc2Endpoint string
	switch region {
	case "cn-hangzhou-finance":
		fc2Endpoint = fmt.Sprintf("%s.fc.aliyuncs.com", region)
	default:
		fc2Endpoint = fmt.Sprintf("fc.%s.aliyuncs.com", region)
	}

	fc2Config := &aliyunOpen.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
		Endpoint:        tea.String(fc2Endpoint),
	}
	fc2Client, err := aliyunFc2.NewClient(fc2Config)
	if err != nil {
		return nil, err
	}

	// 接入点一览 https://api.aliyun.com/product/FC-Open
	fc3Endpoint := fmt.Sprintf("fcv3.%s.aliyuncs.com", region)
	fc3Config := &aliyunOpen.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
		Endpoint:        tea.String(fc3Endpoint),
	}
	fc3Client, err := aliyunFc3.NewClient(fc3Config)
	if err != nil {
		return nil, err
	}

	return &wSdkClients{
		fc2: fc2Client,
		fc3: fc3Client,
	}, nil
}
