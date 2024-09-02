package applicant

import (
	"certimate/internal/domain"
	"encoding/json"
	"os"

	cf "github.com/go-acme/lego/v4/providers/dns/cloudflare"
)

type cloudflare struct {
	option *ApplyOption
}

func NewCloudflare(option *ApplyOption) Applicant {
	return &cloudflare{
		option: option,
	}
}

func (c *cloudflare) Apply() (*Certificate, error) {
	access := &domain.CloudflareAccess{}
	json.Unmarshal([]byte(c.option.Access), access)

	os.Setenv("CLOUDFLARE_DNS_API_TOKEN", access.CloudflareDnsApiToken)

	provider, err := cf.NewDNSProvider()
	if err != nil {
		return nil, err
	}

	return apply(c.option, provider)
}
