package tencentcloudscf

import (
	"context"
	"log/slog"

	xerrors "github.com/pkg/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tcScf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/scf/v20180416"

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
	sdkClient   *tcScf.Client
	sslUploader uploader.Uploader
}

var _ deployer.Deployer = (*DeployerProvider)(nil)

func NewDeployer(config *DeployerConfig) (*DeployerProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	client, err := createSdkClient(config.SecretId, config.SecretKey, config.Region)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk client")
	}

	uploader, err := uploadersp.NewUploader(&uploadersp.UploaderConfig{
		SecretId:  config.SecretId,
		SecretKey: config.SecretKey,
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
	// 查看云函数自定义域名详情
	// REF: https://cloud.tencent.com/document/product/583/111924
	getCustomDomainReq := tcScf.NewGetCustomDomainRequest()
	getCustomDomainReq.Domain = common.StringPtr(d.config.Domain)
	getCustomDomainResp, err := d.sdkClient.GetCustomDomain(getCustomDomainReq)
	d.logger.Debug("sdk request 'scf.GetCustomDomain'", slog.Any("request", getCustomDomainReq), slog.Any("response", getCustomDomainResp))
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'scf.GetCustomDomain'")
	}

	// 上传证书到 SSL
	upres, err := d.sslUploader.Upload(ctx, certPem, privkeyPem)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to upload certificate file")
	} else {
		d.logger.Info("ssl certificate uploaded", slog.Any("result", upres))
	}

	// 更新云函数自定义域名
	// REF: https://cloud.tencent.com/document/product/583/111922
	updateCustomDomainReq := tcScf.NewUpdateCustomDomainRequest()
	updateCustomDomainReq.Domain = common.StringPtr(d.config.Domain)
	updateCustomDomainReq.CertConfig = &tcScf.CertConf{
		CertificateId: common.StringPtr(upres.CertId),
	}
	updateCustomDomainReq.Protocol = getCustomDomainResp.Response.Protocol
	updateCustomDomainResp, err := d.sdkClient.UpdateCustomDomain(updateCustomDomainReq)
	d.logger.Debug("sdk request 'scf.UpdateCustomDomain'", slog.Any("request", updateCustomDomainReq), slog.Any("response", updateCustomDomainResp))
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'scf.UpdateCustomDomain'")
	}

	return &deployer.DeployResult{}, nil
}

func createSdkClient(secretId, secretKey, region string) (*tcScf.Client, error) {
	credential := common.NewCredential(secretId, secretKey)
	client, err := tcScf.NewClient(credential, region, profile.NewClientProfile())
	if err != nil {
		return nil, err
	}

	return client, nil
}
