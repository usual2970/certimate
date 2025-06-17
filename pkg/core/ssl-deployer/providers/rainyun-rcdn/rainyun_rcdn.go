package rainyunrcdn

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/certimate-go/certimate/pkg/core"
	sslmgrsp "github.com/certimate-go/certimate/pkg/core/ssl-manager/providers/rainyun-sslcenter"
	rainyunsdk "github.com/certimate-go/certimate/pkg/sdk3rd/rainyun"
)

type SSLDeployerProviderConfig struct {
	// 雨云 API 密钥。
	ApiKey string `json:"apiKey"`
	// RCDN 实例 ID。
	InstanceId int32 `json:"instanceId"`
	// 加速域名（支持泛域名）。
	Domain string `json:"domain"`
}

type SSLDeployerProvider struct {
	config     *SSLDeployerProviderConfig
	logger     *slog.Logger
	sdkClient  *rainyunsdk.Client
	sslManager core.SSLManager
}

var _ core.SSLDeployer = (*SSLDeployerProvider)(nil)

func NewSSLDeployerProvider(config *SSLDeployerProviderConfig) (*SSLDeployerProvider, error) {
	if config == nil {
		return nil, errors.New("the configuration of the ssl deployer provider is nil")
	}

	client, err := createSDKClient(config.ApiKey)
	if err != nil {
		return nil, fmt.Errorf("could not create sdk client: %w", err)
	}

	sslmgr, err := sslmgrsp.NewSSLManagerProvider(&sslmgrsp.SSLManagerProviderConfig{
		ApiKey: config.ApiKey,
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
		return nil, fmt.Errorf("config `domain` is required")
	}

	// 上传证书
	upres, err := d.sslManager.Upload(ctx, certPEM, privkeyPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to upload certificate file: %w", err)
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
		return nil, fmt.Errorf("failed to execute sdk request 'rcdn.InstanceSslBind': %w", err)
	}

	return &core.SSLDeployResult{}, nil
}

func createSDKClient(apiKey string) (*rainyunsdk.Client, error) {
	return rainyunsdk.NewClient(apiKey)
}
