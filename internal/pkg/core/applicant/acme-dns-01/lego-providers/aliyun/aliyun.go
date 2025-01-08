package aliyun

import (
	"errors"
	"time"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/alidns"
)

type AliyunApplicantConfig struct {
	AccessKeyId        string `json:"accessKeyId"`
	AccessKeySecret    string `json:"accessKeySecret"`
	PropagationTimeout int32  `json:"propagationTimeout,omitempty"`
}

func NewChallengeProvider(config *AliyunApplicantConfig) (challenge.Provider, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}

	providerConfig := alidns.NewDefaultConfig()
	providerConfig.APIKey = config.AccessKeyId
	providerConfig.SecretKey = config.AccessKeySecret
	if config.PropagationTimeout != 0 {
		providerConfig.PropagationTimeout = time.Duration(config.PropagationTimeout) * time.Second
	}

	provider, err := alidns.NewDNSProviderConfig(providerConfig)
	if err != nil {
		return nil, err
	}

	return provider, nil
}
