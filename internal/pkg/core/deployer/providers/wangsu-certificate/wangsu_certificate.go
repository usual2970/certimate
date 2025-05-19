package wangsucertificate

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	uploadersp "github.com/usual2970/certimate/internal/pkg/core/uploader/providers/wangsu-certificate"
	wangsusdk "github.com/usual2970/certimate/internal/pkg/sdk3rd/wangsu/certificate"
	typeutil "github.com/usual2970/certimate/internal/pkg/utils/type"
)

type DeployerConfig struct {
	// 网宿云 AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// 网宿云 AccessKeySecret。
	AccessKeySecret string `json:"accessKeySecret"`
	// 证书 ID。
	// 选填。零值时表示新建证书；否则表示更新证书。
	CertificateId string `json:"certificateId,omitempty"`
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
	if d.config.CertificateId == "" {
		// 上传证书到证书管理
		upres, err := d.sslUploader.Upload(ctx, certPEM, privkeyPEM)
		if err != nil {
			return nil, fmt.Errorf("failed to upload certificate file: %w", err)
		} else {
			d.logger.Info("ssl certificate uploaded", slog.Any("result", upres))
		}
	} else {
		// 修改证书
		// REF: https://www.wangsu.com/document/api-doc/25568?productCode=certificatemanagement
		updateCertificateReq := &wangsusdk.UpdateCertificateRequest{
			Name:        typeutil.ToPtr(fmt.Sprintf("certimate_%d", time.Now().UnixMilli())),
			Certificate: typeutil.ToPtr(certPEM),
			PrivateKey:  typeutil.ToPtr(privkeyPEM),
			Comment:     typeutil.ToPtr("upload from certimate"),
		}
		updateCertificateResp, err := d.sdkClient.UpdateCertificate(d.config.CertificateId, updateCertificateReq)
		d.logger.Debug("sdk request 'certificatemanagement.UpdateCertificate'", slog.Any("request", updateCertificateReq), slog.Any("response", updateCertificateResp))
		if err != nil {
			return nil, fmt.Errorf("failed to execute sdk request 'certificatemanagement.CreateCertificate': %w", err)
		}
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
