package acmehttpreq

import (
	"net/url"
	"time"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/httpreq"
)

type ChallengeProviderConfig struct {
	Endpoint              string `json:"endpoint"`
	Mode                  string `json:"mode"`
	Username              string `json:"username"`
	Password              string `json:"password"`
	DnsPropagationTimeout int32  `json:"dnsPropagationTimeout,omitempty"`
}

func NewChallengeProvider(config *ChallengeProviderConfig) (challenge.Provider, error) {
	if config == nil {
		panic("config is nil")
	}

	endpoint, _ := url.Parse(config.Endpoint)
	providerConfig := httpreq.NewDefaultConfig()
	providerConfig.Endpoint = endpoint
	providerConfig.Mode = config.Mode
	providerConfig.Username = config.Username
	providerConfig.Password = config.Password
	if config.DnsPropagationTimeout != 0 {
		providerConfig.PropagationTimeout = time.Duration(config.DnsPropagationTimeout) * time.Second
	}

	provider, err := httpreq.NewDNSProviderConfig(providerConfig)
	if err != nil {
		return nil, err
	}

	return provider, nil
}
