package azuredns

import (
	"time"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/azuredns"

	azcommon "github.com/usual2970/certimate/internal/pkg/sdk3rd/azure/common"
)

type ChallengeProviderConfig struct {
	TenantId              string `json:"tenantId"`
	ClientId              string `json:"clientId"`
	ClientSecret          string `json:"clientSecret"`
	CloudName             string `json:"cloudName,omitempty"`
	DnsPropagationTimeout int32  `json:"dnsPropagationTimeout,omitempty"`
	DnsTTL                int32  `json:"dnsTTL,omitempty"`
}

func NewChallengeProvider(config *ChallengeProviderConfig) (challenge.Provider, error) {
	if config == nil {
		panic("config is nil")
	}

	providerConfig := azuredns.NewDefaultConfig()
	providerConfig.TenantID = config.TenantId
	providerConfig.ClientID = config.ClientId
	providerConfig.ClientSecret = config.ClientSecret
	if config.CloudName != "" {
		env, err := azcommon.GetCloudEnvironmentConfiguration(config.CloudName)
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
