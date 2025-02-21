package dnsla

import (
	"time"

	"github.com/go-acme/lego/v4/challenge"

	internal "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/dnsla/internal"
)

type ChallengeProviderConfig struct {
	ApiId                 string `json:"apiId"`
	ApiSecret             string `json:"apiSecret"`
	DnsPropagationTimeout int32  `json:"dnsPropagationTimeout,omitempty"`
	DnsTTL                int32  `json:"dnsTTL,omitempty"`
}

func NewChallengeProvider(config *ChallengeProviderConfig) (challenge.Provider, error) {
	if config == nil {
		panic("config is nil")
	}

	providerConfig := internal.NewDefaultConfig()
	providerConfig.APIId = config.ApiId
	providerConfig.APISecret = config.ApiSecret
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
