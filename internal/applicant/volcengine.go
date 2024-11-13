package applicant

import (
	"encoding/json"
	"fmt"
	volcengineDns "github.com/go-acme/lego/v4/providers/dns/volcengine"
	"github.com/usual2970/certimate/internal/domain"
	"os"
)

type volcengine struct {
	option *ApplyOption
}

func NewVolcengine(option *ApplyOption) Applicant {
	return &volcengine{
		option: option,
	}
}

func (a *volcengine) Apply() (*Certificate, error) {
	access := &domain.VolcengineAccess{}
	json.Unmarshal([]byte(a.option.Access), access)

	os.Setenv("VOLC_ACCESSKEY", access.AccessKeyID)
	os.Setenv("VOLC_SECRETKEY", access.SecretAccessKey)
	os.Setenv("VOLC_PROPAGATION_TIMEOUT", fmt.Sprintf("%d", a.option.Timeout))
	dnsProvider, err := volcengineDns.NewDNSProvider()
	if err != nil {
		return nil, err
	}

	return apply(a.option, dnsProvider)
}
