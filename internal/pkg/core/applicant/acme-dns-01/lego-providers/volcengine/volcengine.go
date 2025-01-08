package volcengine

import (
	"errors"
	"time"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/volcengine"
)

type VolcEngineApplicantConfig struct {
	AccessKeyId        string `json:"accessKeyId"`
	SecretAccessKey    string `json:"secretAccessKey"`
	PropagationTimeout int32  `json:"propagationTimeout,omitempty"`
}

func NewChallengeProvider(config *VolcEngineApplicantConfig) (challenge.Provider, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}

	providerConfig := volcengine.NewDefaultConfig()
	providerConfig.AccessKey = config.AccessKeyId
	providerConfig.SecretKey = config.SecretAccessKey
	if config.PropagationTimeout != 0 {
		providerConfig.PropagationTimeout = time.Duration(config.PropagationTimeout) * time.Second
	}

	provider, err := volcengine.NewDNSProviderConfig(providerConfig)
	if err != nil {
		return nil, err
	}

	return provider, nil
}
