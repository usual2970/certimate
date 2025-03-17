package tencentcloudcdn

import (
	"context"
	"log/slog"
	"strings"

	xerrors "github.com/pkg/errors"
	tcCdn "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn/v20180606"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tcSsl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"
	"golang.org/x/exp/slices"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	uploadersp "github.com/usual2970/certimate/internal/pkg/core/uploader/providers/tencentcloud-ssl"
)

type DeployerConfig struct {
	// 腾讯云 SecretId。
	SecretId string `json:"secretId"`
	// 腾讯云 SecretKey。
	SecretKey string `json:"secretKey"`
	// 加速域名（支持泛域名）。
	Domain string `json:"domain"`
}

type DeployerProvider struct {
	config      *DeployerConfig
	logger      *slog.Logger
	sdkClients  *wSdkClients
	sslUploader uploader.Uploader
}

var _ deployer.Deployer = (*DeployerProvider)(nil)

type wSdkClients struct {
	ssl *tcSsl.Client
	cdn *tcCdn.Client
}

func NewDeployer(config *DeployerConfig) (*DeployerProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	clients, err := createSdkClients(config.SecretId, config.SecretKey)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk clients")
	}

	uploader, err := uploadersp.NewUploader(&uploadersp.UploaderConfig{
		SecretId:  config.SecretId,
		SecretKey: config.SecretKey,
	})
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create ssl uploader")
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
		d.logger = slog.Default()
	} else {
		d.logger = logger
	}
	d.sslUploader.WithLogger(logger)
	return d
}

func (d *DeployerProvider) Deploy(ctx context.Context, certPem string, privkeyPem string) (*deployer.DeployResult, error) {
	// 上传证书到 SSL
	upres, err := d.sslUploader.Upload(ctx, certPem, privkeyPem)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to upload certificate file")
	} else {
		d.logger.Info("ssl certificate uploaded", slog.Any("result", upres))
	}

	// 获取待部署的 CDN 实例
	// 如果是泛域名，根据证书匹配 CDN 实例
	instanceIds := make([]string, 0)
	if strings.HasPrefix(d.config.Domain, "*.") {
		domains, err := d.getDomainsByCertificateId(upres.CertId)
		if err != nil {
			return nil, err
		}

		instanceIds = domains
	} else {
		instanceIds = append(instanceIds, d.config.Domain)
	}

	// 跳过已部署的 CDN 实例
	if len(instanceIds) > 0 {
		deployedDomains, err := d.getDeployedDomainsByCertificateId(upres.CertId)
		if err != nil {
			return nil, err
		}

		temp := make([]string, 0)
		for _, instanceId := range instanceIds {
			if !slices.Contains(deployedDomains, instanceId) {
				temp = append(temp, instanceId)
			}
		}
		instanceIds = temp
	}

	if len(instanceIds) == 0 {
		d.logger.Info("no cdn instances to deploy")
	} else {
		d.logger.Info("found cdn instances to deploy", slog.Any("instanceIds", instanceIds))

		// 证书部署到 CDN 实例
		// REF: https://cloud.tencent.com/document/product/400/91667
		deployCertificateInstanceReq := tcSsl.NewDeployCertificateInstanceRequest()
		deployCertificateInstanceReq.CertificateId = common.StringPtr(upres.CertId)
		deployCertificateInstanceReq.ResourceType = common.StringPtr("cdn")
		deployCertificateInstanceReq.Status = common.Int64Ptr(1)
		deployCertificateInstanceReq.InstanceIdList = common.StringPtrs(instanceIds)
		deployCertificateInstanceResp, err := d.sdkClients.ssl.DeployCertificateInstance(deployCertificateInstanceReq)
		d.logger.Debug("sdk request 'ssl.DeployCertificateInstance'", slog.Any("request", deployCertificateInstanceReq), slog.Any("response", deployCertificateInstanceResp))
		if err != nil {
			return nil, xerrors.Wrap(err, "failed to execute sdk request 'ssl.DeployCertificateInstance'")
		}
	}

	return &deployer.DeployResult{}, nil
}

func (d *DeployerProvider) getDomainsByCertificateId(cloudCertId string) ([]string, error) {
	// 获取证书中的可用域名
	// REF: https://cloud.tencent.com/document/product/228/42491
	describeCertDomainsReq := tcCdn.NewDescribeCertDomainsRequest()
	describeCertDomainsReq.CertId = common.StringPtr(cloudCertId)
	describeCertDomainsReq.Product = common.StringPtr("cdn")
	describeCertDomainsResp, err := d.sdkClients.cdn.DescribeCertDomains(describeCertDomainsReq)
	d.logger.Debug("sdk request 'cdn.DescribeCertDomains'", slog.Any("request", describeCertDomainsReq), slog.Any("response", describeCertDomainsResp))
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'cdn.DescribeCertDomains'")
	}

	domains := make([]string, 0)
	if describeCertDomainsResp.Response.Domains != nil {
		for _, domain := range describeCertDomainsResp.Response.Domains {
			domains = append(domains, *domain)
		}
	}

	return domains, nil
}

func (d *DeployerProvider) getDeployedDomainsByCertificateId(cloudCertId string) ([]string, error) {
	// 根据证书查询关联 CDN 域名
	// REF: https://cloud.tencent.com/document/product/400/62674
	describeDeployedResourcesReq := tcSsl.NewDescribeDeployedResourcesRequest()
	describeDeployedResourcesReq.CertificateIds = common.StringPtrs([]string{cloudCertId})
	describeDeployedResourcesReq.ResourceType = common.StringPtr("cdn")
	describeDeployedResourcesResp, err := d.sdkClients.ssl.DescribeDeployedResources(describeDeployedResourcesReq)
	d.logger.Debug("sdk request 'cdn.DescribeDeployedResources'", slog.Any("request", describeDeployedResourcesReq), slog.Any("response", describeDeployedResourcesResp))
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'cdn.DescribeDeployedResources'")
	}

	domains := make([]string, 0)
	if describeDeployedResourcesResp.Response.DeployedResources != nil {
		for _, deployedResource := range describeDeployedResourcesResp.Response.DeployedResources {
			for _, resource := range deployedResource.Resources {
				domains = append(domains, *resource)
			}
		}
	}

	return domains, nil
}

func createSdkClients(secretId, secretKey string) (*wSdkClients, error) {
	credential := common.NewCredential(secretId, secretKey)

	sslClient, err := tcSsl.NewClient(credential, "", profile.NewClientProfile())
	if err != nil {
		return nil, err
	}

	cdnClient, err := tcCdn.NewClient(credential, "", profile.NewClientProfile())
	if err != nil {
		return nil, err
	}

	return &wSdkClients{
		ssl: sslClient,
		cdn: cdnClient,
	}, nil
}
