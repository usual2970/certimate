package applicant

import (
	"encoding/json"
	"fmt"
	"os"

	godaddyProvider "github.com/go-acme/lego/v4/providers/dns/godaddy"

	"certimate/internal/domain"
)

type godaddy struct {
	option *ApplyOption
}

func NewGodaddy(option *ApplyOption) Applicant {
	return &godaddy{
		option: option,
	}
}

func (a *godaddy) Apply() (*Certificate, error) {
	access := &domain.GodaddyAccess{}
	json.Unmarshal([]byte(a.option.Access), access)

	os.Setenv("GODADDY_API_KEY", access.ApiKey)
	os.Setenv("GODADDY_API_SECRET", access.ApiSecret)
	os.Setenv("GODADDY_PROPAGATION_TIMEOUT", fmt.Sprintf("%d", a.option.Timeout))

	dnsProvider, err := godaddyProvider.NewDNSProvider()
	if err != nil {
		return nil, err
	}

	return apply(a.option, dnsProvider)
}
