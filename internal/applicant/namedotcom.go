package applicant

import (
	"encoding/json"
	"time"

	"github.com/go-acme/lego/v4/providers/dns/namedotcom"
	"github.com/usual2970/certimate/internal/domain"
)

type nameDotComApplicant struct {
	option *ApplyOption
}

func NewNameDotComApplicant(option *ApplyOption) Applicant {
	return &nameDotComApplicant{
		option: option,
	}
}

func (a *nameDotComApplicant) Apply() (*Certificate, error) {
	access := &domain.NameDotComAccess{}
	json.Unmarshal([]byte(a.option.Access), access)

	config := namedotcom.NewDefaultConfig()
	config.Username = access.Username
	config.APIToken = access.ApiToken
	if a.option.Timeout != 0 {
		config.PropagationTimeout = time.Duration(a.option.Timeout) * time.Second
	}

	provider, err := namedotcom.NewDNSProviderConfig(config)
	if err != nil {
		return nil, err
	}

	return apply(a.option, provider)
}
