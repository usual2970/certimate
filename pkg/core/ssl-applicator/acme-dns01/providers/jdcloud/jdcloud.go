package jdcloud

import (
	"errors"
	"time"

	"github.com/certimate-go/certimate/pkg/core"
	"github.com/certimate-go/certimate/pkg/core/ssl-applicator/acme-dns01/providers/jdcloud/internal"
)

type ChallengeProviderConfig struct {
	AccessKeyId           string `json:"accessKeyId"`
	AccessKeySecret       string `json:"accessKeySecret"`
	RegionId              string `json:"regionId"`
	DnsPropagationTimeout int32  `json:"dnsPropagationTimeout,omitempty"`
	DnsTTL                int32  `json:"dnsTTL,omitempty"`
}

func NewChallengeProvider(config *ChallengeProviderConfig) (core.ACMEChallenger, error) {
	if config == nil {
		return nil, errors.New("the configuration of the acme challenge provider is nil")
	}

	regionId := config.RegionId
	if regionId == "" {
		// 京东云的 SDK 要求必须传一个区域，实际上 DNS-01 流程里用不到，但不传会报错
		regionId = "cn-north-1"
	}

	providerConfig := internal.NewDefaultConfig()
	providerConfig.AccessKeyID = config.AccessKeyId
	providerConfig.AccessKeySecret = config.AccessKeySecret
	providerConfig.RegionId = regionId
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
