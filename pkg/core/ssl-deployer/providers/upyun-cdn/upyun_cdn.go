package upyuncdn

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"golang.org/x/exp/slices"

	"github.com/certimate-go/certimate/pkg/core"
	sslmgrsp "github.com/certimate-go/certimate/pkg/core/ssl-manager/providers/upyun-ssl"
	upyunsdk "github.com/certimate-go/certimate/pkg/sdk3rd/upyun/console"
)

type SSLDeployerProviderConfig struct {
	// 又拍云账号用户名。
	Username string `json:"username"`
	// 又拍云账号密码。
	Password string `json:"password"`
	// 加速域名（支持泛域名）。
	Domain string `json:"domain"`
}

type SSLDeployerProvider struct {
	config     *SSLDeployerProviderConfig
	logger     *slog.Logger
	sdkClient  *upyunsdk.Client
	sslManager core.SSLManager
}

var _ core.SSLDeployer = (*SSLDeployerProvider)(nil)

func NewSSLDeployerProvider(config *SSLDeployerProviderConfig) (*SSLDeployerProvider, error) {
	if config == nil {
		return nil, errors.New("the configuration of the ssl deployer provider is nil")
	}

	client, err := createSDKClient(config.Username, config.Password)
	if err != nil {
		return nil, fmt.Errorf("could not create sdk client: %w", err)
	}

	sslmgr, err := sslmgrsp.NewSSLManagerProvider(&sslmgrsp.SSLManagerProviderConfig{
		Username: config.Username,
		Password: config.Password,
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

	// 获取域名证书配置
	getHttpsServiceManagerResp, err := d.sdkClient.GetHttpsServiceManager(d.config.Domain)
	d.logger.Debug("sdk request 'console.GetHttpsServiceManager'", slog.String("request.domain", d.config.Domain), slog.Any("response", getHttpsServiceManagerResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'console.GetHttpsServiceManager': %w", err)
	}

	// 判断域名是否已启用 HTTPS。如果已启用，迁移域名证书；否则，设置新证书
	lastCertIndex := slices.IndexFunc(getHttpsServiceManagerResp.Data.Domains, func(item upyunsdk.HttpsServiceManagerDomain) bool {
		return item.Https
	})
	if lastCertIndex == -1 {
		updateHttpsCertificateManagerReq := &upyunsdk.UpdateHttpsCertificateManagerRequest{
			CertificateId: upres.CertId,
			Domain:        d.config.Domain,
			Https:         true,
			ForceHttps:    true,
		}
		updateHttpsCertificateManagerResp, err := d.sdkClient.UpdateHttpsCertificateManager(updateHttpsCertificateManagerReq)
		d.logger.Debug("sdk request 'console.EnableDomainHttps'", slog.Any("request", updateHttpsCertificateManagerReq), slog.Any("response", updateHttpsCertificateManagerResp))
		if err != nil {
			return nil, fmt.Errorf("failed to execute sdk request 'console.UpdateHttpsCertificateManager': %w", err)
		}
	} else if getHttpsServiceManagerResp.Data.Domains[lastCertIndex].CertificateId != upres.CertId {
		migrateHttpsDomainReq := &upyunsdk.MigrateHttpsDomainRequest{
			CertificateId: upres.CertId,
			Domain:        d.config.Domain,
		}
		migrateHttpsDomainResp, err := d.sdkClient.MigrateHttpsDomain(migrateHttpsDomainReq)
		d.logger.Debug("sdk request 'console.MigrateHttpsDomain'", slog.Any("request", migrateHttpsDomainReq), slog.Any("response", migrateHttpsDomainResp))
		if err != nil {
			return nil, fmt.Errorf("failed to execute sdk request 'console.MigrateHttpsDomain': %w", err)
		}
	}

	return &core.SSLDeployResult{}, nil
}

func createSDKClient(username, password string) (*upyunsdk.Client, error) {
	return upyunsdk.NewClient(username, password)
}
