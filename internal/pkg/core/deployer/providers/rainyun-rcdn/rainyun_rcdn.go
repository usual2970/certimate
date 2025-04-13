package rainyunrcdn

import (
	"context"
	"errors"
	"log/slog"
	"strconv"

	xerrors "github.com/pkg/errors"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	uploadersp "github.com/usual2970/certimate/internal/pkg/core/uploader/providers/rainyun-sslcenter"
	rainyunsdk "github.com/usual2970/certimate/internal/pkg/vendors/rainyun-sdk"
)

type DeployerConfig struct {
	// 雨云 API 密钥。
	ApiKey string `json:"apiKey"`
	// RCDN 实例 ID。
	InstanceId int32 `json:"instanceId"`
	// 加速域名（支持泛域名）。
	Domain string `json:"domain"`
}

type DeployerProvider struct {
	config      *DeployerConfig
	logger      *slog.Logger
	sdkClient   *rainyunsdk.Client
	sslUploader uploader.Uploader
}

var _ deployer.Deployer = (*DeployerProvider)(nil)

func NewDeployer(config *DeployerConfig) (*DeployerProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	client, err := createSdkClient(config.ApiKey)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk client")
	}

	uploader, err := uploadersp.NewUploader(&uploadersp.UploaderConfig{
		ApiKey: config.ApiKey,
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
	// 上传证书到 SSL 证书
	upres, err := d.sslUploader.Upload(ctx, certPem, privkeyPem)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to upload certificate file")
	} else {
		d.logger.Info("ssl certificate uploaded", slog.Any("result", upres))
	}

	// RCDN SSL 绑定域名
	// REF: https://apifox.com/apidoc/shared/a4595cc8-44c5-4678-a2a3-eed7738dab03/api-184214120
	certId, _ := strconv.Atoi(upres.CertId)
	rcdnInstanceSslBindReq := &rainyunsdk.RcdnInstanceSslBindRequest{
		CertId:  int32(certId),
		Domains: []string{d.config.Domain},
	}
	rcdnInstanceSslBindResp, err := d.sdkClient.RcdnInstanceSslBind(d.config.InstanceId, rcdnInstanceSslBindReq)
	d.logger.Debug("sdk request 'rcdn.InstanceSslBind'", slog.Any("instanceId", d.config.InstanceId), slog.Any("request", rcdnInstanceSslBindReq), slog.Any("response", rcdnInstanceSslBindResp))
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'rcdn.InstanceSslBind'")
	}

	return &deployer.DeployResult{}, nil
}

func createSdkClient(apiKey string) (*rainyunsdk.Client, error) {
	if apiKey == "" {
		return nil, errors.New("invalid rainyun api key")
	}

	client := rainyunsdk.NewClient(apiKey)
	return client, nil
}
