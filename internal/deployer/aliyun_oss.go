package deployer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"

	"github.com/usual2970/certimate/internal/domain"
)

type AliyunOSSDeployer struct {
	option *DeployerOption
	infos  []string

	sdkClient *oss.Client
}

func NewAliyunOSSDeployer(option *DeployerOption) (Deployer, error) {
	access := &domain.AliyunAccess{}
	json.Unmarshal([]byte(option.Access), access)

	client, err := (&AliyunOSSDeployer{}).createSdkClient(
		access.AccessKeyId,
		access.AccessKeySecret,
		option.DeployConfig.GetConfigAsString("endpoint"),
	)
	if err != nil {
		return nil, err
	}

	return &AliyunOSSDeployer{
		option:    option,
		infos:     make([]string, 0),
		sdkClient: client,
	}, nil
}

func (d *AliyunOSSDeployer) GetID() string {
	return fmt.Sprintf("%s-%s", d.option.AccessRecord.GetString("name"), d.option.AccessRecord.Id)
}

func (d *AliyunOSSDeployer) GetInfo() []string {
	return d.infos
}

func (d *AliyunOSSDeployer) Deploy(ctx context.Context) error {
	aliBucket := d.option.DeployConfig.GetConfigAsString("bucket")
	if aliBucket == "" {
		return errors.New("`bucket` is required")
	}

	// 为存储空间绑定自定义域名
	// REF: https://help.aliyun.com/zh/oss/developer-reference/putcname
	err := d.sdkClient.PutBucketCnameWithCertificate(aliBucket, oss.PutBucketCname{
		Cname: d.option.DeployConfig.GetConfigAsString("domain"),
		CertificateConfiguration: &oss.CertificateConfiguration{
			Certificate: d.option.Certificate.Certificate,
			PrivateKey:  d.option.Certificate.PrivateKey,
			Force:       true,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to execute sdk request 'oss.PutBucketCnameWithCertificate': %w", err)
	}

	return nil
}

func (d *AliyunOSSDeployer) createSdkClient(accessKeyId, accessKeySecret, endpoint string) (*oss.Client, error) {
	if endpoint == "" {
		endpoint = "oss.aliyuncs.com"
	}

	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		return nil, err
	}

	return client, nil
}
