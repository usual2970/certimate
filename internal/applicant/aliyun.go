package applicant

import (
	"encoding/json"
	"os"

	"github.com/go-acme/lego/v4/providers/dns/alidns"
)

type aliyunAccess struct {
	AccessKeyId     string `json:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret"`
}

type aliyun struct {
	option *ApplyOption
}

func NewAliyun(option *ApplyOption) Applicant {
	return &aliyun{
		option: option,
	}
}

func (a *aliyun) Apply() (*Certificate, error) {

	access := &aliyunAccess{}
	json.Unmarshal([]byte(a.option.Access), access)

	os.Setenv("ALICLOUD_ACCESS_KEY", access.AccessKeyId)
	os.Setenv("ALICLOUD_SECRET_KEY", access.AccessKeySecret)
	dnsProvider, err := alidns.NewDNSProvider()
	if err != nil {
		return nil, err
	}

	return apply(a.option, dnsProvider)
}
