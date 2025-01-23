package volcenginetos

import (
	"context"
	"errors"
	"fmt"

	xerrors "github.com/pkg/errors"
	veTos "github.com/volcengine/ve-tos-golang-sdk/v2/tos"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/logger"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	uploaderp "github.com/usual2970/certimate/internal/pkg/core/uploader/providers/volcengine-certcenter"
)

type VolcEngineTOSDeployerConfig struct {
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

type VolcEngineTOSDeployer struct {
	config      *VolcEngineTOSDeployerConfig
	logger      logger.Logger
	sdkClient   *veTos.ClientV2
	sslUploader uploader.Uploader
}

var _ deployer.Deployer = (*VolcEngineTOSDeployer)(nil)

func New(config *VolcEngineTOSDeployerConfig) (*VolcEngineTOSDeployer, error) {
	return NewWithLogger(config, logger.NewNilLogger())
}

func NewWithLogger(config *VolcEngineTOSDeployerConfig, logger logger.Logger) (*VolcEngineTOSDeployer, error) {
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

	uploader, err := uploaderp.New(&uploaderp.VolcEngineCertCenterUploaderConfig{
		AccessKeyId:     config.AccessKeyId,
		AccessKeySecret: config.AccessKeySecret,
		Region:          config.Region,
	})
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create ssl uploader")
	}

	return &VolcEngineTOSDeployer{
		logger:      logger,
		config:      config,
		sdkClient:   client,
		sslUploader: uploader,
	}, nil
}

func (d *VolcEngineTOSDeployer) Deploy(ctx context.Context, certPem string, privkeyPem string) (*deployer.DeployResult, error) {
	if d.config.Bucket == "" {
		return nil, errors.New("config `bucket` is required")
	}
	if d.config.Domain == "" {
		return nil, errors.New("config `domain` is required")
	}

	// 上传证书到证书中心
	upres, err := d.sslUploader.Upload(ctx, certPem, privkeyPem)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to upload certificate file")
	}

	d.logger.Logt("certificate file uploaded", upres)

	// 设置自定义域名
	// REF: https://www.volcengine.com/docs/6559/1250189
	putBucketCustomDomainReq := &veTos.PutBucketCustomDomainInput{
		Bucket: d.config.Bucket,
		Rule: veTos.CustomDomainRule{
			Domain: d.config.Domain,
			CertID: upres.CertId,
		},
	}
	putBucketCustomDomainResp, err := d.sdkClient.PutBucketCustomDomain(context.TODO(), putBucketCustomDomainReq)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'tos.PutBucketCustomDomain'")
	} else {
		d.logger.Logt("已设置自定义域名", putBucketCustomDomainResp)
	}

	return &deployer.DeployResult{}, nil
}

func createSdkClient(accessKeyId, accessKeySecret, region string) (*veTos.ClientV2, error) {
	endpoint := fmt.Sprintf("tos-%s.ivolces.com", region)

	client, err := veTos.NewClientV2(
		endpoint,
		veTos.WithRegion(region),
		veTos.WithCredentials(veTos.NewStaticCredentials(accessKeyId, accessKeySecret)),
	)
	if err != nil {
		return nil, err
	}

	return client, nil
}
