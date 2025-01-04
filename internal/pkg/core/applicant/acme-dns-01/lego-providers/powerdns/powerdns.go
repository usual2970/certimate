package namesilo

import (
	"errors"
	"net/url"
	"time"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/pdns"
)

type PowerDNSApplicantConfig struct {
	ApiUrl             string `json:"apiUrl"`
	ApiKey             string `json:"apiKey"`
	PropagationTimeout int32  `json:"propagationTimeout,omitempty"`
}

func NewChallengeProvider(config *PowerDNSApplicantConfig) (challenge.Provider, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}

	host, _ := url.Parse(config.ApiUrl)
	providerConfig := pdns.NewDefaultConfig()
	providerConfig.Host = host
	providerConfig.APIKey = config.ApiKey
	if config.PropagationTimeout != 0 {
		providerConfig.PropagationTimeout = time.Duration(config.PropagationTimeout) * time.Second
	}

	provider, err := pdns.NewDNSProviderConfig(providerConfig)
	if err != nil {
		return nil, err
	}

	return provider, nil
}
