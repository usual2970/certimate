package wangsucdn

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	uploadersp "github.com/usual2970/certimate/internal/pkg/core/uploader/providers/wangsu-certificate"
	wangsusdk "github.com/usual2970/certimate/internal/pkg/sdk3rd/wangsu/cdn"
)

type DeployerConfig struct {
	// 网宿云 AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// 网宿云 AccessKeySecret。
	AccessKeySecret string `json:"accessKeySecret"`
	// 加速域名数组。
	Domains []string `json:"domains"`
}

type DeployerProvider struct {
	config      *DeployerConfig
	logger      *slog.Logger
	sdkClient   *wangsusdk.Client
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
		d.logger = slog.Default()
	} else {
		d.logger = logger
	}
	return d
}

func (d *DeployerProvider) Deploy(ctx context.Context, certPEM string, privkeyPEM string) (*deployer.DeployResult, error) {
	// 上传证书到证书管理
	upres, err := d.sslUploader.Upload(ctx, certPEM, privkeyPEM)
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
		DomainNames:   d.config.Domains,
	}
	batchUpdateCertificateConfigResp, err := d.sdkClient.BatchUpdateCertificateConfig(batchUpdateCertificateConfigReq)
	d.logger.Debug("sdk request 'cdn.BatchUpdateCertificateConfig'", slog.Any("request", batchUpdateCertificateConfigReq), slog.Any("response", batchUpdateCertificateConfigResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'cdn.BatchUpdateCertificateConfig': %w", err)
	}

	return &deployer.DeployResult{}, nil
}

func createSdkClient(accessKeyId, accessKeySecret string) (*wangsusdk.Client, error) {
	if accessKeyId == "" {
		return nil, errors.New("invalid wangsu access key id")
	}

	if accessKeySecret == "" {
		return nil, errors.New("invalid wangsu access key secret")
	}

	return wangsusdk.NewClient(accessKeyId, accessKeySecret), nil
}
