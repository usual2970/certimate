package dynv6

import (
	"time"

	"github.com/go-acme/lego/v4/challenge"

	internal "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/dynv6/internal"
)

type ChallengeProviderConfig struct {
	HttpToken             string `json:"httpToken"`
	DnsPropagationTimeout int32  `json:"dnsPropagationTimeout,omitempty"`
	DnsTTL                int32  `json:"dnsTTL,omitempty"`
}

func NewChallengeProvider(config *ChallengeProviderConfig) (challenge.Provider, error) {
	if config == nil {
		panic("config is nil")
	}

	providerConfig := internal.NewDefaultConfig()
	providerConfig.HTTPToken = config.HttpToken
	if config.DnsPropagationTimeout != 0 {
		providerConfig.PropagationTimeout = time.Duration(config.DnsPropagationTimeout) * time.Second
	}
	if config.DnsTTL != 0 {
		providerConfig.TTL = int(config.DnsTTL)
	}

	provider, err := internal.NewDNSProviderConfig(providerConfig)
	if err != nil {
		return nil, err
	}

	return provider, nil
}
