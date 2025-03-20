package upyuncdn

import (
	"context"
	"errors"
	"log/slog"

	xerrors "github.com/pkg/errors"
	"golang.org/x/exp/slices"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	uploadersp "github.com/usual2970/certimate/internal/pkg/core/uploader/providers/upyun-ssl"
	upyunsdk "github.com/usual2970/certimate/internal/pkg/vendors/upyun-sdk/console"
)

type DeployerConfig struct {
	// 又拍云账号用户名。
	Username string `json:"username"`
	// 又拍云账号密码。
	Password string `json:"password"`
	// 加速域名（支持泛域名）。
	Domain string `json:"domain"`
}

type DeployerProvider struct {
	config      *DeployerConfig
	logger      *slog.Logger
	sdkClient   *upyunsdk.Client
	sslUploader uploader.Uploader
}

var _ deployer.Deployer = (*DeployerProvider)(nil)

func NewDeployer(config *DeployerConfig) (*DeployerProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	client, err := createSdkClient(config.Username, config.Password)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk client")
	}

	uploader, err := uploadersp.NewUploader(&uploadersp.UploaderConfig{
		Username: config.Username,
		Password: config.Password,
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

	// 获取域名证书配置
	getHttpsServiceManagerResp, err := d.sdkClient.GetHttpsServiceManager(d.config.Domain)
	d.logger.Debug("sdk request 'console.GetHttpsServiceManager'", slog.String("request.domain", d.config.Domain), slog.Any("response", getHttpsServiceManagerResp))
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'console.GetHttpsServiceManager'")
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
			return nil, xerrors.Wrap(err, "failed to execute sdk request 'console.UpdateHttpsCertificateManager'")
		}
	} else if getHttpsServiceManagerResp.Data.Domains[lastCertIndex].CertificateId != upres.CertId {
		migrateHttpsDomainReq := &upyunsdk.MigrateHttpsDomainRequest{
			CertificateId: upres.CertId,
			Domain:        d.config.Domain,
		}
		migrateHttpsDomainResp, err := d.sdkClient.MigrateHttpsDomain(migrateHttpsDomainReq)
		d.logger.Debug("sdk request 'console.MigrateHttpsDomain'", slog.Any("request", migrateHttpsDomainReq), slog.Any("response", migrateHttpsDomainResp))
		if err != nil {
			return nil, xerrors.Wrap(err, "failed to execute sdk request 'console.MigrateHttpsDomain'")
		}
	}

	return &deployer.DeployResult{}, nil
}

func createSdkClient(username, password string) (*upyunsdk.Client, error) {
	if username == "" {
		return nil, errors.New("invalid upyun username")
	}

	if password == "" {
		return nil, errors.New("invalid upyun password")
	}

	client := upyunsdk.NewClient(username, password)
	return client, nil
}
