package aliyunfc

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	aliopen "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	alifc3 "github.com/alibabacloud-go/fc-20230330/v4/client"
	alifc2 "github.com/alibabacloud-go/fc-open-20210406/v2/client"
	"github.com/alibabacloud-go/tea/tea"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
)

type DeployerConfig struct {
	// 阿里云 AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// 阿里云 AccessKeySecret。
	AccessKeySecret string `json:"accessKeySecret"`
	// 阿里云地域。
	Region string `json:"region"`
	// 服务版本。
	ServiceVersion string `json:"serviceVersion"`
	// 自定义域名（支持泛域名）。
	Domain string `json:"domain"`
}

type DeployerProvider struct {
	config     *DeployerConfig
	logger     *slog.Logger
	sdkClients *wSdkClients
}

var _ deployer.Deployer = (*DeployerProvider)(nil)

type wSdkClients struct {
	FC2 *alifc2.Client
	FC3 *alifc3.Client
}

func NewDeployer(config *DeployerConfig) (*DeployerProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	clients, err := createSdkClients(config.AccessKeyId, config.AccessKeySecret, config.Region)
	if err != nil {
		return nil, fmt.Errorf("failed to create sdk clients: %w", err)
	}

	return &DeployerProvider{
		config:     config,
		logger:     slog.Default(),
		sdkClients: clients,
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

func (d *DeployerProvider) Deploy(ctx context.Context, certPEM string, privkeyPEM string) (*deployer.DeployResult, error) {
	switch d.config.ServiceVersion {
	case "3", "3.0":
		if err := d.deployToFC3(ctx, certPEM, privkeyPEM); err != nil {
			return nil, err
		}

	case "2", "2.0":
		if err := d.deployToFC2(ctx, certPEM, privkeyPEM); err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("unsupported service version '%s'", d.config.ServiceVersion)
	}

	return &deployer.DeployResult{}, nil
}

func (d *DeployerProvider) deployToFC3(ctx context.Context, certPEM string, privkeyPEM string) error {
	// 获取自定义域名
	// REF: https://help.aliyun.com/zh/functioncompute/fc-3-0/developer-reference/api-fc-2023-03-30-getcustomdomain
	getCustomDomainResp, err := d.sdkClients.FC3.GetCustomDomain(tea.String(d.config.Domain))
	d.logger.Debug("sdk request 'fc.GetCustomDomain'", slog.Any("response", getCustomDomainResp))
	if err != nil {
		return fmt.Errorf("failed to execute sdk request 'fc.GetCustomDomain': %w", err)
	}

	// 更新自定义域名
	// REF: https://help.aliyun.com/zh/functioncompute/fc-3-0/developer-reference/api-fc-2023-03-30-updatecustomdomain
	updateCustomDomainReq := &alifc3.UpdateCustomDomainRequest{
		Body: &alifc3.UpdateCustomDomainInput{
			CertConfig: &alifc3.CertConfig{
				CertName:    tea.String(fmt.Sprintf("certimate-%d", time.Now().UnixMilli())),
				Certificate: tea.String(certPEM),
				PrivateKey:  tea.String(privkeyPEM),
			},
			Protocol:  getCustomDomainResp.Body.Protocol,
			TlsConfig: getCustomDomainResp.Body.TlsConfig,
		},
	}
	updateCustomDomainResp, err := d.sdkClients.FC3.UpdateCustomDomain(tea.String(d.config.Domain), updateCustomDomainReq)
	d.logger.Debug("sdk request 'fc.UpdateCustomDomain'", slog.Any("request", updateCustomDomainReq), slog.Any("response", updateCustomDomainResp))
	if err != nil {
		return fmt.Errorf("failed to execute sdk request 'fc.UpdateCustomDomain': %w", err)
	}

	return nil
}

func (d *DeployerProvider) deployToFC2(ctx context.Context, certPEM string, privkeyPEM string) error {
	// 获取自定义域名
	// REF: https://help.aliyun.com/zh/functioncompute/fc-2-0/developer-reference/api-fc-open-2021-04-06-getcustomdomain
	getCustomDomainResp, err := d.sdkClients.FC2.GetCustomDomain(tea.String(d.config.Domain))
	d.logger.Debug("sdk request 'fc.GetCustomDomain'", slog.Any("response", getCustomDomainResp))
	if err != nil {
		return fmt.Errorf("failed to execute sdk request 'fc.GetCustomDomain': %w", err)
	}

	// 更新自定义域名
	// REF: https://help.aliyun.com/zh/functioncompute/fc-2-0/developer-reference/api-fc-open-2021-04-06-updatecustomdomain
	updateCustomDomainReq := &alifc2.UpdateCustomDomainRequest{
		CertConfig: &alifc2.CertConfig{
			CertName:    tea.String(fmt.Sprintf("certimate-%d", time.Now().UnixMilli())),
			Certificate: tea.String(certPEM),
			PrivateKey:  tea.String(privkeyPEM),
		},
		Protocol:  getCustomDomainResp.Body.Protocol,
		TlsConfig: getCustomDomainResp.Body.TlsConfig,
	}
	updateCustomDomainResp, err := d.sdkClients.FC2.UpdateCustomDomain(tea.String(d.config.Domain), updateCustomDomainReq)
	d.logger.Debug("sdk request 'fc.UpdateCustomDomain'", slog.Any("request", updateCustomDomainReq), slog.Any("response", updateCustomDomainResp))
	if err != nil {
		return fmt.Errorf("failed to execute sdk request 'fc.UpdateCustomDomain': %w", err)
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

	fc2Config := &aliopen.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
		Endpoint:        tea.String(fc2Endpoint),
	}
	fc2Client, err := alifc2.NewClient(fc2Config)
	if err != nil {
		return nil, err
	}

	// 接入点一览 https://api.aliyun.com/product/FC-Open
	fc3Endpoint := fmt.Sprintf("fcv3.%s.aliyuncs.com", region)
	fc3Config := &aliopen.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
		Endpoint:        tea.String(fc3Endpoint),
	}
	fc3Client, err := alifc3.NewClient(fc3Config)
	if err != nil {
		return nil, err
	}

	return &wSdkClients{
		FC2: fc2Client,
		FC3: fc3Client,
	}, nil
}
