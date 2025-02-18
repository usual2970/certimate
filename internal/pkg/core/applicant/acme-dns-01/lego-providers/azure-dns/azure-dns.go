package azuredns

import (
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/azuredns"
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
		switch strings.ToLower(config.CloudName) {
		case "default", "public", "cloud", "azurecloud":
			providerConfig.Environment = cloud.AzurePublic
		case "usgovernment", "azureusgovernment":
			providerConfig.Environment = cloud.AzureGovernment
		case "china", "chinacloud", "azurechina", "azurechinacloud":
			providerConfig.Environment = cloud.AzureChina
		default:
			return nil, fmt.Errorf("azuredns: unknown environment %s", config.CloudName)
		}
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
