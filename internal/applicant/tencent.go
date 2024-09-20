package applicant

import (
	"certimate/internal/domain"
	"encoding/json"
	"os"

	"github.com/go-acme/lego/v4/providers/dns/tencentcloud"
)

type tencent struct {
	option *ApplyOption
}

func NewTencent(option *ApplyOption) Applicant {
	return &tencent{
		option: option,
	}
}

func (t *tencent) Apply() (*Certificate, error) {

	access := &domain.TencentAccess{}
	json.Unmarshal([]byte(t.option.Access), access)

	os.Setenv("TENCENTCLOUD_SECRET_ID", access.SecretId)
	os.Setenv("TENCENTCLOUD_SECRET_KEY", access.SecretKey)
	dnsProvider, err := tencentcloud.NewDNSProvider()
	if err != nil {
		return nil, err
	}

	switch t.option.SSLprovider {
	case "letsencrypt":
		return applyLetsencrypt(t.option, dnsProvider)
	case "zerossl":
		return applyZeroSSL(t.option, dnsProvider)
	default:
		return applyLetsencrypt(t.option, dnsProvider)
	}
}
