package jdcloudcdn

import (
	"context"
	"log/slog"

	jdCore "github.com/jdcloud-api/jdcloud-sdk-go/core"
	jdCdnApi "github.com/jdcloud-api/jdcloud-sdk-go/services/cdn/apis"
	jdCdnClient "github.com/jdcloud-api/jdcloud-sdk-go/services/cdn/client"
	xerrors "github.com/pkg/errors"

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
	sdkClient   *jdCdnClient.CdnClient
	sslUploader uploader.Uploader
}

var _ deployer.Deployer = (*DeployerProvider)(nil)

func NewDeployer(config *DeployerConfig) (*DeployerProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	client, err := createSdkClient(config.AccessKeyId, config.AccessKeySecret)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk client")
	}

	uploader, err := uploadersp.NewUploader(&uploadersp.UploaderConfig{
		AccessKeyId:     config.AccessKeyId,
		AccessKeySecret: config.AccessKeySecret,
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
	// 查询域名配置信息
	// REF: https://docs.jdcloud.com/cn/cdn/api/querydomainconfig
	queryDomainConfigReq := jdCdnApi.NewQueryDomainConfigRequest(d.config.Domain)
	queryDomainConfigResp, err := d.sdkClient.QueryDomainConfig(queryDomainConfigReq)
	d.logger.Debug("sdk request 'cdn.QueryDomainConfig'", slog.Any("request", queryDomainConfigReq), slog.Any("response", queryDomainConfigResp))
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'cdn.QueryDomainConfig'")
	}

	// 上传证书到 SSL
	upres, err := d.sslUploader.Upload(ctx, certPem, privkeyPem)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to upload certificate file")
	} else {
		d.logger.Info("ssl certificate uploaded", slog.Any("result", upres))
	}

	// 设置通讯协议
	// REF: https://docs.jdcloud.com/cn/cdn/api/sethttptype
	setHttpTypeReq := jdCdnApi.NewSetHttpTypeRequest(d.config.Domain)
	setHttpTypeReq.SetHttpType("https")
	setHttpTypeReq.SetCertificate(certPem)
	setHttpTypeReq.SetRsaKey(privkeyPem)
	setHttpTypeReq.SetCertFrom("ssl")
	setHttpTypeReq.SetSslCertId(upres.CertId)
	setHttpTypeReq.SetJumpType(queryDomainConfigResp.Result.HttpsJumpType)
	setHttpTypeResp, err := d.sdkClient.SetHttpType(setHttpTypeReq)
	d.logger.Debug("sdk request 'cdn.QueryDomainConfig'", slog.Any("request", setHttpTypeReq), slog.Any("response", setHttpTypeResp))
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'cdn.SetHttpType'")
	}

	return &deployer.DeployResult{}, nil
}

func createSdkClient(accessKeyId, accessKeySecret string) (*jdCdnClient.CdnClient, error) {
	clientCredentials := jdCore.NewCredentials(accessKeyId, accessKeySecret)
	client := jdCdnClient.NewCdnClient(clientCredentials)
	client.SetLogger(jdCore.NewDefaultLogger(jdCore.LogWarn))
	return client, nil
}
