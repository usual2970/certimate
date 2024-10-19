package applicant

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-acme/lego/v4/providers/dns/route53"

	"certimate/internal/domain"
)

type aws struct {
	option *ApplyOption
}

func NewAws(option *ApplyOption) Applicant {
	return &aws{
		option: option,
	}
}

func (t *aws) Apply() (*Certificate, error) {
	access := &domain.AwsAccess{}
	json.Unmarshal([]byte(t.option.Access), access)

	os.Setenv("AWS_REGION", access.Region)
	os.Setenv("AWS_ACCESS_KEY_ID", access.AccessKeyId)
	os.Setenv("AWS_SECRET_ACCESS_KEY", access.SecretAccessKey)
	os.Setenv("AWS_HOSTED_ZONE_ID", access.HostedZoneId)
	os.Setenv("AWS_PROPAGATION_TIMEOUT", fmt.Sprintf("%d", t.option.Timeout))

	dnsProvider, err := route53.NewDNSProvider()
	if err != nil {
		return nil, err
	}

	return apply(t.option, dnsProvider)
}
