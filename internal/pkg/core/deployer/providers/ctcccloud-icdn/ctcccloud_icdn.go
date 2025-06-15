package ctcccloudicdn

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	uploadersp "github.com/usual2970/certimate/internal/pkg/core/uploader/providers/ctcccloud-icdn"
	ctyunicdn "github.com/usual2970/certimate/internal/pkg/sdk3rd/ctyun/icdn"
	xtypes "github.com/usual2970/certimate/internal/pkg/utils/types"
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
	sdkClient   *ctyunicdn.Client
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

	// 上传证书到 ICDN
	upres, err := d.sslUploader.Upload(ctx, certPEM, privkeyPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to upload certificate file: %w", err)
	} else {
		d.logger.Info("ssl certificate uploaded", slog.Any("result", upres))
	}

	// 查询域名配置信息
	// REF: https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=112&api=10849&data=173&isNormal=1&vid=166
	queryDomainDetailReq := &ctyunicdn.QueryDomainDetailRequest{
		Domain: xtypes.ToPtr(d.config.Domain),
	}
	queryDomainDetailResp, err := d.sdkClient.QueryDomainDetail(queryDomainDetailReq)
	d.logger.Debug("sdk request 'icdn.QueryDomainDetail'", slog.Any("request", queryDomainDetailReq), slog.Any("response", queryDomainDetailResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'icdn.QueryDomainDetail': %w", err)
	}

	// 修改域名配置
	// REF: https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=112&api=10853&data=173&isNormal=1&vid=166
	updateDomainReq := &ctyunicdn.UpdateDomainRequest{
		Domain:      xtypes.ToPtr(d.config.Domain),
		HttpsStatus: xtypes.ToPtr("on"),
		CertName:    xtypes.ToPtr(upres.CertName),
	}
	updateDomainResp, err := d.sdkClient.UpdateDomain(updateDomainReq)
	d.logger.Debug("sdk request 'icdn.UpdateDomain'", slog.Any("request", updateDomainReq), slog.Any("response", updateDomainResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'icdn.UpdateDomain': %w", err)
	}

	return &deployer.DeployResult{}, nil
}

func createSdkClient(accessKeyId, secretAccessKey string) (*ctyunicdn.Client, error) {
	return ctyunicdn.NewClient(accessKeyId, secretAccessKey)
}
