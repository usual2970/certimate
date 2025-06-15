package namedotcom

import (
	"time"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/duckdns"
)

type ChallengeProviderConfig struct {
	Token                 string `json:"token"`
	DnsPropagationTimeout int32  `json:"dnsPropagationTimeout,omitempty"`
}

func NewChallengeProvider(config *ChallengeProviderConfig) (challenge.Provider, error) {
	if config == nil {
		panic("config is nil")
	}

	providerConfig := duckdns.NewDefaultConfig()
	providerConfig.Token = config.Token
	if config.DnsPropagationTimeout != 0 {
		providerConfig.PropagationTimeout = time.Duration(config.DnsPropagationTimeout) * time.Second
	}

	provider, err := duckdns.NewDNSProviderConfig(providerConfig)
	if err != nil {
		return nil, err
	}

	return provider, nil
}
