package aliyunoss

import (
	"context"
	"errors"
	"fmt"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	xerrors "github.com/pkg/errors"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/logger"
)

type AliyunOSSDeployerConfig struct {
	// 阿里云 AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// 阿里云 AccessKeySecret。
	AccessKeySecret string `json:"accessKeySecret"`
	// 阿里云地域。
	Region string `json:"region"`
	// 存储桶名。
	Bucket string `json:"bucket"`
	// 自定义域名（不支持泛域名）。
	Domain string `json:"domain"`
}

type AliyunOSSDeployer struct {
	config    *AliyunOSSDeployerConfig
	logger    logger.Logger
	sdkClient *oss.Client
}

var _ deployer.Deployer = (*AliyunOSSDeployer)(nil)

func New(config *AliyunOSSDeployerConfig) (*AliyunOSSDeployer, error) {
	return NewWithLogger(config, logger.NewNilLogger())
}

func NewWithLogger(config *AliyunOSSDeployerConfig, logger logger.Logger) (*AliyunOSSDeployer, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}

	if logger == nil {
		return nil, errors.New("logger is nil")
	}

	client, err := createSdkClient(config.AccessKeyId, config.AccessKeySecret, config.Region)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk client")
	}

	return &AliyunOSSDeployer{
		logger:    logger,
		config:    config,
		sdkClient: client,
	}, nil
}

func (d *AliyunOSSDeployer) Deploy(ctx context.Context, certPem string, privkeyPem string) (*deployer.DeployResult, error) {
	if d.config.Bucket == "" {
		return nil, errors.New("config `bucket` is required")
	}
	if d.config.Domain == "" {
		return nil, errors.New("config `domain` is required")
	}

	// 为存储空间绑定自定义域名
	// REF: https://help.aliyun.com/zh/oss/developer-reference/putcname
	err := d.sdkClient.PutBucketCnameWithCertificate(d.config.Bucket, oss.PutBucketCname{
		Cname: d.config.Domain,
		CertificateConfiguration: &oss.CertificateConfiguration{
			Certificate: certPem,
			PrivateKey:  privkeyPem,
			Force:       true,
		},
	})
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'oss.PutBucketCnameWithCertificate'")
	}

	return &deployer.DeployResult{}, nil
}

func createSdkClient(accessKeyId, accessKeySecret, region string) (*oss.Client, error) {
	// 接入点一览 https://help.aliyun.com/zh/oss/user-guide/regions-and-endpoints
	var endpoint string
	switch region {
	case "":
		endpoint = "oss.aliyuncs.com"
	case
		"cn-hzjbp",
		"cn-hzjbp-a",
		"cn-hzjbp-b":
		endpoint = "oss-cn-hzjbp-a-internal.aliyuncs.com"
	case
		"cn-shanghai-finance-1",
		"cn-shenzhen-finance-1",
		"cn-beijing-finance-1",
		"cn-north-2-gov-1":
		endpoint = fmt.Sprintf("oss-%s-internal.aliyuncs.com", region)
	default:
		endpoint = fmt.Sprintf("oss-%s.aliyuncs.com", region)
	}

	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		return nil, err
	}

	return client, nil
}
