package deployer

import (
	"fmt"
	"strings"

	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	p1PanelConsole "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/1panel-console"
	p1PanelSite "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/1panel-site"
	pAliyunALB "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/aliyun-alb"
	pAliyunAPIGW "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/aliyun-apigw"
	pAliyunCAS "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/aliyun-cas"
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
	pAWSACM "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/aws-acm"
	pAWSCloudFront "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/aws-cloudfront"
	pAzureKeyVault "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/azure-keyvault"
	pBaiduCloudAppBLB "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/baiducloud-appblb"
	pBaiduCloudBLB "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/baiducloud-blb"
	pBaiduCloudCDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/baiducloud-cdn"
	pBaiduCloudCert "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/baiducloud-cert"
	pBaishanCDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/baishan-cdn"
	pBaotaPanelConsole "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/baotapanel-console"
	pBaotaPanelSite "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/baotapanel-site"
	pBunnyCDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/bunny-cdn"
	pBytePlusCDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/byteplus-cdn"
	pCacheFly "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/cachefly"
	pCdnfly "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/cdnfly"
	pDogeCDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/dogecloud-cdn"
	pEdgioApplications "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/edgio-applications"
	pGcoreCDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/gcore-cdn"
	pHuaweiCloudCDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/huaweicloud-cdn"
	pHuaweiCloudELB "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/huaweicloud-elb"
	pHuaweiCloudSCM "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/huaweicloud-scm"
	pHuaweiCloudWAF "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/huaweicloud-waf"
	pJDCloudALB "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/jdcloud-alb"
	pJDCloudCDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/jdcloud-cdn"
	pJDCloudLive "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/jdcloud-live"
	pJDCloudVOD "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/jdcloud-vod"
	pK8sSecret "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/k8s-secret"
	pLocal "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/local"
	pQiniuCDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/qiniu-cdn"
	pQiniuPili "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/qiniu-pili"
	pRainYunRCDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/rainyun-rcdn"
	pSafeLine "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/safeline"
	pSSH "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/ssh"
	pTencentCloudCDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/tencentcloud-cdn"
	pTencentCloudCLB "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/tencentcloud-clb"
	pTencentCloudCOS "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/tencentcloud-cos"
	pTencentCloudCSS "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/tencentcloud-css"
	pTencentCloudECDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/tencentcloud-ecdn"
	pTencentCloudEO "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/tencentcloud-eo"
	pTencentCloudSCF "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/tencentcloud-scf"
	pTencentCloudSSL "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/tencentcloud-ssl"
	pTencentCloudSSLDeploy "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/tencentcloud-ssl-deploy"
	pTencentCloudVOD "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/tencentcloud-vod"
	pTencentCloudWAF "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/tencentcloud-waf"
	pUCloudUCDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/ucloud-ucdn"
	pUCloudUS3 "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/ucloud-us3"
	pUpyunCDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/upyun-cdn"
	pVolcEngineALB "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/volcengine-alb"
	pVolcEngineCDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/volcengine-cdn"
	pVolcEngineCertCenter "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/volcengine-certcenter"
	pVolcEngineCLB "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/volcengine-clb"
	pVolcEngineDCDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/volcengine-dcdn"
	pVolcEngineImageX "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/volcengine-imagex"
	pVolcEngineLive "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/volcengine-live"
	pVolcEngineTOS "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/volcengine-tos"
	pWangsuCDNPro "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/wangsu-cdnpro"
	pWebhook "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/webhook"
	"github.com/usual2970/certimate/internal/pkg/utils/maputil"
	"github.com/usual2970/certimate/internal/pkg/utils/sliceutil"
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
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeployProviderType1PanelConsole:
				deployer, err := p1PanelConsole.NewDeployer(&p1PanelConsole.DeployerConfig{
					ApiUrl:                   access.ApiUrl,
					ApiKey:                   access.ApiKey,
					AllowInsecureConnections: access.AllowInsecureConnections,
					AutoRestart:              maputil.GetBool(options.ProviderDeployConfig, "autoRestart"),
				})
				return deployer, err

			case domain.DeployProviderType1PanelSite:
				deployer, err := p1PanelSite.NewDeployer(&p1PanelSite.DeployerConfig{
					ApiUrl:                   access.ApiUrl,
					ApiKey:                   access.ApiKey,
					AllowInsecureConnections: access.AllowInsecureConnections,
					ResourceType:             p1PanelSite.ResourceType(maputil.GetOrDefaultString(options.ProviderDeployConfig, "resourceType", string(p1PanelSite.RESOURCE_TYPE_WEBSITE))),
					WebsiteId:                maputil.GetInt64(options.ProviderDeployConfig, "websiteId"),
					CertificateId:            maputil.GetInt64(options.ProviderDeployConfig, "certificateId"),
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeployProviderTypeAliyunALB, domain.DeployProviderTypeAliyunAPIGW, domain.DeployProviderTypeAliyunCAS, domain.DeployProviderTypeAliyunCASDeploy, domain.DeployProviderTypeAliyunCDN, domain.DeployProviderTypeAliyunCLB, domain.DeployProviderTypeAliyunDCDN, domain.DeployProviderTypeAliyunESA, domain.DeployProviderTypeAliyunFC, domain.DeployProviderTypeAliyunLive, domain.DeployProviderTypeAliyunNLB, domain.DeployProviderTypeAliyunOSS, domain.DeployProviderTypeAliyunVOD, domain.DeployProviderTypeAliyunWAF:
		{
			access := domain.AccessConfigForAliyun{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeployProviderTypeAliyunALB:
				deployer, err := pAliyunALB.NewDeployer(&pAliyunALB.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Region:          maputil.GetString(options.ProviderDeployConfig, "region"),
					ResourceType:    pAliyunALB.ResourceType(maputil.GetString(options.ProviderDeployConfig, "resourceType")),
					LoadbalancerId:  maputil.GetString(options.ProviderDeployConfig, "loadbalancerId"),
					ListenerId:      maputil.GetString(options.ProviderDeployConfig, "listenerId"),
					Domain:          maputil.GetString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			case domain.DeployProviderTypeAliyunAPIGW:
				deployer, err := pAliyunAPIGW.NewDeployer(&pAliyunAPIGW.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Region:          maputil.GetString(options.ProviderDeployConfig, "region"),
					ServiceType:     pAliyunAPIGW.ServiceType(maputil.GetString(options.ProviderDeployConfig, "serviceType")),
					GatewayId:       maputil.GetString(options.ProviderDeployConfig, "gatewayId"),
					GroupId:         maputil.GetString(options.ProviderDeployConfig, "groupId"),
					Domain:          maputil.GetString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			case domain.DeployProviderTypeAliyunCAS:
				deployer, err := pAliyunCAS.NewDeployer(&pAliyunCAS.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Region:          maputil.GetString(options.ProviderDeployConfig, "region"),
				})
				return deployer, err

			case domain.DeployProviderTypeAliyunCASDeploy:
				deployer, err := pAliyunCASDeploy.NewDeployer(&pAliyunCASDeploy.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Region:          maputil.GetString(options.ProviderDeployConfig, "region"),
					ResourceIds:     sliceutil.Filter(strings.Split(maputil.GetString(options.ProviderDeployConfig, "resourceIds"), ";"), func(s string) bool { return s != "" }),
					ContactIds:      sliceutil.Filter(strings.Split(maputil.GetString(options.ProviderDeployConfig, "contactIds"), ";"), func(s string) bool { return s != "" }),
				})
				return deployer, err

			case domain.DeployProviderTypeAliyunCDN:
				deployer, err := pAliyunCDN.NewDeployer(&pAliyunCDN.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Domain:          maputil.GetString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			case domain.DeployProviderTypeAliyunCLB:
				deployer, err := pAliyunCLB.NewDeployer(&pAliyunCLB.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Region:          maputil.GetString(options.ProviderDeployConfig, "region"),
					ResourceType:    pAliyunCLB.ResourceType(maputil.GetString(options.ProviderDeployConfig, "resourceType")),
					LoadbalancerId:  maputil.GetString(options.ProviderDeployConfig, "loadbalancerId"),
					ListenerPort:    maputil.GetOrDefaultInt32(options.ProviderDeployConfig, "listenerPort", 443),
					Domain:          maputil.GetString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			case domain.DeployProviderTypeAliyunDCDN:
				deployer, err := pAliyunDCDN.NewDeployer(&pAliyunDCDN.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Domain:          maputil.GetString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			case domain.DeployProviderTypeAliyunESA:
				deployer, err := pAliyunESA.NewDeployer(&pAliyunESA.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Region:          maputil.GetString(options.ProviderDeployConfig, "region"),
					SiteId:          maputil.GetInt64(options.ProviderDeployConfig, "siteId"),
				})
				return deployer, err

			case domain.DeployProviderTypeAliyunFC:
				deployer, err := pAliyunFC.NewDeployer(&pAliyunFC.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Region:          maputil.GetString(options.ProviderDeployConfig, "region"),
					ServiceVersion:  maputil.GetOrDefaultString(options.ProviderDeployConfig, "serviceVersion", "3.0"),
					Domain:          maputil.GetString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			case domain.DeployProviderTypeAliyunLive:
				deployer, err := pAliyunLive.NewDeployer(&pAliyunLive.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Region:          maputil.GetString(options.ProviderDeployConfig, "region"),
					Domain:          maputil.GetString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			case domain.DeployProviderTypeAliyunNLB:
				deployer, err := pAliyunNLB.NewDeployer(&pAliyunNLB.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Region:          maputil.GetString(options.ProviderDeployConfig, "region"),
					ResourceType:    pAliyunNLB.ResourceType(maputil.GetString(options.ProviderDeployConfig, "resourceType")),
					LoadbalancerId:  maputil.GetString(options.ProviderDeployConfig, "loadbalancerId"),
					ListenerId:      maputil.GetString(options.ProviderDeployConfig, "listenerId"),
				})
				return deployer, err

			case domain.DeployProviderTypeAliyunOSS:
				deployer, err := pAliyunOSS.NewDeployer(&pAliyunOSS.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Region:          maputil.GetString(options.ProviderDeployConfig, "region"),
					Bucket:          maputil.GetString(options.ProviderDeployConfig, "bucket"),
					Domain:          maputil.GetString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			case domain.DeployProviderTypeAliyunVOD:
				deployer, err := pAliyunVOD.NewDeployer(&pAliyunVOD.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Region:          maputil.GetString(options.ProviderDeployConfig, "region"),
					Domain:          maputil.GetString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			case domain.DeployProviderTypeAliyunWAF:
				deployer, err := pAliyunWAF.NewDeployer(&pAliyunWAF.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Region:          maputil.GetString(options.ProviderDeployConfig, "region"),
					ServiceVersion:  maputil.GetOrDefaultString(options.ProviderDeployConfig, "serviceVersion", "3.0"),
					InstanceId:      maputil.GetString(options.ProviderDeployConfig, "instanceId"),
					Domain:          maputil.GetString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeployProviderTypeAWSACM, domain.DeployProviderTypeAWSCloudFront:
		{
			access := domain.AccessConfigForAWS{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeployProviderTypeAWSACM:
				deployer, err := pAWSACM.NewDeployer(&pAWSACM.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					SecretAccessKey: access.SecretAccessKey,
					Region:          maputil.GetString(options.ProviderDeployConfig, "region"),
				})
				return deployer, err

			case domain.DeployProviderTypeAWSCloudFront:
				deployer, err := pAWSCloudFront.NewDeployer(&pAWSCloudFront.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					SecretAccessKey: access.SecretAccessKey,
					Region:          maputil.GetString(options.ProviderDeployConfig, "region"),
					DistributionId:  maputil.GetString(options.ProviderDeployConfig, "distributionId"),
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeployProviderTypeAzureKeyVault:
		{
			access := domain.AccessConfigForAzure{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeployProviderTypeAzureKeyVault:
				deployer, err := pAzureKeyVault.NewDeployer(&pAzureKeyVault.DeployerConfig{
					TenantId:        access.TenantId,
					ClientId:        access.ClientId,
					ClientSecret:    access.ClientSecret,
					CloudName:       access.CloudName,
					KeyVaultName:    maputil.GetString(options.ProviderDeployConfig, "keyvaultName"),
					CertificateName: maputil.GetString(options.ProviderDeployConfig, "certificateName"),
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeployProviderTypeBaiduCloudAppBLB, domain.DeployProviderTypeBaiduCloudBLB, domain.DeployProviderTypeBaiduCloudCDN, domain.DeployProviderTypeBaiduCloudCert:
		{
			access := domain.AccessConfigForBaiduCloud{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeployProviderTypeBaiduCloudAppBLB:
				deployer, err := pBaiduCloudAppBLB.NewDeployer(&pBaiduCloudAppBLB.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					SecretAccessKey: access.SecretAccessKey,
					Region:          maputil.GetString(options.ProviderDeployConfig, "region"),
					ResourceType:    pBaiduCloudAppBLB.ResourceType(maputil.GetString(options.ProviderDeployConfig, "resourceType")),
					LoadbalancerId:  maputil.GetString(options.ProviderDeployConfig, "loadbalancerId"),
					ListenerPort:    maputil.GetInt32(options.ProviderDeployConfig, "listenerPort"),
					Domain:          maputil.GetString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			case domain.DeployProviderTypeBaiduCloudBLB:
				deployer, err := pBaiduCloudBLB.NewDeployer(&pBaiduCloudBLB.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					SecretAccessKey: access.SecretAccessKey,
					Region:          maputil.GetString(options.ProviderDeployConfig, "region"),
					ResourceType:    pBaiduCloudBLB.ResourceType(maputil.GetString(options.ProviderDeployConfig, "resourceType")),
					LoadbalancerId:  maputil.GetString(options.ProviderDeployConfig, "loadbalancerId"),
					ListenerPort:    maputil.GetInt32(options.ProviderDeployConfig, "listenerPort"),
					Domain:          maputil.GetString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			case domain.DeployProviderTypeBaiduCloudCDN:
				deployer, err := pBaiduCloudCDN.NewDeployer(&pBaiduCloudCDN.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					SecretAccessKey: access.SecretAccessKey,
					Domain:          maputil.GetString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			case domain.DeployProviderTypeBaiduCloudCert:
				deployer, err := pBaiduCloudCert.NewDeployer(&pBaiduCloudCert.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					SecretAccessKey: access.SecretAccessKey,
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeployProviderTypeBaishanCDN:
		{
			access := domain.AccessConfigForBaishan{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeployProviderTypeBaishanCDN:
				deployer, err := pBaishanCDN.NewDeployer(&pBaishanCDN.DeployerConfig{
					ApiToken:      access.ApiToken,
					Domain:        maputil.GetString(options.ProviderDeployConfig, "domain"),
					CertificateId: maputil.GetString(options.ProviderDeployConfig, "certificateId"),
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeployProviderTypeBaotaPanelConsole, domain.DeployProviderTypeBaotaPanelSite:
		{
			access := domain.AccessConfigForBaotaPanel{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeployProviderTypeBaotaPanelConsole:
				deployer, err := pBaotaPanelConsole.NewDeployer(&pBaotaPanelConsole.DeployerConfig{
					ApiUrl:                   access.ApiUrl,
					ApiKey:                   access.ApiKey,
					AllowInsecureConnections: access.AllowInsecureConnections,
					AutoRestart:              maputil.GetBool(options.ProviderDeployConfig, "autoRestart"),
				})
				return deployer, err

			case domain.DeployProviderTypeBaotaPanelSite:
				deployer, err := pBaotaPanelSite.NewDeployer(&pBaotaPanelSite.DeployerConfig{
					ApiUrl:                   access.ApiUrl,
					ApiKey:                   access.ApiKey,
					AllowInsecureConnections: access.AllowInsecureConnections,
					SiteType:                 maputil.GetOrDefaultString(options.ProviderDeployConfig, "siteType", "other"),
					SiteName:                 maputil.GetString(options.ProviderDeployConfig, "siteName"),
					SiteNames:                sliceutil.Filter(strings.Split(maputil.GetString(options.ProviderDeployConfig, "siteNames"), ";"), func(s string) bool { return s != "" }),
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeployProviderTypeBunnyCDN:
		{
			access := domain.AccessConfigForBunny{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			deployer, err := pBunnyCDN.NewDeployer(&pBunnyCDN.DeployerConfig{
				ApiKey:     access.ApiKey,
				PullZoneId: maputil.GetString(options.ProviderDeployConfig, "pullZoneId"),
				HostName:   maputil.GetString(options.ProviderDeployConfig, "hostName"),
			})
			return deployer, err
		}

	case domain.DeployProviderTypeBytePlusCDN:
		{
			access := domain.AccessConfigForBytePlus{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeployProviderTypeBytePlusCDN:
				deployer, err := pBytePlusCDN.NewDeployer(&pBytePlusCDN.DeployerConfig{
					AccessKey: access.AccessKey,
					SecretKey: access.SecretKey,
					Domain:    maputil.GetString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeployProviderTypeCacheFly:
		{
			access := domain.AccessConfigForCacheFly{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
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
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			deployer, err := pCdnfly.NewDeployer(&pCdnfly.DeployerConfig{
				ApiUrl:        access.ApiUrl,
				ApiKey:        access.ApiKey,
				ApiSecret:     access.ApiSecret,
				ResourceType:  pCdnfly.ResourceType(maputil.GetOrDefaultString(options.ProviderDeployConfig, "resourceType", string(pCdnfly.RESOURCE_TYPE_SITE))),
				SiteId:        maputil.GetString(options.ProviderDeployConfig, "siteId"),
				CertificateId: maputil.GetString(options.ProviderDeployConfig, "certificateId"),
			})
			return deployer, err
		}

	case domain.DeployProviderTypeDogeCloudCDN:
		{
			access := domain.AccessConfigForDogeCloud{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			deployer, err := pDogeCDN.NewDeployer(&pDogeCDN.DeployerConfig{
				AccessKey: access.AccessKey,
				SecretKey: access.SecretKey,
				Domain:    maputil.GetString(options.ProviderDeployConfig, "domain"),
			})
			return deployer, err
		}

	case domain.DeployProviderTypeEdgioApplications:
		{
			access := domain.AccessConfigForEdgio{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			deployer, err := pEdgioApplications.NewDeployer(&pEdgioApplications.DeployerConfig{
				ClientId:      access.ClientId,
				ClientSecret:  access.ClientSecret,
				EnvironmentId: maputil.GetString(options.ProviderDeployConfig, "environmentId"),
			})
			return deployer, err
		}

	case domain.DeployProviderTypeGcoreCDN:
		{
			access := domain.AccessConfigForGcore{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeployProviderTypeGcoreCDN:
				deployer, err := pGcoreCDN.NewDeployer(&pGcoreCDN.DeployerConfig{
					ApiToken:   access.ApiToken,
					ResourceId: maputil.GetInt64(options.ProviderDeployConfig, "resourceId"),
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeployProviderTypeHuaweiCloudCDN, domain.DeployProviderTypeHuaweiCloudELB, domain.DeployProviderTypeHuaweiCloudSCM, domain.DeployProviderTypeHuaweiCloudWAF:
		{
			access := domain.AccessConfigForHuaweiCloud{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeployProviderTypeHuaweiCloudCDN:
				deployer, err := pHuaweiCloudCDN.NewDeployer(&pHuaweiCloudCDN.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					SecretAccessKey: access.SecretAccessKey,
					Region:          maputil.GetString(options.ProviderDeployConfig, "region"),
					Domain:          maputil.GetString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			case domain.DeployProviderTypeHuaweiCloudELB:
				deployer, err := pHuaweiCloudELB.NewDeployer(&pHuaweiCloudELB.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					SecretAccessKey: access.SecretAccessKey,
					Region:          maputil.GetString(options.ProviderDeployConfig, "region"),
					ResourceType:    pHuaweiCloudELB.ResourceType(maputil.GetString(options.ProviderDeployConfig, "resourceType")),
					CertificateId:   maputil.GetString(options.ProviderDeployConfig, "certificateId"),
					LoadbalancerId:  maputil.GetString(options.ProviderDeployConfig, "loadbalancerId"),
					ListenerId:      maputil.GetString(options.ProviderDeployConfig, "listenerId"),
				})
				return deployer, err

			case domain.DeployProviderTypeHuaweiCloudSCM:
				deployer, err := pHuaweiCloudSCM.NewDeployer(&pHuaweiCloudSCM.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					SecretAccessKey: access.SecretAccessKey,
				})
				return deployer, err

			case domain.DeployProviderTypeHuaweiCloudWAF:
				deployer, err := pHuaweiCloudWAF.NewDeployer(&pHuaweiCloudWAF.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					SecretAccessKey: access.SecretAccessKey,
					Region:          maputil.GetString(options.ProviderDeployConfig, "region"),
					ResourceType:    pHuaweiCloudWAF.ResourceType(maputil.GetString(options.ProviderDeployConfig, "resourceType")),
					CertificateId:   maputil.GetString(options.ProviderDeployConfig, "certificateId"),
					Domain:          maputil.GetString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeployProviderTypeJDCloudALB, domain.DeployProviderTypeJDCloudCDN, domain.DeployProviderTypeJDCloudLive, domain.DeployProviderTypeJDCloudVOD:
		{
			access := domain.AccessConfigForJDCloud{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeployProviderTypeJDCloudALB:
				deployer, err := pJDCloudALB.NewDeployer(&pJDCloudALB.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					RegionId:        maputil.GetString(options.ProviderDeployConfig, "regionId"),
					ResourceType:    pJDCloudALB.ResourceType(maputil.GetString(options.ProviderDeployConfig, "resourceType")),
					LoadbalancerId:  maputil.GetString(options.ProviderDeployConfig, "loadbalancerId"),
					ListenerId:      maputil.GetString(options.ProviderDeployConfig, "listenerId"),
				})
				return deployer, err

			case domain.DeployProviderTypeJDCloudCDN:
				deployer, err := pJDCloudCDN.NewDeployer(&pJDCloudCDN.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Domain:          maputil.GetString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			case domain.DeployProviderTypeJDCloudLive:
				deployer, err := pJDCloudLive.NewDeployer(&pJDCloudLive.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Domain:          maputil.GetString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			case domain.DeployProviderTypeJDCloudVOD:
				deployer, err := pJDCloudVOD.NewDeployer(&pJDCloudVOD.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Domain:          maputil.GetString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeployProviderTypeLocal:
		{
			deployer, err := pLocal.NewDeployer(&pLocal.DeployerConfig{
				ShellEnv:       pLocal.ShellEnvType(maputil.GetString(options.ProviderDeployConfig, "shellEnv")),
				PreCommand:     maputil.GetString(options.ProviderDeployConfig, "preCommand"),
				PostCommand:    maputil.GetString(options.ProviderDeployConfig, "postCommand"),
				OutputFormat:   pLocal.OutputFormatType(maputil.GetOrDefaultString(options.ProviderDeployConfig, "format", string(pLocal.OUTPUT_FORMAT_PEM))),
				OutputCertPath: maputil.GetString(options.ProviderDeployConfig, "certPath"),
				OutputKeyPath:  maputil.GetString(options.ProviderDeployConfig, "keyPath"),
				PfxPassword:    maputil.GetString(options.ProviderDeployConfig, "pfxPassword"),
				JksAlias:       maputil.GetString(options.ProviderDeployConfig, "jksAlias"),
				JksKeypass:     maputil.GetString(options.ProviderDeployConfig, "jksKeypass"),
				JksStorepass:   maputil.GetString(options.ProviderDeployConfig, "jksStorepass"),
			})
			return deployer, err
		}

	case domain.DeployProviderTypeKubernetesSecret:
		{
			access := domain.AccessConfigForKubernetes{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			deployer, err := pK8sSecret.NewDeployer(&pK8sSecret.DeployerConfig{
				KubeConfig:          access.KubeConfig,
				Namespace:           maputil.GetOrDefaultString(options.ProviderDeployConfig, "namespace", "default"),
				SecretName:          maputil.GetString(options.ProviderDeployConfig, "secretName"),
				SecretType:          maputil.GetOrDefaultString(options.ProviderDeployConfig, "secretType", "kubernetes.io/tls"),
				SecretDataKeyForCrt: maputil.GetOrDefaultString(options.ProviderDeployConfig, "secretDataKeyForCrt", "tls.crt"),
				SecretDataKeyForKey: maputil.GetOrDefaultString(options.ProviderDeployConfig, "secretDataKeyForKey", "tls.key"),
			})
			return deployer, err
		}

	case domain.DeployProviderTypeQiniuCDN, domain.DeployProviderTypeQiniuKodo, domain.DeployProviderTypeQiniuPili:
		{
			access := domain.AccessConfigForQiniu{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeployProviderTypeQiniuCDN, domain.DeployProviderTypeQiniuKodo:
				deployer, err := pQiniuCDN.NewDeployer(&pQiniuCDN.DeployerConfig{
					AccessKey: access.AccessKey,
					SecretKey: access.SecretKey,
					Domain:    maputil.GetString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			case domain.DeployProviderTypeQiniuPili:
				deployer, err := pQiniuPili.NewDeployer(&pQiniuPili.DeployerConfig{
					AccessKey: access.AccessKey,
					SecretKey: access.SecretKey,
					Hub:       maputil.GetString(options.ProviderDeployConfig, "hub"),
					Domain:    maputil.GetString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeployProviderTypeRainYunRCDN:
		{
			access := domain.AccessConfigForRainYun{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeployProviderTypeTencentCloudCDN:
				deployer, err := pRainYunRCDN.NewDeployer(&pRainYunRCDN.DeployerConfig{
					ApiKey:     access.ApiKey,
					InstanceId: maputil.GetInt32(options.ProviderDeployConfig, "instanceId"),
					Domain:     maputil.GetString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeployProviderTypeSafeLine:
		{
			access := domain.AccessConfigForSafeLine{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			deployer, err := pSafeLine.NewDeployer(&pSafeLine.DeployerConfig{
				ApiUrl:                   access.ApiUrl,
				ApiToken:                 access.ApiToken,
				AllowInsecureConnections: access.AllowInsecureConnections,
				ResourceType:             pSafeLine.ResourceType(maputil.GetString(options.ProviderDeployConfig, "resourceType")),
				CertificateId:            maputil.GetInt32(options.ProviderDeployConfig, "certificateId"),
			})
			return deployer, err
		}

	case domain.DeployProviderTypeSSH:
		{
			access := domain.AccessConfigForSSH{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			deployer, err := pSSH.NewDeployer(&pSSH.DeployerConfig{
				SshHost:          access.Host,
				SshPort:          access.Port,
				SshUsername:      access.Username,
				SshPassword:      access.Password,
				SshKey:           access.Key,
				SshKeyPassphrase: access.KeyPassphrase,
				UseSCP:           maputil.GetBool(options.ProviderDeployConfig, "useSCP"),
				PreCommand:       maputil.GetString(options.ProviderDeployConfig, "preCommand"),
				PostCommand:      maputil.GetString(options.ProviderDeployConfig, "postCommand"),
				OutputFormat:     pSSH.OutputFormatType(maputil.GetOrDefaultString(options.ProviderDeployConfig, "format", string(pSSH.OUTPUT_FORMAT_PEM))),
				OutputCertPath:   maputil.GetString(options.ProviderDeployConfig, "certPath"),
				OutputKeyPath:    maputil.GetString(options.ProviderDeployConfig, "keyPath"),
				PfxPassword:      maputil.GetString(options.ProviderDeployConfig, "pfxPassword"),
				JksAlias:         maputil.GetString(options.ProviderDeployConfig, "jksAlias"),
				JksKeypass:       maputil.GetString(options.ProviderDeployConfig, "jksKeypass"),
				JksStorepass:     maputil.GetString(options.ProviderDeployConfig, "jksStorepass"),
			})
			return deployer, err
		}

	case domain.DeployProviderTypeTencentCloudCDN, domain.DeployProviderTypeTencentCloudCLB, domain.DeployProviderTypeTencentCloudCOS, domain.DeployProviderTypeTencentCloudCSS, domain.DeployProviderTypeTencentCloudECDN, domain.DeployProviderTypeTencentCloudEO, domain.DeployProviderTypeTencentCloudSCF, domain.DeployProviderTypeTencentCloudSSL, domain.DeployProviderTypeTencentCloudSSLDeploy, domain.DeployProviderTypeTencentCloudVOD, domain.DeployProviderTypeTencentCloudWAF:
		{
			access := domain.AccessConfigForTencentCloud{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeployProviderTypeTencentCloudCDN:
				deployer, err := pTencentCloudCDN.NewDeployer(&pTencentCloudCDN.DeployerConfig{
					SecretId:  access.SecretId,
					SecretKey: access.SecretKey,
					Domain:    maputil.GetString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			case domain.DeployProviderTypeTencentCloudCLB:
				deployer, err := pTencentCloudCLB.NewDeployer(&pTencentCloudCLB.DeployerConfig{
					SecretId:       access.SecretId,
					SecretKey:      access.SecretKey,
					Region:         maputil.GetString(options.ProviderDeployConfig, "region"),
					ResourceType:   pTencentCloudCLB.ResourceType(maputil.GetString(options.ProviderDeployConfig, "resourceType")),
					LoadbalancerId: maputil.GetString(options.ProviderDeployConfig, "loadbalancerId"),
					ListenerId:     maputil.GetString(options.ProviderDeployConfig, "listenerId"),
					Domain:         maputil.GetString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			case domain.DeployProviderTypeTencentCloudCOS:
				deployer, err := pTencentCloudCOS.NewDeployer(&pTencentCloudCOS.DeployerConfig{
					SecretId:  access.SecretId,
					SecretKey: access.SecretKey,
					Region:    maputil.GetString(options.ProviderDeployConfig, "region"),
					Bucket:    maputil.GetString(options.ProviderDeployConfig, "bucket"),
					Domain:    maputil.GetString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			case domain.DeployProviderTypeTencentCloudCSS:
				deployer, err := pTencentCloudCSS.NewDeployer(&pTencentCloudCSS.DeployerConfig{
					SecretId:  access.SecretId,
					SecretKey: access.SecretKey,
					Domain:    maputil.GetString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			case domain.DeployProviderTypeTencentCloudECDN:
				deployer, err := pTencentCloudECDN.NewDeployer(&pTencentCloudECDN.DeployerConfig{
					SecretId:  access.SecretId,
					SecretKey: access.SecretKey,
					Domain:    maputil.GetString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			case domain.DeployProviderTypeTencentCloudEO:
				deployer, err := pTencentCloudEO.NewDeployer(&pTencentCloudEO.DeployerConfig{
					SecretId:  access.SecretId,
					SecretKey: access.SecretKey,
					ZoneId:    maputil.GetString(options.ProviderDeployConfig, "zoneId"),
					Domain:    maputil.GetString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			case domain.DeployProviderTypeTencentCloudSCF:
				deployer, err := pTencentCloudSCF.NewDeployer(&pTencentCloudSCF.DeployerConfig{
					SecretId:  access.SecretId,
					SecretKey: access.SecretKey,
					Region:    maputil.GetString(options.ProviderDeployConfig, "region"),
					Domain:    maputil.GetString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			case domain.DeployProviderTypeTencentCloudSSL:
				deployer, err := pTencentCloudSSL.NewDeployer(&pTencentCloudSSL.DeployerConfig{
					SecretId:  access.SecretId,
					SecretKey: access.SecretKey,
				})
				return deployer, err

			case domain.DeployProviderTypeTencentCloudSSLDeploy:
				deployer, err := pTencentCloudSSLDeploy.NewDeployer(&pTencentCloudSSLDeploy.DeployerConfig{
					SecretId:     access.SecretId,
					SecretKey:    access.SecretKey,
					Region:       maputil.GetString(options.ProviderDeployConfig, "region"),
					ResourceType: maputil.GetString(options.ProviderDeployConfig, "resourceType"),
					ResourceIds:  sliceutil.Filter(strings.Split(maputil.GetString(options.ProviderDeployConfig, "resourceIds"), ";"), func(s string) bool { return s != "" }),
				})
				return deployer, err

			case domain.DeployProviderTypeTencentCloudVOD:
				deployer, err := pTencentCloudVOD.NewDeployer(&pTencentCloudVOD.DeployerConfig{
					SecretId:  access.SecretId,
					SecretKey: access.SecretKey,
					SubAppId:  maputil.GetInt64(options.ProviderDeployConfig, "subAppId"),
					Domain:    maputil.GetString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			case domain.DeployProviderTypeTencentCloudWAF:
				deployer, err := pTencentCloudWAF.NewDeployer(&pTencentCloudWAF.DeployerConfig{
					SecretId:   access.SecretId,
					SecretKey:  access.SecretKey,
					Domain:     maputil.GetString(options.ProviderDeployConfig, "domain"),
					DomainId:   maputil.GetString(options.ProviderDeployConfig, "domainId"),
					InstanceId: maputil.GetString(options.ProviderDeployConfig, "instanceId"),
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeployProviderTypeUCloudUCDN, domain.DeployProviderTypeUCloudUS3:
		{
			access := domain.AccessConfigForUCloud{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeployProviderTypeUCloudUCDN:
				deployer, err := pUCloudUCDN.NewDeployer(&pUCloudUCDN.DeployerConfig{
					PrivateKey: access.PrivateKey,
					PublicKey:  access.PublicKey,
					ProjectId:  access.ProjectId,
					DomainId:   maputil.GetString(options.ProviderDeployConfig, "domainId"),
				})
				return deployer, err

			case domain.DeployProviderTypeUCloudUS3:
				deployer, err := pUCloudUS3.NewDeployer(&pUCloudUS3.DeployerConfig{
					PrivateKey: access.PrivateKey,
					PublicKey:  access.PublicKey,
					ProjectId:  access.ProjectId,
					Region:     maputil.GetString(options.ProviderDeployConfig, "region"),
					Bucket:     maputil.GetString(options.ProviderDeployConfig, "bucket"),
					Domain:     maputil.GetString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeployProviderTypeUpyunCDN, domain.DeployProviderTypeUpyunFile:
		{
			access := domain.AccessConfigForUpyun{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeployProviderTypeUpyunCDN, domain.DeployProviderTypeUpyunFile:
				deployer, err := pUpyunCDN.NewDeployer(&pUpyunCDN.DeployerConfig{
					Username: access.Username,
					Password: access.Password,
					Domain:   maputil.GetString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeployProviderTypeVolcEngineALB, domain.DeployProviderTypeVolcEngineCDN, domain.DeployProviderTypeVolcEngineCertCenter, domain.DeployProviderTypeVolcEngineCLB, domain.DeployProviderTypeVolcEngineDCDN, domain.DeployProviderTypeVolcEngineImageX, domain.DeployProviderTypeVolcEngineLive, domain.DeployProviderTypeVolcEngineTOS:
		{
			access := domain.AccessConfigForVolcEngine{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeployProviderTypeVolcEngineALB:
				deployer, err := pVolcEngineALB.NewDeployer(&pVolcEngineALB.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.SecretAccessKey,
					Region:          maputil.GetString(options.ProviderDeployConfig, "region"),
					ResourceType:    pVolcEngineALB.ResourceType(maputil.GetString(options.ProviderDeployConfig, "resourceType")),
					LoadbalancerId:  maputil.GetString(options.ProviderDeployConfig, "loadbalancerId"),
					ListenerId:      maputil.GetString(options.ProviderDeployConfig, "listenerId"),
					Domain:          maputil.GetString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			case domain.DeployProviderTypeVolcEngineCDN:
				deployer, err := pVolcEngineCDN.NewDeployer(&pVolcEngineCDN.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.SecretAccessKey,
					Domain:          maputil.GetString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			case domain.DeployProviderTypeVolcEngineCertCenter:
				deployer, err := pVolcEngineCertCenter.NewDeployer(&pVolcEngineCertCenter.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.SecretAccessKey,
					Region:          maputil.GetString(options.ProviderDeployConfig, "region"),
				})
				return deployer, err

			case domain.DeployProviderTypeVolcEngineCLB:
				deployer, err := pVolcEngineCLB.NewDeployer(&pVolcEngineCLB.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.SecretAccessKey,
					Region:          maputil.GetString(options.ProviderDeployConfig, "region"),
					ResourceType:    pVolcEngineCLB.ResourceType(maputil.GetString(options.ProviderDeployConfig, "resourceType")),
					LoadbalancerId:  maputil.GetString(options.ProviderDeployConfig, "loadbalancerId"),
					ListenerId:      maputil.GetString(options.ProviderDeployConfig, "listenerId"),
				})
				return deployer, err

			case domain.DeployProviderTypeVolcEngineDCDN:
				deployer, err := pVolcEngineDCDN.NewDeployer(&pVolcEngineDCDN.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.SecretAccessKey,
					Domain:          maputil.GetString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			case domain.DeployProviderTypeVolcEngineImageX:
				deployer, err := pVolcEngineImageX.NewDeployer(&pVolcEngineImageX.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.SecretAccessKey,
					Region:          maputil.GetString(options.ProviderDeployConfig, "region"),
					ServiceId:       maputil.GetString(options.ProviderDeployConfig, "serviceId"),
					Domain:          maputil.GetString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			case domain.DeployProviderTypeVolcEngineLive:
				deployer, err := pVolcEngineLive.NewDeployer(&pVolcEngineLive.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.SecretAccessKey,
					Domain:          maputil.GetString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			case domain.DeployProviderTypeVolcEngineTOS:
				deployer, err := pVolcEngineTOS.NewDeployer(&pVolcEngineTOS.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.SecretAccessKey,
					Region:          maputil.GetString(options.ProviderDeployConfig, "region"),
					Bucket:          maputil.GetString(options.ProviderDeployConfig, "bucket"),
					Domain:          maputil.GetString(options.ProviderDeployConfig, "domain"),
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeployProviderTypeWangsuCDNPro:
		{
			access := domain.AccessConfigForWangsu{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeployProviderTypeWangsuCDNPro:
				deployer, err := pWangsuCDNPro.NewDeployer(&pWangsuCDNPro.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					ApiKey:          access.ApiKey,
					Environment:     maputil.GetOrDefaultString(options.ProviderDeployConfig, "environment", "production"),
					Domain:          maputil.GetString(options.ProviderDeployConfig, "domain"),
					CertificateId:   maputil.GetString(options.ProviderDeployConfig, "certificateId"),
					WebhookId:       maputil.GetString(options.ProviderDeployConfig, "webhookId"),
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeployProviderTypeWebhook:
		{
			access := domain.AccessConfigForWebhook{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			deployer, err := pWebhook.NewDeployer(&pWebhook.DeployerConfig{
				WebhookUrl:               access.Url,
				WebhookData:              maputil.GetString(options.ProviderDeployConfig, "webhookData"),
				AllowInsecureConnections: access.AllowInsecureConnections,
			})
			return deployer, err
		}
	}

	return nil, fmt.Errorf("unsupported deployer provider: %s", string(options.Provider))
}
