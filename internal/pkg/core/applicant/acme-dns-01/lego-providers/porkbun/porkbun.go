package porkbun

import (
	"time"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/porkbun"
)

type ChallengeProviderConfig struct {
	ApiKey                string `json:"apiKey"`
	SecretApiKey          string `json:"secretApiKey"`
	DnsPropagationTimeout int32  `json:"dnsPropagationTimeout,omitempty"`
	DnsTTL                int32  `json:"dnsTTL,omitempty"`
}

func NewChallengeProvider(config *ChallengeProviderConfig) (challenge.Provider, error) {
	if config == nil {
		panic("config is nil")
	}

	providerConfig := porkbun.NewDefaultConfig()
	providerConfig.APIKey = config.ApiKey
	providerConfig.SecretAPIKey = config.SecretApiKey
	if config.DnsPropagationTimeout != 0 {
		providerConfig.PropagationTimeout = time.Duration(config.DnsPropagationTimeout) * time.Second
	}
	if config.DnsTTL != 0 {
		providerConfig.TTL = int(config.DnsTTL)
	}

	provider, err := porkbun.NewDNSProviderConfig(providerConfig)
	if err != nil {
		return nil, err
	}

	return provider, nil
}
