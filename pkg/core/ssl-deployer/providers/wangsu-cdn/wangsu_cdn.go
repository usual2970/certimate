package wangsucdn

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/certimate-go/certimate/pkg/core"
	sslmgrsp "github.com/certimate-go/certimate/pkg/core/ssl-manager/providers/wangsu-certificate"
	wangsusdk "github.com/certimate-go/certimate/pkg/sdk3rd/wangsu/cdn"
	xslices "github.com/certimate-go/certimate/pkg/utils/slices"
)

type SSLDeployerProviderConfig struct {
	// 网宿云 AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// 网宿云 AccessKeySecret。
	AccessKeySecret string `json:"accessKeySecret"`
	// 加速域名数组（支持泛域名）。
	Domains []string `json:"domains"`
}

type SSLDeployerProvider struct {
	config     *SSLDeployerProviderConfig
	logger     *slog.Logger
	sdkClient  *wangsusdk.Client
	sslManager core.SSLManager
}

var _ core.SSLDeployer = (*SSLDeployerProvider)(nil)

func NewSSLDeployerProvider(config *SSLDeployerProviderConfig) (*SSLDeployerProvider, error) {
	if config == nil {
		return nil, errors.New("the configuration of the ssl deployer provider is nil")
	}

	client, err := createSDKClient(config.AccessKeyId, config.AccessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("could not create sdk client: %w", err)
	}

	sslmgr, err := sslmgrsp.NewSSLManagerProvider(&sslmgrsp.SSLManagerProviderConfig{
		AccessKeyId:     config.AccessKeyId,
		AccessKeySecret: config.AccessKeySecret,
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
}

func (d *SSLDeployerProvider) Deploy(ctx context.Context, certPEM string, privkeyPEM string) (*core.SSLDeployResult, error) {
	if len(d.config.Domains) == 0 {
		return nil, errors.New("config `domains` is required")
	}

	// 上传证书
	upres, err := d.sslManager.Upload(ctx, certPEM, privkeyPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to upload certificate file: %w", err)
	} else {
		d.logger.Info("ssl certificate uploaded", slog.Any("result", upres))
	}

	// 批量修改域名证书配置
	// REF: https://www.wangsu.com/document/api-doc/37447
	certId, _ := strconv.ParseInt(upres.CertId, 10, 64)
	batchUpdateCertificateConfigReq := &wangsusdk.BatchUpdateCertificateConfigRequest{
		CertificateId: certId,
		DomainNames: xslices.Map(d.config.Domains, func(domain string) string {
			// "*.example.com" → ".example.com"，适配网宿云 CDN 要求的泛域名格式
			return strings.TrimPrefix(domain, "*")
		}),
	}
	batchUpdateCertificateConfigResp, err := d.sdkClient.BatchUpdateCertificateConfig(batchUpdateCertificateConfigReq)
	d.logger.Debug("sdk request 'cdn.BatchUpdateCertificateConfig'", slog.Any("request", batchUpdateCertificateConfigReq), slog.Any("response", batchUpdateCertificateConfigResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'cdn.BatchUpdateCertificateConfig': %w", err)
	}

	return &core.SSLDeployResult{}, nil
}

func createSDKClient(accessKeyId, accessKeySecret string) (*wangsusdk.Client, error) {
	return wangsusdk.NewClient(accessKeyId, accessKeySecret)
}
