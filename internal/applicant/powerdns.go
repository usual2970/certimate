package applicant

import (
	"encoding/json"
	"net/url"
	"time"

	"github.com/go-acme/lego/v4/providers/dns/pdns"

	"github.com/usual2970/certimate/internal/domain"
)

type powerdnsApplicant struct {
	option *ApplyOption
}

func NewPowerDNSApplicant(option *ApplyOption) Applicant {
	return &powerdnsApplicant{
		option: option,
	}
}

func (a *powerdnsApplicant) Apply() (*Certificate, error) {
	access := &domain.PowerDNSAccessConfig{}
	json.Unmarshal([]byte(a.option.AccessConfig), access)

	config := pdns.NewDefaultConfig()
	host, _ := url.Parse(access.ApiUrl)
	config.Host = host
	config.APIKey = access.ApiKey
	if a.option.PropagationTimeout != 0 {
		config.PropagationTimeout = time.Duration(a.option.PropagationTimeout) * time.Second
	}

	provider, err := pdns.NewDNSProviderConfig(config)
	if err != nil {
		return nil, err
	}

	return apply(a.option, provider)
}
