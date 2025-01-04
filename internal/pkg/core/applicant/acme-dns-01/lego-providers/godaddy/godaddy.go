package godaddy

import (
	"errors"
	"time"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/godaddy"
)

type GoDaddyApplicantConfig struct {
	ApiKey             string `json:"apiKey"`
	ApiSecret          string `json:"apiSecret"`
	PropagationTimeout int32  `json:"propagationTimeout,omitempty"`
}

func NewChallengeProvider(config *GoDaddyApplicantConfig) (challenge.Provider, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}

	providerConfig := godaddy.NewDefaultConfig()
	providerConfig.APIKey = config.ApiKey
	providerConfig.APISecret = config.ApiSecret
	if config.PropagationTimeout != 0 {
		providerConfig.PropagationTimeout = time.Duration(config.PropagationTimeout) * time.Second
	}

	provider, err := godaddy.NewDNSProviderConfig(providerConfig)
	if err != nil {
		return nil, err
	}

	return provider, nil
}
