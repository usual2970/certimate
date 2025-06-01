package jdcloudcdn

import (
	"context"
	"fmt"
	"log/slog"

	jdcore "github.com/jdcloud-api/jdcloud-sdk-go/core"
	jdcdnapi "github.com/jdcloud-api/jdcloud-sdk-go/services/cdn/apis"
	jdcdnclient "github.com/jdcloud-api/jdcloud-sdk-go/services/cdn/client"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	uploadersp "github.com/usual2970/certimate/internal/pkg/core/uploader/providers/jdcloud-ssl"
)

type DeployerConfig struct {
	// 京东云 AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// 京东云 AccessKeySecret。
	AccessKeySecret string `json:"accessKeySecret"`
	// 加速域名（支持泛域名）。
	Domain string `json:"domain"`
}

type DeployerProvider struct {
	config      *DeployerConfig
	logger      *slog.Logger
	sdkClient   *jdcdnclient.CdnClient
	sslUploader uploader.Uploader
}

var _ deployer.Deployer = (*DeployerProvider)(nil)

func NewDeployer(config *DeployerConfig) (*DeployerProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	client, err := createSdkClient(config.AccessKeyId, config.AccessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("failed to create sdk client: %w", err)
	}

	uploader, err := uploadersp.NewUploader(&uploadersp.UploaderConfig{
		AccessKeyId:     config.AccessKeyId,
		AccessKeySecret: config.AccessKeySecret,
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
	// 查询域名配置信息
	// REF: https://docs.jdcloud.com/cn/cdn/api/querydomainconfig
	queryDomainConfigReq := jdcdnapi.NewQueryDomainConfigRequest(d.config.Domain)
	queryDomainConfigResp, err := d.sdkClient.QueryDomainConfig(queryDomainConfigReq)
	d.logger.Debug("sdk request 'cdn.QueryDomainConfig'", slog.Any("request", queryDomainConfigReq), slog.Any("response", queryDomainConfigResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'cdn.QueryDomainConfig': %w", err)
	}

	// 上传证书到 SSL
	upres, err := d.sslUploader.Upload(ctx, certPEM, privkeyPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to upload certificate file: %w", err)
	} else {
		d.logger.Info("ssl certificate uploaded", slog.Any("result", upres))
	}

	// 设置通讯协议
	// REF: https://docs.jdcloud.com/cn/cdn/api/sethttptype
	setHttpTypeReq := jdcdnapi.NewSetHttpTypeRequest(d.config.Domain)
	setHttpTypeReq.SetHttpType("https")
	setHttpTypeReq.SetCertificate(certPEM)
	setHttpTypeReq.SetRsaKey(privkeyPEM)
	setHttpTypeReq.SetCertFrom("ssl")
	setHttpTypeReq.SetSslCertId(upres.CertId)
	setHttpTypeReq.SetJumpType(queryDomainConfigResp.Result.HttpsJumpType)
	setHttpTypeResp, err := d.sdkClient.SetHttpType(setHttpTypeReq)
	d.logger.Debug("sdk request 'cdn.QueryDomainConfig'", slog.Any("request", setHttpTypeReq), slog.Any("response", setHttpTypeResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'cdn.SetHttpType': %w", err)
	}

	return &deployer.DeployResult{}, nil
}

func createSdkClient(accessKeyId, accessKeySecret string) (*jdcdnclient.CdnClient, error) {
	clientCredentials := jdcore.NewCredentials(accessKeyId, accessKeySecret)
	client := jdcdnclient.NewCdnClient(clientCredentials)
	client.SetLogger(jdcore.NewDefaultLogger(jdcore.LogWarn))
	return client, nil
}
