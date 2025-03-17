package huaweicloudcdn

import (
	"context"
	"log/slog"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/global"
	hcCdn "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cdn/v2"
	hcCdnModel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cdn/v2/model"
	hcCdnRegion "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cdn/v2/region"
	xerrors "github.com/pkg/errors"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	uploadersp "github.com/usual2970/certimate/internal/pkg/core/uploader/providers/huaweicloud-scm"
	hwsdk "github.com/usual2970/certimate/internal/pkg/vendors/huaweicloud-sdk"
)

type DeployerConfig struct {
	// 华为云 AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// 华为云 SecretAccessKey。
	SecretAccessKey string `json:"secretAccessKey"`
	// 华为云区域。
	Region string `json:"region"`
	// 加速域名（不支持泛域名）。
	Domain string `json:"domain"`
}

type DeployerProvider struct {
	config      *DeployerConfig
	logger      *slog.Logger
	sdkClient   *hcCdn.CdnClient
	sslUploader uploader.Uploader
}

var _ deployer.Deployer = (*DeployerProvider)(nil)

func NewDeployer(config *DeployerConfig) (*DeployerProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	client, err := createSdkClient(
		config.AccessKeyId,
		config.SecretAccessKey,
		config.Region,
	)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk client")
	}

	uploader, err := uploadersp.NewUploader(&uploadersp.UploaderConfig{
		AccessKeyId:     config.AccessKeyId,
		SecretAccessKey: config.SecretAccessKey,
	})
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create ssl uploader")
	}

	return &DeployerProvider{
		config:      config,
		logger:      slog.Default(),
		sdkClient:   client,
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
	// 上传证书到 SCM
	upres, err := d.sslUploader.Upload(ctx, certPem, privkeyPem)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to upload certificate file")
	} else {
		d.logger.Info("ssl certificate uploaded", slog.Any("result", upres))
	}

	// 查询加速域名配置
	// REF: https://support.huaweicloud.com/api-cdn/ShowDomainFullConfig.html
	showDomainFullConfigReq := &hcCdnModel.ShowDomainFullConfigRequest{
		DomainName: d.config.Domain,
	}
	showDomainFullConfigResp, err := d.sdkClient.ShowDomainFullConfig(showDomainFullConfigReq)
	d.logger.Debug("sdk request 'cdn.ShowDomainFullConfig'", slog.Any("request", showDomainFullConfigReq), slog.Any("response", showDomainFullConfigResp))
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'cdn.ShowDomainFullConfig'")
	}

	// 更新加速域名配置
	// REF: https://support.huaweicloud.com/api-cdn/UpdateDomainMultiCertificates.html
	// REF: https://support.huaweicloud.com/usermanual-cdn/cdn_01_0306.html
	updateDomainMultiCertificatesReqBodyContent := &hcCdnModel.UpdateDomainMultiCertificatesRequestBodyContent{}
	updateDomainMultiCertificatesReqBodyContent.DomainName = d.config.Domain
	updateDomainMultiCertificatesReqBodyContent.HttpsSwitch = 1
	updateDomainMultiCertificatesReqBodyContent.CertificateType = hwsdk.Int32Ptr(2)
	updateDomainMultiCertificatesReqBodyContent.ScmCertificateId = hwsdk.StringPtr(upres.CertId)
	updateDomainMultiCertificatesReqBodyContent.CertName = hwsdk.StringPtr(upres.CertName)
	updateDomainMultiCertificatesReqBodyContent = assign(updateDomainMultiCertificatesReqBodyContent, showDomainFullConfigResp.Configs)
	updateDomainMultiCertificatesReq := &hcCdnModel.UpdateDomainMultiCertificatesRequest{
		Body: &hcCdnModel.UpdateDomainMultiCertificatesRequestBody{
			Https: updateDomainMultiCertificatesReqBodyContent,
		},
	}
	updateDomainMultiCertificatesResp, err := d.sdkClient.UpdateDomainMultiCertificates(updateDomainMultiCertificatesReq)
	d.logger.Debug("sdk request 'cdn.UploadDomainMultiCertificates'", slog.Any("request", updateDomainMultiCertificatesReq), slog.Any("response", updateDomainMultiCertificatesResp))
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'cdn.UploadDomainMultiCertificates'")
	}

	return &deployer.DeployResult{}, nil
}

func createSdkClient(accessKeyId, secretAccessKey, region string) (*hcCdn.CdnClient, error) {
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

	client := hcCdn.NewCdnClient(hcClient)
	return client, nil
}

func assign(reqContent *hcCdnModel.UpdateDomainMultiCertificatesRequestBodyContent, target *hcCdnModel.ConfigsGetBody) *hcCdnModel.UpdateDomainMultiCertificatesRequestBodyContent {
	if target == nil {
		return reqContent
	}

	// 华为云 API 中不传的字段表示使用默认值、而非保留原值，因此这里需要把原配置中的参数重新赋值回去。
	// 而且蛋疼的是查询接口返回的数据结构和更新接口传入的参数结构不一致，需要做很多转化。

	if *target.OriginProtocol == "follow" {
		reqContent.AccessOriginWay = hwsdk.Int32Ptr(1)
	} else if *target.OriginProtocol == "http" {
		reqContent.AccessOriginWay = hwsdk.Int32Ptr(2)
	} else if *target.OriginProtocol == "https" {
		reqContent.AccessOriginWay = hwsdk.Int32Ptr(3)
	}

	if target.ForceRedirect != nil {
		reqContent.ForceRedirectConfig = &hcCdnModel.ForceRedirect{}

		if target.ForceRedirect.Status == "on" {
			reqContent.ForceRedirectConfig.Switch = 1
			reqContent.ForceRedirectConfig.RedirectType = target.ForceRedirect.Type
		} else {
			reqContent.ForceRedirectConfig.Switch = 0
		}
	}

	if target.Https != nil {
		if *target.Https.Http2Status == "on" {
			reqContent.Http2 = hwsdk.Int32Ptr(1)
		}
	}

	return reqContent
}
