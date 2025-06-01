package volcengineimagex

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	vebase "github.com/volcengine/volc-sdk-golang/base"
	veimagex "github.com/volcengine/volc-sdk-golang/service/imagex/v2"

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
	sdkClient   *veimagex.Imagex
	sslUploader uploader.Uploader
}

var _ deployer.Deployer = (*DeployerProvider)(nil)

func NewDeployer(config *DeployerConfig) (*DeployerProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	client, err := createSdkClient(config.AccessKeyId, config.AccessKeySecret, config.Region)
	if err != nil {
		return nil, fmt.Errorf("failed to create sdk client: %w", err)
	}

	uploader, err := uploadersp.NewUploader(&uploadersp.UploaderConfig{
		AccessKeyId:     config.AccessKeyId,
		AccessKeySecret: config.AccessKeySecret,
		Region:          config.Region,
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
		d.logger = slog.New(slog.DiscardHandler)
	} else {
		d.logger = logger
	}
	d.sslUploader.WithLogger(logger)
	return d
}

func (d *DeployerProvider) Deploy(ctx context.Context, certPEM string, privkeyPEM string) (*deployer.DeployResult, error) {
	if d.config.ServiceId == "" {
		return nil, errors.New("config `serviceId` is required")
	}
	if d.config.Domain == "" {
		return nil, errors.New("config `domain` is required")
	}

	// 上传证书到证书中心
	upres, err := d.sslUploader.Upload(ctx, certPEM, privkeyPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to upload certificate file: %w", err)
	} else {
		d.logger.Info("ssl certificate uploaded", slog.Any("result", upres))
	}

	// 获取域名配置
	// REF: https://www.volcengine.com/docs/508/9366
	getDomainConfigReq := &veimagex.GetDomainConfigQuery{
		ServiceID:  d.config.ServiceId,
		DomainName: d.config.Domain,
	}
	getDomainConfigResp, err := d.sdkClient.GetDomainConfig(context.TODO(), getDomainConfigReq)
	d.logger.Debug("sdk request 'imagex.GetDomainConfig'", slog.Any("request", getDomainConfigReq), slog.Any("response", getDomainConfigResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'imagex.GetDomainConfig': %w", err)
	}

	// 更新 HTTPS 配置
	// REF: https://www.volcengine.com/docs/508/66012
	updateHttpsReq := &veimagex.UpdateHTTPSReq{
		UpdateHTTPSQuery: &veimagex.UpdateHTTPSQuery{
			ServiceID: d.config.ServiceId,
		},
		UpdateHTTPSBody: &veimagex.UpdateHTTPSBody{
			Domain: d.config.Domain,
			HTTPS: &veimagex.UpdateHTTPSBodyHTTPS{
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
		return nil, fmt.Errorf("failed to execute sdk request 'imagex.UpdateHttps': %w", err)
	}

	return &deployer.DeployResult{}, nil
}

func createSdkClient(accessKeyId, accessKeySecret, region string) (*veimagex.Imagex, error) {
	var instance *veimagex.Imagex
	if region == "" {
		instance = veimagex.NewInstance()
	} else {
		instance = veimagex.NewInstanceWithRegion(region)
	}

	instance.SetCredential(vebase.Credentials{
		AccessKeyID:     accessKeyId,
		SecretAccessKey: accessKeySecret,
	})

	return instance, nil
}
