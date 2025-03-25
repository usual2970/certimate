package tencentcloudcss

import (
	"context"
	"log/slog"

	xerrors "github.com/pkg/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tclive "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	uploadersp "github.com/usual2970/certimate/internal/pkg/core/uploader/providers/tencentcloud-ssl"
)

type DeployerConfig struct {
	// 腾讯云 SecretId。
	SecretId string `json:"secretId"`
	// 腾讯云 SecretKey。
	SecretKey string `json:"secretKey"`
	// 直播播放域名（不支持泛域名）。
	Domain string `json:"domain"`
}

type DeployerProvider struct {
	config      *DeployerConfig
	logger      *slog.Logger
	sdkClient   *tclive.Client
	sslUploader uploader.Uploader
}

var _ deployer.Deployer = (*DeployerProvider)(nil)

func NewDeployer(config *DeployerConfig) (*DeployerProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	client, err := createSdkClient(config.SecretId, config.SecretKey)
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
	// 上传证书到 SSL
	upres, err := d.sslUploader.Upload(ctx, certPem, privkeyPem)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to upload certificate file")
	} else {
		d.logger.Info("ssl certificate uploaded", slog.Any("result", upres))
	}

	// 绑定证书对应的播放域名
	// REF: https://cloud.tencent.com/document/product/267/78655
	modifyLiveDomainCertBindingsReq := tclive.NewModifyLiveDomainCertBindingsRequest()
	modifyLiveDomainCertBindingsReq.DomainInfos = []*tclive.LiveCertDomainInfo{
		{
			DomainName: common.StringPtr(d.config.Domain),
			Status:     common.Int64Ptr(1),
		},
	}
	modifyLiveDomainCertBindingsReq.CloudCertId = common.StringPtr(upres.CertId)
	modifyLiveDomainCertBindingsResp, err := d.sdkClient.ModifyLiveDomainCertBindings(modifyLiveDomainCertBindingsReq)
	d.logger.Debug("sdk request 'live.ModifyLiveDomainCertBindings'", slog.Any("request", modifyLiveDomainCertBindingsReq), slog.Any("response", modifyLiveDomainCertBindingsResp))
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'live.ModifyLiveDomainCertBindings'")
	}

	return &deployer.DeployResult{}, nil
}

func createSdkClient(secretId, secretKey string) (*tclive.Client, error) {
	credential := common.NewCredential(secretId, secretKey)

	client, err := tclive.NewClient(credential, "", profile.NewClientProfile())
	if err != nil {
		return nil, err
	}

	return client, nil
}
