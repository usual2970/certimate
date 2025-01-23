package tencentcloudcss

import (
	"context"
	"errors"

	xerrors "github.com/pkg/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tcLive "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/logger"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	uploaderp "github.com/usual2970/certimate/internal/pkg/core/uploader/providers/tencentcloud-ssl"
)

type TencentCloudCSSDeployerConfig struct {
	// 腾讯云 SecretId。
	SecretId string `json:"secretId"`
	// 腾讯云 SecretKey。
	SecretKey string `json:"secretKey"`
	// 直播播放域名（不支持泛域名）。
	Domain string `json:"domain"`
}

type TencentCloudCSSDeployer struct {
	config      *TencentCloudCSSDeployerConfig
	logger      logger.Logger
	sdkClient   *tcLive.Client
	sslUploader uploader.Uploader
}

var _ deployer.Deployer = (*TencentCloudCSSDeployer)(nil)

func New(config *TencentCloudCSSDeployerConfig) (*TencentCloudCSSDeployer, error) {
	return NewWithLogger(config, logger.NewNilLogger())
}

func NewWithLogger(config *TencentCloudCSSDeployerConfig, logger logger.Logger) (*TencentCloudCSSDeployer, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}

	if logger == nil {
		return nil, errors.New("logger is nil")
	}

	client, err := createSdkClient(config.SecretId, config.SecretKey)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk client")
	}

	uploader, err := uploaderp.New(&uploaderp.TencentCloudSSLUploaderConfig{
		SecretId:  config.SecretId,
		SecretKey: config.SecretKey,
	})
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create ssl uploader")
	}

	return &TencentCloudCSSDeployer{
		logger:      logger,
		config:      config,
		sdkClient:   client,
		sslUploader: uploader,
	}, nil
}

func (d *TencentCloudCSSDeployer) Deploy(ctx context.Context, certPem string, privkeyPem string) (*deployer.DeployResult, error) {
	// 上传证书到 SSL
	upres, err := d.sslUploader.Upload(ctx, certPem, privkeyPem)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to upload certificate file")
	}

	d.logger.Logt("certificate file uploaded", upres)

	// 绑定证书对应的播放域名
	// REF: https://cloud.tencent.com/document/product/267/78655
	modifyLiveDomainCertBindingsReq := &tcLive.ModifyLiveDomainCertBindingsRequest{
		DomainInfos: []*tcLive.LiveCertDomainInfo{
			{
				DomainName: common.StringPtr(d.config.Domain),
				Status:     common.Int64Ptr(1),
			},
		},
		CloudCertId: common.StringPtr(upres.CertId),
	}
	modifyLiveDomainCertBindingsResp, err := d.sdkClient.ModifyLiveDomainCertBindings(modifyLiveDomainCertBindingsReq)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'live.ModifyLiveDomainCertBindings'")
	}

	d.logger.Logt("已部署证书到云资源实例", modifyLiveDomainCertBindingsResp.Response)

	return &deployer.DeployResult{}, nil
}

func createSdkClient(secretId, secretKey string) (*tcLive.Client, error) {
	credential := common.NewCredential(secretId, secretKey)

	client, err := tcLive.NewClient(credential, "", profile.NewClientProfile())
	if err != nil {
		return nil, err
	}

	return client, nil
}
