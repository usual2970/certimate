package awscloudfront

import (
	"context"
	"errors"
	"log/slog"

	aws "github.com/aws/aws-sdk-go-v2/aws"
	awscfg "github.com/aws/aws-sdk-go-v2/config"
	awscred "github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	xerrors "github.com/pkg/errors"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	uploadersp "github.com/usual2970/certimate/internal/pkg/core/uploader/providers/aws-acm"
)

type DeployerConfig struct {
	// AWS AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// AWS SecretAccessKey。
	SecretAccessKey string `json:"secretAccessKey"`
	// AWS 区域。
	Region string `json:"region"`
	// AWS CloudFront 分配 ID。
	DistributionId string `json:"distributionId"`
}

type DeployerProvider struct {
	config      *DeployerConfig
	logger      *slog.Logger
	sdkClient   *cloudfront.Client
	sslUploader uploader.Uploader
}

var _ deployer.Deployer = (*DeployerProvider)(nil)

func NewDeployer(config *DeployerConfig) (*DeployerProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	client, err := createSdkClient(config.AccessKeyId, config.SecretAccessKey, config.Region)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk client")
	}

	uploader, err := uploadersp.NewUploader(&uploadersp.UploaderConfig{
		AccessKeyId:     config.AccessKeyId,
		SecretAccessKey: config.SecretAccessKey,
		Region:          config.Region,
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
	if d.config.DistributionId == "" {
		return nil, errors.New("config `distribuitionId` is required")
	}

	// 上传证书到 ACM
	upres, err := d.sslUploader.Upload(ctx, certPem, privkeyPem)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to upload certificate file")
	} else {
		d.logger.Info("ssl certificate uploaded", slog.Any("result", upres))
	}

	// 获取分配配置
	// REF: https://docs.aws.amazon.com/en_us/cloudfront/latest/APIReference/API_GetDistributionConfig.html
	getDistributionConfigReq := &cloudfront.GetDistributionConfigInput{
		Id: aws.String(d.config.DistributionId),
	}
	getDistributionConfigResp, err := d.sdkClient.GetDistributionConfig(context.TODO(), getDistributionConfigReq)
	d.logger.Debug("sdk request 'cloudfront.GetDistributionConfig'", slog.Any("request", getDistributionConfigReq), slog.Any("response", getDistributionConfigResp))
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'cloudfront.GetDistributionConfig'")
	}

	// 更新分配配置
	// REF: https://docs.aws.amazon.com/zh_cn/cloudfront/latest/APIReference/API_UpdateDistribution.html
	updateDistributionReq := &cloudfront.UpdateDistributionInput{
		Id:                 aws.String(d.config.DistributionId),
		DistributionConfig: getDistributionConfigResp.DistributionConfig,
		IfMatch:            getDistributionConfigResp.ETag,
	}
	if updateDistributionReq.DistributionConfig.ViewerCertificate == nil {
		updateDistributionReq.DistributionConfig.ViewerCertificate = &types.ViewerCertificate{}
	}
	updateDistributionReq.DistributionConfig.ViewerCertificate.CloudFrontDefaultCertificate = aws.Bool(false)
	updateDistributionReq.DistributionConfig.ViewerCertificate.ACMCertificateArn = aws.String(upres.CertId)
	updateDistributionResp, err := d.sdkClient.UpdateDistribution(context.TODO(), updateDistributionReq)
	d.logger.Debug("sdk request 'cloudfront.UpdateDistribution'", slog.Any("request", updateDistributionReq), slog.Any("response", updateDistributionResp))
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'cloudfront.UpdateDistribution'")
	}

	return &deployer.DeployResult{}, nil
}

func createSdkClient(accessKeyId, secretAccessKey, region string) (*cloudfront.Client, error) {
	cfg, err := awscfg.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	client := cloudfront.NewFromConfig(cfg, func(o *cloudfront.Options) {
		o.Region = region
		o.Credentials = aws.NewCredentialsCache(awscred.NewStaticCredentialsProvider(accessKeyId, secretAccessKey, ""))
	})
	return client, nil
}
