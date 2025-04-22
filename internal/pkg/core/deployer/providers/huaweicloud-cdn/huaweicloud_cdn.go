package huaweicloudcdn

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/global"
	hccdn "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cdn/v2"
	hccdnmodel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cdn/v2/model"
	hccdnregion "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cdn/v2/region"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	uploadersp "github.com/usual2970/certimate/internal/pkg/core/uploader/providers/huaweicloud-scm"
	typeutil "github.com/usual2970/certimate/internal/pkg/utils/type"
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
	sdkClient   *hccdn.CdnClient
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
		return nil, fmt.Errorf("failed to create sdk client: %w", err)
	}

	uploader, err := uploadersp.NewUploader(&uploadersp.UploaderConfig{
		AccessKeyId:     config.AccessKeyId,
		SecretAccessKey: config.SecretAccessKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create ssl uploader: %w", err)
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

func (d *DeployerProvider) Deploy(ctx context.Context, certPEM string, privkeyPEM string) (*deployer.DeployResult, error) {
	// 上传证书到 SCM
	upres, err := d.sslUploader.Upload(ctx, certPEM, privkeyPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to upload certificate file: %w", err)
	} else {
		d.logger.Info("ssl certificate uploaded", slog.Any("result", upres))
	}

	// 查询加速域名配置
	// REF: https://support.huaweicloud.com/api-cdn/ShowDomainFullConfig.html
	showDomainFullConfigReq := &hccdnmodel.ShowDomainFullConfigRequest{
		DomainName: d.config.Domain,
	}
	showDomainFullConfigResp, err := d.sdkClient.ShowDomainFullConfig(showDomainFullConfigReq)
	d.logger.Debug("sdk request 'cdn.ShowDomainFullConfig'", slog.Any("request", showDomainFullConfigReq), slog.Any("response", showDomainFullConfigResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'cdn.ShowDomainFullConfig': %w", err)
	}

	// 更新加速域名配置
	// REF: https://support.huaweicloud.com/api-cdn/UpdateDomainMultiCertificates.html
	// REF: https://support.huaweicloud.com/usermanual-cdn/cdn_01_0306.html
	updateDomainMultiCertificatesReqBodyContent := &hccdnmodel.UpdateDomainMultiCertificatesRequestBodyContent{}
	updateDomainMultiCertificatesReqBodyContent.DomainName = d.config.Domain
	updateDomainMultiCertificatesReqBodyContent.HttpsSwitch = 1
	updateDomainMultiCertificatesReqBodyContent.CertificateType = typeutil.ToPtr(int32(2))
	updateDomainMultiCertificatesReqBodyContent.ScmCertificateId = typeutil.ToPtr(upres.CertId)
	updateDomainMultiCertificatesReqBodyContent.CertName = typeutil.ToPtr(upres.CertName)
	updateDomainMultiCertificatesReqBodyContent = assign(updateDomainMultiCertificatesReqBodyContent, showDomainFullConfigResp.Configs)
	updateDomainMultiCertificatesReq := &hccdnmodel.UpdateDomainMultiCertificatesRequest{
		Body: &hccdnmodel.UpdateDomainMultiCertificatesRequestBody{
			Https: updateDomainMultiCertificatesReqBodyContent,
		},
	}
	updateDomainMultiCertificatesResp, err := d.sdkClient.UpdateDomainMultiCertificates(updateDomainMultiCertificatesReq)
	d.logger.Debug("sdk request 'cdn.UploadDomainMultiCertificates'", slog.Any("request", updateDomainMultiCertificatesReq), slog.Any("response", updateDomainMultiCertificatesResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'cdn.UploadDomainMultiCertificates': %w", err)
	}

	return &deployer.DeployResult{}, nil
}

func createSdkClient(accessKeyId, secretAccessKey, region string) (*hccdn.CdnClient, error) {
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

	hcRegion, err := hccdnregion.SafeValueOf(region)
	if err != nil {
		return nil, err
	}

	hcClient, err := hccdn.CdnClientBuilder().
		WithRegion(hcRegion).
		WithCredential(auth).
		SafeBuild()
	if err != nil {
		return nil, err
	}

	client := hccdn.NewCdnClient(hcClient)
	return client, nil
}

func assign(source *hccdnmodel.UpdateDomainMultiCertificatesRequestBodyContent, target *hccdnmodel.ConfigsGetBody) *hccdnmodel.UpdateDomainMultiCertificatesRequestBodyContent {
	// `UpdateDomainMultiCertificates` 中不传的字段表示使用默认值、而非保留原值，
	// 因此这里需要把原配置中的参数重新赋值回去。

	if target == nil {
		return source
	}

	if *target.OriginProtocol == "follow" {
		source.AccessOriginWay = typeutil.ToPtr(int32(1))
	} else if *target.OriginProtocol == "http" {
		source.AccessOriginWay = typeutil.ToPtr(int32(2))
	} else if *target.OriginProtocol == "https" {
		source.AccessOriginWay = typeutil.ToPtr(int32(3))
	}

	if target.ForceRedirect != nil {
		if source.ForceRedirectConfig == nil {
			source.ForceRedirectConfig = &hccdnmodel.ForceRedirect{}
		}

		if target.ForceRedirect.Status == "on" {
			source.ForceRedirectConfig.Switch = 1
			source.ForceRedirectConfig.RedirectType = target.ForceRedirect.Type
		} else {
			source.ForceRedirectConfig.Switch = 0
		}
	}

	if target.Https != nil {
		if *target.Https.Http2Status == "on" {
			source.Http2 = typeutil.ToPtr(int32(1))
		}
	}

	return source
}
