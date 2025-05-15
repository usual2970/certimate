package deployer

import (
	"fmt"
	"net/http"
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
	pAliyunDDoS "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/aliyun-ddos"
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
	pGoEdge "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/goedge"
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
	pNetlifySite "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/netlify-site"
	pProxmoxVE "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/proxmoxve"
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
	httputil "github.com/usual2970/certimate/internal/pkg/utils/http"
	maputil "github.com/usual2970/certimate/internal/pkg/utils/map"
	sliceutil "github.com/usual2970/certimate/internal/pkg/utils/slice"
)

type deployerProviderOptions struct {
	Provider               domain.DeploymentProviderType
	ProviderAccessConfig   map[string]any
	ProviderExtendedConfig map[string]any
}

func createDeployerProvider(options *deployerProviderOptions) (deployer.Deployer, error) {
	/*
	  注意：如果追加新的常量值，请保持以 ASCII 排序。
	  NOTICE: If you add new constant, please keep ASCII order.
	*/
	switch options.Provider {
	case domain.DeploymentProviderType1PanelConsole, domain.DeploymentProviderType1PanelSite:
		{
			access := domain.AccessConfigFor1Panel{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeploymentProviderType1PanelConsole:
				deployer, err := p1PanelConsole.NewDeployer(&p1PanelConsole.DeployerConfig{
					ApiUrl:                   access.ApiUrl,
					ApiKey:                   access.ApiKey,
					AllowInsecureConnections: access.AllowInsecureConnections,
					AutoRestart:              maputil.GetBool(options.ProviderExtendedConfig, "autoRestart"),
				})
				return deployer, err

			case domain.DeploymentProviderType1PanelSite:
				deployer, err := p1PanelSite.NewDeployer(&p1PanelSite.DeployerConfig{
					ApiUrl:                   access.ApiUrl,
					ApiKey:                   access.ApiKey,
					AllowInsecureConnections: access.AllowInsecureConnections,
					ResourceType:             p1PanelSite.ResourceType(maputil.GetOrDefaultString(options.ProviderExtendedConfig, "resourceType", string(p1PanelSite.RESOURCE_TYPE_WEBSITE))),
					WebsiteId:                maputil.GetInt64(options.ProviderExtendedConfig, "websiteId"),
					CertificateId:            maputil.GetInt64(options.ProviderExtendedConfig, "certificateId"),
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeploymentProviderTypeAliyunALB, domain.DeploymentProviderTypeAliyunAPIGW, domain.DeploymentProviderTypeAliyunCAS, domain.DeploymentProviderTypeAliyunCASDeploy, domain.DeploymentProviderTypeAliyunCDN, domain.DeploymentProviderTypeAliyunCLB, domain.DeploymentProviderTypeAliyunDCDN, domain.DeploymentProviderTypeAliyunDDoS, domain.DeploymentProviderTypeAliyunESA, domain.DeploymentProviderTypeAliyunFC, domain.DeploymentProviderTypeAliyunLive, domain.DeploymentProviderTypeAliyunNLB, domain.DeploymentProviderTypeAliyunOSS, domain.DeploymentProviderTypeAliyunVOD, domain.DeploymentProviderTypeAliyunWAF:
		{
			access := domain.AccessConfigForAliyun{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeploymentProviderTypeAliyunALB:
				deployer, err := pAliyunALB.NewDeployer(&pAliyunALB.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Region:          maputil.GetString(options.ProviderExtendedConfig, "region"),
					ResourceType:    pAliyunALB.ResourceType(maputil.GetString(options.ProviderExtendedConfig, "resourceType")),
					LoadbalancerId:  maputil.GetString(options.ProviderExtendedConfig, "loadbalancerId"),
					ListenerId:      maputil.GetString(options.ProviderExtendedConfig, "listenerId"),
					Domain:          maputil.GetString(options.ProviderExtendedConfig, "domain"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeAliyunAPIGW:
				deployer, err := pAliyunAPIGW.NewDeployer(&pAliyunAPIGW.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Region:          maputil.GetString(options.ProviderExtendedConfig, "region"),
					ServiceType:     pAliyunAPIGW.ServiceType(maputil.GetString(options.ProviderExtendedConfig, "serviceType")),
					GatewayId:       maputil.GetString(options.ProviderExtendedConfig, "gatewayId"),
					GroupId:         maputil.GetString(options.ProviderExtendedConfig, "groupId"),
					Domain:          maputil.GetString(options.ProviderExtendedConfig, "domain"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeAliyunCAS:
				deployer, err := pAliyunCAS.NewDeployer(&pAliyunCAS.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Region:          maputil.GetString(options.ProviderExtendedConfig, "region"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeAliyunCASDeploy:
				deployer, err := pAliyunCASDeploy.NewDeployer(&pAliyunCASDeploy.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Region:          maputil.GetString(options.ProviderExtendedConfig, "region"),
					ResourceIds:     sliceutil.Filter(strings.Split(maputil.GetString(options.ProviderExtendedConfig, "resourceIds"), ";"), func(s string) bool { return s != "" }),
					ContactIds:      sliceutil.Filter(strings.Split(maputil.GetString(options.ProviderExtendedConfig, "contactIds"), ";"), func(s string) bool { return s != "" }),
				})
				return deployer, err

			case domain.DeploymentProviderTypeAliyunCDN:
				deployer, err := pAliyunCDN.NewDeployer(&pAliyunCDN.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Domain:          maputil.GetString(options.ProviderExtendedConfig, "domain"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeAliyunCLB:
				deployer, err := pAliyunCLB.NewDeployer(&pAliyunCLB.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Region:          maputil.GetString(options.ProviderExtendedConfig, "region"),
					ResourceType:    pAliyunCLB.ResourceType(maputil.GetString(options.ProviderExtendedConfig, "resourceType")),
					LoadbalancerId:  maputil.GetString(options.ProviderExtendedConfig, "loadbalancerId"),
					ListenerPort:    maputil.GetOrDefaultInt32(options.ProviderExtendedConfig, "listenerPort", 443),
					Domain:          maputil.GetString(options.ProviderExtendedConfig, "domain"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeAliyunDCDN:
				deployer, err := pAliyunDCDN.NewDeployer(&pAliyunDCDN.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Domain:          maputil.GetString(options.ProviderExtendedConfig, "domain"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeAliyunDDoS:
				deployer, err := pAliyunDDoS.NewDeployer(&pAliyunDDoS.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Region:          maputil.GetString(options.ProviderExtendedConfig, "region"),
					Domain:          maputil.GetString(options.ProviderExtendedConfig, "domain"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeAliyunESA:
				deployer, err := pAliyunESA.NewDeployer(&pAliyunESA.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Region:          maputil.GetString(options.ProviderExtendedConfig, "region"),
					SiteId:          maputil.GetInt64(options.ProviderExtendedConfig, "siteId"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeAliyunFC:
				deployer, err := pAliyunFC.NewDeployer(&pAliyunFC.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Region:          maputil.GetString(options.ProviderExtendedConfig, "region"),
					ServiceVersion:  maputil.GetOrDefaultString(options.ProviderExtendedConfig, "serviceVersion", "3.0"),
					Domain:          maputil.GetString(options.ProviderExtendedConfig, "domain"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeAliyunLive:
				deployer, err := pAliyunLive.NewDeployer(&pAliyunLive.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Region:          maputil.GetString(options.ProviderExtendedConfig, "region"),
					Domain:          maputil.GetString(options.ProviderExtendedConfig, "domain"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeAliyunNLB:
				deployer, err := pAliyunNLB.NewDeployer(&pAliyunNLB.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Region:          maputil.GetString(options.ProviderExtendedConfig, "region"),
					ResourceType:    pAliyunNLB.ResourceType(maputil.GetString(options.ProviderExtendedConfig, "resourceType")),
					LoadbalancerId:  maputil.GetString(options.ProviderExtendedConfig, "loadbalancerId"),
					ListenerId:      maputil.GetString(options.ProviderExtendedConfig, "listenerId"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeAliyunOSS:
				deployer, err := pAliyunOSS.NewDeployer(&pAliyunOSS.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Region:          maputil.GetString(options.ProviderExtendedConfig, "region"),
					Bucket:          maputil.GetString(options.ProviderExtendedConfig, "bucket"),
					Domain:          maputil.GetString(options.ProviderExtendedConfig, "domain"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeAliyunVOD:
				deployer, err := pAliyunVOD.NewDeployer(&pAliyunVOD.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Region:          maputil.GetString(options.ProviderExtendedConfig, "region"),
					Domain:          maputil.GetString(options.ProviderExtendedConfig, "domain"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeAliyunWAF:
				deployer, err := pAliyunWAF.NewDeployer(&pAliyunWAF.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Region:          maputil.GetString(options.ProviderExtendedConfig, "region"),
					ServiceVersion:  maputil.GetOrDefaultString(options.ProviderExtendedConfig, "serviceVersion", "3.0"),
					InstanceId:      maputil.GetString(options.ProviderExtendedConfig, "instanceId"),
					Domain:          maputil.GetString(options.ProviderExtendedConfig, "domain"),
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeploymentProviderTypeAWSACM, domain.DeploymentProviderTypeAWSCloudFront:
		{
			access := domain.AccessConfigForAWS{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeploymentProviderTypeAWSACM:
				deployer, err := pAWSACM.NewDeployer(&pAWSACM.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					SecretAccessKey: access.SecretAccessKey,
					Region:          maputil.GetString(options.ProviderExtendedConfig, "region"),
					CertificateArn:  maputil.GetString(options.ProviderExtendedConfig, "certificateArn"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeAWSCloudFront:
				deployer, err := pAWSCloudFront.NewDeployer(&pAWSCloudFront.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					SecretAccessKey: access.SecretAccessKey,
					Region:          maputil.GetString(options.ProviderExtendedConfig, "region"),
					DistributionId:  maputil.GetString(options.ProviderExtendedConfig, "distributionId"),
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeploymentProviderTypeAzureKeyVault:
		{
			access := domain.AccessConfigForAzure{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeploymentProviderTypeAzureKeyVault:
				deployer, err := pAzureKeyVault.NewDeployer(&pAzureKeyVault.DeployerConfig{
					TenantId:        access.TenantId,
					ClientId:        access.ClientId,
					ClientSecret:    access.ClientSecret,
					CloudName:       access.CloudName,
					KeyVaultName:    maputil.GetString(options.ProviderExtendedConfig, "keyvaultName"),
					CertificateName: maputil.GetString(options.ProviderExtendedConfig, "certificateName"),
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeploymentProviderTypeBaiduCloudAppBLB, domain.DeploymentProviderTypeBaiduCloudBLB, domain.DeploymentProviderTypeBaiduCloudCDN, domain.DeploymentProviderTypeBaiduCloudCert:
		{
			access := domain.AccessConfigForBaiduCloud{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeploymentProviderTypeBaiduCloudAppBLB:
				deployer, err := pBaiduCloudAppBLB.NewDeployer(&pBaiduCloudAppBLB.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					SecretAccessKey: access.SecretAccessKey,
					Region:          maputil.GetString(options.ProviderExtendedConfig, "region"),
					ResourceType:    pBaiduCloudAppBLB.ResourceType(maputil.GetString(options.ProviderExtendedConfig, "resourceType")),
					LoadbalancerId:  maputil.GetString(options.ProviderExtendedConfig, "loadbalancerId"),
					ListenerPort:    maputil.GetInt32(options.ProviderExtendedConfig, "listenerPort"),
					Domain:          maputil.GetString(options.ProviderExtendedConfig, "domain"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeBaiduCloudBLB:
				deployer, err := pBaiduCloudBLB.NewDeployer(&pBaiduCloudBLB.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					SecretAccessKey: access.SecretAccessKey,
					Region:          maputil.GetString(options.ProviderExtendedConfig, "region"),
					ResourceType:    pBaiduCloudBLB.ResourceType(maputil.GetString(options.ProviderExtendedConfig, "resourceType")),
					LoadbalancerId:  maputil.GetString(options.ProviderExtendedConfig, "loadbalancerId"),
					ListenerPort:    maputil.GetInt32(options.ProviderExtendedConfig, "listenerPort"),
					Domain:          maputil.GetString(options.ProviderExtendedConfig, "domain"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeBaiduCloudCDN:
				deployer, err := pBaiduCloudCDN.NewDeployer(&pBaiduCloudCDN.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					SecretAccessKey: access.SecretAccessKey,
					Domain:          maputil.GetString(options.ProviderExtendedConfig, "domain"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeBaiduCloudCert:
				deployer, err := pBaiduCloudCert.NewDeployer(&pBaiduCloudCert.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					SecretAccessKey: access.SecretAccessKey,
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeploymentProviderTypeBaishanCDN:
		{
			access := domain.AccessConfigForBaishan{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeploymentProviderTypeBaishanCDN:
				deployer, err := pBaishanCDN.NewDeployer(&pBaishanCDN.DeployerConfig{
					ApiToken:      access.ApiToken,
					Domain:        maputil.GetString(options.ProviderExtendedConfig, "domain"),
					CertificateId: maputil.GetString(options.ProviderExtendedConfig, "certificateId"),
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeploymentProviderTypeBaotaPanelConsole, domain.DeploymentProviderTypeBaotaPanelSite:
		{
			access := domain.AccessConfigForBaotaPanel{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeploymentProviderTypeBaotaPanelConsole:
				deployer, err := pBaotaPanelConsole.NewDeployer(&pBaotaPanelConsole.DeployerConfig{
					ApiUrl:                   access.ApiUrl,
					ApiKey:                   access.ApiKey,
					AllowInsecureConnections: access.AllowInsecureConnections,
					AutoRestart:              maputil.GetBool(options.ProviderExtendedConfig, "autoRestart"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeBaotaPanelSite:
				deployer, err := pBaotaPanelSite.NewDeployer(&pBaotaPanelSite.DeployerConfig{
					ApiUrl:                   access.ApiUrl,
					ApiKey:                   access.ApiKey,
					AllowInsecureConnections: access.AllowInsecureConnections,
					SiteType:                 maputil.GetOrDefaultString(options.ProviderExtendedConfig, "siteType", "other"),
					SiteName:                 maputil.GetString(options.ProviderExtendedConfig, "siteName"),
					SiteNames:                sliceutil.Filter(strings.Split(maputil.GetString(options.ProviderExtendedConfig, "siteNames"), ";"), func(s string) bool { return s != "" }),
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeploymentProviderTypeBunnyCDN:
		{
			access := domain.AccessConfigForBunny{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			deployer, err := pBunnyCDN.NewDeployer(&pBunnyCDN.DeployerConfig{
				ApiKey:     access.ApiKey,
				PullZoneId: maputil.GetString(options.ProviderExtendedConfig, "pullZoneId"),
				Hostname:   maputil.GetString(options.ProviderExtendedConfig, "hostname"),
			})
			return deployer, err
		}

	case domain.DeploymentProviderTypeBytePlusCDN:
		{
			access := domain.AccessConfigForBytePlus{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeploymentProviderTypeBytePlusCDN:
				deployer, err := pBytePlusCDN.NewDeployer(&pBytePlusCDN.DeployerConfig{
					AccessKey: access.AccessKey,
					SecretKey: access.SecretKey,
					Domain:    maputil.GetString(options.ProviderExtendedConfig, "domain"),
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeploymentProviderTypeCacheFly:
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

	case domain.DeploymentProviderTypeCdnfly:
		{
			access := domain.AccessConfigForCdnfly{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			deployer, err := pCdnfly.NewDeployer(&pCdnfly.DeployerConfig{
				ApiUrl:                   access.ApiUrl,
				ApiKey:                   access.ApiKey,
				ApiSecret:                access.ApiSecret,
				AllowInsecureConnections: access.AllowInsecureConnections,
				ResourceType:             pCdnfly.ResourceType(maputil.GetOrDefaultString(options.ProviderExtendedConfig, "resourceType", string(pCdnfly.RESOURCE_TYPE_SITE))),
				SiteId:                   maputil.GetString(options.ProviderExtendedConfig, "siteId"),
				CertificateId:            maputil.GetString(options.ProviderExtendedConfig, "certificateId"),
			})
			return deployer, err
		}

	case domain.DeploymentProviderTypeDogeCloudCDN:
		{
			access := domain.AccessConfigForDogeCloud{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			deployer, err := pDogeCDN.NewDeployer(&pDogeCDN.DeployerConfig{
				AccessKey: access.AccessKey,
				SecretKey: access.SecretKey,
				Domain:    maputil.GetString(options.ProviderExtendedConfig, "domain"),
			})
			return deployer, err
		}

	case domain.DeploymentProviderTypeEdgioApplications:
		{
			access := domain.AccessConfigForEdgio{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			deployer, err := pEdgioApplications.NewDeployer(&pEdgioApplications.DeployerConfig{
				ClientId:      access.ClientId,
				ClientSecret:  access.ClientSecret,
				EnvironmentId: maputil.GetString(options.ProviderExtendedConfig, "environmentId"),
			})
			return deployer, err
		}

	case domain.DeploymentProviderTypeGcoreCDN:
		{
			access := domain.AccessConfigForGcore{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeploymentProviderTypeGcoreCDN:
				deployer, err := pGcoreCDN.NewDeployer(&pGcoreCDN.DeployerConfig{
					ApiToken:      access.ApiToken,
					ResourceId:    maputil.GetInt64(options.ProviderExtendedConfig, "resourceId"),
					CertificateId: maputil.GetInt64(options.ProviderExtendedConfig, "certificateId"),
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeploymentProviderTypeGoEdge:
		{
			access := domain.AccessConfigForGoEdge{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			deployer, err := pGoEdge.NewDeployer(&pGoEdge.DeployerConfig{
				ApiUrl:                   access.ApiUrl,
				ApiRole:                  access.ApiRole,
				AccessKeyId:              access.AccessKeyId,
				AccessKey:                access.AccessKey,
				AllowInsecureConnections: access.AllowInsecureConnections,
				ResourceType:             pGoEdge.ResourceType(maputil.GetString(options.ProviderExtendedConfig, "resourceType")),
				CertificateId:            maputil.GetInt64(options.ProviderExtendedConfig, "certificateId"),
			})
			return deployer, err
		}

	case domain.DeploymentProviderTypeHuaweiCloudCDN, domain.DeploymentProviderTypeHuaweiCloudELB, domain.DeploymentProviderTypeHuaweiCloudSCM, domain.DeploymentProviderTypeHuaweiCloudWAF:
		{
			access := domain.AccessConfigForHuaweiCloud{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeploymentProviderTypeHuaweiCloudCDN:
				deployer, err := pHuaweiCloudCDN.NewDeployer(&pHuaweiCloudCDN.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					SecretAccessKey: access.SecretAccessKey,
					Region:          maputil.GetString(options.ProviderExtendedConfig, "region"),
					Domain:          maputil.GetString(options.ProviderExtendedConfig, "domain"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeHuaweiCloudELB:
				deployer, err := pHuaweiCloudELB.NewDeployer(&pHuaweiCloudELB.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					SecretAccessKey: access.SecretAccessKey,
					Region:          maputil.GetString(options.ProviderExtendedConfig, "region"),
					ResourceType:    pHuaweiCloudELB.ResourceType(maputil.GetString(options.ProviderExtendedConfig, "resourceType")),
					CertificateId:   maputil.GetString(options.ProviderExtendedConfig, "certificateId"),
					LoadbalancerId:  maputil.GetString(options.ProviderExtendedConfig, "loadbalancerId"),
					ListenerId:      maputil.GetString(options.ProviderExtendedConfig, "listenerId"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeHuaweiCloudSCM:
				deployer, err := pHuaweiCloudSCM.NewDeployer(&pHuaweiCloudSCM.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					SecretAccessKey: access.SecretAccessKey,
				})
				return deployer, err

			case domain.DeploymentProviderTypeHuaweiCloudWAF:
				deployer, err := pHuaweiCloudWAF.NewDeployer(&pHuaweiCloudWAF.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					SecretAccessKey: access.SecretAccessKey,
					Region:          maputil.GetString(options.ProviderExtendedConfig, "region"),
					ResourceType:    pHuaweiCloudWAF.ResourceType(maputil.GetString(options.ProviderExtendedConfig, "resourceType")),
					CertificateId:   maputil.GetString(options.ProviderExtendedConfig, "certificateId"),
					Domain:          maputil.GetString(options.ProviderExtendedConfig, "domain"),
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeploymentProviderTypeJDCloudALB, domain.DeploymentProviderTypeJDCloudCDN, domain.DeploymentProviderTypeJDCloudLive, domain.DeploymentProviderTypeJDCloudVOD:
		{
			access := domain.AccessConfigForJDCloud{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeploymentProviderTypeJDCloudALB:
				deployer, err := pJDCloudALB.NewDeployer(&pJDCloudALB.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					RegionId:        maputil.GetString(options.ProviderExtendedConfig, "regionId"),
					ResourceType:    pJDCloudALB.ResourceType(maputil.GetString(options.ProviderExtendedConfig, "resourceType")),
					LoadbalancerId:  maputil.GetString(options.ProviderExtendedConfig, "loadbalancerId"),
					ListenerId:      maputil.GetString(options.ProviderExtendedConfig, "listenerId"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeJDCloudCDN:
				deployer, err := pJDCloudCDN.NewDeployer(&pJDCloudCDN.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Domain:          maputil.GetString(options.ProviderExtendedConfig, "domain"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeJDCloudLive:
				deployer, err := pJDCloudLive.NewDeployer(&pJDCloudLive.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Domain:          maputil.GetString(options.ProviderExtendedConfig, "domain"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeJDCloudVOD:
				deployer, err := pJDCloudVOD.NewDeployer(&pJDCloudVOD.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Domain:          maputil.GetString(options.ProviderExtendedConfig, "domain"),
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeploymentProviderTypeLocal:
		{
			deployer, err := pLocal.NewDeployer(&pLocal.DeployerConfig{
				ShellEnv:                 pLocal.ShellEnvType(maputil.GetString(options.ProviderExtendedConfig, "shellEnv")),
				PreCommand:               maputil.GetString(options.ProviderExtendedConfig, "preCommand"),
				PostCommand:              maputil.GetString(options.ProviderExtendedConfig, "postCommand"),
				OutputFormat:             pLocal.OutputFormatType(maputil.GetOrDefaultString(options.ProviderExtendedConfig, "format", string(pLocal.OUTPUT_FORMAT_PEM))),
				OutputCertPath:           maputil.GetString(options.ProviderExtendedConfig, "certPath"),
				OutputServerCertPath:     maputil.GetString(options.ProviderExtendedConfig, "certPathForServerOnly"),
				OutputIntermediaCertPath: maputil.GetString(options.ProviderExtendedConfig, "certPathForIntermediaOnly"),
				OutputKeyPath:            maputil.GetString(options.ProviderExtendedConfig, "keyPath"),
				PfxPassword:              maputil.GetString(options.ProviderExtendedConfig, "pfxPassword"),
				JksAlias:                 maputil.GetString(options.ProviderExtendedConfig, "jksAlias"),
				JksKeypass:               maputil.GetString(options.ProviderExtendedConfig, "jksKeypass"),
				JksStorepass:             maputil.GetString(options.ProviderExtendedConfig, "jksStorepass"),
			})
			return deployer, err
		}

	case domain.DeploymentProviderTypeKubernetesSecret:
		{
			access := domain.AccessConfigForKubernetes{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			deployer, err := pK8sSecret.NewDeployer(&pK8sSecret.DeployerConfig{
				KubeConfig:          access.KubeConfig,
				Namespace:           maputil.GetOrDefaultString(options.ProviderExtendedConfig, "namespace", "default"),
				SecretName:          maputil.GetString(options.ProviderExtendedConfig, "secretName"),
				SecretType:          maputil.GetOrDefaultString(options.ProviderExtendedConfig, "secretType", "kubernetes.io/tls"),
				SecretDataKeyForCrt: maputil.GetOrDefaultString(options.ProviderExtendedConfig, "secretDataKeyForCrt", "tls.crt"),
				SecretDataKeyForKey: maputil.GetOrDefaultString(options.ProviderExtendedConfig, "secretDataKeyForKey", "tls.key"),
			})
			return deployer, err
		}

	case domain.DeploymentProviderTypeNetlifySite:
		{
			access := domain.AccessConfigForNetlify{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			deployer, err := pNetlifySite.NewDeployer(&pNetlifySite.DeployerConfig{
				ApiToken: access.ApiToken,
				SiteId:   maputil.GetString(options.ProviderExtendedConfig, "siteId"),
			})
			return deployer, err
		}

	case domain.DeploymentProviderTypeProxmoxVE:
		{
			access := domain.AccessConfigForProxmoxVE{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			deployer, err := pProxmoxVE.NewDeployer(&pProxmoxVE.DeployerConfig{
				ApiUrl:                   access.ApiUrl,
				ApiToken:                 access.ApiToken,
				ApiTokenSecret:           access.ApiTokenSecret,
				AllowInsecureConnections: access.AllowInsecureConnections,
				NodeName:                 maputil.GetString(options.ProviderExtendedConfig, "nodeName"),
				AutoRestart:              maputil.GetBool(options.ProviderExtendedConfig, "autoRestart"),
			})
			return deployer, err
		}

	case domain.DeploymentProviderTypeQiniuCDN, domain.DeploymentProviderTypeQiniuKodo, domain.DeploymentProviderTypeQiniuPili:
		{
			access := domain.AccessConfigForQiniu{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeploymentProviderTypeQiniuCDN, domain.DeploymentProviderTypeQiniuKodo:
				deployer, err := pQiniuCDN.NewDeployer(&pQiniuCDN.DeployerConfig{
					AccessKey: access.AccessKey,
					SecretKey: access.SecretKey,
					Domain:    maputil.GetString(options.ProviderExtendedConfig, "domain"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeQiniuPili:
				deployer, err := pQiniuPili.NewDeployer(&pQiniuPili.DeployerConfig{
					AccessKey: access.AccessKey,
					SecretKey: access.SecretKey,
					Hub:       maputil.GetString(options.ProviderExtendedConfig, "hub"),
					Domain:    maputil.GetString(options.ProviderExtendedConfig, "domain"),
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeploymentProviderTypeRainYunRCDN:
		{
			access := domain.AccessConfigForRainYun{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeploymentProviderTypeTencentCloudCDN:
				deployer, err := pRainYunRCDN.NewDeployer(&pRainYunRCDN.DeployerConfig{
					ApiKey:     access.ApiKey,
					InstanceId: maputil.GetInt32(options.ProviderExtendedConfig, "instanceId"),
					Domain:     maputil.GetString(options.ProviderExtendedConfig, "domain"),
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeploymentProviderTypeSafeLine:
		{
			access := domain.AccessConfigForSafeLine{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			deployer, err := pSafeLine.NewDeployer(&pSafeLine.DeployerConfig{
				ApiUrl:                   access.ApiUrl,
				ApiToken:                 access.ApiToken,
				AllowInsecureConnections: access.AllowInsecureConnections,
				ResourceType:             pSafeLine.ResourceType(maputil.GetString(options.ProviderExtendedConfig, "resourceType")),
				CertificateId:            maputil.GetInt32(options.ProviderExtendedConfig, "certificateId"),
			})
			return deployer, err
		}

	case domain.DeploymentProviderTypeSSH:
		{
			access := domain.AccessConfigForSSH{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			deployer, err := pSSH.NewDeployer(&pSSH.DeployerConfig{
				SshHost:                  access.Host,
				SshPort:                  access.Port,
				SshUsername:              access.Username,
				SshPassword:              access.Password,
				SshKey:                   access.Key,
				SshKeyPassphrase:         access.KeyPassphrase,
				UseSCP:                   maputil.GetBool(options.ProviderExtendedConfig, "useSCP"),
				PreCommand:               maputil.GetString(options.ProviderExtendedConfig, "preCommand"),
				PostCommand:              maputil.GetString(options.ProviderExtendedConfig, "postCommand"),
				OutputFormat:             pSSH.OutputFormatType(maputil.GetOrDefaultString(options.ProviderExtendedConfig, "format", string(pSSH.OUTPUT_FORMAT_PEM))),
				OutputCertPath:           maputil.GetString(options.ProviderExtendedConfig, "certPath"),
				OutputServerCertPath:     maputil.GetString(options.ProviderExtendedConfig, "certPathForServerOnly"),
				OutputIntermediaCertPath: maputil.GetString(options.ProviderExtendedConfig, "certPathForIntermediaOnly"),
				OutputKeyPath:            maputil.GetString(options.ProviderExtendedConfig, "keyPath"),
				PfxPassword:              maputil.GetString(options.ProviderExtendedConfig, "pfxPassword"),
				JksAlias:                 maputil.GetString(options.ProviderExtendedConfig, "jksAlias"),
				JksKeypass:               maputil.GetString(options.ProviderExtendedConfig, "jksKeypass"),
				JksStorepass:             maputil.GetString(options.ProviderExtendedConfig, "jksStorepass"),
			})
			return deployer, err
		}

	case domain.DeploymentProviderTypeTencentCloudCDN, domain.DeploymentProviderTypeTencentCloudCLB, domain.DeploymentProviderTypeTencentCloudCOS, domain.DeploymentProviderTypeTencentCloudCSS, domain.DeploymentProviderTypeTencentCloudECDN, domain.DeploymentProviderTypeTencentCloudEO, domain.DeploymentProviderTypeTencentCloudSCF, domain.DeploymentProviderTypeTencentCloudSSL, domain.DeploymentProviderTypeTencentCloudSSLDeploy, domain.DeploymentProviderTypeTencentCloudVOD, domain.DeploymentProviderTypeTencentCloudWAF:
		{
			access := domain.AccessConfigForTencentCloud{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeploymentProviderTypeTencentCloudCDN:
				deployer, err := pTencentCloudCDN.NewDeployer(&pTencentCloudCDN.DeployerConfig{
					SecretId:  access.SecretId,
					SecretKey: access.SecretKey,
					Domain:    maputil.GetString(options.ProviderExtendedConfig, "domain"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeTencentCloudCLB:
				deployer, err := pTencentCloudCLB.NewDeployer(&pTencentCloudCLB.DeployerConfig{
					SecretId:       access.SecretId,
					SecretKey:      access.SecretKey,
					Region:         maputil.GetString(options.ProviderExtendedConfig, "region"),
					ResourceType:   pTencentCloudCLB.ResourceType(maputil.GetString(options.ProviderExtendedConfig, "resourceType")),
					LoadbalancerId: maputil.GetString(options.ProviderExtendedConfig, "loadbalancerId"),
					ListenerId:     maputil.GetString(options.ProviderExtendedConfig, "listenerId"),
					Domain:         maputil.GetString(options.ProviderExtendedConfig, "domain"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeTencentCloudCOS:
				deployer, err := pTencentCloudCOS.NewDeployer(&pTencentCloudCOS.DeployerConfig{
					SecretId:  access.SecretId,
					SecretKey: access.SecretKey,
					Region:    maputil.GetString(options.ProviderExtendedConfig, "region"),
					Bucket:    maputil.GetString(options.ProviderExtendedConfig, "bucket"),
					Domain:    maputil.GetString(options.ProviderExtendedConfig, "domain"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeTencentCloudCSS:
				deployer, err := pTencentCloudCSS.NewDeployer(&pTencentCloudCSS.DeployerConfig{
					SecretId:  access.SecretId,
					SecretKey: access.SecretKey,
					Domain:    maputil.GetString(options.ProviderExtendedConfig, "domain"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeTencentCloudECDN:
				deployer, err := pTencentCloudECDN.NewDeployer(&pTencentCloudECDN.DeployerConfig{
					SecretId:  access.SecretId,
					SecretKey: access.SecretKey,
					Domain:    maputil.GetString(options.ProviderExtendedConfig, "domain"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeTencentCloudEO:
				deployer, err := pTencentCloudEO.NewDeployer(&pTencentCloudEO.DeployerConfig{
					SecretId:  access.SecretId,
					SecretKey: access.SecretKey,
					ZoneId:    maputil.GetString(options.ProviderExtendedConfig, "zoneId"),
					Domain:    maputil.GetString(options.ProviderExtendedConfig, "domain"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeTencentCloudSCF:
				deployer, err := pTencentCloudSCF.NewDeployer(&pTencentCloudSCF.DeployerConfig{
					SecretId:  access.SecretId,
					SecretKey: access.SecretKey,
					Region:    maputil.GetString(options.ProviderExtendedConfig, "region"),
					Domain:    maputil.GetString(options.ProviderExtendedConfig, "domain"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeTencentCloudSSL:
				deployer, err := pTencentCloudSSL.NewDeployer(&pTencentCloudSSL.DeployerConfig{
					SecretId:  access.SecretId,
					SecretKey: access.SecretKey,
				})
				return deployer, err

			case domain.DeploymentProviderTypeTencentCloudSSLDeploy:
				deployer, err := pTencentCloudSSLDeploy.NewDeployer(&pTencentCloudSSLDeploy.DeployerConfig{
					SecretId:     access.SecretId,
					SecretKey:    access.SecretKey,
					Region:       maputil.GetString(options.ProviderExtendedConfig, "region"),
					ResourceType: maputil.GetString(options.ProviderExtendedConfig, "resourceType"),
					ResourceIds:  sliceutil.Filter(strings.Split(maputil.GetString(options.ProviderExtendedConfig, "resourceIds"), ";"), func(s string) bool { return s != "" }),
				})
				return deployer, err

			case domain.DeploymentProviderTypeTencentCloudVOD:
				deployer, err := pTencentCloudVOD.NewDeployer(&pTencentCloudVOD.DeployerConfig{
					SecretId:  access.SecretId,
					SecretKey: access.SecretKey,
					SubAppId:  maputil.GetInt64(options.ProviderExtendedConfig, "subAppId"),
					Domain:    maputil.GetString(options.ProviderExtendedConfig, "domain"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeTencentCloudWAF:
				deployer, err := pTencentCloudWAF.NewDeployer(&pTencentCloudWAF.DeployerConfig{
					SecretId:   access.SecretId,
					SecretKey:  access.SecretKey,
					Domain:     maputil.GetString(options.ProviderExtendedConfig, "domain"),
					DomainId:   maputil.GetString(options.ProviderExtendedConfig, "domainId"),
					InstanceId: maputil.GetString(options.ProviderExtendedConfig, "instanceId"),
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeploymentProviderTypeUCloudUCDN, domain.DeploymentProviderTypeUCloudUS3:
		{
			access := domain.AccessConfigForUCloud{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeploymentProviderTypeUCloudUCDN:
				deployer, err := pUCloudUCDN.NewDeployer(&pUCloudUCDN.DeployerConfig{
					PrivateKey: access.PrivateKey,
					PublicKey:  access.PublicKey,
					ProjectId:  access.ProjectId,
					DomainId:   maputil.GetString(options.ProviderExtendedConfig, "domainId"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeUCloudUS3:
				deployer, err := pUCloudUS3.NewDeployer(&pUCloudUS3.DeployerConfig{
					PrivateKey: access.PrivateKey,
					PublicKey:  access.PublicKey,
					ProjectId:  access.ProjectId,
					Region:     maputil.GetString(options.ProviderExtendedConfig, "region"),
					Bucket:     maputil.GetString(options.ProviderExtendedConfig, "bucket"),
					Domain:     maputil.GetString(options.ProviderExtendedConfig, "domain"),
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeploymentProviderTypeUpyunCDN, domain.DeploymentProviderTypeUpyunFile:
		{
			access := domain.AccessConfigForUpyun{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeploymentProviderTypeUpyunCDN, domain.DeploymentProviderTypeUpyunFile:
				deployer, err := pUpyunCDN.NewDeployer(&pUpyunCDN.DeployerConfig{
					Username: access.Username,
					Password: access.Password,
					Domain:   maputil.GetString(options.ProviderExtendedConfig, "domain"),
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeploymentProviderTypeVolcEngineALB, domain.DeploymentProviderTypeVolcEngineCDN, domain.DeploymentProviderTypeVolcEngineCertCenter, domain.DeploymentProviderTypeVolcEngineCLB, domain.DeploymentProviderTypeVolcEngineDCDN, domain.DeploymentProviderTypeVolcEngineImageX, domain.DeploymentProviderTypeVolcEngineLive, domain.DeploymentProviderTypeVolcEngineTOS:
		{
			access := domain.AccessConfigForVolcEngine{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeploymentProviderTypeVolcEngineALB:
				deployer, err := pVolcEngineALB.NewDeployer(&pVolcEngineALB.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.SecretAccessKey,
					Region:          maputil.GetString(options.ProviderExtendedConfig, "region"),
					ResourceType:    pVolcEngineALB.ResourceType(maputil.GetString(options.ProviderExtendedConfig, "resourceType")),
					LoadbalancerId:  maputil.GetString(options.ProviderExtendedConfig, "loadbalancerId"),
					ListenerId:      maputil.GetString(options.ProviderExtendedConfig, "listenerId"),
					Domain:          maputil.GetString(options.ProviderExtendedConfig, "domain"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeVolcEngineCDN:
				deployer, err := pVolcEngineCDN.NewDeployer(&pVolcEngineCDN.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.SecretAccessKey,
					Domain:          maputil.GetString(options.ProviderExtendedConfig, "domain"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeVolcEngineCertCenter:
				deployer, err := pVolcEngineCertCenter.NewDeployer(&pVolcEngineCertCenter.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.SecretAccessKey,
					Region:          maputil.GetString(options.ProviderExtendedConfig, "region"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeVolcEngineCLB:
				deployer, err := pVolcEngineCLB.NewDeployer(&pVolcEngineCLB.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.SecretAccessKey,
					Region:          maputil.GetString(options.ProviderExtendedConfig, "region"),
					ResourceType:    pVolcEngineCLB.ResourceType(maputil.GetString(options.ProviderExtendedConfig, "resourceType")),
					LoadbalancerId:  maputil.GetString(options.ProviderExtendedConfig, "loadbalancerId"),
					ListenerId:      maputil.GetString(options.ProviderExtendedConfig, "listenerId"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeVolcEngineDCDN:
				deployer, err := pVolcEngineDCDN.NewDeployer(&pVolcEngineDCDN.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.SecretAccessKey,
					Domain:          maputil.GetString(options.ProviderExtendedConfig, "domain"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeVolcEngineImageX:
				deployer, err := pVolcEngineImageX.NewDeployer(&pVolcEngineImageX.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.SecretAccessKey,
					Region:          maputil.GetString(options.ProviderExtendedConfig, "region"),
					ServiceId:       maputil.GetString(options.ProviderExtendedConfig, "serviceId"),
					Domain:          maputil.GetString(options.ProviderExtendedConfig, "domain"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeVolcEngineLive:
				deployer, err := pVolcEngineLive.NewDeployer(&pVolcEngineLive.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.SecretAccessKey,
					Domain:          maputil.GetString(options.ProviderExtendedConfig, "domain"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeVolcEngineTOS:
				deployer, err := pVolcEngineTOS.NewDeployer(&pVolcEngineTOS.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.SecretAccessKey,
					Region:          maputil.GetString(options.ProviderExtendedConfig, "region"),
					Bucket:          maputil.GetString(options.ProviderExtendedConfig, "bucket"),
					Domain:          maputil.GetString(options.ProviderExtendedConfig, "domain"),
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeploymentProviderTypeWangsuCDNPro:
		{
			access := domain.AccessConfigForWangsu{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeploymentProviderTypeWangsuCDNPro:
				deployer, err := pWangsuCDNPro.NewDeployer(&pWangsuCDNPro.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					ApiKey:          access.ApiKey,
					Environment:     maputil.GetOrDefaultString(options.ProviderExtendedConfig, "environment", "production"),
					Domain:          maputil.GetString(options.ProviderExtendedConfig, "domain"),
					CertificateId:   maputil.GetString(options.ProviderExtendedConfig, "certificateId"),
					WebhookId:       maputil.GetString(options.ProviderExtendedConfig, "webhookId"),
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeploymentProviderTypeWebhook:
		{
			access := domain.AccessConfigForWebhook{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			mergedHeaders := make(map[string]string)
			if defaultHeadersString := access.HeadersString; defaultHeadersString != "" {
				h, err := httputil.ParseHeaders(defaultHeadersString)
				if err != nil {
					return nil, fmt.Errorf("failed to parse webhook headers: %w", err)
				}
				for key := range h {
					mergedHeaders[http.CanonicalHeaderKey(key)] = h.Get(key)
				}
			}
			if extendedHeadersString := maputil.GetString(options.ProviderExtendedConfig, "headers"); extendedHeadersString != "" {
				h, err := httputil.ParseHeaders(extendedHeadersString)
				if err != nil {
					return nil, fmt.Errorf("failed to parse webhook headers: %w", err)
				}
				for key := range h {
					mergedHeaders[http.CanonicalHeaderKey(key)] = h.Get(key)
				}
			}

			deployer, err := pWebhook.NewDeployer(&pWebhook.DeployerConfig{
				WebhookUrl:               access.Url,
				WebhookData:              maputil.GetOrDefaultString(options.ProviderExtendedConfig, "webhookData", access.DefaultDataForDeployment),
				Method:                   access.Method,
				Headers:                  mergedHeaders,
				AllowInsecureConnections: access.AllowInsecureConnections,
			})
			return deployer, err
		}
	}

	return nil, fmt.Errorf("unsupported deployer provider '%s'", string(options.Provider))
}
