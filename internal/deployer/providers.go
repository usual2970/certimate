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
	pAliyunGA "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/aliyun-ga"
	pAliyunLive "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/aliyun-live"
	pAliyunNLB "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/aliyun-nlb"
	pAliyunOSS "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/aliyun-oss"
	pAliyunVOD "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/aliyun-vod"
	pAliyunWAF "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/aliyun-waf"
	pAPISIX "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/apisix"
	pAWSACM "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/aws-acm"
	pAWSCloudFront "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/aws-cloudfront"
	pAWSIAM "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/aws-iam"
	pAzureKeyVault "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/azure-keyvault"
	pBaiduCloudAppBLB "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/baiducloud-appblb"
	pBaiduCloudBLB "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/baiducloud-blb"
	pBaiduCloudCDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/baiducloud-cdn"
	pBaiduCloudCert "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/baiducloud-cert"
	pBaishanCDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/baishan-cdn"
	pBaotaPanelConsole "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/baotapanel-console"
	pBaotaPanelSite "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/baotapanel-site"
	pBaotaWAFConsole "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/baotawaf-console"
	pBaotaWAFSite "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/baotawaf-site"
	pBunnyCDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/bunny-cdn"
	pBytePlusCDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/byteplus-cdn"
	pCacheFly "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/cachefly"
	pCdnfly "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/cdnfly"
	pCTCCCloudAO "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/ctcccloud-ao"
	pCTCCCloudCDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/ctcccloud-cdn"
	pCTCCCloudCMS "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/ctcccloud-cms"
	pCTCCCloudELB "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/ctcccloud-elb"
	pCTCCCloudICDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/ctcccloud-icdn"
	pCTCCCloudLVDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/ctcccloud-lvdn"
	pDogeCDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/dogecloud-cdn"
	pEdgioApplications "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/edgio-applications"
	pFlexCDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/flexcdn"
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
	pLeCDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/lecdn"
	pLocal "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/local"
	pNetlifySite "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/netlify-site"
	pProxmoxVE "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/proxmoxve"
	pQiniuCDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/qiniu-cdn"
	pQiniuPili "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/qiniu-pili"
	pRainYunRCDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/rainyun-rcdn"
	pRatPanelConsole "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/ratpanel-console"
	pRatPanelSite "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/ratpanel-site"
	pSafeLine "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/safeline"
	pSSH "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/ssh"
	pTencentCloudCDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/tencentcloud-cdn"
	pTencentCloudCLB "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/tencentcloud-clb"
	pTencentCloudCOS "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/tencentcloud-cos"
	pTencentCloudCSS "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/tencentcloud-css"
	pTencentCloudECDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/tencentcloud-ecdn"
	pTencentCloudEO "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/tencentcloud-eo"
	pTencentCloudGAAP "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/tencentcloud-gaap"
	pTencentCloudSCF "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/tencentcloud-scf"
	pTencentCloudSSL "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/tencentcloud-ssl"
	pTencentCloudSSLDeploy "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/tencentcloud-ssl-deploy"
	pTencentCloudVOD "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/tencentcloud-vod"
	pTencentCloudWAF "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/tencentcloud-waf"
	pUCloudUCDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/ucloud-ucdn"
	pUCloudUS3 "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/ucloud-us3"
	pUniCloudWebHost "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/unicloud-webhost"
	pUpyunCDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/upyun-cdn"
	pVolcEngineALB "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/volcengine-alb"
	pVolcEngineCDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/volcengine-cdn"
	pVolcEngineCertCenter "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/volcengine-certcenter"
	pVolcEngineCLB "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/volcengine-clb"
	pVolcEngineDCDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/volcengine-dcdn"
	pVolcEngineImageX "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/volcengine-imagex"
	pVolcEngineLive "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/volcengine-live"
	pVolcEngineTOS "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/volcengine-tos"
	pWangsuCDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/wangsu-cdn"
	pWangsuCDNPro "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/wangsu-cdnpro"
	pWangsuCertificate "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/wangsu-certificate"
	pWebhook "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/webhook"
	xhttp "github.com/usual2970/certimate/internal/pkg/utils/http"
	xmaps "github.com/usual2970/certimate/internal/pkg/utils/maps"
	xslices "github.com/usual2970/certimate/internal/pkg/utils/slices"
)

type deployerProviderOptions struct {
	Provider              domain.DeploymentProviderType
	ProviderAccessConfig  map[string]any
	ProviderServiceConfig map[string]any
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
			if err := xmaps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeploymentProviderType1PanelConsole:
				deployer, err := p1PanelConsole.NewDeployer(&p1PanelConsole.DeployerConfig{
					ServerUrl:                access.ServerUrl,
					ApiVersion:               access.ApiVersion,
					ApiKey:                   access.ApiKey,
					AllowInsecureConnections: access.AllowInsecureConnections,
					AutoRestart:              xmaps.GetBool(options.ProviderServiceConfig, "autoRestart"),
				})
				return deployer, err

			case domain.DeploymentProviderType1PanelSite:
				deployer, err := p1PanelSite.NewDeployer(&p1PanelSite.DeployerConfig{
					ServerUrl:                access.ServerUrl,
					ApiVersion:               access.ApiVersion,
					ApiKey:                   access.ApiKey,
					AllowInsecureConnections: access.AllowInsecureConnections,
					ResourceType:             p1PanelSite.ResourceType(xmaps.GetOrDefaultString(options.ProviderServiceConfig, "resourceType", string(p1PanelSite.RESOURCE_TYPE_WEBSITE))),
					WebsiteId:                xmaps.GetInt64(options.ProviderServiceConfig, "websiteId"),
					CertificateId:            xmaps.GetInt64(options.ProviderServiceConfig, "certificateId"),
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeploymentProviderTypeAliyunALB, domain.DeploymentProviderTypeAliyunAPIGW, domain.DeploymentProviderTypeAliyunCAS, domain.DeploymentProviderTypeAliyunCASDeploy, domain.DeploymentProviderTypeAliyunCDN, domain.DeploymentProviderTypeAliyunCLB, domain.DeploymentProviderTypeAliyunDCDN, domain.DeploymentProviderTypeAliyunDDoS, domain.DeploymentProviderTypeAliyunESA, domain.DeploymentProviderTypeAliyunFC, domain.DeploymentProviderTypeAliyunGA, domain.DeploymentProviderTypeAliyunLive, domain.DeploymentProviderTypeAliyunNLB, domain.DeploymentProviderTypeAliyunOSS, domain.DeploymentProviderTypeAliyunVOD, domain.DeploymentProviderTypeAliyunWAF:
		{
			access := domain.AccessConfigForAliyun{}
			if err := xmaps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeploymentProviderTypeAliyunALB:
				deployer, err := pAliyunALB.NewDeployer(&pAliyunALB.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					ResourceGroupId: access.ResourceGroupId,
					Region:          xmaps.GetString(options.ProviderServiceConfig, "region"),
					ResourceType:    pAliyunALB.ResourceType(xmaps.GetString(options.ProviderServiceConfig, "resourceType")),
					LoadbalancerId:  xmaps.GetString(options.ProviderServiceConfig, "loadbalancerId"),
					ListenerId:      xmaps.GetString(options.ProviderServiceConfig, "listenerId"),
					Domain:          xmaps.GetString(options.ProviderServiceConfig, "domain"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeAliyunAPIGW:
				deployer, err := pAliyunAPIGW.NewDeployer(&pAliyunAPIGW.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					ResourceGroupId: access.ResourceGroupId,
					Region:          xmaps.GetString(options.ProviderServiceConfig, "region"),
					ServiceType:     pAliyunAPIGW.ServiceType(xmaps.GetString(options.ProviderServiceConfig, "serviceType")),
					GatewayId:       xmaps.GetString(options.ProviderServiceConfig, "gatewayId"),
					GroupId:         xmaps.GetString(options.ProviderServiceConfig, "groupId"),
					Domain:          xmaps.GetString(options.ProviderServiceConfig, "domain"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeAliyunCAS:
				deployer, err := pAliyunCAS.NewDeployer(&pAliyunCAS.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					ResourceGroupId: access.ResourceGroupId,
					Region:          xmaps.GetString(options.ProviderServiceConfig, "region"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeAliyunCASDeploy:
				deployer, err := pAliyunCASDeploy.NewDeployer(&pAliyunCASDeploy.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					ResourceGroupId: access.ResourceGroupId,
					Region:          xmaps.GetString(options.ProviderServiceConfig, "region"),
					ResourceIds:     xslices.Filter(strings.Split(xmaps.GetString(options.ProviderServiceConfig, "resourceIds"), ";"), func(s string) bool { return s != "" }),
					ContactIds:      xslices.Filter(strings.Split(xmaps.GetString(options.ProviderServiceConfig, "contactIds"), ";"), func(s string) bool { return s != "" }),
				})
				return deployer, err

			case domain.DeploymentProviderTypeAliyunCDN:
				deployer, err := pAliyunCDN.NewDeployer(&pAliyunCDN.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					ResourceGroupId: access.ResourceGroupId,
					Domain:          xmaps.GetString(options.ProviderServiceConfig, "domain"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeAliyunCLB:
				deployer, err := pAliyunCLB.NewDeployer(&pAliyunCLB.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					ResourceGroupId: access.ResourceGroupId,
					Region:          xmaps.GetString(options.ProviderServiceConfig, "region"),
					ResourceType:    pAliyunCLB.ResourceType(xmaps.GetString(options.ProviderServiceConfig, "resourceType")),
					LoadbalancerId:  xmaps.GetString(options.ProviderServiceConfig, "loadbalancerId"),
					ListenerPort:    xmaps.GetOrDefaultInt32(options.ProviderServiceConfig, "listenerPort", 443),
					Domain:          xmaps.GetString(options.ProviderServiceConfig, "domain"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeAliyunDCDN:
				deployer, err := pAliyunDCDN.NewDeployer(&pAliyunDCDN.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					ResourceGroupId: access.ResourceGroupId,
					Domain:          xmaps.GetString(options.ProviderServiceConfig, "domain"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeAliyunDDoS:
				deployer, err := pAliyunDDoS.NewDeployer(&pAliyunDDoS.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					ResourceGroupId: access.ResourceGroupId,
					Region:          xmaps.GetString(options.ProviderServiceConfig, "region"),
					Domain:          xmaps.GetString(options.ProviderServiceConfig, "domain"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeAliyunESA:
				deployer, err := pAliyunESA.NewDeployer(&pAliyunESA.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Region:          xmaps.GetString(options.ProviderServiceConfig, "region"),
					SiteId:          xmaps.GetInt64(options.ProviderServiceConfig, "siteId"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeAliyunFC:
				deployer, err := pAliyunFC.NewDeployer(&pAliyunFC.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					ResourceGroupId: access.ResourceGroupId,
					Region:          xmaps.GetString(options.ProviderServiceConfig, "region"),
					ServiceVersion:  xmaps.GetOrDefaultString(options.ProviderServiceConfig, "serviceVersion", "3.0"),
					Domain:          xmaps.GetString(options.ProviderServiceConfig, "domain"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeAliyunGA:
				deployer, err := pAliyunGA.NewDeployer(&pAliyunGA.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					ResourceGroupId: access.ResourceGroupId,
					ResourceType:    pAliyunGA.ResourceType(xmaps.GetString(options.ProviderServiceConfig, "resourceType")),
					AcceleratorId:   xmaps.GetString(options.ProviderServiceConfig, "acceleratorId"),
					ListenerId:      xmaps.GetString(options.ProviderServiceConfig, "listenerId"),
					Domain:          xmaps.GetString(options.ProviderServiceConfig, "domain"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeAliyunLive:
				deployer, err := pAliyunLive.NewDeployer(&pAliyunLive.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Region:          xmaps.GetString(options.ProviderServiceConfig, "region"),
					Domain:          xmaps.GetString(options.ProviderServiceConfig, "domain"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeAliyunNLB:
				deployer, err := pAliyunNLB.NewDeployer(&pAliyunNLB.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					ResourceGroupId: access.ResourceGroupId,
					Region:          xmaps.GetString(options.ProviderServiceConfig, "region"),
					ResourceType:    pAliyunNLB.ResourceType(xmaps.GetString(options.ProviderServiceConfig, "resourceType")),
					LoadbalancerId:  xmaps.GetString(options.ProviderServiceConfig, "loadbalancerId"),
					ListenerId:      xmaps.GetString(options.ProviderServiceConfig, "listenerId"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeAliyunOSS:
				deployer, err := pAliyunOSS.NewDeployer(&pAliyunOSS.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					ResourceGroupId: access.ResourceGroupId,
					Region:          xmaps.GetString(options.ProviderServiceConfig, "region"),
					Bucket:          xmaps.GetString(options.ProviderServiceConfig, "bucket"),
					Domain:          xmaps.GetString(options.ProviderServiceConfig, "domain"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeAliyunVOD:
				deployer, err := pAliyunVOD.NewDeployer(&pAliyunVOD.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					ResourceGroupId: access.ResourceGroupId,
					Region:          xmaps.GetString(options.ProviderServiceConfig, "region"),
					Domain:          xmaps.GetString(options.ProviderServiceConfig, "domain"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeAliyunWAF:
				deployer, err := pAliyunWAF.NewDeployer(&pAliyunWAF.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					ResourceGroupId: access.ResourceGroupId,
					Region:          xmaps.GetString(options.ProviderServiceConfig, "region"),
					ServiceVersion:  xmaps.GetOrDefaultString(options.ProviderServiceConfig, "serviceVersion", "3.0"),
					InstanceId:      xmaps.GetString(options.ProviderServiceConfig, "instanceId"),
					Domain:          xmaps.GetString(options.ProviderServiceConfig, "domain"),
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeploymentProviderTypeAPISIX:
		{
			access := domain.AccessConfigForAPISIX{}
			if err := xmaps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			deployer, err := pAPISIX.NewDeployer(&pAPISIX.DeployerConfig{
				ServerUrl:                access.ServerUrl,
				ApiKey:                   access.ApiKey,
				AllowInsecureConnections: access.AllowInsecureConnections,
				ResourceType:             pAPISIX.ResourceType(xmaps.GetString(options.ProviderServiceConfig, "resourceType")),
				CertificateId:            xmaps.GetString(options.ProviderServiceConfig, "certificateId"),
			})
			return deployer, err
		}

	case domain.DeploymentProviderTypeAWSACM, domain.DeploymentProviderTypeAWSCloudFront, domain.DeploymentProviderTypeAWSIAM:
		{
			access := domain.AccessConfigForAWS{}
			if err := xmaps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeploymentProviderTypeAWSACM:
				deployer, err := pAWSACM.NewDeployer(&pAWSACM.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					SecretAccessKey: access.SecretAccessKey,
					Region:          xmaps.GetString(options.ProviderServiceConfig, "region"),
					CertificateArn:  xmaps.GetString(options.ProviderServiceConfig, "certificateArn"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeAWSCloudFront:
				deployer, err := pAWSCloudFront.NewDeployer(&pAWSCloudFront.DeployerConfig{
					AccessKeyId:       access.AccessKeyId,
					SecretAccessKey:   access.SecretAccessKey,
					Region:            xmaps.GetString(options.ProviderServiceConfig, "region"),
					DistributionId:    xmaps.GetString(options.ProviderServiceConfig, "distributionId"),
					CertificateSource: xmaps.GetOrDefaultString(options.ProviderServiceConfig, "certificateSource", "ACM"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeAWSIAM:
				deployer, err := pAWSIAM.NewDeployer(&pAWSIAM.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					SecretAccessKey: access.SecretAccessKey,
					Region:          xmaps.GetString(options.ProviderServiceConfig, "region"),
					CertificatePath: xmaps.GetOrDefaultString(options.ProviderServiceConfig, "certificatePath", "/"),
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeploymentProviderTypeAzureKeyVault:
		{
			access := domain.AccessConfigForAzure{}
			if err := xmaps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeploymentProviderTypeAzureKeyVault:
				deployer, err := pAzureKeyVault.NewDeployer(&pAzureKeyVault.DeployerConfig{
					TenantId:        access.TenantId,
					ClientId:        access.ClientId,
					ClientSecret:    access.ClientSecret,
					CloudName:       access.CloudName,
					KeyVaultName:    xmaps.GetString(options.ProviderServiceConfig, "keyvaultName"),
					CertificateName: xmaps.GetString(options.ProviderServiceConfig, "certificateName"),
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeploymentProviderTypeBaiduCloudAppBLB, domain.DeploymentProviderTypeBaiduCloudBLB, domain.DeploymentProviderTypeBaiduCloudCDN, domain.DeploymentProviderTypeBaiduCloudCert:
		{
			access := domain.AccessConfigForBaiduCloud{}
			if err := xmaps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeploymentProviderTypeBaiduCloudAppBLB:
				deployer, err := pBaiduCloudAppBLB.NewDeployer(&pBaiduCloudAppBLB.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					SecretAccessKey: access.SecretAccessKey,
					Region:          xmaps.GetString(options.ProviderServiceConfig, "region"),
					ResourceType:    pBaiduCloudAppBLB.ResourceType(xmaps.GetString(options.ProviderServiceConfig, "resourceType")),
					LoadbalancerId:  xmaps.GetString(options.ProviderServiceConfig, "loadbalancerId"),
					ListenerPort:    xmaps.GetInt32(options.ProviderServiceConfig, "listenerPort"),
					Domain:          xmaps.GetString(options.ProviderServiceConfig, "domain"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeBaiduCloudBLB:
				deployer, err := pBaiduCloudBLB.NewDeployer(&pBaiduCloudBLB.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					SecretAccessKey: access.SecretAccessKey,
					Region:          xmaps.GetString(options.ProviderServiceConfig, "region"),
					ResourceType:    pBaiduCloudBLB.ResourceType(xmaps.GetString(options.ProviderServiceConfig, "resourceType")),
					LoadbalancerId:  xmaps.GetString(options.ProviderServiceConfig, "loadbalancerId"),
					ListenerPort:    xmaps.GetInt32(options.ProviderServiceConfig, "listenerPort"),
					Domain:          xmaps.GetString(options.ProviderServiceConfig, "domain"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeBaiduCloudCDN:
				deployer, err := pBaiduCloudCDN.NewDeployer(&pBaiduCloudCDN.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					SecretAccessKey: access.SecretAccessKey,
					Domain:          xmaps.GetString(options.ProviderServiceConfig, "domain"),
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
			if err := xmaps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeploymentProviderTypeBaishanCDN:
				deployer, err := pBaishanCDN.NewDeployer(&pBaishanCDN.DeployerConfig{
					ApiToken:      access.ApiToken,
					Domain:        xmaps.GetString(options.ProviderServiceConfig, "domain"),
					CertificateId: xmaps.GetString(options.ProviderServiceConfig, "certificateId"),
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeploymentProviderTypeBaotaPanelConsole, domain.DeploymentProviderTypeBaotaPanelSite:
		{
			access := domain.AccessConfigForBaotaPanel{}
			if err := xmaps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeploymentProviderTypeBaotaPanelConsole:
				deployer, err := pBaotaPanelConsole.NewDeployer(&pBaotaPanelConsole.DeployerConfig{
					ServerUrl:                access.ServerUrl,
					ApiKey:                   access.ApiKey,
					AllowInsecureConnections: access.AllowInsecureConnections,
					AutoRestart:              xmaps.GetBool(options.ProviderServiceConfig, "autoRestart"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeBaotaPanelSite:
				deployer, err := pBaotaPanelSite.NewDeployer(&pBaotaPanelSite.DeployerConfig{
					ServerUrl:                access.ServerUrl,
					ApiKey:                   access.ApiKey,
					AllowInsecureConnections: access.AllowInsecureConnections,
					SiteType:                 xmaps.GetOrDefaultString(options.ProviderServiceConfig, "siteType", "other"),
					SiteName:                 xmaps.GetString(options.ProviderServiceConfig, "siteName"),
					SiteNames:                xslices.Filter(strings.Split(xmaps.GetString(options.ProviderServiceConfig, "siteNames"), ";"), func(s string) bool { return s != "" }),
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeploymentProviderTypeBaotaWAFConsole, domain.DeploymentProviderTypeBaotaWAFSite:
		{
			access := domain.AccessConfigForBaotaWAF{}
			if err := xmaps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeploymentProviderTypeBaotaWAFConsole:
				deployer, err := pBaotaWAFConsole.NewDeployer(&pBaotaWAFConsole.DeployerConfig{
					ServerUrl:                access.ServerUrl,
					ApiKey:                   access.ApiKey,
					AllowInsecureConnections: access.AllowInsecureConnections,
				})
				return deployer, err

			case domain.DeploymentProviderTypeBaotaWAFSite:
				deployer, err := pBaotaWAFSite.NewDeployer(&pBaotaWAFSite.DeployerConfig{
					ServerUrl:                access.ServerUrl,
					ApiKey:                   access.ApiKey,
					AllowInsecureConnections: access.AllowInsecureConnections,
					SiteName:                 xmaps.GetString(options.ProviderServiceConfig, "siteName"),
					SitePort:                 xmaps.GetOrDefaultInt32(options.ProviderServiceConfig, "sitePort", 443),
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeploymentProviderTypeBunnyCDN:
		{
			access := domain.AccessConfigForBunny{}
			if err := xmaps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			deployer, err := pBunnyCDN.NewDeployer(&pBunnyCDN.DeployerConfig{
				ApiKey:     access.ApiKey,
				PullZoneId: xmaps.GetString(options.ProviderServiceConfig, "pullZoneId"),
				Hostname:   xmaps.GetString(options.ProviderServiceConfig, "hostname"),
			})
			return deployer, err
		}

	case domain.DeploymentProviderTypeBytePlusCDN:
		{
			access := domain.AccessConfigForBytePlus{}
			if err := xmaps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeploymentProviderTypeBytePlusCDN:
				deployer, err := pBytePlusCDN.NewDeployer(&pBytePlusCDN.DeployerConfig{
					AccessKey: access.AccessKey,
					SecretKey: access.SecretKey,
					Domain:    xmaps.GetString(options.ProviderServiceConfig, "domain"),
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeploymentProviderTypeCacheFly:
		{
			access := domain.AccessConfigForCacheFly{}
			if err := xmaps.Populate(options.ProviderAccessConfig, &access); err != nil {
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
			if err := xmaps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			deployer, err := pCdnfly.NewDeployer(&pCdnfly.DeployerConfig{
				ServerUrl:                access.ServerUrl,
				ApiKey:                   access.ApiKey,
				ApiSecret:                access.ApiSecret,
				AllowInsecureConnections: access.AllowInsecureConnections,
				ResourceType:             pCdnfly.ResourceType(xmaps.GetOrDefaultString(options.ProviderServiceConfig, "resourceType", string(pCdnfly.RESOURCE_TYPE_SITE))),
				SiteId:                   xmaps.GetString(options.ProviderServiceConfig, "siteId"),
				CertificateId:            xmaps.GetString(options.ProviderServiceConfig, "certificateId"),
			})
			return deployer, err
		}

	case domain.DeploymentProviderTypeCTCCCloudAO, domain.DeploymentProviderTypeCTCCCloudCDN, domain.DeploymentProviderTypeCTCCCloudCMS, domain.DeploymentProviderTypeCTCCCloudELB, domain.DeploymentProviderTypeCTCCCloudICDN, domain.DeploymentProviderTypeCTCCCloudLVDN:
		{
			access := domain.AccessConfigForCTCCCloud{}
			if err := xmaps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeploymentProviderTypeCTCCCloudAO:
				deployer, err := pCTCCCloudAO.NewDeployer(&pCTCCCloudAO.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					SecretAccessKey: access.SecretAccessKey,
					Domain:          xmaps.GetString(options.ProviderServiceConfig, "domain"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeCTCCCloudCDN:
				deployer, err := pCTCCCloudCDN.NewDeployer(&pCTCCCloudCDN.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					SecretAccessKey: access.SecretAccessKey,
					Domain:          xmaps.GetString(options.ProviderServiceConfig, "domain"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeCTCCCloudCMS:
				deployer, err := pCTCCCloudCMS.NewDeployer(&pCTCCCloudCMS.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					SecretAccessKey: access.SecretAccessKey,
				})
				return deployer, err

			case domain.DeploymentProviderTypeCTCCCloudELB:
				deployer, err := pCTCCCloudELB.NewDeployer(&pCTCCCloudELB.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					SecretAccessKey: access.SecretAccessKey,
					RegionId:        xmaps.GetString(options.ProviderServiceConfig, "regionId"),
					ResourceType:    pCTCCCloudELB.ResourceType(xmaps.GetString(options.ProviderServiceConfig, "resourceType")),
					LoadbalancerId:  xmaps.GetString(options.ProviderServiceConfig, "loadbalancerId"),
					ListenerId:      xmaps.GetString(options.ProviderServiceConfig, "listenerId"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeCTCCCloudICDN:
				deployer, err := pCTCCCloudICDN.NewDeployer(&pCTCCCloudICDN.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					SecretAccessKey: access.SecretAccessKey,
					Domain:          xmaps.GetString(options.ProviderServiceConfig, "domain"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeCTCCCloudLVDN:
				deployer, err := pCTCCCloudLVDN.NewDeployer(&pCTCCCloudLVDN.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					SecretAccessKey: access.SecretAccessKey,
					Domain:          xmaps.GetString(options.ProviderServiceConfig, "domain"),
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeploymentProviderTypeDogeCloudCDN:
		{
			access := domain.AccessConfigForDogeCloud{}
			if err := xmaps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			deployer, err := pDogeCDN.NewDeployer(&pDogeCDN.DeployerConfig{
				AccessKey: access.AccessKey,
				SecretKey: access.SecretKey,
				Domain:    xmaps.GetString(options.ProviderServiceConfig, "domain"),
			})
			return deployer, err
		}

	case domain.DeploymentProviderTypeEdgioApplications:
		{
			access := domain.AccessConfigForEdgio{}
			if err := xmaps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			deployer, err := pEdgioApplications.NewDeployer(&pEdgioApplications.DeployerConfig{
				ClientId:      access.ClientId,
				ClientSecret:  access.ClientSecret,
				EnvironmentId: xmaps.GetString(options.ProviderServiceConfig, "environmentId"),
			})
			return deployer, err
		}

	case domain.DeploymentProviderTypeFlexCDN:
		{
			access := domain.AccessConfigForFlexCDN{}
			if err := xmaps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			deployer, err := pFlexCDN.NewDeployer(&pFlexCDN.DeployerConfig{
				ServerUrl:                access.ServerUrl,
				ApiRole:                  access.ApiRole,
				AccessKeyId:              access.AccessKeyId,
				AccessKey:                access.AccessKey,
				AllowInsecureConnections: access.AllowInsecureConnections,
				ResourceType:             pFlexCDN.ResourceType(xmaps.GetString(options.ProviderServiceConfig, "resourceType")),
				CertificateId:            xmaps.GetInt64(options.ProviderServiceConfig, "certificateId"),
			})
			return deployer, err
		}

	case domain.DeploymentProviderTypeGcoreCDN:
		{
			access := domain.AccessConfigForGcore{}
			if err := xmaps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeploymentProviderTypeGcoreCDN:
				deployer, err := pGcoreCDN.NewDeployer(&pGcoreCDN.DeployerConfig{
					ApiToken:      access.ApiToken,
					ResourceId:    xmaps.GetInt64(options.ProviderServiceConfig, "resourceId"),
					CertificateId: xmaps.GetInt64(options.ProviderServiceConfig, "certificateId"),
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeploymentProviderTypeGoEdge:
		{
			access := domain.AccessConfigForGoEdge{}
			if err := xmaps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			deployer, err := pGoEdge.NewDeployer(&pGoEdge.DeployerConfig{
				ServerUrl:                access.ServerUrl,
				ApiRole:                  access.ApiRole,
				AccessKeyId:              access.AccessKeyId,
				AccessKey:                access.AccessKey,
				AllowInsecureConnections: access.AllowInsecureConnections,
				ResourceType:             pGoEdge.ResourceType(xmaps.GetString(options.ProviderServiceConfig, "resourceType")),
				CertificateId:            xmaps.GetInt64(options.ProviderServiceConfig, "certificateId"),
			})
			return deployer, err
		}

	case domain.DeploymentProviderTypeHuaweiCloudCDN, domain.DeploymentProviderTypeHuaweiCloudELB, domain.DeploymentProviderTypeHuaweiCloudSCM, domain.DeploymentProviderTypeHuaweiCloudWAF:
		{
			access := domain.AccessConfigForHuaweiCloud{}
			if err := xmaps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeploymentProviderTypeHuaweiCloudCDN:
				deployer, err := pHuaweiCloudCDN.NewDeployer(&pHuaweiCloudCDN.DeployerConfig{
					AccessKeyId:         access.AccessKeyId,
					SecretAccessKey:     access.SecretAccessKey,
					EnterpriseProjectId: access.EnterpriseProjectId,
					Region:              xmaps.GetString(options.ProviderServiceConfig, "region"),
					Domain:              xmaps.GetString(options.ProviderServiceConfig, "domain"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeHuaweiCloudELB:
				deployer, err := pHuaweiCloudELB.NewDeployer(&pHuaweiCloudELB.DeployerConfig{
					AccessKeyId:         access.AccessKeyId,
					SecretAccessKey:     access.SecretAccessKey,
					EnterpriseProjectId: access.EnterpriseProjectId,
					Region:              xmaps.GetString(options.ProviderServiceConfig, "region"),
					ResourceType:        pHuaweiCloudELB.ResourceType(xmaps.GetString(options.ProviderServiceConfig, "resourceType")),
					CertificateId:       xmaps.GetString(options.ProviderServiceConfig, "certificateId"),
					LoadbalancerId:      xmaps.GetString(options.ProviderServiceConfig, "loadbalancerId"),
					ListenerId:          xmaps.GetString(options.ProviderServiceConfig, "listenerId"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeHuaweiCloudSCM:
				deployer, err := pHuaweiCloudSCM.NewDeployer(&pHuaweiCloudSCM.DeployerConfig{
					AccessKeyId:         access.AccessKeyId,
					SecretAccessKey:     access.SecretAccessKey,
					EnterpriseProjectId: access.EnterpriseProjectId,
				})
				return deployer, err

			case domain.DeploymentProviderTypeHuaweiCloudWAF:
				deployer, err := pHuaweiCloudWAF.NewDeployer(&pHuaweiCloudWAF.DeployerConfig{
					AccessKeyId:         access.AccessKeyId,
					SecretAccessKey:     access.SecretAccessKey,
					EnterpriseProjectId: access.EnterpriseProjectId,
					Region:              xmaps.GetString(options.ProviderServiceConfig, "region"),
					ResourceType:        pHuaweiCloudWAF.ResourceType(xmaps.GetString(options.ProviderServiceConfig, "resourceType")),
					CertificateId:       xmaps.GetString(options.ProviderServiceConfig, "certificateId"),
					Domain:              xmaps.GetString(options.ProviderServiceConfig, "domain"),
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeploymentProviderTypeJDCloudALB, domain.DeploymentProviderTypeJDCloudCDN, domain.DeploymentProviderTypeJDCloudLive, domain.DeploymentProviderTypeJDCloudVOD:
		{
			access := domain.AccessConfigForJDCloud{}
			if err := xmaps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeploymentProviderTypeJDCloudALB:
				deployer, err := pJDCloudALB.NewDeployer(&pJDCloudALB.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					RegionId:        xmaps.GetString(options.ProviderServiceConfig, "regionId"),
					ResourceType:    pJDCloudALB.ResourceType(xmaps.GetString(options.ProviderServiceConfig, "resourceType")),
					LoadbalancerId:  xmaps.GetString(options.ProviderServiceConfig, "loadbalancerId"),
					ListenerId:      xmaps.GetString(options.ProviderServiceConfig, "listenerId"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeJDCloudCDN:
				deployer, err := pJDCloudCDN.NewDeployer(&pJDCloudCDN.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Domain:          xmaps.GetString(options.ProviderServiceConfig, "domain"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeJDCloudLive:
				deployer, err := pJDCloudLive.NewDeployer(&pJDCloudLive.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Domain:          xmaps.GetString(options.ProviderServiceConfig, "domain"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeJDCloudVOD:
				deployer, err := pJDCloudVOD.NewDeployer(&pJDCloudVOD.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Domain:          xmaps.GetString(options.ProviderServiceConfig, "domain"),
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeploymentProviderTypeLeCDN:
		{
			access := domain.AccessConfigForLeCDN{}
			if err := xmaps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			deployer, err := pLeCDN.NewDeployer(&pLeCDN.DeployerConfig{
				ServerUrl:                access.ServerUrl,
				ApiVersion:               access.ApiVersion,
				ApiRole:                  access.ApiRole,
				Username:                 access.Username,
				Password:                 access.Password,
				AllowInsecureConnections: access.AllowInsecureConnections,
				ResourceType:             pLeCDN.ResourceType(xmaps.GetString(options.ProviderServiceConfig, "resourceType")),
				CertificateId:            xmaps.GetInt64(options.ProviderServiceConfig, "certificateId"),
				ClientId:                 xmaps.GetInt64(options.ProviderServiceConfig, "clientId"),
			})
			return deployer, err
		}

	case domain.DeploymentProviderTypeLocal:
		{
			deployer, err := pLocal.NewDeployer(&pLocal.DeployerConfig{
				ShellEnv:                 pLocal.ShellEnvType(xmaps.GetString(options.ProviderServiceConfig, "shellEnv")),
				PreCommand:               xmaps.GetString(options.ProviderServiceConfig, "preCommand"),
				PostCommand:              xmaps.GetString(options.ProviderServiceConfig, "postCommand"),
				OutputFormat:             pLocal.OutputFormatType(xmaps.GetOrDefaultString(options.ProviderServiceConfig, "format", string(pLocal.OUTPUT_FORMAT_PEM))),
				OutputCertPath:           xmaps.GetString(options.ProviderServiceConfig, "certPath"),
				OutputServerCertPath:     xmaps.GetString(options.ProviderServiceConfig, "certPathForServerOnly"),
				OutputIntermediaCertPath: xmaps.GetString(options.ProviderServiceConfig, "certPathForIntermediaOnly"),
				OutputKeyPath:            xmaps.GetString(options.ProviderServiceConfig, "keyPath"),
				PfxPassword:              xmaps.GetString(options.ProviderServiceConfig, "pfxPassword"),
				JksAlias:                 xmaps.GetString(options.ProviderServiceConfig, "jksAlias"),
				JksKeypass:               xmaps.GetString(options.ProviderServiceConfig, "jksKeypass"),
				JksStorepass:             xmaps.GetString(options.ProviderServiceConfig, "jksStorepass"),
			})
			return deployer, err
		}

	case domain.DeploymentProviderTypeKubernetesSecret:
		{
			access := domain.AccessConfigForKubernetes{}
			if err := xmaps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			deployer, err := pK8sSecret.NewDeployer(&pK8sSecret.DeployerConfig{
				KubeConfig:          access.KubeConfig,
				Namespace:           xmaps.GetOrDefaultString(options.ProviderServiceConfig, "namespace", "default"),
				SecretName:          xmaps.GetString(options.ProviderServiceConfig, "secretName"),
				SecretType:          xmaps.GetOrDefaultString(options.ProviderServiceConfig, "secretType", "kubernetes.io/tls"),
				SecretDataKeyForCrt: xmaps.GetOrDefaultString(options.ProviderServiceConfig, "secretDataKeyForCrt", "tls.crt"),
				SecretDataKeyForKey: xmaps.GetOrDefaultString(options.ProviderServiceConfig, "secretDataKeyForKey", "tls.key"),
			})
			return deployer, err
		}

	case domain.DeploymentProviderTypeNetlifySite:
		{
			access := domain.AccessConfigForNetlify{}
			if err := xmaps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			deployer, err := pNetlifySite.NewDeployer(&pNetlifySite.DeployerConfig{
				ApiToken: access.ApiToken,
				SiteId:   xmaps.GetString(options.ProviderServiceConfig, "siteId"),
			})
			return deployer, err
		}

	case domain.DeploymentProviderTypeProxmoxVE:
		{
			access := domain.AccessConfigForProxmoxVE{}
			if err := xmaps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			deployer, err := pProxmoxVE.NewDeployer(&pProxmoxVE.DeployerConfig{
				ServerUrl:                access.ServerUrl,
				ApiToken:                 access.ApiToken,
				ApiTokenSecret:           access.ApiTokenSecret,
				AllowInsecureConnections: access.AllowInsecureConnections,
				NodeName:                 xmaps.GetString(options.ProviderServiceConfig, "nodeName"),
				AutoRestart:              xmaps.GetBool(options.ProviderServiceConfig, "autoRestart"),
			})
			return deployer, err
		}

	case domain.DeploymentProviderTypeQiniuCDN, domain.DeploymentProviderTypeQiniuKodo, domain.DeploymentProviderTypeQiniuPili:
		{
			access := domain.AccessConfigForQiniu{}
			if err := xmaps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeploymentProviderTypeQiniuCDN, domain.DeploymentProviderTypeQiniuKodo:
				deployer, err := pQiniuCDN.NewDeployer(&pQiniuCDN.DeployerConfig{
					AccessKey: access.AccessKey,
					SecretKey: access.SecretKey,
					Domain:    xmaps.GetString(options.ProviderServiceConfig, "domain"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeQiniuPili:
				deployer, err := pQiniuPili.NewDeployer(&pQiniuPili.DeployerConfig{
					AccessKey: access.AccessKey,
					SecretKey: access.SecretKey,
					Hub:       xmaps.GetString(options.ProviderServiceConfig, "hub"),
					Domain:    xmaps.GetString(options.ProviderServiceConfig, "domain"),
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeploymentProviderTypeRainYunRCDN:
		{
			access := domain.AccessConfigForRainYun{}
			if err := xmaps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeploymentProviderTypeTencentCloudCDN:
				deployer, err := pRainYunRCDN.NewDeployer(&pRainYunRCDN.DeployerConfig{
					ApiKey:     access.ApiKey,
					InstanceId: xmaps.GetInt32(options.ProviderServiceConfig, "instanceId"),
					Domain:     xmaps.GetString(options.ProviderServiceConfig, "domain"),
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeploymentProviderTypeRatPanelConsole, domain.DeploymentProviderTypeRatPanelSite:
		{
			access := domain.AccessConfigForRatPanel{}
			if err := xmaps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeploymentProviderTypeRatPanelConsole:
				deployer, err := pRatPanelConsole.NewDeployer(&pRatPanelConsole.DeployerConfig{
					ServerUrl:                access.ServerUrl,
					AccessTokenId:            access.AccessTokenId,
					AccessToken:              access.AccessToken,
					AllowInsecureConnections: access.AllowInsecureConnections,
				})
				return deployer, err

			case domain.DeploymentProviderTypeRatPanelSite:
				deployer, err := pRatPanelSite.NewDeployer(&pRatPanelSite.DeployerConfig{
					ServerUrl:                access.ServerUrl,
					AccessTokenId:            access.AccessTokenId,
					AccessToken:              access.AccessToken,
					AllowInsecureConnections: access.AllowInsecureConnections,
					SiteName:                 xmaps.GetString(options.ProviderServiceConfig, "siteName"),
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeploymentProviderTypeSafeLine:
		{
			access := domain.AccessConfigForSafeLine{}
			if err := xmaps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			deployer, err := pSafeLine.NewDeployer(&pSafeLine.DeployerConfig{
				ServerUrl:                access.ServerUrl,
				ApiToken:                 access.ApiToken,
				AllowInsecureConnections: access.AllowInsecureConnections,
				ResourceType:             pSafeLine.ResourceType(xmaps.GetString(options.ProviderServiceConfig, "resourceType")),
				CertificateId:            xmaps.GetInt32(options.ProviderServiceConfig, "certificateId"),
			})
			return deployer, err
		}

	case domain.DeploymentProviderTypeSSH:
		{
			access := domain.AccessConfigForSSH{}
			if err := xmaps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			jumpServers := make([]pSSH.JumpServerConfig, len(access.JumpServers))
			for i, jumpServer := range access.JumpServers {
				jumpServers[i] = pSSH.JumpServerConfig{
					SshHost:          jumpServer.Host,
					SshPort:          jumpServer.Port,
					SshAuthMethod:    jumpServer.AuthMethod,
					SshUsername:      jumpServer.Username,
					SshPassword:      jumpServer.Password,
					SshKey:           jumpServer.Key,
					SshKeyPassphrase: jumpServer.KeyPassphrase,
				}
			}

			deployer, err := pSSH.NewDeployer(&pSSH.DeployerConfig{
				SshHost:                  access.Host,
				SshPort:                  access.Port,
				SshAuthMethod:            access.AuthMethod,
				SshUsername:              access.Username,
				SshPassword:              access.Password,
				SshKey:                   access.Key,
				SshKeyPassphrase:         access.KeyPassphrase,
				JumpServers:              jumpServers,
				UseSCP:                   xmaps.GetBool(options.ProviderServiceConfig, "useSCP"),
				PreCommand:               xmaps.GetString(options.ProviderServiceConfig, "preCommand"),
				PostCommand:              xmaps.GetString(options.ProviderServiceConfig, "postCommand"),
				OutputFormat:             pSSH.OutputFormatType(xmaps.GetOrDefaultString(options.ProviderServiceConfig, "format", string(pSSH.OUTPUT_FORMAT_PEM))),
				OutputCertPath:           xmaps.GetString(options.ProviderServiceConfig, "certPath"),
				OutputServerCertPath:     xmaps.GetString(options.ProviderServiceConfig, "certPathForServerOnly"),
				OutputIntermediaCertPath: xmaps.GetString(options.ProviderServiceConfig, "certPathForIntermediaOnly"),
				OutputKeyPath:            xmaps.GetString(options.ProviderServiceConfig, "keyPath"),
				PfxPassword:              xmaps.GetString(options.ProviderServiceConfig, "pfxPassword"),
				JksAlias:                 xmaps.GetString(options.ProviderServiceConfig, "jksAlias"),
				JksKeypass:               xmaps.GetString(options.ProviderServiceConfig, "jksKeypass"),
				JksStorepass:             xmaps.GetString(options.ProviderServiceConfig, "jksStorepass"),
			})
			return deployer, err
		}

	case domain.DeploymentProviderTypeTencentCloudCDN, domain.DeploymentProviderTypeTencentCloudCLB, domain.DeploymentProviderTypeTencentCloudCOS, domain.DeploymentProviderTypeTencentCloudCSS, domain.DeploymentProviderTypeTencentCloudECDN, domain.DeploymentProviderTypeTencentCloudEO, domain.DeploymentProviderTypeTencentCloudGAAP, domain.DeploymentProviderTypeTencentCloudSCF, domain.DeploymentProviderTypeTencentCloudSSL, domain.DeploymentProviderTypeTencentCloudSSLDeploy, domain.DeploymentProviderTypeTencentCloudVOD, domain.DeploymentProviderTypeTencentCloudWAF:
		{
			access := domain.AccessConfigForTencentCloud{}
			if err := xmaps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeploymentProviderTypeTencentCloudCDN:
				deployer, err := pTencentCloudCDN.NewDeployer(&pTencentCloudCDN.DeployerConfig{
					SecretId:  access.SecretId,
					SecretKey: access.SecretKey,
					Domain:    xmaps.GetString(options.ProviderServiceConfig, "domain"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeTencentCloudCLB:
				deployer, err := pTencentCloudCLB.NewDeployer(&pTencentCloudCLB.DeployerConfig{
					SecretId:       access.SecretId,
					SecretKey:      access.SecretKey,
					Region:         xmaps.GetString(options.ProviderServiceConfig, "region"),
					ResourceType:   pTencentCloudCLB.ResourceType(xmaps.GetString(options.ProviderServiceConfig, "resourceType")),
					LoadbalancerId: xmaps.GetString(options.ProviderServiceConfig, "loadbalancerId"),
					ListenerId:     xmaps.GetString(options.ProviderServiceConfig, "listenerId"),
					Domain:         xmaps.GetString(options.ProviderServiceConfig, "domain"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeTencentCloudCOS:
				deployer, err := pTencentCloudCOS.NewDeployer(&pTencentCloudCOS.DeployerConfig{
					SecretId:  access.SecretId,
					SecretKey: access.SecretKey,
					Region:    xmaps.GetString(options.ProviderServiceConfig, "region"),
					Bucket:    xmaps.GetString(options.ProviderServiceConfig, "bucket"),
					Domain:    xmaps.GetString(options.ProviderServiceConfig, "domain"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeTencentCloudCSS:
				deployer, err := pTencentCloudCSS.NewDeployer(&pTencentCloudCSS.DeployerConfig{
					SecretId:  access.SecretId,
					SecretKey: access.SecretKey,
					Domain:    xmaps.GetString(options.ProviderServiceConfig, "domain"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeTencentCloudECDN:
				deployer, err := pTencentCloudECDN.NewDeployer(&pTencentCloudECDN.DeployerConfig{
					SecretId:  access.SecretId,
					SecretKey: access.SecretKey,
					Domain:    xmaps.GetString(options.ProviderServiceConfig, "domain"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeTencentCloudEO:
				deployer, err := pTencentCloudEO.NewDeployer(&pTencentCloudEO.DeployerConfig{
					SecretId:  access.SecretId,
					SecretKey: access.SecretKey,
					ZoneId:    xmaps.GetString(options.ProviderServiceConfig, "zoneId"),
					Domain:    xmaps.GetString(options.ProviderServiceConfig, "domain"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeTencentCloudGAAP:
				deployer, err := pTencentCloudGAAP.NewDeployer(&pTencentCloudGAAP.DeployerConfig{
					SecretId:     access.SecretId,
					SecretKey:    access.SecretKey,
					ResourceType: pTencentCloudGAAP.ResourceType(xmaps.GetString(options.ProviderServiceConfig, "resourceType")),
					ProxyId:      xmaps.GetString(options.ProviderServiceConfig, "proxyId"),
					ListenerId:   xmaps.GetString(options.ProviderServiceConfig, "listenerId"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeTencentCloudSCF:
				deployer, err := pTencentCloudSCF.NewDeployer(&pTencentCloudSCF.DeployerConfig{
					SecretId:  access.SecretId,
					SecretKey: access.SecretKey,
					Region:    xmaps.GetString(options.ProviderServiceConfig, "region"),
					Domain:    xmaps.GetString(options.ProviderServiceConfig, "domain"),
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
					Region:       xmaps.GetString(options.ProviderServiceConfig, "region"),
					ResourceType: xmaps.GetString(options.ProviderServiceConfig, "resourceType"),
					ResourceIds:  xslices.Filter(strings.Split(xmaps.GetString(options.ProviderServiceConfig, "resourceIds"), ";"), func(s string) bool { return s != "" }),
				})
				return deployer, err

			case domain.DeploymentProviderTypeTencentCloudVOD:
				deployer, err := pTencentCloudVOD.NewDeployer(&pTencentCloudVOD.DeployerConfig{
					SecretId:  access.SecretId,
					SecretKey: access.SecretKey,
					SubAppId:  xmaps.GetInt64(options.ProviderServiceConfig, "subAppId"),
					Domain:    xmaps.GetString(options.ProviderServiceConfig, "domain"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeTencentCloudWAF:
				deployer, err := pTencentCloudWAF.NewDeployer(&pTencentCloudWAF.DeployerConfig{
					SecretId:   access.SecretId,
					SecretKey:  access.SecretKey,
					Domain:     xmaps.GetString(options.ProviderServiceConfig, "domain"),
					DomainId:   xmaps.GetString(options.ProviderServiceConfig, "domainId"),
					InstanceId: xmaps.GetString(options.ProviderServiceConfig, "instanceId"),
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeploymentProviderTypeUCloudUCDN, domain.DeploymentProviderTypeUCloudUS3:
		{
			access := domain.AccessConfigForUCloud{}
			if err := xmaps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeploymentProviderTypeUCloudUCDN:
				deployer, err := pUCloudUCDN.NewDeployer(&pUCloudUCDN.DeployerConfig{
					PrivateKey: access.PrivateKey,
					PublicKey:  access.PublicKey,
					ProjectId:  access.ProjectId,
					DomainId:   xmaps.GetString(options.ProviderServiceConfig, "domainId"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeUCloudUS3:
				deployer, err := pUCloudUS3.NewDeployer(&pUCloudUS3.DeployerConfig{
					PrivateKey: access.PrivateKey,
					PublicKey:  access.PublicKey,
					ProjectId:  access.ProjectId,
					Region:     xmaps.GetString(options.ProviderServiceConfig, "region"),
					Bucket:     xmaps.GetString(options.ProviderServiceConfig, "bucket"),
					Domain:     xmaps.GetString(options.ProviderServiceConfig, "domain"),
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeploymentProviderTypeUniCloudWebHost:
		{
			access := domain.AccessConfigForUniCloud{}
			if err := xmaps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			deployer, err := pUniCloudWebHost.NewDeployer(&pUniCloudWebHost.DeployerConfig{
				Username:      access.Username,
				Password:      access.Password,
				SpaceProvider: xmaps.GetString(options.ProviderServiceConfig, "spaceProvider"),
				SpaceId:       xmaps.GetString(options.ProviderServiceConfig, "spaceId"),
				Domain:        xmaps.GetString(options.ProviderServiceConfig, "domain"),
			})
			return deployer, err
		}

	case domain.DeploymentProviderTypeUpyunCDN, domain.DeploymentProviderTypeUpyunFile:
		{
			access := domain.AccessConfigForUpyun{}
			if err := xmaps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeploymentProviderTypeUpyunCDN, domain.DeploymentProviderTypeUpyunFile:
				deployer, err := pUpyunCDN.NewDeployer(&pUpyunCDN.DeployerConfig{
					Username: access.Username,
					Password: access.Password,
					Domain:   xmaps.GetString(options.ProviderServiceConfig, "domain"),
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeploymentProviderTypeVolcEngineALB, domain.DeploymentProviderTypeVolcEngineCDN, domain.DeploymentProviderTypeVolcEngineCertCenter, domain.DeploymentProviderTypeVolcEngineCLB, domain.DeploymentProviderTypeVolcEngineDCDN, domain.DeploymentProviderTypeVolcEngineImageX, domain.DeploymentProviderTypeVolcEngineLive, domain.DeploymentProviderTypeVolcEngineTOS:
		{
			access := domain.AccessConfigForVolcEngine{}
			if err := xmaps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeploymentProviderTypeVolcEngineALB:
				deployer, err := pVolcEngineALB.NewDeployer(&pVolcEngineALB.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.SecretAccessKey,
					Region:          xmaps.GetString(options.ProviderServiceConfig, "region"),
					ResourceType:    pVolcEngineALB.ResourceType(xmaps.GetString(options.ProviderServiceConfig, "resourceType")),
					LoadbalancerId:  xmaps.GetString(options.ProviderServiceConfig, "loadbalancerId"),
					ListenerId:      xmaps.GetString(options.ProviderServiceConfig, "listenerId"),
					Domain:          xmaps.GetString(options.ProviderServiceConfig, "domain"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeVolcEngineCDN:
				deployer, err := pVolcEngineCDN.NewDeployer(&pVolcEngineCDN.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.SecretAccessKey,
					Domain:          xmaps.GetString(options.ProviderServiceConfig, "domain"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeVolcEngineCertCenter:
				deployer, err := pVolcEngineCertCenter.NewDeployer(&pVolcEngineCertCenter.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.SecretAccessKey,
					Region:          xmaps.GetString(options.ProviderServiceConfig, "region"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeVolcEngineCLB:
				deployer, err := pVolcEngineCLB.NewDeployer(&pVolcEngineCLB.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.SecretAccessKey,
					Region:          xmaps.GetString(options.ProviderServiceConfig, "region"),
					ResourceType:    pVolcEngineCLB.ResourceType(xmaps.GetString(options.ProviderServiceConfig, "resourceType")),
					LoadbalancerId:  xmaps.GetString(options.ProviderServiceConfig, "loadbalancerId"),
					ListenerId:      xmaps.GetString(options.ProviderServiceConfig, "listenerId"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeVolcEngineDCDN:
				deployer, err := pVolcEngineDCDN.NewDeployer(&pVolcEngineDCDN.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.SecretAccessKey,
					Domain:          xmaps.GetString(options.ProviderServiceConfig, "domain"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeVolcEngineImageX:
				deployer, err := pVolcEngineImageX.NewDeployer(&pVolcEngineImageX.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.SecretAccessKey,
					Region:          xmaps.GetString(options.ProviderServiceConfig, "region"),
					ServiceId:       xmaps.GetString(options.ProviderServiceConfig, "serviceId"),
					Domain:          xmaps.GetString(options.ProviderServiceConfig, "domain"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeVolcEngineLive:
				deployer, err := pVolcEngineLive.NewDeployer(&pVolcEngineLive.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.SecretAccessKey,
					Domain:          xmaps.GetString(options.ProviderServiceConfig, "domain"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeVolcEngineTOS:
				deployer, err := pVolcEngineTOS.NewDeployer(&pVolcEngineTOS.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.SecretAccessKey,
					Region:          xmaps.GetString(options.ProviderServiceConfig, "region"),
					Bucket:          xmaps.GetString(options.ProviderServiceConfig, "bucket"),
					Domain:          xmaps.GetString(options.ProviderServiceConfig, "domain"),
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeploymentProviderTypeWangsuCDN, domain.DeploymentProviderTypeWangsuCDNPro, domain.DeploymentProviderTypeWangsuCertificate:
		{
			access := domain.AccessConfigForWangsu{}
			if err := xmaps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeploymentProviderTypeWangsuCDN:
				deployer, err := pWangsuCDN.NewDeployer(&pWangsuCDN.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Domains:         xslices.Filter(strings.Split(xmaps.GetString(options.ProviderServiceConfig, "domains"), ";"), func(s string) bool { return s != "" }),
				})
				return deployer, err

			case domain.DeploymentProviderTypeWangsuCDNPro:
				deployer, err := pWangsuCDNPro.NewDeployer(&pWangsuCDNPro.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					ApiKey:          access.ApiKey,
					Environment:     xmaps.GetOrDefaultString(options.ProviderServiceConfig, "environment", "production"),
					Domain:          xmaps.GetString(options.ProviderServiceConfig, "domain"),
					CertificateId:   xmaps.GetString(options.ProviderServiceConfig, "certificateId"),
					WebhookId:       xmaps.GetString(options.ProviderServiceConfig, "webhookId"),
				})
				return deployer, err

			case domain.DeploymentProviderTypeWangsuCertificate:
				deployer, err := pWangsuCertificate.NewDeployer(&pWangsuCertificate.DeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					CertificateId:   xmaps.GetString(options.ProviderServiceConfig, "certificateId"),
				})
				return deployer, err

			default:
				break
			}
		}

	case domain.DeploymentProviderTypeWebhook:
		{
			access := domain.AccessConfigForWebhook{}
			if err := xmaps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			mergedHeaders := make(map[string]string)
			if defaultHeadersString := access.HeadersString; defaultHeadersString != "" {
				h, err := xhttp.ParseHeaders(defaultHeadersString)
				if err != nil {
					return nil, fmt.Errorf("failed to parse webhook headers: %w", err)
				}
				for key := range h {
					mergedHeaders[http.CanonicalHeaderKey(key)] = h.Get(key)
				}
			}
			if extendedHeadersString := xmaps.GetString(options.ProviderServiceConfig, "headers"); extendedHeadersString != "" {
				h, err := xhttp.ParseHeaders(extendedHeadersString)
				if err != nil {
					return nil, fmt.Errorf("failed to parse webhook headers: %w", err)
				}
				for key := range h {
					mergedHeaders[http.CanonicalHeaderKey(key)] = h.Get(key)
				}
			}

			deployer, err := pWebhook.NewDeployer(&pWebhook.DeployerConfig{
				WebhookUrl:               access.Url,
				WebhookData:              xmaps.GetOrDefaultString(options.ProviderServiceConfig, "webhookData", access.DefaultDataForDeployment),
				Method:                   access.Method,
				Headers:                  mergedHeaders,
				AllowInsecureConnections: access.AllowInsecureConnections,
			})
			return deployer, err
		}
	}

	return nil, fmt.Errorf("unsupported deployer provider '%s'", string(options.Provider))
}
