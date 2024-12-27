package applicant

import (
	"encoding/json"
	"time"

	"github.com/go-acme/lego/v4/providers/dns/volcengine"
	"github.com/usual2970/certimate/internal/domain"
)

type volcengineApplicant struct {
	option *ApplyOption
}

func NewVolcEngineApplicant(option *ApplyOption) Applicant {
	return &volcengineApplicant{
		option: option,
	}
}

func (a *volcengineApplicant) Apply() (*Certificate, error) {
	access := &domain.VolcEngineAccessConfig{}
	json.Unmarshal([]byte(a.option.AccessConfig), access)

	config := volcengine.NewDefaultConfig()
	config.AccessKey = access.AccessKeyId
	config.SecretKey = access.SecretAccessKey
	if a.option.PropagationTimeout != 0 {
		config.PropagationTimeout = time.Duration(a.option.PropagationTimeout) * time.Second
	}

	provider, err := volcengine.NewDNSProviderConfig(config)
	if err != nil {
		return nil, err
	}

	return apply(a.option, provider)
}
