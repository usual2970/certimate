package aws

import (
	"errors"
	"time"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/route53"
)

type AWSApplicantConfig struct {
	AccessKeyId        string `json:"accessKeyId"`
	SecretAccessKey    string `json:"secretAccessKey"`
	Region             string `json:"region"`
	HostedZoneId       string `json:"hostedZoneId"`
	PropagationTimeout int32  `json:"propagationTimeout,omitempty"`
}

func NewChallengeProvider(config *AWSApplicantConfig) (challenge.Provider, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}

	providerConfig := route53.NewDefaultConfig()
	providerConfig.AccessKeyID = config.AccessKeyId
	providerConfig.SecretAccessKey = config.SecretAccessKey
	providerConfig.Region = config.Region
	providerConfig.HostedZoneID = config.HostedZoneId
	if config.PropagationTimeout != 0 {
		providerConfig.PropagationTimeout = time.Duration(config.PropagationTimeout) * time.Second
	}

	provider, err := route53.NewDNSProviderConfig(providerConfig)
	if err != nil {
		return nil, err
	}

	return provider, nil
}
