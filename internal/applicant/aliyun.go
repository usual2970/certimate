package applicant

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-acme/lego/v4/providers/dns/alidns"

	"certimate/internal/domain"
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
	os.Setenv("ALICLOUD_PROPAGATION_TIMEOUT", fmt.Sprintf("%d", a.option.Timeout))
	dnsProvider, err := alidns.NewDNSProvider()
	if err != nil {
		return nil, err
	}

	return apply(a.option, dnsProvider)
}
