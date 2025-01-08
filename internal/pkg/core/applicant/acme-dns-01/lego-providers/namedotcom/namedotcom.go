package namedotcom

import (
	"errors"
	"time"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/namedotcom"
)

type NameDotComApplicantConfig struct {
	Username           string `json:"username"`
	ApiToken           string `json:"apiToken"`
	PropagationTimeout int32  `json:"propagationTimeout,omitempty"`
}

func NewChallengeProvider(config *NameDotComApplicantConfig) (challenge.Provider, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}

	providerConfig := namedotcom.NewDefaultConfig()
	providerConfig.Username = config.Username
	providerConfig.APIToken = config.ApiToken
	if config.PropagationTimeout != 0 {
		providerConfig.PropagationTimeout = time.Duration(config.PropagationTimeout) * time.Second
	}

	provider, err := namedotcom.NewDNSProviderConfig(providerConfig)
	if err != nil {
		return nil, err
	}

	return provider, nil
}
