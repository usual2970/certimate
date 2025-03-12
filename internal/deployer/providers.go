package deployer

import (
	"fmt"
	"strings"

	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	p1PanelConsole "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/1panel-console"
	p1PanelSite "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/1panel-site"
	pAliyunALB "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/aliyun-alb"
	pAliyunCASDeploy "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/aliyun-cas-deploy"
	pAliyunCDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/aliyun-cdn"
	pAliyunCLB "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/aliyun-clb"
	pAliyunDCDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/aliyun-dcdn"
	pAliyunESA "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/aliyun-esa"
	pAliyunFC "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/aliyun-fc"
	pAliyunLive "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/aliyun-live"
	pAliyunNLB "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/aliyun-nlb"
	pAliyunOSS "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/aliyun-oss"
	pAliyunVOD "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/aliyun-vod"
	pAliyunWAF "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/aliyun-waf"
	pAWSCloudFront "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/aws-cloudfront"
	pBaiduCloudCDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/baiducloud-cdn"
	pBaishanCDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/baishan-cdn"
	pBaotaPanelConsole "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/baotapanel-console"
	pBaotaPanelSite "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/baotapanel-site"
	pBytePlusCDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/byteplus-cdn"
	pCacheFly "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/cachefly"
	pCdnfly "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/cdnfly"
	pDogeCDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/dogecloud-cdn"
	pEdgioApplications "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/edgio-applications"
	pGcoreCDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/gcore-cdn"
	pHuaweiCloudCDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/huaweicloud-cdn"
	pHuaweiCloudELB "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/huaweicloud-elb"
	pHuaweiCloudWAF "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/huaweicloud-waf"
	pJDCloudALB "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/jdcloud-alb"
	pJDCloudCDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/jdcloud-cdn"
	pJDCloudLive "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/jdcloud-live"
	pJDCloudVOD "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/jdcloud-vod"
	pK8sSecret "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/k8s-secret"
	pLocal "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/local"
	pQiniuCDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/qiniu-cdn"
	pQiniuPili "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/qiniu-pili"
	pSafeLine "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/safeline"
	pSSH "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/ssh"
	pTencentCloudCDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/tencentcloud-cdn"
	pTencentCloudCLB "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/tencentcloud-clb"
	pTencentCloudCOS "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/tencentcloud-cos"
	pTencentCloudCSS "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/tencentcloud-css"
	pTencentCloudECDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/tencentcloud-ecdn"
	pTencentCloudEO "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/tencentcloud-eo"
	pTencentCloudSCF "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/tencentcloud-scf"
	pTencentCloudSSLDeploy "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/tencentcloud-ssl-deploy"
	pTencentCloudVOD "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/tencentcloud-vod"
	pTencentCloudWAF "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/tencentcloud-waf"
	pUCloudUCDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/ucloud-ucdn"
	pUCloudUS3 "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/ucloud-us3"
	pVolcEngineCDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/volcengine-cdn"
	pVolcEngineCLB "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/volcengine-clb"
	pVolcEngineDCDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/volcengine-dcdn"
	pVolcEngineImageX "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/volcengine-imagex"
	pVolcEngineLive "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/volcengine-live"
	pVolcEngineTOS "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/volcengine-tos"
	pWebhook "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/webhook"
	"github.com/usual2970/certimate/internal/pkg/utils/maps"
	"github.com/usual2970/certimate/internal/pkg/utils/slices"
)

func createDeployer(options *deployerOptions) (deployer.Deployer, error) {
	/*
	  注意：如果追加新的常量值，请保持以 ASCII 排序。
	  NOTICE: If you add new constant, please keep ASCII order.
	*/
	switch options.Provider {
	case domain.DeployProviderType1PanelConsole, domain.DeployProviderType1PanelSite:
		{
			access := domain.AccessConfigFor1Panel{}
			if err := maps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeployProviderType1PanelConsole:
				deployer, err := p1PanelConsole.NewDeployer(&p1PanelConsole.DeployerConfig{
					ApiUrl:                   access.ApiUrl,
					ApiKey:                   access.ApiKey,
					AllowInsecureConnections: access.AllowInsecureConnections,
					AutoRestart:              maps.GetValueAsBool(options.ProviderDeployConfig, "autoRestart"),
				})
				return deployer, err

			case domain.DeployProviderType1PanelSite:
				deployer, err := p1PanelSite.NewDeployer(&p1PanelSite.DeployerConfig{
					ApiUrl:                   access.ApiUrl,
					ApiKey:                   access.ApiKey,
					AllowInsecureConnections: access.AllowInsecureConnections,
					WebsiteId:                maps.GetValueAsInt64(options.ProviderDeployConfig, "websiteId"),
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeployProviderTypeAliyunALB, domain.DeployProviderTypeAliyunCASDeploy, domain.DeployProviderTypeAliyunCDN, domain.DeployProviderTypeAliyunCLB, domain.DeployProviderTypeAliyunDCDN, domain.DeployProviderTypeAliyunESA, domain.DeployProviderTypeAliyunFC, domain.DeployProviderTypeAliyunLive, domain.DeployProviderTypeAliyunNLB, domain.DeployProviderTypeAliyunOSS, domain.DeployProviderTypeAliyunVOD, domain.DeployProviderTypeAliyunWAF:
		{
			access := domain.AccessConfigForAliyun{}
			if err := maps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeployProviderTypeAliyunALB:
				deployer, err := pAliyunALB.NewDeployer(&pAliyunALB.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Region:          maps.GetValueAsString(options.ProviderDeployConfig, "region"),
					ResourceType:    pAliyunALB.ResourceType(maps.GetValueAsString(options.ProviderDeployConfig, "resourceType")),
					LoadbalancerId:  maps.GetValueAsString(options.ProviderDeployConfig, "loadbalancerId"),
					ListenerId:      maps.GetValueAsString(options.ProviderDeployConfig, "listenerId"),
					Domain:          maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			case domain.DeployProviderTypeAliyunCASDeploy:
				deployer, err := pAliyunCASDeploy.NewDeployer(&pAliyunCASDeploy.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Region:          maps.GetValueAsString(options.ProviderDeployConfig, "region"),
					ResourceIds:     slices.Filter(strings.Split(maps.GetValueAsString(options.ProviderDeployConfig, "resourceIds"), ";"), func(s string) bool { return s != "" }),
					ContactIds:      slices.Filter(strings.Split(maps.GetValueAsString(options.ProviderDeployConfig, "contactIds"), ";"), func(s string) bool { return s != "" }),
				})
				return deployer, err

			case domain.DeployProviderTypeAliyunCDN:
				deployer, err := pAliyunCDN.NewDeployer(&pAliyunCDN.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Domain:          maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			case domain.DeployProviderTypeAliyunCLB:
				deployer, err := pAliyunCLB.NewDeployer(&pAliyunCLB.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Region:          maps.GetValueAsString(options.ProviderDeployConfig, "region"),
					ResourceType:    pAliyunCLB.ResourceType(maps.GetValueAsString(options.ProviderDeployConfig, "resourceType")),
					LoadbalancerId:  maps.GetValueAsString(options.ProviderDeployConfig, "loadbalancerId"),
					ListenerPort:    maps.GetValueOrDefaultAsInt32(options.ProviderDeployConfig, "listenerPort", 443),
					Domain:          maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			case domain.DeployProviderTypeAliyunDCDN:
				deployer, err := pAliyunDCDN.NewDeployer(&pAliyunDCDN.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Domain:          maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			case domain.DeployProviderTypeAliyunESA:
				deployer, err := pAliyunESA.NewDeployer(&pAliyunESA.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Region:          maps.GetValueAsString(options.ProviderDeployConfig, "region"),
					SiteId:          maps.GetValueAsInt64(options.ProviderDeployConfig, "siteId"),
				})
				return deployer, err

			case domain.DeployProviderTypeAliyunFC:
				deployer, err := pAliyunFC.NewDeployer(&pAliyunFC.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Region:          maps.GetValueAsString(options.ProviderDeployConfig, "region"),
					ServiceVersion:  maps.GetValueAsString(options.ProviderDeployConfig, "serviceVersion"),
					Domain:          maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			case domain.DeployProviderTypeAliyunLive:
				deployer, err := pAliyunLive.NewDeployer(&pAliyunLive.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Region:          maps.GetValueAsString(options.ProviderDeployConfig, "region"),
					Domain:          maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			case domain.DeployProviderTypeAliyunNLB:
				deployer, err := pAliyunNLB.NewDeployer(&pAliyunNLB.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Region:          maps.GetValueAsString(options.ProviderDeployConfig, "region"),
					ResourceType:    pAliyunNLB.ResourceType(maps.GetValueAsString(options.ProviderDeployConfig, "resourceType")),
					LoadbalancerId:  maps.GetValueAsString(options.ProviderDeployConfig, "loadbalancerId"),
					ListenerId:      maps.GetValueAsString(options.ProviderDeployConfig, "listenerId"),
				})
				return deployer, err

			case domain.DeployProviderTypeAliyunOSS:
				deployer, err := pAliyunOSS.NewDeployer(&pAliyunOSS.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Region:          maps.GetValueAsString(options.ProviderDeployConfig, "region"),
					Bucket:          maps.GetValueAsString(options.ProviderDeployConfig, "bucket"),
					Domain:          maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			case domain.DeployProviderTypeAliyunVOD:
				deployer, err := pAliyunVOD.NewDeployer(&pAliyunVOD.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Region:          maps.GetValueAsString(options.ProviderDeployConfig, "region"),
					Domain:          maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			case domain.DeployProviderTypeAliyunWAF:
				deployer, err := pAliyunWAF.NewDeployer(&pAliyunWAF.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Region:          maps.GetValueAsString(options.ProviderDeployConfig, "region"),
					InstanceId:      maps.GetValueAsString(options.ProviderDeployConfig, "instanceId"),
					Domain:          maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeployProviderTypeAWSCloudFront:
		{
			access := domain.AccessConfigForAWS{}
			if err := maps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeployProviderTypeAWSCloudFront:
				deployer, err := pAWSCloudFront.NewDeployer(&pAWSCloudFront.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					SecretAccessKey: access.SecretAccessKey,
					Region:          maps.GetValueAsString(options.ProviderDeployConfig, "region"),
					DistributionId:  maps.GetValueAsString(options.ProviderDeployConfig, "distributionId"),
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeployProviderTypeBaiduCloudCDN:
		{
			access := domain.AccessConfigForBaiduCloud{}
			if err := maps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeployProviderTypeBaiduCloudCDN:
				deployer, err := pBaiduCloudCDN.NewDeployer(&pBaiduCloudCDN.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					SecretAccessKey: access.SecretAccessKey,
					Domain:          maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeployProviderTypeBaishanCDN:
		{
			access := domain.AccessConfigForBaishan{}
			if err := maps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeployProviderTypeBaishanCDN:
				deployer, err := pBaishanCDN.NewDeployer(&pBaishanCDN.DeployerConfig{
					ApiToken: access.ApiToken,
					Domain:   maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeployProviderTypeBaotaPanelConsole, domain.DeployProviderTypeBaotaPanelSite:
		{
			access := domain.AccessConfigForBaotaPanel{}
			if err := maps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeployProviderTypeBaotaPanelConsole:
				deployer, err := pBaotaPanelConsole.NewDeployer(&pBaotaPanelConsole.DeployerConfig{
					ApiUrl:                   access.ApiUrl,
					ApiKey:                   access.ApiKey,
					AllowInsecureConnections: access.AllowInsecureConnections,
					AutoRestart:              maps.GetValueAsBool(options.ProviderDeployConfig, "autoRestart"),
				})
				return deployer, err

			case domain.DeployProviderTypeBaotaPanelSite:
				deployer, err := pBaotaPanelSite.NewDeployer(&pBaotaPanelSite.DeployerConfig{
					ApiUrl:                   access.ApiUrl,
					ApiKey:                   access.ApiKey,
					AllowInsecureConnections: access.AllowInsecureConnections,
					SiteType:                 maps.GetValueOrDefaultAsString(options.ProviderDeployConfig, "siteType", "other"),
					SiteName:                 maps.GetValueAsString(options.ProviderDeployConfig, "siteName"),
					SiteNames:                slices.Filter(strings.Split(maps.GetValueAsString(options.ProviderDeployConfig, "siteNames"), ";"), func(s string) bool { return s != "" }),
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeployProviderTypeBytePlusCDN:
		{
			access := domain.AccessConfigForBytePlus{}
			if err := maps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeployProviderTypeBytePlusCDN:
				deployer, err := pBytePlusCDN.NewDeployer(&pBytePlusCDN.DeployerConfig{
					AccessKey: access.AccessKey,
					SecretKey: access.SecretKey,
					Domain:    maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeployProviderTypeCacheFly:
		{
			access := domain.AccessConfigForCacheFly{}
			if err := maps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			deployer, err := pCacheFly.NewDeployer(&pCacheFly.DeployerConfig{
				ApiToken: access.ApiToken,
			})
			return deployer, err
		}

	case domain.DeployProviderTypeCdnfly:
		{
			access := domain.AccessConfigForCdnfly{}
			if err := maps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			deployer, err := pCdnfly.NewDeployer(&pCdnfly.DeployerConfig{
				ApiUrl:        access.ApiUrl,
				ApiKey:        access.ApiKey,
				ApiSecret:     access.ApiSecret,
				ResourceType:  pCdnfly.ResourceType(maps.GetValueAsString(options.ProviderDeployConfig, "resourceType")),
				SiteId:        maps.GetValueAsString(options.ProviderDeployConfig, "siteId"),
				CertificateId: maps.GetValueAsString(options.ProviderDeployConfig, "certificateId"),
			})
			return deployer, err
		}

	case domain.DeployProviderTypeDogeCloudCDN:
		{
			access := domain.AccessConfigForDogeCloud{}
			if err := maps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			deployer, err := pDogeCDN.NewDeployer(&pDogeCDN.DeployerConfig{
				AccessKey: access.AccessKey,
				SecretKey: access.SecretKey,
				Domain:    maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
			})
			return deployer, err
		}

	case domain.DeployProviderTypeEdgioApplications:
		{
			access := domain.AccessConfigForEdgio{}
			if err := maps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			deployer, err := pEdgioApplications.NewDeployer(&pEdgioApplications.DeployerConfig{
				ClientId:      access.ClientId,
				ClientSecret:  access.ClientSecret,
				EnvironmentId: maps.GetValueAsString(options.ProviderDeployConfig, "environmentId"),
			})
			return deployer, err
		}

	case domain.DeployProviderTypeGcoreCDN:
		{
			access := domain.AccessConfigForGcore{}
			if err := maps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeployProviderTypeGcoreCDN:
				deployer, err := pGcoreCDN.NewDeployer(&pGcoreCDN.DeployerConfig{
					ApiToken:   access.ApiToken,
					ResourceId: maps.GetValueAsInt64(options.ProviderDeployConfig, "resourceId"),
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeployProviderTypeHuaweiCloudCDN, domain.DeployProviderTypeHuaweiCloudELB, domain.DeployProviderTypeHuaweiCloudWAF:
		{
			access := domain.AccessConfigForHuaweiCloud{}
			if err := maps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeployProviderTypeHuaweiCloudCDN:
				deployer, err := pHuaweiCloudCDN.NewDeployer(&pHuaweiCloudCDN.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					SecretAccessKey: access.SecretAccessKey,
					Region:          maps.GetValueAsString(options.ProviderDeployConfig, "region"),
					Domain:          maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			case domain.DeployProviderTypeHuaweiCloudELB:
				deployer, err := pHuaweiCloudELB.NewDeployer(&pHuaweiCloudELB.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					SecretAccessKey: access.SecretAccessKey,
					Region:          maps.GetValueAsString(options.ProviderDeployConfig, "region"),
					ResourceType:    pHuaweiCloudELB.ResourceType(maps.GetValueAsString(options.ProviderDeployConfig, "resourceType")),
					CertificateId:   maps.GetValueAsString(options.ProviderDeployConfig, "certificateId"),
					LoadbalancerId:  maps.GetValueAsString(options.ProviderDeployConfig, "loadbalancerId"),
					ListenerId:      maps.GetValueAsString(options.ProviderDeployConfig, "listenerId"),
				})
				return deployer, err

			case domain.DeployProviderTypeHuaweiCloudWAF:
				deployer, err := pHuaweiCloudWAF.NewDeployer(&pHuaweiCloudWAF.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					SecretAccessKey: access.SecretAccessKey,
					Region:          maps.GetValueAsString(options.ProviderDeployConfig, "region"),
					ResourceType:    pHuaweiCloudWAF.ResourceType(maps.GetValueAsString(options.ProviderDeployConfig, "resourceType")),
					CertificateId:   maps.GetValueAsString(options.ProviderDeployConfig, "certificateId"),
					Domain:          maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeployProviderTypeJDCloudALB, domain.DeployProviderTypeJDCloudCDN, domain.DeployProviderTypeJDCloudLive, domain.DeployProviderTypeJDCloudVOD:
		{
			access := domain.AccessConfigForJDCloud{}
			if err := maps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeployProviderTypeJDCloudALB:
				deployer, err := pJDCloudALB.NewDeployer(&pJDCloudALB.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					RegionId:        maps.GetValueAsString(options.ProviderDeployConfig, "regionId"),
					ResourceType:    pJDCloudALB.ResourceType(maps.GetValueAsString(options.ProviderDeployConfig, "resourceType")),
					LoadbalancerId:  maps.GetValueAsString(options.ProviderDeployConfig, "loadbalancerId"),
					ListenerId:      maps.GetValueAsString(options.ProviderDeployConfig, "listenerId"),
				})
				return deployer, err

			case domain.DeployProviderTypeJDCloudCDN:
				deployer, err := pJDCloudCDN.NewDeployer(&pJDCloudCDN.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Domain:          maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			case domain.DeployProviderTypeJDCloudLive:
				deployer, err := pJDCloudLive.NewDeployer(&pJDCloudLive.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Domain:          maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			case domain.DeployProviderTypeJDCloudVOD:
				deployer, err := pJDCloudVOD.NewDeployer(&pJDCloudVOD.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Domain:          maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeployProviderTypeLocal:
		{
			deployer, err := pLocal.NewDeployer(&pLocal.DeployerConfig{
				ShellEnv:       pLocal.ShellEnvType(maps.GetValueAsString(options.ProviderDeployConfig, "shellEnv")),
				PreCommand:     maps.GetValueAsString(options.ProviderDeployConfig, "preCommand"),
				PostCommand:    maps.GetValueAsString(options.ProviderDeployConfig, "postCommand"),
				OutputFormat:   pLocal.OutputFormatType(maps.GetValueOrDefaultAsString(options.ProviderDeployConfig, "format", string(pLocal.OUTPUT_FORMAT_PEM))),
				OutputCertPath: maps.GetValueAsString(options.ProviderDeployConfig, "certPath"),
				OutputKeyPath:  maps.GetValueAsString(options.ProviderDeployConfig, "keyPath"),
				PfxPassword:    maps.GetValueAsString(options.ProviderDeployConfig, "pfxPassword"),
				JksAlias:       maps.GetValueAsString(options.ProviderDeployConfig, "jksAlias"),
				JksKeypass:     maps.GetValueAsString(options.ProviderDeployConfig, "jksKeypass"),
				JksStorepass:   maps.GetValueAsString(options.ProviderDeployConfig, "jksStorepass"),
			})
			return deployer, err
		}

	case domain.DeployProviderTypeKubernetesSecret:
		{
			access := domain.AccessConfigForKubernetes{}
			if err := maps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			deployer, err := pK8sSecret.NewDeployer(&pK8sSecret.DeployerConfig{
				KubeConfig:          access.KubeConfig,
				Namespace:           maps.GetValueOrDefaultAsString(options.ProviderDeployConfig, "namespace", "default"),
				SecretName:          maps.GetValueAsString(options.ProviderDeployConfig, "secretName"),
				SecretType:          maps.GetValueOrDefaultAsString(options.ProviderDeployConfig, "secretType", "kubernetes.io/tls"),
				SecretDataKeyForCrt: maps.GetValueOrDefaultAsString(options.ProviderDeployConfig, "secretDataKeyForCrt", "tls.crt"),
				SecretDataKeyForKey: maps.GetValueOrDefaultAsString(options.ProviderDeployConfig, "secretDataKeyForKey", "tls.key"),
			})
			return deployer, err
		}

	case domain.DeployProviderTypeQiniuCDN, domain.DeployProviderTypeQiniuPili:
		{
			access := domain.AccessConfigForQiniu{}
			if err := maps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeployProviderTypeQiniuCDN:
				deployer, err := pQiniuCDN.NewDeployer(&pQiniuCDN.DeployerConfig{
					AccessKey: access.AccessKey,
					SecretKey: access.SecretKey,
					Domain:    maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			case domain.DeployProviderTypeQiniuPili:
				deployer, err := pQiniuPili.NewDeployer(&pQiniuPili.DeployerConfig{
					AccessKey: access.AccessKey,
					SecretKey: access.SecretKey,
					Hub:       maps.GetValueAsString(options.ProviderDeployConfig, "hub"),
					Domain:    maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeployProviderTypeSafeLine:
		{
			access := domain.AccessConfigForSafeLine{}
			if err := maps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			deployer, err := pSafeLine.NewDeployer(&pSafeLine.DeployerConfig{
				ApiUrl:                   access.ApiUrl,
				ApiToken:                 access.ApiToken,
				AllowInsecureConnections: access.AllowInsecureConnections,
				ResourceType:             pSafeLine.ResourceType(maps.GetValueAsString(options.ProviderDeployConfig, "resourceType")),
				CertificateId:            maps.GetValueAsInt32(options.ProviderDeployConfig, "certificateId"),
			})
			return deployer, err
		}

	case domain.DeployProviderTypeSSH:
		{
			access := domain.AccessConfigForSSH{}
			if err := maps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			deployer, err := pSSH.NewDeployer(&pSSH.DeployerConfig{
				SshHost:          access.Host,
				SshPort:          access.Port,
				SshUsername:      access.Username,
				SshPassword:      access.Password,
				SshKey:           access.Key,
				SshKeyPassphrase: access.KeyPassphrase,
				UseSCP:           maps.GetValueAsBool(options.ProviderDeployConfig, "useSCP"),
				PreCommand:       maps.GetValueAsString(options.ProviderDeployConfig, "preCommand"),
				PostCommand:      maps.GetValueAsString(options.ProviderDeployConfig, "postCommand"),
				OutputFormat:     pSSH.OutputFormatType(maps.GetValueOrDefaultAsString(options.ProviderDeployConfig, "format", string(pSSH.OUTPUT_FORMAT_PEM))),
				OutputCertPath:   maps.GetValueAsString(options.ProviderDeployConfig, "certPath"),
				OutputKeyPath:    maps.GetValueAsString(options.ProviderDeployConfig, "keyPath"),
				PfxPassword:      maps.GetValueAsString(options.ProviderDeployConfig, "pfxPassword"),
				JksAlias:         maps.GetValueAsString(options.ProviderDeployConfig, "jksAlias"),
				JksKeypass:       maps.GetValueAsString(options.ProviderDeployConfig, "jksKeypass"),
				JksStorepass:     maps.GetValueAsString(options.ProviderDeployConfig, "jksStorepass"),
			})
			return deployer, err
		}

	case domain.DeployProviderTypeTencentCloudCDN, domain.DeployProviderTypeTencentCloudCLB, domain.DeployProviderTypeTencentCloudCOS, domain.DeployProviderTypeTencentCloudCSS, domain.DeployProviderTypeTencentCloudECDN, domain.DeployProviderTypeTencentCloudEO, domain.DeployProviderTypeTencentCloudSCF, domain.DeployProviderTypeTencentCloudSSLDeploy, domain.DeployProviderTypeTencentCloudVOD, domain.DeployProviderTypeTencentCloudWAF:
		{
			access := domain.AccessConfigForTencentCloud{}
			if err := maps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeployProviderTypeTencentCloudCDN:
				deployer, err := pTencentCloudCDN.NewDeployer(&pTencentCloudCDN.DeployerConfig{
					SecretId:  access.SecretId,
					SecretKey: access.SecretKey,
					Domain:    maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			case domain.DeployProviderTypeTencentCloudCLB:
				deployer, err := pTencentCloudCLB.NewDeployer(&pTencentCloudCLB.DeployerConfig{
					SecretId:       access.SecretId,
					SecretKey:      access.SecretKey,
					Region:         maps.GetValueAsString(options.ProviderDeployConfig, "region"),
					ResourceType:   pTencentCloudCLB.ResourceType(maps.GetValueAsString(options.ProviderDeployConfig, "resourceType")),
					LoadbalancerId: maps.GetValueAsString(options.ProviderDeployConfig, "loadbalancerId"),
					ListenerId:     maps.GetValueAsString(options.ProviderDeployConfig, "listenerId"),
					Domain:         maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			case domain.DeployProviderTypeTencentCloudCOS:
				deployer, err := pTencentCloudCOS.NewDeployer(&pTencentCloudCOS.DeployerConfig{
					SecretId:  access.SecretId,
					SecretKey: access.SecretKey,
					Region:    maps.GetValueAsString(options.ProviderDeployConfig, "region"),
					Bucket:    maps.GetValueAsString(options.ProviderDeployConfig, "bucket"),
					Domain:    maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			case domain.DeployProviderTypeTencentCloudCSS:
				deployer, err := pTencentCloudCSS.NewDeployer(&pTencentCloudCSS.DeployerConfig{
					SecretId:  access.SecretId,
					SecretKey: access.SecretKey,
					Domain:    maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			case domain.DeployProviderTypeTencentCloudECDN:
				deployer, err := pTencentCloudECDN.NewDeployer(&pTencentCloudECDN.DeployerConfig{
					SecretId:  access.SecretId,
					SecretKey: access.SecretKey,
					Domain:    maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			case domain.DeployProviderTypeTencentCloudEO:
				deployer, err := pTencentCloudEO.NewDeployer(&pTencentCloudEO.DeployerConfig{
					SecretId:  access.SecretId,
					SecretKey: access.SecretKey,
					ZoneId:    maps.GetValueAsString(options.ProviderDeployConfig, "zoneId"),
					Domain:    maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			case domain.DeployProviderTypeTencentCloudSCF:
				deployer, err := pTencentCloudSCF.NewDeployer(&pTencentCloudSCF.DeployerConfig{
					SecretId:  access.SecretId,
					SecretKey: access.SecretKey,
					Region:    maps.GetValueAsString(options.ProviderDeployConfig, "region"),
					Domain:    maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			case domain.DeployProviderTypeTencentCloudSSLDeploy:
				deployer, err := pTencentCloudSSLDeploy.NewDeployer(&pTencentCloudSSLDeploy.DeployerConfig{
					SecretId:     access.SecretId,
					SecretKey:    access.SecretKey,
					Region:       maps.GetValueAsString(options.ProviderDeployConfig, "region"),
					ResourceType: maps.GetValueAsString(options.ProviderDeployConfig, "resourceType"),
					ResourceIds:  slices.Filter(strings.Split(maps.GetValueAsString(options.ProviderDeployConfig, "resourceIds"), ";"), func(s string) bool { return s != "" }),
				})
				return deployer, err

			case domain.DeployProviderTypeTencentCloudVOD:
				deployer, err := pTencentCloudVOD.NewDeployer(&pTencentCloudVOD.DeployerConfig{
					SecretId:  access.SecretId,
					SecretKey: access.SecretKey,
					SubAppId:  maps.GetValueAsInt64(options.ProviderDeployConfig, "subAppId"),
					Domain:    maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			case domain.DeployProviderTypeTencentCloudWAF:
				deployer, err := pTencentCloudWAF.NewDeployer(&pTencentCloudWAF.DeployerConfig{
					SecretId:   access.SecretId,
					SecretKey:  access.SecretKey,
					Domain:     maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
					DomainId:   maps.GetValueAsString(options.ProviderDeployConfig, "domainId"),
					InstanceId: maps.GetValueAsString(options.ProviderDeployConfig, "instanceId"),
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeployProviderTypeUCloudUCDN, domain.DeployProviderTypeUCloudUS3:
		{
			access := domain.AccessConfigForUCloud{}
			if err := maps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeployProviderTypeUCloudUCDN:
				deployer, err := pUCloudUCDN.NewDeployer(&pUCloudUCDN.DeployerConfig{
					PrivateKey: access.PrivateKey,
					PublicKey:  access.PublicKey,
					ProjectId:  access.ProjectId,
					DomainId:   maps.GetValueAsString(options.ProviderDeployConfig, "domainId"),
				})
				return deployer, err

			case domain.DeployProviderTypeUCloudUS3:
				deployer, err := pUCloudUS3.NewDeployer(&pUCloudUS3.DeployerConfig{
					PrivateKey: access.PrivateKey,
					PublicKey:  access.PublicKey,
					ProjectId:  access.ProjectId,
					Region:     maps.GetValueAsString(options.ProviderDeployConfig, "region"),
					Bucket:     maps.GetValueAsString(options.ProviderDeployConfig, "bucket"),
					Domain:     maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeployProviderTypeVolcEngineCDN, domain.DeployProviderTypeVolcEngineCLB, domain.DeployProviderTypeVolcEngineDCDN, domain.DeployProviderTypeVolcEngineImageX, domain.DeployProviderTypeVolcEngineLive, domain.DeployProviderTypeVolcEngineTOS:
		{
			access := domain.AccessConfigForVolcEngine{}
			if err := maps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeployProviderTypeVolcEngineCDN:
				deployer, err := pVolcEngineCDN.NewDeployer(&pVolcEngineCDN.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.SecretAccessKey,
					Domain:          maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			case domain.DeployProviderTypeVolcEngineCLB:
				deployer, err := pVolcEngineCLB.NewDeployer(&pVolcEngineCLB.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.SecretAccessKey,
					Region:          maps.GetValueAsString(options.ProviderDeployConfig, "region"),
					ResourceType:    pVolcEngineCLB.ResourceType(maps.GetValueAsString(options.ProviderDeployConfig, "resourceType")),
					ListenerId:      maps.GetValueAsString(options.ProviderDeployConfig, "listenerId"),
				})
				return deployer, err

			case domain.DeployProviderTypeVolcEngineDCDN:
				deployer, err := pVolcEngineDCDN.NewDeployer(&pVolcEngineDCDN.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.SecretAccessKey,
					Domain:          maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			case domain.DeployProviderTypeVolcEngineImageX:
				deployer, err := pVolcEngineImageX.NewDeployer(&pVolcEngineImageX.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.SecretAccessKey,
					Region:          maps.GetValueAsString(options.ProviderDeployConfig, "region"),
					ServiceId:       maps.GetValueAsString(options.ProviderDeployConfig, "serviceId"),
					Domain:          maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			case domain.DeployProviderTypeVolcEngineLive:
				deployer, err := pVolcEngineLive.NewDeployer(&pVolcEngineLive.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.SecretAccessKey,
					Domain:          maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			case domain.DeployProviderTypeVolcEngineTOS:
				deployer, err := pVolcEngineTOS.NewDeployer(&pVolcEngineTOS.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.SecretAccessKey,
					Region:          maps.GetValueAsString(options.ProviderDeployConfig, "region"),
					Bucket:          maps.GetValueAsString(options.ProviderDeployConfig, "bucket"),
					Domain:          maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeployProviderTypeWebhook:
		{
			access := domain.AccessConfigForWebhook{}
			if err := maps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			deployer, err := pWebhook.NewDeployer(&pWebhook.DeployerConfig{
				WebhookUrl:               access.Url,
				WebhookData:              maps.GetValueAsString(options.ProviderDeployConfig, "webhookData"),
				AllowInsecureConnections: access.AllowInsecureConnections,
			})
			return deployer, err
		}
	}

	return nil, fmt.Errorf("unsupported deployer provider: %s", string(options.Provider))
}
