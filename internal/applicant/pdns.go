package applicant

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-acme/lego/v4/providers/dns/pdns"

	"github.com/usual2970/certimate/internal/domain"
)

type powerdns struct {
	option *ApplyOption
}

func NewPdns(option *ApplyOption) Applicant {
	return &powerdns{
		option: option,
	}
}

func (a *powerdns) Apply() (*Certificate, error) {
	access := &domain.PdnsAccess{}
	json.Unmarshal([]byte(a.option.Access), access)

	os.Setenv("PDNS_API_URL", access.ApiUrl)
	os.Setenv("PDNS_API_KEY", access.ApiKey)
	os.Setenv("PDNS_HTTP_TIMEOUT", fmt.Sprintf("%d", a.option.Timeout))
	dnsProvider, err := pdns.NewDNSProvider()
	if err != nil {
		return nil, err
	}

	return apply(a.option, dnsProvider)
}
