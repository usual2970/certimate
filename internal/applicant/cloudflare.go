package applicant

import (
	"encoding/json"
	"fmt"
	"os"

	cf "github.com/go-acme/lego/v4/providers/dns/cloudflare"

	"certimate/internal/domain"
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

	os.Setenv("CLOUDFLARE_DNS_API_TOKEN", access.DnsApiToken)
	os.Setenv("CLOUDFLARE_PROPAGATION_TIMEOUT", fmt.Sprintf("%d", c.option.Timeout))

	provider, err := cf.NewDNSProvider()
	if err != nil {
		return nil, err
	}

	return apply(c.option, provider)
}
