package applicant

import (
	"encoding/json"
	"time"

	"github.com/go-acme/lego/v4/providers/dns/namedotcom"
	"github.com/usual2970/certimate/internal/domain"
)

type nameDotCom struct {
	option *ApplyOption
}

func NewNameDotCom(option *ApplyOption) Applicant {
	return &nameDotCom{
		option: option,
	}
}

func (n *nameDotCom) Apply() (*Certificate, error) {
	access := &domain.NameDotComAccess{}
	json.Unmarshal([]byte(n.option.Access), access)

	config := namedotcom.NewDefaultConfig()
	config.Username = access.Username
	config.APIToken = access.ApiToken
	config.PropagationTimeout = time.Duration(n.option.Timeout) * time.Second

	dnsProvider, err := namedotcom.NewDNSProviderConfig(config)
	if err != nil {
		return nil, err
	}

	return apply(n.option, dnsProvider)
}
