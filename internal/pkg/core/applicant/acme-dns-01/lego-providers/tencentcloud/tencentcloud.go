package tencentcloud

import (
	"errors"
	"time"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/tencentcloud"
)

type TencentCloudApplicantConfig struct {
	SecretId           string `json:"secretId"`
	SecretKey          string `json:"secretKey"`
	PropagationTimeout int32  `json:"propagationTimeout,omitempty"`
}

func NewChallengeProvider(config *TencentCloudApplicantConfig) (challenge.Provider, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}

	providerConfig := tencentcloud.NewDefaultConfig()
	providerConfig.SecretID = config.SecretId
	providerConfig.SecretKey = config.SecretKey
	if config.PropagationTimeout != 0 {
		providerConfig.PropagationTimeout = time.Duration(config.PropagationTimeout) * time.Second
	}

	provider, err := tencentcloud.NewDNSProviderConfig(providerConfig)
	if err != nil {
		return nil, err
	}

	return provider, nil
}
