package applicant

import (
	"encoding/json"
	"time"

	namesilo "github.com/go-acme/lego/v4/providers/dns/namesilo"

	"github.com/usual2970/certimate/internal/domain"
)

type namesiloApplicant struct {
	option *ApplyOption
}

func NewNamesiloApplicant(option *ApplyOption) Applicant {
	return &namesiloApplicant{
		option: option,
	}
}

func (a *namesiloApplicant) Apply() (*Certificate, error) {
	access := &domain.NameSiloAccess{}
	json.Unmarshal([]byte(a.option.Access), access)

	config := namesilo.NewDefaultConfig()
	config.APIKey = access.ApiKey
	if a.option.PropagationTimeout != 0 {
		config.PropagationTimeout = time.Duration(a.option.PropagationTimeout) * time.Second
	}

	provider, err := namesilo.NewDNSProviderConfig(config)
	if err != nil {
		return nil, err
	}

	return apply(a.option, provider)
}
