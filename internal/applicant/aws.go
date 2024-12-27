package applicant

import (
	"encoding/json"
	"time"

	"github.com/go-acme/lego/v4/providers/dns/route53"

	"github.com/usual2970/certimate/internal/domain"
)

type awsApplicant struct {
	option *ApplyOption
}

func NewAWSApplicant(option *ApplyOption) Applicant {
	return &awsApplicant{
		option: option,
	}
}

func (a *awsApplicant) Apply() (*Certificate, error) {
	access := &domain.AWSAccessConfig{}
	json.Unmarshal([]byte(a.option.AccessConfig), access)

	config := route53.NewDefaultConfig()
	config.AccessKeyID = access.AccessKeyId
	config.SecretAccessKey = access.SecretAccessKey
	config.Region = access.Region
	config.HostedZoneID = access.HostedZoneId
	if a.option.PropagationTimeout != 0 {
		config.PropagationTimeout = time.Duration(a.option.PropagationTimeout) * time.Second
	}

	provider, err := route53.NewDNSProviderConfig(config)
	if err != nil {
		return nil, err
	}

	return apply(a.option, provider)
}
