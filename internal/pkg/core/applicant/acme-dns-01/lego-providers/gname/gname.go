package gname

import (
	"errors"
	"time"

	"github.com/go-acme/lego/v4/challenge"

	internal "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/gname/internal"
)

type GnameApplicantConfig struct {
	AppId                 string `json:"appId"`
	AppKey                string `json:"appKey"`
	DnsPropagationTimeout int32  `json:"dnsPropagationTimeout,omitempty"`
	DnsTTL                int32  `json:"dnsTTL,omitempty"`
}

func NewChallengeProvider(config *GnameApplicantConfig) (challenge.Provider, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}

	providerConfig := internal.NewDefaultConfig()
	providerConfig.AppID = config.AppId
	providerConfig.AppKey = config.AppKey
	if config.DnsPropagationTimeout != 0 {
		providerConfig.PropagationTimeout = time.Duration(config.DnsPropagationTimeout) * time.Second
	}
	if config.DnsTTL != 0 {
		providerConfig.TTL = int(config.DnsTTL)
	}

	provider, err := internal.NewDNSProviderConfig(providerConfig)
	if err != nil {
		return nil, err
	}

	return provider, nil
}
