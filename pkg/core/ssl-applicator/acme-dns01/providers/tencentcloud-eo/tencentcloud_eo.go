package tencentcloudeo

import (
	"errors"
	"time"

	"github.com/certimate-go/certimate/pkg/core"
	"github.com/certimate-go/certimate/pkg/core/ssl-applicator/acme-dns01/providers/tencentcloud-eo/internal"
)

type ChallengeProviderConfig struct {
	SecretId              string `json:"secretId"`
	SecretKey             string `json:"secretKey"`
	ZoneId                string `json:"zoneId"`
	DnsPropagationTimeout int32  `json:"dnsPropagationTimeout,omitempty"`
	DnsTTL                int32  `json:"dnsTTL,omitempty"`
}

func NewChallengeProvider(config *ChallengeProviderConfig) (core.ACMEChallenger, error) {
	if config == nil {
		return nil, errors.New("the configuration of the acme challenge provider is nil")
	}

	providerConfig := internal.NewDefaultConfig()
	providerConfig.SecretID = config.SecretId
	providerConfig.SecretKey = config.SecretKey
	providerConfig.ZoneID = config.ZoneId
	if config.DnsPropagationTimeout != 0 {
		providerConfig.PropagationTimeout = time.Duration(config.DnsPropagationTimeout) * time.Second
	}
	if config.DnsTTL != 0 {
		providerConfig.TTL = config.DnsTTL
	}

	provider, err := internal.NewDNSProviderConfig(providerConfig)
	if err != nil {
		return nil, err
	}

	return provider, nil
}
