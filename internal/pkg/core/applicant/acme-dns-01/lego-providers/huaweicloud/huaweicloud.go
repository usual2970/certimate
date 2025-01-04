package huaweicloud

import (
	"errors"
	"time"

	"github.com/go-acme/lego/v4/challenge"
	hwc "github.com/go-acme/lego/v4/providers/dns/huaweicloud"
)

type HuaweiCloudApplicantConfig struct {
	AccessKeyId        string `json:"accessKeyId"`
	SecretAccessKey    string `json:"secretAccessKey"`
	Region             string `json:"region"`
	PropagationTimeout int32  `json:"propagationTimeout,omitempty"`
}

func NewChallengeProvider(config *HuaweiCloudApplicantConfig) (challenge.Provider, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}

	region := config.Region
	if region == "" {
		// 华为云的 SDK 要求必须传一个区域，实际上 DNS-01 流程里用不到，但不传会报错
		region = "cn-north-1"
	}

	providerConfig := hwc.NewDefaultConfig()
	providerConfig.AccessKeyID = config.AccessKeyId
	providerConfig.SecretAccessKey = config.SecretAccessKey
	providerConfig.Region = region
	if config.PropagationTimeout != 0 {
		providerConfig.PropagationTimeout = time.Duration(config.PropagationTimeout) * time.Second
	}

	provider, err := hwc.NewDNSProviderConfig(providerConfig)
	if err != nil {
		return nil, err
	}

	return provider, nil
}
