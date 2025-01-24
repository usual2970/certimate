package tencentcloudecdn

import (
	"context"
	"errors"
	"strings"

	xerrors "github.com/pkg/errors"
	tcCdn "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn/v20180606"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tcSsl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/logger"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	uploaderp "github.com/usual2970/certimate/internal/pkg/core/uploader/providers/tencentcloud-ssl"
)

type TencentCloudECDNDeployerConfig struct {
	// 腾讯云 SecretId。
	SecretId string `json:"secretId"`
	// 腾讯云 SecretKey。
	SecretKey string `json:"secretKey"`
	// 加速域名（支持泛域名）。
	Domain string `json:"domain"`
}

type TencentCloudECDNDeployer struct {
	config      *TencentCloudECDNDeployerConfig
	logger      logger.Logger
	sdkClients  *wSdkClients
	sslUploader uploader.Uploader
}

var _ deployer.Deployer = (*TencentCloudECDNDeployer)(nil)

type wSdkClients struct {
	ssl *tcSsl.Client
	cdn *tcCdn.Client
}

func New(config *TencentCloudECDNDeployerConfig) (*TencentCloudECDNDeployer, error) {
	return NewWithLogger(config, logger.NewNilLogger())
}

func NewWithLogger(config *TencentCloudECDNDeployerConfig, logger logger.Logger) (*TencentCloudECDNDeployer, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}

	if logger == nil {
		return nil, errors.New("logger is nil")
	}

	clients, err := createSdkClients(config.SecretId, config.SecretKey)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk clients")
	}

	uploader, err := uploaderp.New(&uploaderp.TencentCloudSSLUploaderConfig{
		SecretId:  config.SecretId,
		SecretKey: config.SecretKey,
	})
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create ssl uploader")
	}

	return &TencentCloudECDNDeployer{
		logger:      logger,
		config:      config,
		sdkClients:  clients,
		sslUploader: uploader,
	}, nil
}

func (d *TencentCloudECDNDeployer) Deploy(ctx context.Context, certPem string, privkeyPem string) (*deployer.DeployResult, error) {
	// 上传证书到 SSL
	upres, err := d.sslUploader.Upload(ctx, certPem, privkeyPem)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to upload certificate file")
	}

	d.logger.Logt("certificate file uploaded", upres)

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

	if len(instanceIds) == 0 {
		d.logger.Logt("已部署过或没有要部署的 ECDN 实例")
	} else {
		// 证书部署到 ECDN 实例
		// REF: https://cloud.tencent.com/document/product/400/91667
		deployCertificateInstanceReq := tcSsl.NewDeployCertificateInstanceRequest()
		deployCertificateInstanceReq.CertificateId = common.StringPtr(upres.CertId)
		deployCertificateInstanceReq.ResourceType = common.StringPtr("ecdn")
		deployCertificateInstanceReq.Status = common.Int64Ptr(1)
		deployCertificateInstanceReq.InstanceIdList = common.StringPtrs(instanceIds)
		deployCertificateInstanceResp, err := d.sdkClients.ssl.DeployCertificateInstance(deployCertificateInstanceReq)
		if err != nil {
			return nil, xerrors.Wrap(err, "failed to execute sdk request 'ssl.DeployCertificateInstance'")
		}

		d.logger.Logt("已部署证书到云资源实例", deployCertificateInstanceResp.Response)
	}

	return &deployer.DeployResult{}, nil
}

func (d *TencentCloudECDNDeployer) getDomainsByCertificateId(cloudCertId string) ([]string, error) {
	// 获取证书中的可用域名
	// REF: https://cloud.tencent.com/document/product/228/42491
	describeCertDomainsReq := tcCdn.NewDescribeCertDomainsRequest()
	describeCertDomainsReq.CertId = common.StringPtr(cloudCertId)
	describeCertDomainsReq.Product = common.StringPtr("ecdn")
	describeCertDomainsResp, err := d.sdkClients.cdn.DescribeCertDomains(describeCertDomainsReq)
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
