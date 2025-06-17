package azuredns

import (
	"errors"
	"time"

	"github.com/go-acme/lego/v4/providers/dns/azuredns"

	"github.com/certimate-go/certimate/pkg/core"
	azenv "github.com/certimate-go/certimate/pkg/sdk3rd/azure/env"
)

type ChallengeProviderConfig struct {
	TenantId              string `json:"tenantId"`
	ClientId              string `json:"clientId"`
	ClientSecret          string `json:"clientSecret"`
	CloudName             string `json:"cloudName,omitempty"`
	DnsPropagationTimeout int32  `json:"dnsPropagationTimeout,omitempty"`
	DnsTTL                int32  `json:"dnsTTL,omitempty"`
}

func NewChallengeProvider(config *ChallengeProviderConfig) (core.ACMEChallenger, error) {
	if config == nil {
		return nil, errors.New("the configuration of the acme challenge provider is nil")
	}

	providerConfig := azuredns.NewDefaultConfig()
	providerConfig.TenantID = config.TenantId
	providerConfig.ClientID = config.ClientId
	providerConfig.ClientSecret = config.ClientSecret
	if config.CloudName != "" {
		env, err := azenv.GetCloudEnvConfiguration(config.CloudName)
		if err != nil {
			return nil, err
		}
		providerConfig.Environment = env
	}
	if config.DnsPropagationTimeout != 0 {
		providerConfig.PropagationTimeout = time.Duration(config.DnsPropagationTimeout) * time.Second
	}
	if config.DnsTTL != 0 {
		providerConfig.TTL = int(config.DnsTTL)
	}

	provider, err := azuredns.NewDNSProviderConfig(providerConfig)
	if err != nil {
		return nil, err
	}

	return provider, nil
}
