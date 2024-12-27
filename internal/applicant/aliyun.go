package applicant

import (
	"encoding/json"
	"time"

	"github.com/go-acme/lego/v4/providers/dns/alidns"

	"github.com/usual2970/certimate/internal/domain"
)

type aliyunApplicant struct {
	option *ApplyOption
}

func NewAliyunApplicant(option *ApplyOption) Applicant {
	return &aliyunApplicant{
		option: option,
	}
}

func (a *aliyunApplicant) Apply() (*Certificate, error) {
	access := &domain.AliyunAccessConfig{}
	json.Unmarshal([]byte(a.option.AccessConfig), access)

	config := alidns.NewDefaultConfig()
	config.APIKey = access.AccessKeyId
	config.SecretKey = access.AccessKeySecret
	if a.option.PropagationTimeout != 0 {
		config.PropagationTimeout = time.Duration(a.option.PropagationTimeout) * time.Second
	}

	provider, err := alidns.NewDNSProviderConfig(config)
	if err != nil {
		return nil, err
	}

	return apply(a.option, provider)
}
