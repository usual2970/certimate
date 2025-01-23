package applicant

import (
	"fmt"

	"github.com/go-acme/lego/v4/challenge"

	"github.com/usual2970/certimate/internal/domain"
	providerACMEHttpReq "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/acmehttpreq"
	providerAliyun "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/aliyun"
	providerAWSRoute53 "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/aws-route53"
	providerAzureDNS "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/azure-dns"
	providerCloudflare "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/cloudflare"
	providerClouDNS "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/cloudns"
	providerGname "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/gname"
	providerGoDaddy "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/godaddy"
	providerHuaweiCloud "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/huaweicloud"
	providerNameDotCom "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/namedotcom"
	providerNameSilo "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/namesilo"
	providerNS1 "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/ns1"
	providerPowerDNS "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/powerdns"
	providerRainYun "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/rainyun"
	providerTencentCloud "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/tencentcloud"
	providerVolcEngine "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/volcengine"
	providerWestcn "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/westcn"
	"github.com/usual2970/certimate/internal/pkg/utils/maps"
)

func createApplicant(options *applicantOptions) (challenge.Provider, error) {
	/*
	  注意：如果追加新的常量值，请保持以 ASCII 排序。
	  NOTICE: If you add new constant, please keep ASCII order.
	*/
	switch options.Provider {
	case domain.ApplyDNSProviderTypeACMEHttpReq:
		{
			access := domain.AccessConfigForACMEHttpReq{}
			if err := maps.Decode(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to decode provider access config: %w", err)
			}

			applicant, err := providerACMEHttpReq.NewChallengeProvider(&providerACMEHttpReq.ACMEHttpReqApplicantConfig{
				Endpoint:              access.Endpoint,
				Mode:                  access.Mode,
				Username:              access.Username,
				Password:              access.Password,
				DnsPropagationTimeout: options.DnsPropagationTimeout,
			})
			return applicant, err
		}

	case domain.ApplyDNSProviderTypeAliyun, domain.ApplyDNSProviderTypeAliyunDNS:
		{
			access := domain.AccessConfigForAliyun{}
			if err := maps.Decode(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to decode provider access config: %w", err)
			}

			applicant, err := providerAliyun.NewChallengeProvider(&providerAliyun.AliyunApplicantConfig{
				AccessKeyId:           access.AccessKeyId,
				AccessKeySecret:       access.AccessKeySecret,
				DnsPropagationTimeout: options.DnsPropagationTimeout,
				DnsTTL:                options.DnsTTL,
			})
			return applicant, err
		}

	case domain.ApplyDNSProviderTypeAWS, domain.ApplyDNSProviderTypeAWSRoute53:
		{
			access := domain.AccessConfigForAWS{}
			if err := maps.Decode(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to decode provider access config: %w", err)
			}

			applicant, err := providerAWSRoute53.NewChallengeProvider(&providerAWSRoute53.AWSRoute53ApplicantConfig{
				AccessKeyId:           access.AccessKeyId,
				SecretAccessKey:       access.SecretAccessKey,
				Region:                maps.GetValueAsString(options.ProviderApplyConfig, "region"),
				HostedZoneId:          maps.GetValueAsString(options.ProviderApplyConfig, "hostedZoneId"),
				DnsPropagationTimeout: options.DnsPropagationTimeout,
				DnsTTL:                options.DnsTTL,
			})
			return applicant, err
		}

	case domain.ApplyDNSProviderTypeAzureDNS:
		{
			access := domain.AccessConfigForAzure{}
			if err := maps.Decode(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to decode provider access config: %w", err)
			}

			applicant, err := providerAzureDNS.NewChallengeProvider(&providerAzureDNS.AzureDNSApplicantConfig{
				TenantId:              access.TenantId,
				ClientId:              access.ClientId,
				ClientSecret:          access.ClientSecret,
				CloudName:             access.CloudName,
				DnsPropagationTimeout: options.DnsPropagationTimeout,
				DnsTTL:                options.DnsTTL,
			})
			return applicant, err
		}

	case domain.ApplyDNSProviderTypeCloudflare:
		{
			access := domain.AccessConfigForCloudflare{}
			if err := maps.Decode(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to decode provider access config: %w", err)
			}

			applicant, err := providerCloudflare.NewChallengeProvider(&providerCloudflare.CloudflareApplicantConfig{
				DnsApiToken:           access.DnsApiToken,
				DnsPropagationTimeout: options.DnsPropagationTimeout,
				DnsTTL:                options.DnsTTL,
			})
			return applicant, err
		}

	case domain.ApplyDNSProviderTypeClouDNS:
		{
			access := domain.AccessConfigForClouDNS{}
			if err := maps.Decode(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to decode provider access config: %w", err)
			}

			applicant, err := providerClouDNS.NewChallengeProvider(&providerClouDNS.ClouDNSApplicantConfig{
				AuthId:                access.AuthId,
				AuthPassword:          access.AuthPassword,
				DnsPropagationTimeout: options.DnsPropagationTimeout,
				DnsTTL:                options.DnsTTL,
			})
			return applicant, err
		}

	case domain.ApplyDNSProviderTypeGname:
		{
			access := domain.AccessConfigForGname{}
			if err := maps.Decode(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to decode provider access config: %w", err)
			}

			applicant, err := providerGname.NewChallengeProvider(&providerGname.GnameApplicantConfig{
				AppId:                 access.AppId,
				AppKey:                access.AppKey,
				DnsPropagationTimeout: options.DnsPropagationTimeout,
				DnsTTL:                options.DnsTTL,
			})
			return applicant, err
		}

	case domain.ApplyDNSProviderTypeGoDaddy:
		{
			access := domain.AccessConfigForGoDaddy{}
			if err := maps.Decode(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to decode provider access config: %w", err)
			}

			applicant, err := providerGoDaddy.NewChallengeProvider(&providerGoDaddy.GoDaddyApplicantConfig{
				ApiKey:                access.ApiKey,
				ApiSecret:             access.ApiSecret,
				DnsPropagationTimeout: options.DnsPropagationTimeout,
				DnsTTL:                options.DnsTTL,
			})
			return applicant, err
		}

	case domain.ApplyDNSProviderTypeHuaweiCloud, domain.ApplyDNSProviderTypeHuaweiCloudDNS:
		{
			access := domain.AccessConfigForHuaweiCloud{}
			if err := maps.Decode(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to decode provider access config: %w", err)
			}

			applicant, err := providerHuaweiCloud.NewChallengeProvider(&providerHuaweiCloud.HuaweiCloudApplicantConfig{
				AccessKeyId:           access.AccessKeyId,
				SecretAccessKey:       access.SecretAccessKey,
				Region:                maps.GetValueAsString(options.ProviderApplyConfig, "region"),
				DnsPropagationTimeout: options.DnsPropagationTimeout,
				DnsTTL:                options.DnsTTL,
			})
			return applicant, err
		}

	case domain.ApplyDNSProviderTypeNameDotCom:
		{
			access := domain.AccessConfigForNameDotCom{}
			if err := maps.Decode(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to decode provider access config: %w", err)
			}

			applicant, err := providerNameDotCom.NewChallengeProvider(&providerNameDotCom.NameDotComApplicantConfig{
				Username:              access.Username,
				ApiToken:              access.ApiToken,
				DnsPropagationTimeout: options.DnsPropagationTimeout,
				DnsTTL:                options.DnsTTL,
			})
			return applicant, err
		}

	case domain.ApplyDNSProviderTypeNameSilo:
		{
			access := domain.AccessConfigForNameSilo{}
			if err := maps.Decode(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to decode provider access config: %w", err)
			}

			applicant, err := providerNameSilo.NewChallengeProvider(&providerNameSilo.NameSiloApplicantConfig{
				ApiKey:                access.ApiKey,
				DnsPropagationTimeout: options.DnsPropagationTimeout,
				DnsTTL:                options.DnsTTL,
			})
			return applicant, err
		}

	case domain.ApplyDNSProviderTypeNS1:
		{
			access := domain.AccessConfigForNS1{}
			if err := maps.Decode(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to decode provider access config: %w", err)
			}

			applicant, err := providerNS1.NewChallengeProvider(&providerNS1.NS1ApplicantConfig{
				ApiKey:                access.ApiKey,
				DnsPropagationTimeout: options.DnsPropagationTimeout,
				DnsTTL:                options.DnsTTL,
			})
			return applicant, err
		}

	case domain.ApplyDNSProviderTypePowerDNS:
		{
			access := domain.AccessConfigForPowerDNS{}
			if err := maps.Decode(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to decode provider access config: %w", err)
			}

			applicant, err := providerPowerDNS.NewChallengeProvider(&providerPowerDNS.PowerDNSApplicantConfig{
				ApiUrl:                access.ApiUrl,
				ApiKey:                access.ApiKey,
				DnsPropagationTimeout: options.DnsPropagationTimeout,
				DnsTTL:                options.DnsTTL,
			})
			return applicant, err
		}

	case domain.ApplyDNSProviderTypeRainYun:
		{
			access := domain.AccessConfigForRainYun{}
			if err := maps.Decode(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to decode provider access config: %w", err)
			}

			applicant, err := providerRainYun.NewChallengeProvider(&providerRainYun.RainYunApplicantConfig{
				ApiKey:                access.ApiKey,
				DnsPropagationTimeout: options.DnsPropagationTimeout,
				DnsTTL:                options.DnsTTL,
			})
			return applicant, err
		}

	case domain.ApplyDNSProviderTypeTencentCloud, domain.ApplyDNSProviderTypeTencentCloudDNS:
		{
			access := domain.AccessConfigForTencentCloud{}
			if err := maps.Decode(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to decode provider access config: %w", err)
			}

			applicant, err := providerTencentCloud.NewChallengeProvider(&providerTencentCloud.TencentCloudApplicantConfig{
				SecretId:              access.SecretId,
				SecretKey:             access.SecretKey,
				DnsPropagationTimeout: options.DnsPropagationTimeout,
				DnsTTL:                options.DnsTTL,
			})
			return applicant, err
		}

	case domain.ApplyDNSProviderTypeVolcEngine, domain.ApplyDNSProviderTypeVolcEngineDNS:
		{
			access := domain.AccessConfigForVolcEngine{}
			if err := maps.Decode(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to decode provider access config: %w", err)
			}

			applicant, err := providerVolcEngine.NewChallengeProvider(&providerVolcEngine.VolcEngineApplicantConfig{
				AccessKeyId:           access.AccessKeyId,
				SecretAccessKey:       access.SecretAccessKey,
				DnsPropagationTimeout: options.DnsPropagationTimeout,
				DnsTTL:                options.DnsTTL,
			})
			return applicant, err
		}

	case domain.ApplyDNSProviderTypeWestcn:
		{
			access := domain.AccessConfigForWestcn{}
			if err := maps.Decode(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to decode provider access config: %w", err)
			}

			applicant, err := providerWestcn.NewChallengeProvider(&providerWestcn.WestcnApplicantConfig{
				Username:              access.Username,
				ApiPassword:           access.ApiPassword,
				DnsPropagationTimeout: options.DnsPropagationTimeout,
				DnsTTL:                options.DnsTTL,
			})
			return applicant, err
		}
	}

	return nil, fmt.Errorf("unsupported applicant provider: %s", string(options.Provider))
}
