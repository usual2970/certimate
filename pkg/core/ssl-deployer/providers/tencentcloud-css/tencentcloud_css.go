package tencentcloudcss

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tclive "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"

	"github.com/certimate-go/certimate/pkg/core"
	sslmgrsp "github.com/certimate-go/certimate/pkg/core/ssl-manager/providers/tencentcloud-ssl"
)

type SSLDeployerProviderConfig struct {
	// 腾讯云 SecretId。
	SecretId string `json:"secretId"`
	// 腾讯云 SecretKey。
	SecretKey string `json:"secretKey"`
	// 直播播放域名（不支持泛域名）。
	Domain string `json:"domain"`
}

type SSLDeployerProvider struct {
	config     *SSLDeployerProviderConfig
	logger     *slog.Logger
	sdkClient  *tclive.Client
	sslManager core.SSLManager
}

var _ core.SSLDeployer = (*SSLDeployerProvider)(nil)

func NewSSLDeployerProvider(config *SSLDeployerProviderConfig) (*SSLDeployerProvider, error) {
	if config == nil {
		return nil, errors.New("the configuration of the ssl deployer provider is nil")
	}

	client, err := createSDKClient(config.SecretId, config.SecretKey)
	if err != nil {
		return nil, fmt.Errorf("could not create sdk client: %w", err)
	}

	sslmgr, err := sslmgrsp.NewSSLManagerProvider(&sslmgrsp.SSLManagerProviderConfig{
		SecretId:  config.SecretId,
		SecretKey: config.SecretKey,
	})
	if err != nil {
		return nil, fmt.Errorf("could not create ssl manager: %w", err)
	}

	return &SSLDeployerProvider{
		config:     config,
		logger:     slog.Default(),
		sdkClient:  client,
		sslManager: sslmgr,
	}, nil
}

func (d *SSLDeployerProvider) SetLogger(logger *slog.Logger) {
	if logger == nil {
		d.logger = slog.New(slog.DiscardHandler)
	} else {
		d.logger = logger
	}

	d.sslManager.SetLogger(logger)
}

func (d *SSLDeployerProvider) Deploy(ctx context.Context, certPEM string, privkeyPEM string) (*core.SSLDeployResult, error) {
	if d.config.Domain == "" {
		return nil, errors.New("config `domain` is required")
	}

	// 上传证书
	upres, err := d.sslManager.Upload(ctx, certPEM, privkeyPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to upload certificate file: %w", err)
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
		return nil, fmt.Errorf("failed to execute sdk request 'live.ModifyLiveDomainCertBindings': %w", err)
	}

	return &core.SSLDeployResult{}, nil
}

func createSDKClient(secretId, secretKey string) (*tclive.Client, error) {
	credential := common.NewCredential(secretId, secretKey)

	client, err := tclive.NewClient(credential, "", profile.NewClientProfile())
	if err != nil {
		return nil, err
	}

	return client, nil
}
