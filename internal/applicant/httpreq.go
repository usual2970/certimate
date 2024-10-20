package applicant

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-acme/lego/v4/providers/dns/httpreq"

	"github.com/usual2970/certimate/internal/domain"
)

type httpReq struct {
	option *ApplyOption
}

func NewHttpreq(option *ApplyOption) Applicant {
	return &httpReq{
		option: option,
	}
}

func (a *httpReq) Apply() (*Certificate, error) {
	access := &domain.HttpreqAccess{}
	json.Unmarshal([]byte(a.option.Access), access)

	os.Setenv("HTTPREQ_ENDPOINT", access.Endpoint)
	os.Setenv("HTTPREQ_MODE", access.Mode)
	os.Setenv("HTTPREQ_USERNAME", access.Username)
	os.Setenv("HTTPREQ_PASSWORD", access.Password)
	os.Setenv("HTTPREQ_HTTP_TIMEOUT", fmt.Sprintf("%d", a.option.Timeout))
	dnsProvider, err := httpreq.NewDNSProvider()
	if err != nil {
		return nil, err
	}

	return apply(a.option, dnsProvider)
}
