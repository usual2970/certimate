package volcenginetos

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/volcengine/ve-tos-golang-sdk/v2/tos"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	uploadersp "github.com/usual2970/certimate/internal/pkg/core/uploader/providers/volcengine-certcenter"
)

type DeployerConfig struct {
	// 火山引擎 AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// 火山引擎 AccessKeySecret。
	AccessKeySecret string `json:"accessKeySecret"`
	// 火山引擎地域。
	Region string `json:"region"`
	// 存储桶名。
	Bucket string `json:"bucket"`
	// 自定义域名（不支持泛域名）。
	Domain string `json:"domain"`
}

type DeployerProvider struct {
	config      *DeployerConfig
	logger      *slog.Logger
	sdkClient   *tos.ClientV2
	sslUploader uploader.Uploader
}

var _ deployer.Deployer = (*DeployerProvider)(nil)

func NewDeployer(config *DeployerConfig) (*DeployerProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	client, err := createSdkClient(config.AccessKeyId, config.AccessKeySecret, config.Region)
	if err != nil {
		return nil, fmt.Errorf("failed to create sdk client: %w", err)
	}

	uploader, err := uploadersp.NewUploader(&uploadersp.UploaderConfig{
		AccessKeyId:     config.AccessKeyId,
		AccessKeySecret: config.AccessKeySecret,
		Region:          config.Region,
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
	d.sslUploader.WithLogger(logger)
	return d
}

func (d *DeployerProvider) Deploy(ctx context.Context, certPEM string, privkeyPEM string) (*deployer.DeployResult, error) {
	if d.config.Bucket == "" {
		return nil, errors.New("config `bucket` is required")
	}
	if d.config.Domain == "" {
		return nil, errors.New("config `domain` is required")
	}

	// 上传证书到证书中心
	upres, err := d.sslUploader.Upload(ctx, certPEM, privkeyPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to upload certificate file: %w", err)
	} else {
		d.logger.Info("ssl certificate uploaded", slog.Any("result", upres))
	}

	// 设置自定义域名
	// REF: https://www.volcengine.com/docs/6559/1250189
	putBucketCustomDomainReq := &tos.PutBucketCustomDomainInput{
		Bucket: d.config.Bucket,
		Rule: tos.CustomDomainRule{
			Domain: d.config.Domain,
			CertID: upres.CertId,
		},
	}
	putBucketCustomDomainResp, err := d.sdkClient.PutBucketCustomDomain(context.TODO(), putBucketCustomDomainReq)
	d.logger.Debug("sdk request 'tos.PutBucketCustomDomain'", slog.Any("request", putBucketCustomDomainReq), slog.Any("response", putBucketCustomDomainResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'tos.PutBucketCustomDomain': %w", err)
	}

	return &deployer.DeployResult{}, nil
}

func createSdkClient(accessKeyId, accessKeySecret, region string) (*tos.ClientV2, error) {
	endpoint := fmt.Sprintf("tos-%s.ivolces.com", region)

	client, err := tos.NewClientV2(
		endpoint,
		tos.WithRegion(region),
		tos.WithCredentials(tos.NewStaticCredentials(accessKeyId, accessKeySecret)),
	)
	if err != nil {
		return nil, err
	}

	return client, nil
}
