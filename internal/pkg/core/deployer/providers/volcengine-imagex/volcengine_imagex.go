package volcengineimagex

import (
	"context"
	"errors"

	xerrors "github.com/pkg/errors"
	veBase "github.com/volcengine/volc-sdk-golang/base"
	veImageX "github.com/volcengine/volc-sdk-golang/service/imagex/v2"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/logger"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	uploadersp "github.com/usual2970/certimate/internal/pkg/core/uploader/providers/volcengine-certcenter"
)

type VolcEngineImageXDeployerConfig struct {
	// 火山引擎 AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// 火山引擎 AccessKeySecret。
	AccessKeySecret string `json:"accessKeySecret"`
	// 火山引擎地域。
	Region string `json:"region"`
	// 服务 ID。
	ServiceId string `json:"serviceId"`
	// 自定义域名（不支持泛域名）。
	Domain string `json:"domain"`
}

type VolcEngineImageXDeployer struct {
	config      *VolcEngineImageXDeployerConfig
	logger      logger.Logger
	sdkClient   *veImageX.Imagex
	sslUploader uploader.Uploader
}

var _ deployer.Deployer = (*VolcEngineImageXDeployer)(nil)

func New(config *VolcEngineImageXDeployerConfig) (*VolcEngineImageXDeployer, error) {
	return NewWithLogger(config, logger.NewNilLogger())
}

func NewWithLogger(config *VolcEngineImageXDeployerConfig, logger logger.Logger) (*VolcEngineImageXDeployer, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}

	if logger == nil {
		return nil, errors.New("logger is nil")
	}

	client, err := createSdkClient(config.AccessKeyId, config.AccessKeySecret, config.Region)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk client")
	}

	uploader, err := uploadersp.New(&uploadersp.VolcEngineCertCenterUploaderConfig{
		AccessKeyId:     config.AccessKeyId,
		AccessKeySecret: config.AccessKeySecret,
		Region:          config.Region,
	})
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create ssl uploader")
	}

	return &VolcEngineImageXDeployer{
		logger:      logger,
		config:      config,
		sdkClient:   client,
		sslUploader: uploader,
	}, nil
}

func (d *VolcEngineImageXDeployer) Deploy(ctx context.Context, certPem string, privkeyPem string) (*deployer.DeployResult, error) {
	if d.config.ServiceId == "" {
		return nil, errors.New("config `serviceId` is required")
	}
	if d.config.Domain == "" {
		return nil, errors.New("config `domain` is required")
	}

	// 上传证书到证书中心
	upres, err := d.sslUploader.Upload(ctx, certPem, privkeyPem)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to upload certificate file")
	} else {
		d.logger.Logt("certificate file uploaded", upres)
	}

	// 获取域名配置
	// REF: https://www.volcengine.com/docs/508/9366
	getDomainConfigReq := &veImageX.GetDomainConfigQuery{
		ServiceID:  d.config.ServiceId,
		DomainName: d.config.Domain,
	}
	getDomainConfigResp, err := d.sdkClient.GetDomainConfig(context.TODO(), getDomainConfigReq)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'imagex.GetDomainConfig'")
	} else {
		d.logger.Logt("已获取域名配置", getDomainConfigResp)
	}

	// 更新 HTTPS 配置
	// REF: https://www.volcengine.com/docs/508/66012
	updateHttpsReq := &veImageX.UpdateHTTPSReq{
		UpdateHTTPSQuery: &veImageX.UpdateHTTPSQuery{
			ServiceID: d.config.ServiceId,
		},
		UpdateHTTPSBody: &veImageX.UpdateHTTPSBody{
			Domain: getDomainConfigResp.Result.Domain,
			HTTPS: &veImageX.UpdateHTTPSBodyHTTPS{
				CertID:              upres.CertId,
				EnableHTTP2:         getDomainConfigResp.Result.HTTPSConfig.EnableHTTP2,
				EnableHTTPS:         getDomainConfigResp.Result.HTTPSConfig.EnableHTTPS,
				EnableOcsp:          getDomainConfigResp.Result.HTTPSConfig.EnableOcsp,
				TLSVersions:         getDomainConfigResp.Result.HTTPSConfig.TLSVersions,
				EnableForceRedirect: getDomainConfigResp.Result.HTTPSConfig.EnableForceRedirect,
				ForceRedirectType:   getDomainConfigResp.Result.HTTPSConfig.ForceRedirectType,
				ForceRedirectCode:   getDomainConfigResp.Result.HTTPSConfig.ForceRedirectCode,
			},
		},
	}
	updateHttpsResp, err := d.sdkClient.UpdateHTTPS(context.TODO(), updateHttpsReq)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'imagex.UpdateHttps'")
	} else {
		d.logger.Logt("已更新 HTTPS 配置", updateHttpsResp)
	}

	return &deployer.DeployResult{}, nil
}

func createSdkClient(accessKeyId, accessKeySecret, region string) (*veImageX.Imagex, error) {
	var instance *veImageX.Imagex
	if region == "" {
		instance = veImageX.NewInstance()
	} else {
		instance = veImageX.NewInstanceWithRegion(region)
	}

	instance.SetCredential(veBase.Credentials{
		AccessKeyID:     accessKeyId,
		SecretAccessKey: accessKeySecret,
	})

	return instance, nil
}
