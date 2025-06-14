package ctcccloudcdn

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	uploadersp "github.com/usual2970/certimate/internal/pkg/core/uploader/providers/ctcccloud-cdn"
	ctyuncdn "github.com/usual2970/certimate/internal/pkg/sdk3rd/ctyun/cdn"
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
	sdkClient   *ctyuncdn.Client
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

	// 上传证书到 CDN
	upres, err := d.sslUploader.Upload(ctx, certPEM, privkeyPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to upload certificate file: %w", err)
	} else {
		d.logger.Info("ssl certificate uploaded", slog.Any("result", upres))
	}

	// 查询域名配置信息
	// REF: https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=108&api=11304&data=161&isNormal=1&vid=154
	queryDomainDetailReq := &ctyuncdn.QueryDomainDetailRequest{
		Domain: typeutil.ToPtr(d.config.Domain),
	}
	queryDomainDetailResp, err := d.sdkClient.QueryDomainDetail(queryDomainDetailReq)
	d.logger.Debug("sdk request 'cdn.QueryDomainDetail'", slog.Any("request", queryDomainDetailReq), slog.Any("response", queryDomainDetailResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'cdn.QueryDomainDetail': %w", err)
	}

	// 修改域名配置
	// REF: https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=108&api=11308&data=161&isNormal=1&vid=154
	updateDomainReq := &ctyuncdn.UpdateDomainRequest{
		Domain:      typeutil.ToPtr(d.config.Domain),
		HttpsStatus: typeutil.ToPtr("on"),
		CertName:    typeutil.ToPtr(upres.CertName),
	}
	updateDomainResp, err := d.sdkClient.UpdateDomain(updateDomainReq)
	d.logger.Debug("sdk request 'cdn.UpdateDomain'", slog.Any("request", updateDomainReq), slog.Any("response", updateDomainResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'cdn.UpdateDomain': %w", err)
	}

	return &deployer.DeployResult{}, nil
}

func createSdkClient(accessKeyId, secretAccessKey string) (*ctyuncdn.Client, error) {
	return ctyuncdn.NewClient(accessKeyId, secretAccessKey)
}
