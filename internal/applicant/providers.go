package applicant

import (
	"fmt"

	"github.com/go-acme/lego/v4/challenge"

	"github.com/usual2970/certimate/internal/domain"
	pACMEHttpReq "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/acmehttpreq"
	pAliyun "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/aliyun"
	pAWSRoute53 "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/aws-route53"
	pAzureDNS "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/azure-dns"
	pBaiduCloud "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/baiducloud"
	pCloudflare "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/cloudflare"
	pClouDNS "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/cloudns"
	pCMCCCloud "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/cmcccloud"
	pDeSEC "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/desec"
	pDNSLA "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/dnsla"
	pDynv6 "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/dynv6"
	pGcore "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/gcore"
	pGname "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/gname"
	pGoDaddy "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/godaddy"
	pHuaweiCloud "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/huaweicloud"
	pJDCloud "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/jdcloud"
	pNamecheap "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/namecheap"
	pNameDotCom "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/namedotcom"
	pNameSilo "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/namesilo"
	pNS1 "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/ns1"
	pPorkbun "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/porkbun"
	pPowerDNS "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/powerdns"
	pRainYun "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/rainyun"
	pTencentCloud "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/tencentcloud"
	pTencentCloudEO "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/tencentcloud-eo"
	pVercel "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/vercel"
	pVolcEngine "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/volcengine"
	pWestcn "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/westcn"
	"github.com/usual2970/certimate/internal/pkg/utils/maputil"
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
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			applicant, err := pACMEHttpReq.NewChallengeProvider(&pACMEHttpReq.ChallengeProviderConfig{
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
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			applicant, err := pAliyun.NewChallengeProvider(&pAliyun.ChallengeProviderConfig{
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
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			applicant, err := pAWSRoute53.NewChallengeProvider(&pAWSRoute53.ChallengeProviderConfig{
				AccessKeyId:           access.AccessKeyId,
				SecretAccessKey:       access.SecretAccessKey,
				Region:                maputil.GetString(options.ProviderApplyConfig, "region"),
				HostedZoneId:          maputil.GetString(options.ProviderApplyConfig, "hostedZoneId"),
				DnsPropagationTimeout: options.DnsPropagationTimeout,
				DnsTTL:                options.DnsTTL,
			})
			return applicant, err
		}

	case domain.ApplyDNSProviderTypeAzure, domain.ApplyDNSProviderTypeAzureDNS:
		{
			access := domain.AccessConfigForAzure{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			applicant, err := pAzureDNS.NewChallengeProvider(&pAzureDNS.ChallengeProviderConfig{
				TenantId:              access.TenantId,
				ClientId:              access.ClientId,
				ClientSecret:          access.ClientSecret,
				CloudName:             access.CloudName,
				DnsPropagationTimeout: options.DnsPropagationTimeout,
				DnsTTL:                options.DnsTTL,
			})
			return applicant, err
		}

	case domain.ApplyDNSProviderTypeBaiduCloud, domain.ApplyDNSProviderTypeBaiduCloudDNS:
		{
			access := domain.AccessConfigForBaiduCloud{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			applicant, err := pBaiduCloud.NewChallengeProvider(&pBaiduCloud.ChallengeProviderConfig{
				AccessKeyId:           access.AccessKeyId,
				SecretAccessKey:       access.SecretAccessKey,
				DnsPropagationTimeout: options.DnsPropagationTimeout,
				DnsTTL:                options.DnsTTL,
			})
			return applicant, err
		}

	case domain.ApplyDNSProviderTypeCloudflare:
		{
			access := domain.AccessConfigForCloudflare{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			applicant, err := pCloudflare.NewChallengeProvider(&pCloudflare.ChallengeProviderConfig{
				DnsApiToken:           access.DnsApiToken,
				DnsPropagationTimeout: options.DnsPropagationTimeout,
				DnsTTL:                options.DnsTTL,
			})
			return applicant, err
		}

	case domain.ApplyDNSProviderTypeClouDNS:
		{
			access := domain.AccessConfigForClouDNS{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			applicant, err := pClouDNS.NewChallengeProvider(&pClouDNS.ChallengeProviderConfig{
				AuthId:                access.AuthId,
				AuthPassword:          access.AuthPassword,
				DnsPropagationTimeout: options.DnsPropagationTimeout,
				DnsTTL:                options.DnsTTL,
			})
			return applicant, err
		}

	case domain.ApplyDNSProviderTypeCMCCCloud:
		{
			access := domain.AccessConfigForCMCCCloud{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			applicant, err := pCMCCCloud.NewChallengeProvider(&pCMCCCloud.ChallengeProviderConfig{
				AccessKeyId:           access.AccessKeyId,
				AccessKeySecret:       access.AccessKeySecret,
				DnsPropagationTimeout: options.DnsPropagationTimeout,
				DnsTTL:                options.DnsTTL,
			})
			return applicant, err
		}

	case domain.ApplyDNSProviderTypeDeSEC:
		{
			access := domain.AccessConfigForDeSEC{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			applicant, err := pDeSEC.NewChallengeProvider(&pDeSEC.ChallengeProviderConfig{
				Token:                 access.Token,
				DnsPropagationTimeout: options.DnsPropagationTimeout,
				DnsTTL:                options.DnsTTL,
			})
			return applicant, err
		}

	case domain.ApplyDNSProviderTypeDNSLA:
		{
			access := domain.AccessConfigForDNSLA{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			applicant, err := pDNSLA.NewChallengeProvider(&pDNSLA.ChallengeProviderConfig{
				ApiId:                 access.ApiId,
				ApiSecret:             access.ApiSecret,
				DnsPropagationTimeout: options.DnsPropagationTimeout,
				DnsTTL:                options.DnsTTL,
			})
			return applicant, err
		}

	case domain.ApplyDNSProviderTypeDynv6:
		{
			access := domain.AccessConfigForDynv6{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			applicant, err := pDynv6.NewChallengeProvider(&pDynv6.ChallengeProviderConfig{
				HttpToken:             access.HttpToken,
				DnsPropagationTimeout: options.DnsPropagationTimeout,
				DnsTTL:                options.DnsTTL,
			})
			return applicant, err
		}

	case domain.ApplyDNSProviderTypeGcore:
		{
			access := domain.AccessConfigForGcore{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			applicant, err := pGcore.NewChallengeProvider(&pGcore.ChallengeProviderConfig{
				ApiToken:              access.ApiToken,
				DnsPropagationTimeout: options.DnsPropagationTimeout,
				DnsTTL:                options.DnsTTL,
			})
			return applicant, err
		}

	case domain.ApplyDNSProviderTypeGname:
		{
			access := domain.AccessConfigForGname{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			applicant, err := pGname.NewChallengeProvider(&pGname.ChallengeProviderConfig{
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
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			applicant, err := pGoDaddy.NewChallengeProvider(&pGoDaddy.ChallengeProviderConfig{
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
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			applicant, err := pHuaweiCloud.NewChallengeProvider(&pHuaweiCloud.ChallengeProviderConfig{
				AccessKeyId:           access.AccessKeyId,
				SecretAccessKey:       access.SecretAccessKey,
				Region:                maputil.GetString(options.ProviderApplyConfig, "region"),
				DnsPropagationTimeout: options.DnsPropagationTimeout,
				DnsTTL:                options.DnsTTL,
			})
			return applicant, err
		}

	case domain.ApplyDNSProviderTypeJDCloud, domain.ApplyDNSProviderTypeJDCloudDNS:
		{
			access := domain.AccessConfigForJDCloud{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			applicant, err := pJDCloud.NewChallengeProvider(&pJDCloud.ChallengeProviderConfig{
				AccessKeyId:           access.AccessKeyId,
				AccessKeySecret:       access.AccessKeySecret,
				RegionId:              maputil.GetString(options.ProviderApplyConfig, "regionId"),
				DnsPropagationTimeout: options.DnsPropagationTimeout,
				DnsTTL:                options.DnsTTL,
			})
			return applicant, err
		}

	case domain.ApplyDNSProviderTypeNamecheap:
		{
			access := domain.AccessConfigForNamecheap{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			applicant, err := pNamecheap.NewChallengeProvider(&pNamecheap.ChallengeProviderConfig{
				Username:              access.Username,
				ApiKey:                access.ApiKey,
				DnsPropagationTimeout: options.DnsPropagationTimeout,
				DnsTTL:                options.DnsTTL,
			})
			return applicant, err
		}

	case domain.ApplyDNSProviderTypeNameDotCom:
		{
			access := domain.AccessConfigForNameDotCom{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			applicant, err := pNameDotCom.NewChallengeProvider(&pNameDotCom.ChallengeProviderConfig{
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
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			applicant, err := pNameSilo.NewChallengeProvider(&pNameSilo.ChallengeProviderConfig{
				ApiKey:                access.ApiKey,
				DnsPropagationTimeout: options.DnsPropagationTimeout,
				DnsTTL:                options.DnsTTL,
			})
			return applicant, err
		}

	case domain.ApplyDNSProviderTypeNS1:
		{
			access := domain.AccessConfigForNS1{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			applicant, err := pNS1.NewChallengeProvider(&pNS1.ChallengeProviderConfig{
				ApiKey:                access.ApiKey,
				DnsPropagationTimeout: options.DnsPropagationTimeout,
				DnsTTL:                options.DnsTTL,
			})
			return applicant, err
		}

	case domain.ApplyDNSProviderTypePorkbun:
		{
			access := domain.AccessConfigForPorkbun{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			applicant, err := pPorkbun.NewChallengeProvider(&pPorkbun.ChallengeProviderConfig{
				ApiKey:                access.ApiKey,
				SecretApiKey:          access.SecretApiKey,
				DnsPropagationTimeout: options.DnsPropagationTimeout,
				DnsTTL:                options.DnsTTL,
			})
			return applicant, err
		}

	case domain.ApplyDNSProviderTypePowerDNS:
		{
			access := domain.AccessConfigForPowerDNS{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			applicant, err := pPowerDNS.NewChallengeProvider(&pPowerDNS.ChallengeProviderConfig{
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
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			applicant, err := pRainYun.NewChallengeProvider(&pRainYun.ChallengeProviderConfig{
				ApiKey:                access.ApiKey,
				DnsPropagationTimeout: options.DnsPropagationTimeout,
				DnsTTL:                options.DnsTTL,
			})
			return applicant, err
		}

	case domain.ApplyDNSProviderTypeTencentCloud, domain.ApplyDNSProviderTypeTencentCloudDNS, domain.ApplyDNSProviderTypeTencentCloudEO:
		{
			access := domain.AccessConfigForTencentCloud{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.ApplyDNSProviderTypeTencentCloud, domain.ApplyDNSProviderTypeTencentCloudDNS:
				applicant, err := pTencentCloud.NewChallengeProvider(&pTencentCloud.ChallengeProviderConfig{
					SecretId:              access.SecretId,
					SecretKey:             access.SecretKey,
					DnsPropagationTimeout: options.DnsPropagationTimeout,
					DnsTTL:                options.DnsTTL,
				})
				return applicant, err

			case domain.ApplyDNSProviderTypeTencentCloudEO:
				applicant, err := pTencentCloudEO.NewChallengeProvider(&pTencentCloudEO.ChallengeProviderConfig{
					SecretId:              access.SecretId,
					SecretKey:             access.SecretKey,
					ZoneId:                maputil.GetString(options.ProviderApplyConfig, "zoneId"),
					DnsPropagationTimeout: options.DnsPropagationTimeout,
					DnsTTL:                options.DnsTTL,
				})
				return applicant, err

			default:
				break
			}
		}

	case domain.ApplyDNSProviderTypeVercel:
		{
			access := domain.AccessConfigForVercel{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			applicant, err := pVercel.NewChallengeProvider(&pVercel.ChallengeProviderConfig{
				ApiAccessToken:        access.ApiAccessToken,
				TeamId:                access.TeamId,
				DnsPropagationTimeout: options.DnsPropagationTimeout,
				DnsTTL:                options.DnsTTL,
			})
			return applicant, err
		}

	case domain.ApplyDNSProviderTypeVolcEngine, domain.ApplyDNSProviderTypeVolcEngineDNS:
		{
			access := domain.AccessConfigForVolcEngine{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			applicant, err := pVolcEngine.NewChallengeProvider(&pVolcEngine.ChallengeProviderConfig{
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
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			applicant, err := pWestcn.NewChallengeProvider(&pWestcn.ChallengeProviderConfig{
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
