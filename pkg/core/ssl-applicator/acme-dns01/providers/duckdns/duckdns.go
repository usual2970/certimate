package namedotcom

import (
	"errors"
	"time"

	"github.com/go-acme/lego/v4/providers/dns/duckdns"

	"github.com/certimate-go/certimate/pkg/core"
)

type ChallengeProviderConfig struct {
	Token                 string `json:"token"`
	DnsPropagationTimeout int32  `json:"dnsPropagationTimeout,omitempty"`
}

func NewChallengeProvider(config *ChallengeProviderConfig) (core.ACMEChallenger, error) {
	if config == nil {
		return nil, errors.New("the configuration of the acme challenge provider is nil")
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
