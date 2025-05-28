package aliyunapigw

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	aliapig "github.com/alibabacloud-go/apig-20240327/v3/client"
	alicloudapi "github.com/alibabacloud-go/cloudapi-20160714/v5/client"
	aliopen "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/alibabacloud-go/tea/tea"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	uploadersp "github.com/usual2970/certimate/internal/pkg/core/uploader/providers/aliyun-cas"
	typeutil "github.com/usual2970/certimate/internal/pkg/utils/type"
)

type DeployerConfig struct {
	// 阿里云 AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// 阿里云 AccessKeySecret。
	AccessKeySecret string `json:"accessKeySecret"`
	// 阿里云资源组 ID。
	ResourceGroupId string `json:"resourceGroupId,omitempty"`
	// 阿里云地域。
	Region string `json:"region"`
	// 服务类型。
	ServiceType ServiceType `json:"serviceType"`
	// API 网关 ID。
	// 服务类型为 [SERVICE_TYPE_CLOUDNATIVE] 时必填。
	GatewayId string `json:"gatewayId,omitempty"`
	// API 分组 ID。
	// 服务类型为 [SERVICE_TYPE_TRADITIONAL] 时必填。
	GroupId string `json:"groupId,omitempty"`
	// 自定义域名（支持泛域名）。
	Domain string `json:"domain"`
}

type DeployerProvider struct {
	config      *DeployerConfig
	logger      *slog.Logger
	sdkClients  *wSdkClients
	sslUploader uploader.Uploader
}

type wSdkClients struct {
	CloudNativeAPIGateway *aliapig.Client
	TraditionalAPIGateway *alicloudapi.Client
}

var _ deployer.Deployer = (*DeployerProvider)(nil)

func NewDeployer(config *DeployerConfig) (*DeployerProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	clients, err := createSdkClients(config.AccessKeyId, config.AccessKeySecret, config.Region)
	if err != nil {
		return nil, fmt.Errorf("failed to create sdk clients: %w", err)
	}

	uploader, err := createSslUploader(config.AccessKeyId, config.AccessKeySecret, config.ResourceGroupId, config.Region)
	if err != nil {
		return nil, fmt.Errorf("failed to create ssl uploader: %w", err)
	}

	return &DeployerProvider{
		config:      config,
		logger:      slog.Default(),
		sdkClients:  clients,
		sslUploader: uploader,
	}, nil
}

func (d *DeployerProvider) WithLogger(logger *slog.Logger) deployer.Deployer {
	if logger == nil {
		d.logger = slog.New(slog.DiscardHandler)
	} else {
		d.logger = logger
	}
	return d
}

func (d *DeployerProvider) Deploy(ctx context.Context, certPEM string, privkeyPEM string) (*deployer.DeployResult, error) {
	switch d.config.ServiceType {
	case SERVICE_TYPE_TRADITIONAL:
		if err := d.deployToTraditional(ctx, certPEM, privkeyPEM); err != nil {
			return nil, err
		}

	case SERVICE_TYPE_CLOUDNATIVE:
		if err := d.deployToCloudNative(ctx, certPEM, privkeyPEM); err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("unsupported service type '%s'", string(d.config.ServiceType))
	}

	return &deployer.DeployResult{}, nil
}

func (d *DeployerProvider) deployToTraditional(ctx context.Context, certPEM string, privkeyPEM string) error {
	if d.config.GroupId == "" {
		return errors.New("config `groupId` is required")
	}
	if d.config.Domain == "" {
		return errors.New("config `domain` is required")
	}

	// 为自定义域名添加 SSL 证书
	// REF: https://help.aliyun.com/zh/api-gateway/traditional-api-gateway/developer-reference/api-cloudapi-2016-07-14-setdomaincertificate
	setDomainCertificateReq := &alicloudapi.SetDomainCertificateRequest{
		GroupId:               tea.String(d.config.GroupId),
		DomainName:            tea.String(d.config.Domain),
		CertificateName:       tea.String(fmt.Sprintf("certimate_%d", time.Now().UnixMilli())),
		CertificateBody:       tea.String(certPEM),
		CertificatePrivateKey: tea.String(privkeyPEM),
	}
	setDomainCertificateResp, err := d.sdkClients.TraditionalAPIGateway.SetDomainCertificate(setDomainCertificateReq)
	d.logger.Debug("sdk request 'apigateway.SetDomainCertificate'", slog.Any("request", setDomainCertificateReq), slog.Any("response", setDomainCertificateResp))
	if err != nil {
		return fmt.Errorf("failed to execute sdk request 'apigateway.SetDomainCertificate': %w", err)
	}

	return nil
}

func (d *DeployerProvider) deployToCloudNative(ctx context.Context, certPEM string, privkeyPEM string) error {
	if d.config.GatewayId == "" {
		return errors.New("config `gatewayId` is required")
	}
	if d.config.Domain == "" {
		return errors.New("config `domain` is required")
	}

	// 遍历查询域名列表，获取域名 ID
	// REF: https://help.aliyun.com/zh/api-gateway/cloud-native-api-gateway/developer-reference/api-apig-2024-03-27-listdomains
	var domainId string
	listDomainsPageNumber := int32(1)
	listDomainsPageSize := int32(10)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		listDomainsReq := &aliapig.ListDomainsRequest{
			ResourceGroupId: typeutil.ToPtrOrZeroNil(d.config.ResourceGroupId),
			GatewayId:       tea.String(d.config.GatewayId),
			NameLike:        tea.String(d.config.Domain),
			PageNumber:      tea.Int32(listDomainsPageNumber),
			PageSize:        tea.Int32(listDomainsPageSize),
		}
		listDomainsResp, err := d.sdkClients.CloudNativeAPIGateway.ListDomains(listDomainsReq)
		d.logger.Debug("sdk request 'apig.ListDomains'", slog.Any("request", listDomainsReq), slog.Any("response", listDomainsResp))
		if err != nil {
			return fmt.Errorf("failed to execute sdk request 'apig.ListDomains': %w", err)
		}

		if listDomainsResp.Body.Data.Items != nil {
			for _, domainInfo := range listDomainsResp.Body.Data.Items {
				if strings.EqualFold(tea.StringValue(domainInfo.Name), d.config.Domain) {
					domainId = tea.StringValue(domainInfo.DomainId)
					break
				}
			}

			if domainId != "" {
				break
			}
		}

		if listDomainsResp.Body.Data.Items == nil || len(listDomainsResp.Body.Data.Items) < int(listDomainsPageSize) {
			break
		} else {
			listDomainsPageNumber++
		}
	}
	if domainId == "" {
		return errors.New("domain not found")
	}

	// 查询域名
	// REF: https://help.aliyun.com/zh/api-gateway/cloud-native-api-gateway/developer-reference/api-apig-2024-03-27-getdomain
	getDomainReq := &aliapig.GetDomainRequest{}
	getDomainResp, err := d.sdkClients.CloudNativeAPIGateway.GetDomain(tea.String(domainId), getDomainReq)
	d.logger.Debug("sdk request 'apig.GetDomain'", slog.Any("domainId", domainId), slog.Any("request", getDomainReq), slog.Any("response", getDomainResp))
	if err != nil {
		return fmt.Errorf("failed to execute sdk request 'apig.GetDomain': %w", err)
	}

	// 上传证书到 CAS
	upres, err := d.sslUploader.Upload(ctx, certPEM, privkeyPEM)
	if err != nil {
		return fmt.Errorf("failed to upload certificate file: %w", err)
	} else {
		d.logger.Info("ssl certificate uploaded", slog.Any("result", upres))
	}

	// 更新域名
	// REF: https://help.aliyun.com/zh/api-gateway/cloud-native-api-gateway/developer-reference/api-apig-2024-03-27-updatedomain
	updateDomainReq := &aliapig.UpdateDomainRequest{
		Protocol:              tea.String("HTTPS"),
		ForceHttps:            getDomainResp.Body.Data.ForceHttps,
		MTLSEnabled:           getDomainResp.Body.Data.MTLSEnabled,
		Http2Option:           getDomainResp.Body.Data.Http2Option,
		TlsMin:                getDomainResp.Body.Data.TlsMin,
		TlsMax:                getDomainResp.Body.Data.TlsMax,
		TlsCipherSuitesConfig: getDomainResp.Body.Data.TlsCipherSuitesConfig,
		CertIdentifier:        tea.String(upres.ExtendedData["certIdentifier"].(string)),
	}
	updateDomainResp, err := d.sdkClients.CloudNativeAPIGateway.UpdateDomain(tea.String(domainId), updateDomainReq)
	d.logger.Debug("sdk request 'apig.UpdateDomain'", slog.Any("domainId", domainId), slog.Any("request", updateDomainReq), slog.Any("response", updateDomainResp))
	if err != nil {
		return fmt.Errorf("failed to execute sdk request 'apig.UpdateDomain': %w", err)
	}

	return nil
}

func createSdkClients(accessKeyId, accessKeySecret, region string) (*wSdkClients, error) {
	// 接入点一览 https://api.aliyun.com/product/APIG
	cloudNativeAPIGEndpoint := strings.ReplaceAll(fmt.Sprintf("apig.%s.aliyuncs.com", region), "..", ".")
	cloudNativeAPIGConfig := &aliopen.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
		Endpoint:        tea.String(cloudNativeAPIGEndpoint),
	}
	cloudNativeAPIGClient, err := aliapig.NewClient(cloudNativeAPIGConfig)
	if err != nil {
		return nil, err
	}

	// 接入点一览 https://api.aliyun.com/product/CloudAPI
	traditionalAPIGEndpoint := strings.ReplaceAll(fmt.Sprintf("apigateway.%s.aliyuncs.com", region), "..", ".")
	traditionalAPIGConfig := &aliopen.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
		Endpoint:        tea.String(traditionalAPIGEndpoint),
	}
	traditionalAPIGClient, err := alicloudapi.NewClient(traditionalAPIGConfig)
	if err != nil {
		return nil, err
	}

	return &wSdkClients{
		CloudNativeAPIGateway: cloudNativeAPIGClient,
		TraditionalAPIGateway: traditionalAPIGClient,
	}, nil
}

func createSslUploader(accessKeyId, accessKeySecret, resourceGroupId, region string) (uploader.Uploader, error) {
	casRegion := region
	if casRegion != "" {
		// 阿里云 CAS 服务接入点是独立于 APIGateway 服务的
		// 国内版固定接入点：华东一杭州
		// 国际版固定接入点：亚太东南一新加坡
		if !strings.HasPrefix(casRegion, "cn-") {
			casRegion = "ap-southeast-1"
		} else {
			casRegion = "cn-hangzhou"
		}
	}

	uploader, err := uploadersp.NewUploader(&uploadersp.UploaderConfig{
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
		ResourceGroupId: resourceGroupId,
		Region:          casRegion,
	})
	return uploader, err
}
