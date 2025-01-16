package ns1

import (
	"errors"
	"time"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/ns1"
)

type NS1ApplicantConfig struct {
	ApiKey             string `json:"apiKey"`
	PropagationTimeout int32  `json:"propagationTimeout,omitempty"`
}

func NewChallengeProvider(config *NS1ApplicantConfig) (challenge.Provider, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}

	providerConfig := ns1.NewDefaultConfig()
	providerConfig.APIKey = config.ApiKey
	if config.PropagationTimeout != 0 {
		providerConfig.PropagationTimeout = time.Duration(config.PropagationTimeout) * time.Second
	}

	provider, err := ns1.NewDNSProviderConfig(providerConfig)
	if err != nil {
		return nil, err
	}

	return provider, nil
}
