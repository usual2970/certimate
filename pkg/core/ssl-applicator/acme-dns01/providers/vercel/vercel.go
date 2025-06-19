package vercel

import (
	"errors"
	"time"

	"github.com/go-acme/lego/v4/providers/dns/vercel"

	"github.com/certimate-go/certimate/pkg/core"
)

type ChallengeProviderConfig struct {
	ApiAccessToken        string `json:"apiAccessToken"`
	TeamId                string `json:"teamId,omitempty"`
	DnsPropagationTimeout int32  `json:"dnsPropagationTimeout,omitempty"`
	DnsTTL                int32  `json:"dnsTTL,omitempty"`
}

func NewChallengeProvider(config *ChallengeProviderConfig) (core.ACMEChallenger, error) {
	if config == nil {
		return nil, errors.New("the configuration of the acme challenge provider is nil")
	}

	providerConfig := vercel.NewDefaultConfig()
	providerConfig.AuthToken = config.ApiAccessToken
	providerConfig.TeamID = config.TeamId
	if config.DnsPropagationTimeout != 0 {
		providerConfig.PropagationTimeout = time.Duration(config.DnsPropagationTimeout) * time.Second
	}
	if config.DnsTTL != 0 {
		providerConfig.TTL = int(config.DnsTTL)
	}

	provider, err := vercel.NewDNSProviderConfig(providerConfig)
	if err != nil {
		return nil, err
	}

	return provider, nil
}
