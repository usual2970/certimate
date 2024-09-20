package applicant

import (
	"certimate/internal/domain"
	"encoding/json"
	"os"

	namesiloProvider "github.com/go-acme/lego/v4/providers/dns/namesilo"
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

	dnsProvider, err := namesiloProvider.NewDNSProvider()
	if err != nil {
		return nil, err
	}

	switch a.option.SSLprovider {
	case "letsencrypt":
		return applyLetsencrypt(a.option, dnsProvider)
	case "zerossl":
		return applyZeroSSL(a.option, dnsProvider)
	default:
		return applyLetsencrypt(a.option, dnsProvider)
	}
}
