package deployer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"

	"certimate/internal/domain"
)

type aliyun struct {
	client *oss.Client
	option *DeployerOption
	infos  []string
}

func NewAliyun(option *DeployerOption) (Deployer, error) {
	access := &domain.AliyunAccess{}
	json.Unmarshal([]byte(option.Access), access)
	a := &aliyun{
		option: option,
		infos:  make([]string, 0),
	}
	client, err := a.createClient(access.AccessKeyId, access.AccessKeySecret)
	if err != nil {
		return nil, err
	}
	a.client = client
	return a, nil
}

func (a *aliyun) GetID() string {
	return fmt.Sprintf("%s-%s", a.option.AceessRecord.GetString("name"), a.option.AceessRecord.Id)
}

func (a *aliyun) GetInfo() []string {
	return a.infos
}

func (a *aliyun) Deploy(ctx context.Context) error {
	err := a.client.PutBucketCnameWithCertificate(getDeployString(a.option.DeployConfig, "bucket"), oss.PutBucketCname{
		Cname: getDeployString(a.option.DeployConfig, "domain"),
		CertificateConfiguration: &oss.CertificateConfiguration{
			Certificate: a.option.Certificate.Certificate,
			PrivateKey:  a.option.Certificate.PrivateKey,
			Force:       true,
		},
	})
	if err != nil {
		return fmt.Errorf("deploy aliyun oss error: %w", err)
	}

	return nil
}

func (a *aliyun) createClient(accessKeyId, accessKeySecret string) (*oss.Client, error) {
	client, err := oss.New(
		getDeployString(a.option.DeployConfig, "endpoint"),
		accessKeyId,
		accessKeySecret,
	)
	if err != nil {
		return nil, fmt.Errorf("create aliyun client error: %w", err)
	}
	return client, nil
}
