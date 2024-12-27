package applicant

import (
	"encoding/json"
	"time"

	huaweicloud "github.com/go-acme/lego/v4/providers/dns/huaweicloud"

	"github.com/usual2970/certimate/internal/domain"
)

type huaweicloudApplicant struct {
	option *ApplyOption
}

func NewHuaweiCloudApplicant(option *ApplyOption) Applicant {
	return &huaweicloudApplicant{
		option: option,
	}
}

func (a *huaweicloudApplicant) Apply() (*Certificate, error) {
	access := &domain.HuaweiCloudAccessConfig{}
	json.Unmarshal([]byte(a.option.AccessConfig), access)

	region := access.Region
	if region == "" {
		// 华为云的 SDK 要求必须传一个区域，实际上 DNS-01 流程里用不到，但不传会报错
		region = "cn-north-1"
	}

	config := huaweicloud.NewDefaultConfig()
	config.AccessKeyID = access.AccessKeyId
	config.SecretAccessKey = access.SecretAccessKey
	config.Region = region
	if a.option.PropagationTimeout != 0 {
		config.PropagationTimeout = time.Duration(a.option.PropagationTimeout) * time.Second
	}

	provider, err := huaweicloud.NewDNSProviderConfig(config)
	if err != nil {
		return nil, err
	}

	return apply(a.option, provider)
}
