package applicant

import (
	"certimate/internal/domain"
	"encoding/json"
	"os"

	"github.com/go-acme/lego/v4/providers/dns/route53"
)

type aws struct {
	option *ApplyOption
}

func NewAws(option *ApplyOption) Applicant {
	return &aws{
		option: option,
	}
}

func (a *aws) Apply() (*Certificate, error) {

	access := &domain.AwsAccess{}
	json.Unmarshal([]byte(a.option.Access), access)

	os.Setenv("AWS_ACCESS_KEY_ID", access.AccessKeyID)
	os.Setenv("AWS_SECRET_ACCESS_KEY", access.SecretAccessKey)
	dnsProvider, err := route53.NewDNSProvider()
	if err != nil {
		return nil, err
	}
	switch access.SSLprovider {
	case "letsencrypt":
		return applyLetsencrypt(a.option, dnsProvider)
	case "zerossl":
		return applyZeroSSL(a.option, dnsProvider)
	default:
		return applyLetsencrypt(a.option, dnsProvider)
	}
}
