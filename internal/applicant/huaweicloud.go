package applicant

import (
	"encoding/json"
	"fmt"
	"os"

	huaweicloudProvider "github.com/go-acme/lego/v4/providers/dns/huaweicloud"

	"certimate/internal/domain"
)

type huaweicloud struct {
	option *ApplyOption
}

func NewHuaweiCloud(option *ApplyOption) Applicant {
	return &huaweicloud{
		option: option,
	}
}

func (t *huaweicloud) Apply() (*Certificate, error) {
	access := &domain.HuaweiCloudAccess{}
	json.Unmarshal([]byte(t.option.Access), access)

	os.Setenv("HUAWEICLOUD_REGION", access.Region) // 华为云的 SDK 要求必须传一个区域，实际上 DNS-01 流程里用不到，但不传会报错
	os.Setenv("HUAWEICLOUD_ACCESS_KEY_ID", access.AccessKeyId)
	os.Setenv("HUAWEICLOUD_SECRET_ACCESS_KEY", access.SecretAccessKey)
	os.Setenv("HUAWEICLOUD_PROPAGATION_TIMEOUT", fmt.Sprintf("%d", t.option.Timeout))

	dnsProvider, err := huaweicloudProvider.NewDNSProvider()
	if err != nil {
		return nil, err
	}

	return apply(t.option, dnsProvider)
}
