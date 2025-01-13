package azuredns

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/azuredns"
)

type AzureDNSApplicantConfig struct {
	TenantId           string `json:"tenantId"`
	ClientId           string `json:"clientId"`
	ClientSecret       string `json:"clientSecret"`
	CloudName          string `json:"cloudName,omitempty"`
	PropagationTimeout int32  `json:"propagationTimeout,omitempty"`
}

func NewChallengeProvider(config *AzureDNSApplicantConfig) (challenge.Provider, error) {
	if config == nil {
		return nil, errors.New("config is nil")
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
	if config.PropagationTimeout != 0 {
		providerConfig.PropagationTimeout = time.Duration(config.PropagationTimeout) * time.Second
	}

	provider, err := azuredns.NewDNSProviderConfig(providerConfig)
	if err != nil {
		return nil, err
	}

	return provider, nil
}
