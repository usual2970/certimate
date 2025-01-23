package awscloudfront

import (
	"context"
	"errors"

	aws "github.com/aws/aws-sdk-go-v2/aws"
	awsCfg "github.com/aws/aws-sdk-go-v2/config"
	awsCred "github.com/aws/aws-sdk-go-v2/credentials"
	awsCf "github.com/aws/aws-sdk-go-v2/service/cloudfront"
	awsCfTypes "github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	xerrors "github.com/pkg/errors"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/logger"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	uploaderp "github.com/usual2970/certimate/internal/pkg/core/uploader/providers/aws-acm"
)

type AWSCloudFrontDeployerConfig struct {
	// AWS AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// AWS SecretAccessKey。
	SecretAccessKey string `json:"secretAccessKey"`
	// AWS 区域。
	Region string `json:"region"`
	// AWS CloudFront 分配 ID。
	DistributionId string `json:"distributionId"`
}

type AWSCloudFrontDeployer struct {
	config      *AWSCloudFrontDeployerConfig
	logger      logger.Logger
	sdkClient   *awsCf.Client
	sslUploader uploader.Uploader
}

var _ deployer.Deployer = (*AWSCloudFrontDeployer)(nil)

func New(config *AWSCloudFrontDeployerConfig) (*AWSCloudFrontDeployer, error) {
	return NewWithLogger(config, logger.NewNilLogger())
}

func NewWithLogger(config *AWSCloudFrontDeployerConfig, logger logger.Logger) (*AWSCloudFrontDeployer, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}

	if logger == nil {
		return nil, errors.New("logger is nil")
	}

	client, err := createSdkClient(config.AccessKeyId, config.SecretAccessKey, config.Region)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk client")
	}

	uploader, err := uploaderp.New(&uploaderp.AWSCertificateManagerUploaderConfig{
		AccessKeyId:     config.AccessKeyId,
		SecretAccessKey: config.SecretAccessKey,
		Region:          config.Region,
	})
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create ssl uploader")
	}

	return &AWSCloudFrontDeployer{
		logger:      logger,
		config:      config,
		sdkClient:   client,
		sslUploader: uploader,
	}, nil
}

func (d *AWSCloudFrontDeployer) Deploy(ctx context.Context, certPem string, privkeyPem string) (*deployer.DeployResult, error) {
	if d.config.DistributionId == "" {
		return nil, errors.New("config `distribuitionId` is required")
	}

	// 上传证书到 ACM
	upres, err := d.sslUploader.Upload(ctx, certPem, privkeyPem)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to upload certificate file")
	}

	d.logger.Logt("certificate file uploaded", upres)

	// 获取分配配置
	// REF: https://docs.aws.amazon.com/en_us/cloudfront/latest/APIReference/API_GetDistributionConfig.html
	getDistributionConfigReq := &awsCf.GetDistributionConfigInput{
		Id: aws.String(d.config.DistributionId),
	}
	getDistributionConfigResp, err := d.sdkClient.GetDistributionConfig(context.TODO(), getDistributionConfigReq)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'cloudfront.GetDistributionConfig'")
	}

	d.logger.Logt("已获取分配配置", getDistributionConfigResp)

	// 更新分配配置
	// REF: https://docs.aws.amazon.com/zh_cn/cloudfront/latest/APIReference/API_UpdateDistribution.html
	updateDistributionReq := &awsCf.UpdateDistributionInput{
		Id:                 aws.String(d.config.DistributionId),
		DistributionConfig: getDistributionConfigResp.DistributionConfig,
		IfMatch:            getDistributionConfigResp.ETag,
	}
	if updateDistributionReq.DistributionConfig.ViewerCertificate == nil {
		updateDistributionReq.DistributionConfig.ViewerCertificate = &awsCfTypes.ViewerCertificate{}
	}
	updateDistributionReq.DistributionConfig.ViewerCertificate.CloudFrontDefaultCertificate = aws.Bool(false)
	updateDistributionReq.DistributionConfig.ViewerCertificate.ACMCertificateArn = aws.String(upres.CertId)
	updateDistributionResp, err := d.sdkClient.UpdateDistribution(context.TODO(), updateDistributionReq)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'cloudfront.UpdateDistribution'")
	}

	d.logger.Logt("已更新分配配置", updateDistributionResp)

	return &deployer.DeployResult{}, nil
}

func createSdkClient(accessKeyId, secretAccessKey, region string) (*awsCf.Client, error) {
	cfg, err := awsCfg.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	client := awsCf.NewFromConfig(cfg, func(o *awsCf.Options) {
		o.Region = region
		o.Credentials = aws.NewCredentialsCache(awsCred.NewStaticCredentialsProvider(accessKeyId, secretAccessKey, ""))
	})
	return client, nil
}
