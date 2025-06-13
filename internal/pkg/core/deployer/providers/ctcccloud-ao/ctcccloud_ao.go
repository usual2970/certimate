package ctcccloudao

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	uploadersp "github.com/usual2970/certimate/internal/pkg/core/uploader/providers/ctcccloud-ao"
	ctyunao "github.com/usual2970/certimate/internal/pkg/sdk3rd/ctyun/ao"
	typeutil "github.com/usual2970/certimate/internal/pkg/utils/type"
)

type DeployerConfig struct {
	// 天翼云 AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// 天翼云 SecretAccessKey。
	SecretAccessKey string `json:"secretAccessKey"`
	// 加速域名（支持泛域名）。
	Domain string `json:"domain"`
}

type DeployerProvider struct {
	config      *DeployerConfig
	logger      *slog.Logger
	sdkClient   *ctyunao.Client
	sslUploader uploader.Uploader
}

var _ deployer.Deployer = (*DeployerProvider)(nil)

func NewDeployer(config *DeployerConfig) (*DeployerProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	client, err := createSdkClient(config.AccessKeyId, config.SecretAccessKey)
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
		d.logger = slog.New(slog.DiscardHandler)
	} else {
		d.logger = logger
	}
	return d
}

func (d *DeployerProvider) Deploy(ctx context.Context, certPEM string, privkeyPEM string) (*deployer.DeployResult, error) {
	if d.config.Domain == "" {
		return nil, errors.New("config `domain` is required")
	}

	// 上传证书到 AccessOne
	upres, err := d.sslUploader.Upload(ctx, certPEM, privkeyPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to upload certificate file: %w", err)
	} else {
		d.logger.Info("ssl certificate uploaded", slog.Any("result", upres))
	}

	// 域名基础及加速配置查询
	// REF: https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=113&api=13412&data=174&isNormal=1&vid=167
	getDomainConfigReq := &ctyunao.GetDomainConfigRequest{
		Domain: typeutil.ToPtr(d.config.Domain),
	}
	getDomainConfigResp, err := d.sdkClient.GetDomainConfig(getDomainConfigReq)
	d.logger.Debug("sdk request 'cdn.GetDomainConfig'", slog.Any("request", getDomainConfigReq), slog.Any("response", getDomainConfigResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'cdn.GetDomainConfig': %w", err)
	}

	// 域名基础及加速配置修改
	// REF: https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=113&api=13413&data=174&isNormal=1&vid=167
	modifyDomainConfigReq := &ctyunao.ModifyDomainConfigRequest{
		Domain:      typeutil.ToPtr(d.config.Domain),
		ProductCode: typeutil.ToPtr(getDomainConfigResp.ReturnObj.ProductCode),
		Origin:      getDomainConfigResp.ReturnObj.Origin,
		HttpsStatus: typeutil.ToPtr("on"),
		CertName:    typeutil.ToPtr(upres.CertName),
	}
	modifyDomainConfigResp, err := d.sdkClient.ModifyDomainConfig(modifyDomainConfigReq)
	d.logger.Debug("sdk request 'cdn.ModifyDomainConfig'", slog.Any("request", modifyDomainConfigReq), slog.Any("response", modifyDomainConfigResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'cdn.ModifyDomainConfig': %w", err)
	}

	return &deployer.DeployResult{}, nil
}

func createSdkClient(accessKeyId, secretAccessKey string) (*ctyunao.Client, error) {
	return ctyunao.NewClient(accessKeyId, secretAccessKey)
}
