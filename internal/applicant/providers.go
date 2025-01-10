package applicant

import (
	"encoding/json"
	"fmt"

	"github.com/go-acme/lego/v4/challenge"

	"github.com/usual2970/certimate/internal/domain"
	providerACMEHttpReq "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/acmehttpreq"
	providerAliyun "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/aliyun"
	providerAWS "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/aws"
	providerCloudflare "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/cloudflare"
	providerGoDaddy "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/godaddy"
	providerHuaweiCloud "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/huaweicloud"
	providerNameDotCom "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/namedotcom"
	providerNameSilo "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/namesilo"
	providerPowerDNS "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/powerdns"
	providerTencentCloud "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/tencentcloud"
	providerVolcEngine "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/volcengine"
)

func createChallengeProvider(provider domain.ApplyDNSProviderType, accessConfig string, applyConfig *applyConfig) (challenge.Provider, error) {
	/*
	  注意：如果追加新的常量值，请保持以 ASCII 排序。
	  NOTICE: If you add new constant, please keep ASCII order.
	*/
	switch provider {
	case domain.ApplyDNSProviderTypeACMEHttpReq:
		{
			access := &domain.AccessConfigForACMEHttpReq{}
			if err := json.Unmarshal([]byte(accessConfig), access); err != nil {
				return nil, fmt.Errorf("failed to unmarshal access config: %w", err)
			}

			applicant, err := providerACMEHttpReq.NewChallengeProvider(&providerACMEHttpReq.ACMEHttpReqApplicantConfig{
				Endpoint:           access.Endpoint,
				Mode:               access.Mode,
				Username:           access.Username,
				Password:           access.Password,
				PropagationTimeout: applyConfig.PropagationTimeout,
			})
			return applicant, err
		}

	case domain.ApplyDNSProviderTypeAliyun, domain.ApplyDNSProviderTypeAliyunDNS:
		{
			access := &domain.AccessConfigForAliyun{}
			if err := json.Unmarshal([]byte(accessConfig), access); err != nil {
				return nil, fmt.Errorf("failed to unmarshal access config: %w", err)
			}

			applicant, err := providerAliyun.NewChallengeProvider(&providerAliyun.AliyunApplicantConfig{
				AccessKeyId:        access.AccessKeyId,
				AccessKeySecret:    access.AccessKeySecret,
				PropagationTimeout: applyConfig.PropagationTimeout,
			})
			return applicant, err
		}

	case domain.ApplyDNSProviderTypeAWS, domain.ApplyDNSProviderTypeAWSRoute53:
		{
			access := &domain.AccessConfigForAWS{}
			if err := json.Unmarshal([]byte(accessConfig), access); err != nil {
				return nil, fmt.Errorf("failed to unmarshal access config: %w", err)
			}

			applicant, err := providerAWS.NewChallengeProvider(&providerAWS.AWSApplicantConfig{
				AccessKeyId:        access.AccessKeyId,
				SecretAccessKey:    access.SecretAccessKey,
				Region:             access.Region,
				HostedZoneId:       access.HostedZoneId,
				PropagationTimeout: applyConfig.PropagationTimeout,
			})
			return applicant, err
		}

	case domain.ApplyDNSProviderTypeCloudflare:
		{
			access := &domain.AccessConfigForCloudflare{}
			if err := json.Unmarshal([]byte(accessConfig), access); err != nil {
				return nil, fmt.Errorf("failed to unmarshal access config: %w", err)
			}

			applicant, err := providerCloudflare.NewChallengeProvider(&providerCloudflare.CloudflareApplicantConfig{
				DnsApiToken:        access.DnsApiToken,
				PropagationTimeout: applyConfig.PropagationTimeout,
			})
			return applicant, err
		}

	case domain.ApplyDNSProviderTypeGoDaddy:
		{
			access := &domain.AccessConfigForGoDaddy{}
			if err := json.Unmarshal([]byte(accessConfig), access); err != nil {
				return nil, fmt.Errorf("failed to unmarshal access config: %w", err)
			}

			applicant, err := providerGoDaddy.NewChallengeProvider(&providerGoDaddy.GoDaddyApplicantConfig{
				ApiKey:             access.ApiKey,
				ApiSecret:          access.ApiSecret,
				PropagationTimeout: applyConfig.PropagationTimeout,
			})
			return applicant, err
		}

	case domain.ApplyDNSProviderTypeHuaweiCloud, domain.ApplyDNSProviderTypeHuaweiCloudDNS:
		{
			access := &domain.AccessConfigForHuaweiCloud{}
			if err := json.Unmarshal([]byte(accessConfig), access); err != nil {
				return nil, fmt.Errorf("failed to unmarshal access config: %w", err)
			}

			applicant, err := providerHuaweiCloud.NewChallengeProvider(&providerHuaweiCloud.HuaweiCloudApplicantConfig{
				AccessKeyId:        access.AccessKeyId,
				SecretAccessKey:    access.SecretAccessKey,
				Region:             access.Region,
				PropagationTimeout: applyConfig.PropagationTimeout,
			})
			return applicant, err
		}

	case domain.ApplyDNSProviderTypeNameDotCom:
		{
			access := &domain.AccessConfigForNameDotCom{}
			if err := json.Unmarshal([]byte(accessConfig), access); err != nil {
				return nil, fmt.Errorf("failed to unmarshal access config: %w", err)
			}

			applicant, err := providerNameDotCom.NewChallengeProvider(&providerNameDotCom.NameDotComApplicantConfig{
				Username:           access.Username,
				ApiToken:           access.ApiToken,
				PropagationTimeout: applyConfig.PropagationTimeout,
			})
			return applicant, err
		}

	case domain.ApplyDNSProviderTypeNameSilo:
		{
			access := &domain.AccessConfigForNameSilo{}
			if err := json.Unmarshal([]byte(accessConfig), access); err != nil {
				return nil, fmt.Errorf("failed to unmarshal access config: %w", err)
			}

			applicant, err := providerNameSilo.NewChallengeProvider(&providerNameSilo.NameSiloApplicantConfig{
				ApiKey:             access.ApiKey,
				PropagationTimeout: applyConfig.PropagationTimeout,
			})
			return applicant, err
		}

	case domain.ApplyDNSProviderTypePowerDNS:
		{
			access := &domain.AccessConfigForPowerDNS{}
			if err := json.Unmarshal([]byte(accessConfig), access); err != nil {
				return nil, fmt.Errorf("failed to unmarshal access config: %w", err)
			}

			applicant, err := providerPowerDNS.NewChallengeProvider(&providerPowerDNS.PowerDNSApplicantConfig{
				ApiUrl:             access.ApiUrl,
				ApiKey:             access.ApiKey,
				PropagationTimeout: applyConfig.PropagationTimeout,
			})
			return applicant, err
		}

	case domain.ApplyDNSProviderTypeTencentCloud, domain.ApplyDNSProviderTypeTencentCloudDNS:
		{
			access := &domain.AccessConfigForTencentCloud{}
			if err := json.Unmarshal([]byte(accessConfig), access); err != nil {
				return nil, fmt.Errorf("failed to unmarshal access config: %w", err)
			}

			applicant, err := providerTencentCloud.NewChallengeProvider(&providerTencentCloud.TencentCloudApplicantConfig{
				SecretId:           access.SecretId,
				SecretKey:          access.SecretKey,
				PropagationTimeout: applyConfig.PropagationTimeout,
			})
			return applicant, err
		}

	case domain.ApplyDNSProviderTypeVolcEngine, domain.ApplyDNSProviderTypeVolcEngineDNS:
		{
			access := &domain.AccessConfigForVolcEngine{}
			if err := json.Unmarshal([]byte(accessConfig), access); err != nil {
				return nil, fmt.Errorf("failed to unmarshal access config: %w", err)
			}

			applicant, err := providerVolcEngine.NewChallengeProvider(&providerVolcEngine.VolcEngineApplicantConfig{
				AccessKeyId:        access.AccessKeyId,
				SecretAccessKey:    access.SecretAccessKey,
				PropagationTimeout: applyConfig.PropagationTimeout,
			})
			return applicant, err
		}
	}

	return nil, fmt.Errorf("unsupported applicant provider: %s", provider)
}
