package huaweicloudcdn

import (
	"context"
	"errors"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/global"
	hcCdn "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cdn/v2"
	hcCdnModel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cdn/v2/model"
	hcCdnRegion "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cdn/v2/region"
	xerrors "github.com/pkg/errors"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/logger"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	providerScm "github.com/usual2970/certimate/internal/pkg/core/uploader/providers/huaweicloud-scm"
	"github.com/usual2970/certimate/internal/pkg/utils/cast"
	huaweicloudsdk "github.com/usual2970/certimate/internal/pkg/vendors/huaweicloud-cdn-sdk"
)

type HuaweiCloudCDNDeployerConfig struct {
	// 华为云 AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// 华为云 SecretAccessKey。
	SecretAccessKey string `json:"secretAccessKey"`
	// 华为云区域。
	Region string `json:"region"`
	// 加速域名（不支持泛域名）。
	Domain string `json:"domain"`
}

type HuaweiCloudCDNDeployer struct {
	config      *HuaweiCloudCDNDeployerConfig
	logger      logger.Logger
	sdkClient   *huaweicloudsdk.Client
	sslUploader uploader.Uploader
}

var _ deployer.Deployer = (*HuaweiCloudCDNDeployer)(nil)

func New(config *HuaweiCloudCDNDeployerConfig) (*HuaweiCloudCDNDeployer, error) {
	return NewWithLogger(config, logger.NewNilLogger())
}

func NewWithLogger(config *HuaweiCloudCDNDeployerConfig, logger logger.Logger) (*HuaweiCloudCDNDeployer, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}

	if logger == nil {
		return nil, errors.New("logger is nil")
	}

	client, err := createSdkClient(
		config.AccessKeyId,
		config.SecretAccessKey,
		config.Region,
	)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk client")
	}

	uploader, err := providerScm.New(&providerScm.HuaweiCloudSCMUploaderConfig{
		AccessKeyId:     config.AccessKeyId,
		SecretAccessKey: config.SecretAccessKey,
	})
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create ssl uploader")
	}

	return &HuaweiCloudCDNDeployer{
		logger:      logger,
		config:      config,
		sdkClient:   client,
		sslUploader: uploader,
	}, nil
}

func (d *HuaweiCloudCDNDeployer) Deploy(ctx context.Context, certPem string, privkeyPem string) (*deployer.DeployResult, error) {
	// 上传证书到 SCM
	upres, err := d.sslUploader.Upload(ctx, certPem, privkeyPem)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to upload certificate file")
	}

	d.logger.Logt("certificate file uploaded", upres)

	// 查询加速域名配置
	// REF: https://support.huaweicloud.com/api-cdn/ShowDomainFullConfig.html
	showDomainFullConfigReq := &hcCdnModel.ShowDomainFullConfigRequest{
		DomainName: d.config.Domain,
	}
	showDomainFullConfigResp, err := d.sdkClient.ShowDomainFullConfig(showDomainFullConfigReq)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'cdn.ShowDomainFullConfig'")
	}

	d.logger.Logt("已查询到加速域名配置", showDomainFullConfigResp)

	// 更新加速域名配置
	// REF: https://support.huaweicloud.com/api-cdn/UpdateDomainMultiCertificates.html
	// REF: https://support.huaweicloud.com/usermanual-cdn/cdn_01_0306.html
	updateDomainMultiCertificatesReqBodyContent := &huaweicloudsdk.UpdateDomainMultiCertificatesExRequestBodyContent{}
	updateDomainMultiCertificatesReqBodyContent.DomainName = d.config.Domain
	updateDomainMultiCertificatesReqBodyContent.HttpsSwitch = 1
	updateDomainMultiCertificatesReqBodyContent.CertificateType = cast.Int32Ptr(2)
	updateDomainMultiCertificatesReqBodyContent.SCMCertificateId = cast.StringPtr(upres.CertId)
	updateDomainMultiCertificatesReqBodyContent.CertName = cast.StringPtr(upres.CertName)
	updateDomainMultiCertificatesReqBodyContent = updateDomainMultiCertificatesReqBodyContent.MergeConfig(showDomainFullConfigResp.Configs)
	updateDomainMultiCertificatesReq := &huaweicloudsdk.UpdateDomainMultiCertificatesExRequest{
		Body: &huaweicloudsdk.UpdateDomainMultiCertificatesExRequestBody{
			Https: updateDomainMultiCertificatesReqBodyContent,
		},
	}
	updateDomainMultiCertificatesResp, err := d.sdkClient.UploadDomainMultiCertificatesEx(updateDomainMultiCertificatesReq)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'cdn.UploadDomainMultiCertificatesEx'")
	}

	d.logger.Logt("已更新加速域名配置", updateDomainMultiCertificatesResp)

	return &deployer.DeployResult{}, nil
}

func createSdkClient(accessKeyId, secretAccessKey, region string) (*huaweicloudsdk.Client, error) {
	if region == "" {
		region = "cn-north-1" // CDN 服务默认区域：华北一北京
	}

	auth, err := global.NewCredentialsBuilder().
		WithAk(accessKeyId).
		WithSk(secretAccessKey).
		SafeBuild()
	if err != nil {
		return nil, err
	}

	hcRegion, err := hcCdnRegion.SafeValueOf(region)
	if err != nil {
		return nil, err
	}

	hcClient, err := hcCdn.CdnClientBuilder().
		WithRegion(hcRegion).
		WithCredential(auth).
		SafeBuild()
	if err != nil {
		return nil, err
	}

	client := huaweicloudsdk.NewClient(hcClient)
	return client, nil
}
