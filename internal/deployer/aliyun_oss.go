package deployer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"

	"github.com/usual2970/certimate/internal/domain"
)

type AliyunOSSDeployer struct {
	client *oss.Client
	option *DeployerOption
	infos  []string
}

func NewAliyunOssDeployer(option *DeployerOption) (Deployer, error) {
	access := &domain.AliyunAccess{}
	json.Unmarshal([]byte(option.Access), access)

	d := &AliyunOSSDeployer{
		option: option,
		infos:  make([]string, 0),
	}

	client, err := d.createClient(access.AccessKeyId, access.AccessKeySecret)
	if err != nil {
		return nil, err
	}
	d.client = client

	return d, nil
}

func (d *AliyunOSSDeployer) GetID() string {
	return fmt.Sprintf("%s-%s", d.option.AceessRecord.GetString("name"), d.option.AceessRecord.Id)
}

func (d *AliyunOSSDeployer) GetInfo() []string {
	return d.infos
}

func (d *AliyunOSSDeployer) Deploy(ctx context.Context) error {
	err := d.client.PutBucketCnameWithCertificate(getDeployString(d.option.DeployConfig, "bucket"), oss.PutBucketCname{
		Cname: getDeployString(d.option.DeployConfig, "domain"),
		CertificateConfiguration: &oss.CertificateConfiguration{
			Certificate: d.option.Certificate.Certificate,
			PrivateKey:  d.option.Certificate.PrivateKey,
			Force:       true,
		},
	})
	if err != nil {
		return fmt.Errorf("deploy aliyun oss error: %w", err)
	}
	return nil
}

func (d *AliyunOSSDeployer) createClient(accessKeyId, accessKeySecret string) (*oss.Client, error) {
	client, err := oss.New(
		getDeployString(d.option.DeployConfig, "endpoint"),
		accessKeyId,
		accessKeySecret,
	)
	if err != nil {
		return nil, fmt.Errorf("create aliyun client error: %w", err)
	}
	return client, nil
}
