package ctcccloud

import (
	"time"

	"github.com/go-acme/lego/v4/challenge"

	"github.com/usual2970/certimate/internal/pkg/core/applicator/acme-dns01/lego-providers/ctcccloud/internal"
)

type ChallengeProviderConfig struct {
	AccessKeyId           string `json:"accessKeyId"`
	SecretAccessKey       string `json:"secretAccessKey"`
	DnsPropagationTimeout int32  `json:"dnsPropagationTimeout,omitempty"`
	DnsTTL                int32  `json:"dnsTTL,omitempty"`
}

func NewChallengeProvider(config *ChallengeProviderConfig) (challenge.Provider, error) {
	if config == nil {
		panic("config is nil")
	}

	providerConfig := internal.NewDefaultConfig()
	providerConfig.AccessKeyId = config.AccessKeyId
	providerConfig.SecretAccessKey = config.SecretAccessKey
	if config.DnsTTL != 0 {
		providerConfig.TTL = int(config.DnsTTL)
	}
	if config.DnsPropagationTimeout != 0 {
		providerConfig.PropagationTimeout = time.Duration(config.DnsPropagationTimeout) * time.Second
	}

	provider, err := internal.NewDNSProviderConfig(providerConfig)
	if err != nil {
		return nil, err
	}

	return provider, nil
}
