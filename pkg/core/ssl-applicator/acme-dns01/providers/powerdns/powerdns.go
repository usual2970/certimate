package powerdns

import (
	"crypto/tls"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/go-acme/lego/v4/providers/dns/pdns"

	"github.com/certimate-go/certimate/pkg/core"
)

type ChallengeProviderConfig struct {
	ServerUrl                string `json:"serverUrl"`
	ApiKey                   string `json:"apiKey"`
	AllowInsecureConnections bool   `json:"allowInsecureConnections,omitempty"`
	DnsPropagationTimeout    int32  `json:"dnsPropagationTimeout,omitempty"`
	DnsTTL                   int32  `json:"dnsTTL,omitempty"`
}

func NewChallengeProvider(config *ChallengeProviderConfig) (core.ACMEChallenger, error) {
	if config == nil {
		return nil, errors.New("the configuration of the acme challenge provider is nil")
	}

	serverUrl, _ := url.Parse(config.ServerUrl)
	providerConfig := pdns.NewDefaultConfig()
	providerConfig.Host = serverUrl
	providerConfig.APIKey = config.ApiKey
	if config.AllowInsecureConnections {
		providerConfig.HTTPClient.Transport = &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}
	if config.DnsPropagationTimeout != 0 {
		providerConfig.PropagationTimeout = time.Duration(config.DnsPropagationTimeout) * time.Second
	}
	if config.DnsTTL != 0 {
		providerConfig.TTL = int(config.DnsTTL)
	}

	provider, err := pdns.NewDNSProviderConfig(providerConfig)
	if err != nil {
		return nil, err
	}

	return provider, nil
}
