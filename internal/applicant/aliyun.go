package applicant

import (
	"certimate/internal/domain"
	"encoding/json"
	"os"

	"github.com/go-acme/lego/v4/providers/dns/alidns"
)

type aliyun struct {
	option *ApplyOption
}

func NewAliyun(option *ApplyOption) Applicant {
	return &aliyun{
		option: option,
	}
}

func (a *aliyun) Apply() (*Certificate, error) {

	access := &domain.AliyunAccess{}
	json.Unmarshal([]byte(a.option.Access), access)

	os.Setenv("ALICLOUD_ACCESS_KEY", access.AccessKeyId)
	os.Setenv("ALICLOUD_SECRET_KEY", access.AccessKeySecret)
	dnsProvider, err := alidns.NewDNSProvider()
	if err != nil {
		return nil, err
	}

	switch a.option.SSLprovider {
	case "letsencrypt":
		return applyLetsencrypt(a.option, dnsProvider)
	case "zerossl":
		return applyZeroSSL(a.option, dnsProvider)
	default:
		return applyLetsencrypt(a.option, dnsProvider)
	}
}
