package applicant

import (
	"encoding/json"
	"fmt"
	"os"

	namesiloProvider "github.com/go-acme/lego/v4/providers/dns/namesilo"

	"certimate/internal/domain"
)

type namesilo struct {
	option *ApplyOption
}

func NewNamesilo(option *ApplyOption) Applicant {
	return &namesilo{
		option: option,
	}
}

func (a *namesilo) Apply() (*Certificate, error) {
	access := &domain.NameSiloAccess{}
	json.Unmarshal([]byte(a.option.Access), access)

	os.Setenv("NAMESILO_API_KEY", access.ApiKey)
	os.Setenv("NAMESILO_PROPAGATION_TIMEOUT", fmt.Sprintf("%d", a.option.Timeout))

	dnsProvider, err := namesiloProvider.NewDNSProvider()
	if err != nil {
		return nil, err
	}

	return apply(a.option, dnsProvider)
}
