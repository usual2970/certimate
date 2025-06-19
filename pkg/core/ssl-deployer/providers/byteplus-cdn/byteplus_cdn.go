package bytepluscdn

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	bpcdn "github.com/byteplus-sdk/byteplus-sdk-golang/service/cdn"

	"github.com/certimate-go/certimate/pkg/core"
	sslmgrsp "github.com/certimate-go/certimate/pkg/core/ssl-manager/providers/byteplus-cdn"
)

type SSLDeployerProviderConfig struct {
	// BytePlus AccessKey。
	AccessKey string `json:"accessKey"`
	// BytePlus SecretKey。
	SecretKey string `json:"secretKey"`
	// 加速域名（支持泛域名）。
	Domain string `json:"domain"`
}

type SSLDeployerProvider struct {
	config     *SSLDeployerProviderConfig
	logger     *slog.Logger
	sdkClient  *bpcdn.CDN
	sslManager core.SSLManager
}

var _ core.SSLDeployer = (*SSLDeployerProvider)(nil)

func NewSSLDeployerProvider(config *SSLDeployerProviderConfig) (*SSLDeployerProvider, error) {
	if config == nil {
		return nil, errors.New("the configuration of the ssl deployer provider is nil")
	}

	client := bpcdn.NewInstance()
	client.Client.SetAccessKey(config.AccessKey)
	client.Client.SetSecretKey(config.SecretKey)

	sslmgr, err := sslmgrsp.NewSSLManagerProvider(&sslmgrsp.SSLManagerProviderConfig{
		AccessKey: config.AccessKey,
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
	// 上传证书
	upres, err := d.sslManager.Upload(ctx, certPEM, privkeyPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to upload certificate file: %w", err)
	} else {
		d.logger.Info("ssl certificate uploaded", slog.Any("result", upres))
	}

	domains := make([]string, 0)
	if strings.HasPrefix(d.config.Domain, "*.") {
		// 获取指定证书可关联的域名
		// REF: https://docs.byteplus.com/en/docs/byteplus-cdn/reference-describecertconfig-9ea17
		describeCertConfigReq := &bpcdn.DescribeCertConfigRequest{
			CertId: upres.CertId,
		}
		describeCertConfigResp, err := d.sdkClient.DescribeCertConfig(describeCertConfigReq)
		d.logger.Debug("sdk request 'cdn.DescribeCertConfig'", slog.Any("request", describeCertConfigReq), slog.Any("response", describeCertConfigResp))
		if err != nil {
			return nil, fmt.Errorf("failed to execute sdk request 'cdn.DescribeCertConfig': %w", err)
		}

		if describeCertConfigResp.Result.CertNotConfig != nil {
			for i := range describeCertConfigResp.Result.CertNotConfig {
				domains = append(domains, describeCertConfigResp.Result.CertNotConfig[i].Domain)
			}
		}

		if describeCertConfigResp.Result.OtherCertConfig != nil {
			for i := range describeCertConfigResp.Result.OtherCertConfig {
				domains = append(domains, describeCertConfigResp.Result.OtherCertConfig[i].Domain)
			}
		}

		if len(domains) == 0 {
			if len(describeCertConfigResp.Result.SpecifiedCertConfig) > 0 {
				// 所有可关联的域名都配置了该证书，跳过部署
				d.logger.Info("no domains to deploy")
			} else {
				return nil, errors.New("domain not found")
			}
		}
	} else {
		domains = append(domains, d.config.Domain)
	}

	if len(domains) > 0 {
		var errs []error

		for _, domain := range domains {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()

			default:
				// 关联证书与加速域名
				// REF: https://docs.byteplus.com/en/docs/byteplus-cdn/reference-batchdeploycert
				batchDeployCertReq := &bpcdn.BatchDeployCertRequest{
					CertId: upres.CertId,
					Domain: domain,
				}
				batchDeployCertResp, err := d.sdkClient.BatchDeployCert(batchDeployCertReq)
				d.logger.Debug("sdk request 'cdn.BatchDeployCert'", slog.Any("request", batchDeployCertReq), slog.Any("response", batchDeployCertResp))
				if err != nil {
					errs = append(errs, err)
				}
			}
		}

		if len(errs) > 0 {
			return nil, errors.Join(errs...)
		}
	}

	return &core.SSLDeployResult{}, nil
}
