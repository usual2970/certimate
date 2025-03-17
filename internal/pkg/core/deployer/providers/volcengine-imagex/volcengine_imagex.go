package volcengineimagex

import (
	"context"
	"errors"
	"log/slog"

	xerrors "github.com/pkg/errors"
	veBase "github.com/volcengine/volc-sdk-golang/base"
	veImageX "github.com/volcengine/volc-sdk-golang/service/imagex/v2"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	uploadersp "github.com/usual2970/certimate/internal/pkg/core/uploader/providers/volcengine-certcenter"
)

type DeployerConfig struct {
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

type DeployerProvider struct {
	config      *DeployerConfig
	logger      *slog.Logger
	sdkClient   *veImageX.Imagex
	sslUploader uploader.Uploader
}

var _ deployer.Deployer = (*DeployerProvider)(nil)

func NewDeployer(config *DeployerConfig) (*DeployerProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	client, err := createSdkClient(config.AccessKeyId, config.AccessKeySecret, config.Region)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk client")
	}

	uploader, err := uploadersp.NewUploader(&uploadersp.UploaderConfig{
		AccessKeyId:     config.AccessKeyId,
		AccessKeySecret: config.AccessKeySecret,
		Region:          config.Region,
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
		d.logger.Info("ssl certificate uploaded", slog.Any("result", upres))
	}

	// 获取域名配置
	// REF: https://www.volcengine.com/docs/508/9366
	getDomainConfigReq := &veImageX.GetDomainConfigQuery{
		ServiceID:  d.config.ServiceId,
		DomainName: d.config.Domain,
	}
	getDomainConfigResp, err := d.sdkClient.GetDomainConfig(context.TODO(), getDomainConfigReq)
	d.logger.Debug("sdk request 'imagex.GetDomainConfig'", slog.Any("request", getDomainConfigReq), slog.Any("response", getDomainConfigResp))
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'imagex.GetDomainConfig'")
	}

	// 更新 HTTPS 配置
	// REF: https://www.volcengine.com/docs/508/66012
	updateHttpsReq := &veImageX.UpdateHTTPSReq{
		UpdateHTTPSQuery: &veImageX.UpdateHTTPSQuery{
			ServiceID: d.config.ServiceId,
		},
		UpdateHTTPSBody: &veImageX.UpdateHTTPSBody{
			Domain: d.config.Domain,
			HTTPS: &veImageX.UpdateHTTPSBodyHTTPS{
				CertID:      upres.CertId,
				EnableHTTPS: true,
			},
		},
	}
	if getDomainConfigResp.Result != nil && getDomainConfigResp.Result.HTTPSConfig != nil {
		updateHttpsReq.UpdateHTTPSBody.HTTPS.EnableHTTPS = getDomainConfigResp.Result.HTTPSConfig.EnableHTTPS
		updateHttpsReq.UpdateHTTPSBody.HTTPS.EnableHTTP2 = getDomainConfigResp.Result.HTTPSConfig.EnableHTTP2
		updateHttpsReq.UpdateHTTPSBody.HTTPS.EnableOcsp = getDomainConfigResp.Result.HTTPSConfig.EnableOcsp
		updateHttpsReq.UpdateHTTPSBody.HTTPS.TLSVersions = getDomainConfigResp.Result.HTTPSConfig.TLSVersions
		updateHttpsReq.UpdateHTTPSBody.HTTPS.EnableForceRedirect = getDomainConfigResp.Result.HTTPSConfig.EnableForceRedirect
		updateHttpsReq.UpdateHTTPSBody.HTTPS.ForceRedirectType = getDomainConfigResp.Result.HTTPSConfig.ForceRedirectType
		updateHttpsReq.UpdateHTTPSBody.HTTPS.ForceRedirectCode = getDomainConfigResp.Result.HTTPSConfig.ForceRedirectCode
	}
	updateHttpsResp, err := d.sdkClient.UpdateHTTPS(context.TODO(), updateHttpsReq)
	d.logger.Debug("sdk request 'imagex.UpdateHttps'", slog.Any("request", updateHttpsReq), slog.Any("response", updateHttpsResp))
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'imagex.UpdateHttps'")
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
