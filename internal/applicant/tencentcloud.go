package applicant

import (
	"encoding/json"
	"time"

	"github.com/go-acme/lego/v4/providers/dns/tencentcloud"

	"github.com/usual2970/certimate/internal/domain"
)

type tencentcloudApplicant struct {
	option *ApplyOption
}

func NewTencentCloudApplicant(option *ApplyOption) Applicant {
	return &tencentcloudApplicant{
		option: option,
	}
}

func (a *tencentcloudApplicant) Apply() (*Certificate, error) {
	access := &domain.TencentCloudAccessConfig{}
	json.Unmarshal([]byte(a.option.AccessConfig), access)

	config := tencentcloud.NewDefaultConfig()
	config.SecretID = access.SecretId
	config.SecretKey = access.SecretKey
	if a.option.PropagationTimeout != 0 {
		config.PropagationTimeout = time.Duration(a.option.PropagationTimeout) * time.Second
	}

	provider, err := tencentcloud.NewDNSProviderConfig(config)
	if err != nil {
		return nil, err
	}

	return apply(a.option, provider)
}
