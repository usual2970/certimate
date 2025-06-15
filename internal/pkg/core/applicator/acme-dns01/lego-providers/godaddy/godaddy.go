package godaddy

import (
	"time"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/godaddy"
)

type ChallengeProviderConfig struct {
	ApiKey                string `json:"apiKey"`
	ApiSecret             string `json:"apiSecret"`
	DnsPropagationTimeout int32  `json:"dnsPropagationTimeout,omitempty"`
	DnsTTL                int32  `json:"dnsTTL,omitempty"`
}

func NewChallengeProvider(config *ChallengeProviderConfig) (challenge.Provider, error) {
	if config == nil {
		panic("config is nil")
	}

	providerConfig := godaddy.NewDefaultConfig()
	providerConfig.APIKey = config.ApiKey
	providerConfig.APISecret = config.ApiSecret
	if config.DnsPropagationTimeout != 0 {
		providerConfig.PropagationTimeout = time.Duration(config.DnsPropagationTimeout) * time.Second
	}
	if config.DnsTTL != 0 {
		providerConfig.TTL = int(config.DnsTTL)
	}

	provider, err := godaddy.NewDNSProviderConfig(providerConfig)
	if err != nil {
		return nil, err
	}

	return provider, nil
}
