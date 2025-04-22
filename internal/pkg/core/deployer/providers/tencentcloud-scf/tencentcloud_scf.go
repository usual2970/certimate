package tencentcloudscf

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tcscf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/scf/v20180416"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	uploadersp "github.com/usual2970/certimate/internal/pkg/core/uploader/providers/tencentcloud-ssl"
)

type DeployerConfig struct {
	// 腾讯云 SecretId。
	SecretId string `json:"secretId"`
	// 腾讯云 SecretKey。
	SecretKey string `json:"secretKey"`
	// 腾讯云地域。
	Region string `json:"region"`
	// 自定义域名（不支持泛域名）。
	Domain string `json:"domain"`
}

type DeployerProvider struct {
	config      *DeployerConfig
	logger      *slog.Logger
	sdkClient   *tcscf.Client
	sslUploader uploader.Uploader
}

var _ deployer.Deployer = (*DeployerProvider)(nil)

func NewDeployer(config *DeployerConfig) (*DeployerProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	client, err := createSdkClient(config.SecretId, config.SecretKey, config.Region)
	if err != nil {
		return nil, fmt.Errorf("failed to create sdk client: %w", err)
	}

	uploader, err := uploadersp.NewUploader(&uploadersp.UploaderConfig{
		SecretId:  config.SecretId,
		SecretKey: config.SecretKey,
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
	// 查看云函数自定义域名详情
	// REF: https://cloud.tencent.com/document/product/583/111924
	getCustomDomainReq := tcscf.NewGetCustomDomainRequest()
	getCustomDomainReq.Domain = common.StringPtr(d.config.Domain)
	getCustomDomainResp, err := d.sdkClient.GetCustomDomain(getCustomDomainReq)
	d.logger.Debug("sdk request 'scf.GetCustomDomain'", slog.Any("request", getCustomDomainReq), slog.Any("response", getCustomDomainResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'scf.GetCustomDomain': %w", err)
	}

	// 上传证书到 SSL
	upres, err := d.sslUploader.Upload(ctx, certPEM, privkeyPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to upload certificate file: %w", err)
	} else {
		d.logger.Info("ssl certificate uploaded", slog.Any("result", upres))
	}

	// 更新云函数自定义域名
	// REF: https://cloud.tencent.com/document/product/583/111922
	updateCustomDomainReq := tcscf.NewUpdateCustomDomainRequest()
	updateCustomDomainReq.Domain = common.StringPtr(d.config.Domain)
	updateCustomDomainReq.CertConfig = &tcscf.CertConf{
		CertificateId: common.StringPtr(upres.CertId),
	}
	updateCustomDomainReq.Protocol = getCustomDomainResp.Response.Protocol
	updateCustomDomainResp, err := d.sdkClient.UpdateCustomDomain(updateCustomDomainReq)
	d.logger.Debug("sdk request 'scf.UpdateCustomDomain'", slog.Any("request", updateCustomDomainReq), slog.Any("response", updateCustomDomainResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'scf.UpdateCustomDomain': %w", err)
	}

	return &deployer.DeployResult{}, nil
}

func createSdkClient(secretId, secretKey, region string) (*tcscf.Client, error) {
	credential := common.NewCredential(secretId, secretKey)
	client, err := tcscf.NewClient(credential, region, profile.NewClientProfile())
	if err != nil {
		return nil, err
	}

	return client, nil
}
