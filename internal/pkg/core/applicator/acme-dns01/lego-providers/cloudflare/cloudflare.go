package cloudflare

import (
	"time"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/cloudflare"
)

type ChallengeProviderConfig struct {
	DnsApiToken           string `json:"dnsApiToken"`
	ZoneApiToken          string `json:"zoneApiToken,omitempty"`
	DnsPropagationTimeout int32  `json:"dnsPropagationTimeout,omitempty"`
	DnsTTL                int32  `json:"dnsTTL,omitempty"`
}

func NewChallengeProvider(config *ChallengeProviderConfig) (challenge.Provider, error) {
	if config == nil {
		panic("config is nil")
	}

	providerConfig := cloudflare.NewDefaultConfig()
	providerConfig.AuthToken = config.DnsApiToken
	providerConfig.ZoneToken = config.ZoneApiToken
	if config.DnsPropagationTimeout != 0 {
		providerConfig.PropagationTimeout = time.Duration(config.DnsPropagationTimeout) * time.Second
	}
	if config.DnsTTL != 0 {
		providerConfig.TTL = int(config.DnsTTL)
	}

	provider, err := cloudflare.NewDNSProviderConfig(providerConfig)
	if err != nil {
		return nil, err
	}

	return provider, nil
}
