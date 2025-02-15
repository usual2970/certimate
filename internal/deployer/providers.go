package deployer

import (
	"fmt"
	"strings"

	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	pAliyunALB "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/aliyun-alb"
	pAliyunCASDeploy "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/aliyun-cas-deploy"
	pAliyunCDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/aliyun-cdn"
	pAliyunCLB "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/aliyun-clb"
	pAliyunDCDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/aliyun-dcdn"
	pAliyunESA "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/aliyun-esa"
	pAliyunLive "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/aliyun-live"
	pAliyunNLB "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/aliyun-nlb"
	pAliyunOSS "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/aliyun-oss"
	pAliyunWAF "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/aliyun-waf"
	pAWSCloudFront "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/aws-cloudfront"
	pBaiduCloudCDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/baiducloud-cdn"
	pBaotaPanelConsole "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/baotapanel-console"
	pBaotaPanelSite "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/baotapanel-site"
	pBytePlusCDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/byteplus-cdn"
	pDogeCDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/dogecloud-cdn"
	pEdgioApplications "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/edgio-applications"
	pHuaweiCloudCDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/huaweicloud-cdn"
	pHuaweiCloudELB "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/huaweicloud-elb"
	pK8sSecret "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/k8s-secret"
	pLocal "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/local"
	pQiniuCDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/qiniu-cdn"
	pQiniuPili "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/qiniu-pili"
	providerSSH "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/ssh"
	pTencentCloudCDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/tencentcloud-cdn"
	pTencentCloudCLB "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/tencentcloud-clb"
	pTencentCloudCOS "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/tencentcloud-cos"
	pTencentCloudCSS "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/tencentcloud-css"
	pTencentCloudECDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/tencentcloud-ecdn"
	pTencentCloudEO "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/tencentcloud-eo"
	pTencentCloudSSLDeploy "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/tencentcloud-ssl-deploy"
	pUCloudUCDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/ucloud-ucdn"
	pUCloudUS3 "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/ucloud-us3"
	pVolcEngineCDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/volcengine-cdn"
	pVolcEngineCLB "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/volcengine-clb"
	pVolcEngineDCDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/volcengine-dcdn"
	pVolcEngineImageX "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/volcengine-imagex"
	pVolcEngineLive "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/volcengine-live"
	pVolcEngineTOS "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/volcengine-tos"
	pWebhook "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/webhook"
	"github.com/usual2970/certimate/internal/pkg/core/logger"
	"github.com/usual2970/certimate/internal/pkg/utils/maps"
	"github.com/usual2970/certimate/internal/pkg/utils/slices"
)

func createDeployer(options *deployerOptions) (deployer.Deployer, logger.Logger, error) {
	logger := logger.NewDefaultLogger()

	/*
	  注意：如果追加新的常量值，请保持以 ASCII 排序。
	  NOTICE: If you add new constant, please keep ASCII order.
	*/
	switch options.Provider {
	case domain.DeployProviderTypeAliyunALB, domain.DeployProviderTypeAliyunCASDeploy, domain.DeployProviderTypeAliyunCDN, domain.DeployProviderTypeAliyunCLB, domain.DeployProviderTypeAliyunDCDN, domain.DeployProviderTypeAliyunESA, domain.DeployProviderTypeAliyunLive, domain.DeployProviderTypeAliyunNLB, domain.DeployProviderTypeAliyunOSS, domain.DeployProviderTypeAliyunWAF:
		{
			access := domain.AccessConfigForAliyun{}
			if err := maps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeployProviderTypeAliyunALB:
				deployer, err := pAliyunALB.NewWithLogger(&pAliyunALB.AliyunALBDeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Region:          maps.GetValueAsString(options.ProviderDeployConfig, "region"),
					ResourceType:    pAliyunALB.DeployResourceType(maps.GetValueAsString(options.ProviderDeployConfig, "resourceType")),
					LoadbalancerId:  maps.GetValueAsString(options.ProviderDeployConfig, "loadbalancerId"),
					ListenerId:      maps.GetValueAsString(options.ProviderDeployConfig, "listenerId"),
					Domain:          maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
				}, logger)
				return deployer, logger, err

			case domain.DeployProviderTypeAliyunCASDeploy:
				deployer, err := pAliyunCASDeploy.NewWithLogger(&pAliyunCASDeploy.AliyunCASDeployDeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Region:          maps.GetValueAsString(options.ProviderDeployConfig, "region"),
					ResourceIds:     slices.Filter(strings.Split(maps.GetValueAsString(options.ProviderDeployConfig, "resourceIds"), ";"), func(s string) bool { return s != "" }),
					ContactIds:      slices.Filter(strings.Split(maps.GetValueAsString(options.ProviderDeployConfig, "contactIds"), ";"), func(s string) bool { return s != "" }),
				}, logger)
				return deployer, logger, err

			case domain.DeployProviderTypeAliyunCDN:
				deployer, err := pAliyunCDN.NewWithLogger(&pAliyunCDN.AliyunCDNDeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Domain:          maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
				}, logger)
				return deployer, logger, err

			case domain.DeployProviderTypeAliyunCLB:
				deployer, err := pAliyunCLB.NewWithLogger(&pAliyunCLB.AliyunCLBDeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Region:          maps.GetValueAsString(options.ProviderDeployConfig, "region"),
					ResourceType:    pAliyunCLB.DeployResourceType(maps.GetValueAsString(options.ProviderDeployConfig, "resourceType")),
					LoadbalancerId:  maps.GetValueAsString(options.ProviderDeployConfig, "loadbalancerId"),
					ListenerPort:    maps.GetValueOrDefaultAsInt32(options.ProviderDeployConfig, "listenerPort", 443),
					Domain:          maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
				}, logger)
				return deployer, logger, err

			case domain.DeployProviderTypeAliyunDCDN:
				deployer, err := pAliyunDCDN.NewWithLogger(&pAliyunDCDN.AliyunDCDNDeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Domain:          maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
				}, logger)
				return deployer, logger, err

			case domain.DeployProviderTypeAliyunESA:
				deployer, err := pAliyunESA.NewWithLogger(&pAliyunESA.AliyunESADeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Region:          maps.GetValueAsString(options.ProviderDeployConfig, "region"),
					SiteId:          maps.GetValueAsInt64(options.ProviderDeployConfig, "siteId"),
				}, logger)
				return deployer, logger, err

			case domain.DeployProviderTypeAliyunLive:
				deployer, err := pAliyunLive.NewWithLogger(&pAliyunLive.AliyunLiveDeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Region:          maps.GetValueAsString(options.ProviderDeployConfig, "region"),
					Domain:          maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
				}, logger)
				return deployer, logger, err

			case domain.DeployProviderTypeAliyunNLB:
				deployer, err := pAliyunNLB.NewWithLogger(&pAliyunNLB.AliyunNLBDeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Region:          maps.GetValueAsString(options.ProviderDeployConfig, "region"),
					ResourceType:    pAliyunNLB.DeployResourceType(maps.GetValueAsString(options.ProviderDeployConfig, "resourceType")),
					LoadbalancerId:  maps.GetValueAsString(options.ProviderDeployConfig, "loadbalancerId"),
					ListenerId:      maps.GetValueAsString(options.ProviderDeployConfig, "listenerId"),
				}, logger)
				return deployer, logger, err

			case domain.DeployProviderTypeAliyunOSS:
				deployer, err := pAliyunOSS.NewWithLogger(&pAliyunOSS.AliyunOSSDeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Region:          maps.GetValueAsString(options.ProviderDeployConfig, "region"),
					Bucket:          maps.GetValueAsString(options.ProviderDeployConfig, "bucket"),
					Domain:          maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
				}, logger)
				return deployer, logger, err

			case domain.DeployProviderTypeAliyunWAF:
				deployer, err := pAliyunWAF.NewWithLogger(&pAliyunWAF.AliyunWAFDeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Region:          maps.GetValueAsString(options.ProviderDeployConfig, "region"),
					InstanceId:      maps.GetValueAsString(options.ProviderDeployConfig, "instanceId"),
				}, logger)
				return deployer, logger, err

			default:
				break
			}
		}

	case domain.DeployProviderTypeAWSCloudFront:
		{
			access := domain.AccessConfigForAWS{}
			if err := maps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeployProviderTypeAWSCloudFront:
				deployer, err := pAWSCloudFront.NewWithLogger(&pAWSCloudFront.AWSCloudFrontDeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					SecretAccessKey: access.SecretAccessKey,
					Region:          maps.GetValueAsString(options.ProviderDeployConfig, "region"),
					DistributionId:  maps.GetValueAsString(options.ProviderDeployConfig, "distributionId"),
				}, logger)
				return deployer, logger, err

			default:
				break
			}
		}

	case domain.DeployProviderTypeBaiduCloudCDN:
		{
			access := domain.AccessConfigForBaiduCloud{}
			if err := maps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeployProviderTypeBaiduCloudCDN:
				deployer, err := pBaiduCloudCDN.NewWithLogger(&pBaiduCloudCDN.BaiduCloudCDNDeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					SecretAccessKey: access.SecretAccessKey,
					Domain:          maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
				}, logger)
				return deployer, logger, err

			default:
				break
			}
		}

	case domain.DeployProviderTypeBaotaPanelConsole, domain.DeployProviderTypeBaotaPanelSite:
		{
			access := domain.AccessConfigForBaotaPanel{}
			if err := maps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeployProviderTypeBaotaPanelConsole:
				deployer, err := pBaotaPanelConsole.NewWithLogger(&pBaotaPanelConsole.BaotaPanelConsoleDeployerConfig{
					ApiUrl:      access.ApiUrl,
					ApiKey:      access.ApiKey,
					AutoRestart: maps.GetValueAsBool(options.ProviderDeployConfig, "autoRestart"),
				}, logger)
				return deployer, logger, err

			case domain.DeployProviderTypeBaotaPanelSite:
				deployer, err := pBaotaPanelSite.NewWithLogger(&pBaotaPanelSite.BaotaPanelSiteDeployerConfig{
					ApiUrl:   access.ApiUrl,
					ApiKey:   access.ApiKey,
					SiteName: maps.GetValueAsString(options.ProviderDeployConfig, "siteName"),
				}, logger)
				return deployer, logger, err

			default:
				break
			}
		}

	case domain.DeployProviderTypeBytePlusCDN:
		{
			access := domain.AccessConfigForBytePlus{}
			if err := maps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeployProviderTypeBytePlusCDN:
				deployer, err := pBytePlusCDN.NewWithLogger(&pBytePlusCDN.BytePlusCDNDeployerConfig{
					AccessKey: access.AccessKey,
					SecretKey: access.SecretKey,
					Domain:    maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
				}, logger)
				return deployer, logger, err

			default:
				break
			}
		}

	case domain.DeployProviderTypeDogeCloudCDN:
		{
			access := domain.AccessConfigForDogeCloud{}
			if err := maps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			deployer, err := pDogeCDN.NewWithLogger(&pDogeCDN.DogeCloudCDNDeployerConfig{
				AccessKey: access.AccessKey,
				SecretKey: access.SecretKey,
				Domain:    maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
			}, logger)
			return deployer, logger, err
		}

	case domain.DeployProviderTypeEdgioApplications:
		{
			access := domain.AccessConfigForEdgio{}
			if err := maps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			deployer, err := pEdgioApplications.NewWithLogger(&pEdgioApplications.EdgioApplicationsDeployerConfig{
				ClientId:      access.ClientId,
				ClientSecret:  access.ClientSecret,
				EnvironmentId: maps.GetValueAsString(options.ProviderDeployConfig, "environmentId"),
			}, logger)
			return deployer, logger, err
		}

	case domain.DeployProviderTypeHuaweiCloudCDN, domain.DeployProviderTypeHuaweiCloudELB:
		{
			access := domain.AccessConfigForHuaweiCloud{}
			if err := maps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeployProviderTypeHuaweiCloudCDN:
				deployer, err := pHuaweiCloudCDN.NewWithLogger(&pHuaweiCloudCDN.HuaweiCloudCDNDeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					SecretAccessKey: access.SecretAccessKey,
					Region:          maps.GetValueAsString(options.ProviderDeployConfig, "region"),
					Domain:          maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
				}, logger)
				return deployer, logger, err

			case domain.DeployProviderTypeHuaweiCloudELB:
				deployer, err := pHuaweiCloudELB.NewWithLogger(&pHuaweiCloudELB.HuaweiCloudELBDeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					SecretAccessKey: access.SecretAccessKey,
					Region:          maps.GetValueAsString(options.ProviderDeployConfig, "region"),
					ResourceType:    pHuaweiCloudELB.DeployResourceType(maps.GetValueAsString(options.ProviderDeployConfig, "resourceType")),
					CertificateId:   maps.GetValueAsString(options.ProviderDeployConfig, "certificateId"),
					LoadbalancerId:  maps.GetValueAsString(options.ProviderDeployConfig, "loadbalancerId"),
					ListenerId:      maps.GetValueAsString(options.ProviderDeployConfig, "listenerId"),
				}, logger)
				return deployer, logger, err

			default:
				break
			}
		}

	case domain.DeployProviderTypeLocal:
		{
			deployer, err := pLocal.NewWithLogger(&pLocal.LocalDeployerConfig{
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
			}, logger)
			return deployer, logger, err
		}

	case domain.DeployProviderTypeKubernetesSecret:
		{
			access := domain.AccessConfigForKubernetes{}
			if err := maps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			deployer, err := pK8sSecret.NewWithLogger(&pK8sSecret.K8sSecretDeployerConfig{
				KubeConfig:          access.KubeConfig,
				Namespace:           maps.GetValueOrDefaultAsString(options.ProviderDeployConfig, "namespace", "default"),
				SecretName:          maps.GetValueAsString(options.ProviderDeployConfig, "secretName"),
				SecretType:          maps.GetValueOrDefaultAsString(options.ProviderDeployConfig, "secretType", "kubernetes.io/tls"),
				SecretDataKeyForCrt: maps.GetValueOrDefaultAsString(options.ProviderDeployConfig, "secretDataKeyForCrt", "tls.crt"),
				SecretDataKeyForKey: maps.GetValueOrDefaultAsString(options.ProviderDeployConfig, "secretDataKeyForKey", "tls.key"),
			}, logger)
			return deployer, logger, err
		}

	case domain.DeployProviderTypeQiniuCDN, domain.DeployProviderTypeQiniuPili:
		{
			access := domain.AccessConfigForQiniu{}
			if err := maps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeployProviderTypeQiniuCDN:
				deployer, err := pQiniuCDN.NewWithLogger(&pQiniuCDN.QiniuCDNDeployerConfig{
					AccessKey: access.AccessKey,
					SecretKey: access.SecretKey,
					Domain:    maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
				}, logger)
				return deployer, logger, err

			case domain.DeployProviderTypeQiniuPili:
				deployer, err := pQiniuPili.NewWithLogger(&pQiniuPili.QiniuPiliDeployerConfig{
					AccessKey: access.AccessKey,
					SecretKey: access.SecretKey,
					Hub:       maps.GetValueAsString(options.ProviderDeployConfig, "hub"),
					Domain:    maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
				}, logger)
				return deployer, logger, err

			default:
				break
			}
		}

	case domain.DeployProviderTypeSSH:
		{
			access := domain.AccessConfigForSSH{}
			if err := maps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			deployer, err := providerSSH.NewWithLogger(&providerSSH.SshDeployerConfig{
				SshHost:          access.Host,
				SshPort:          access.Port,
				SshUsername:      access.Username,
				SshPassword:      access.Password,
				SshKey:           access.Key,
				SshKeyPassphrase: access.KeyPassphrase,
				UseSCP:           maps.GetValueAsBool(options.ProviderDeployConfig, "useSCP"),
				PreCommand:       maps.GetValueAsString(options.ProviderDeployConfig, "preCommand"),
				PostCommand:      maps.GetValueAsString(options.ProviderDeployConfig, "postCommand"),
				OutputFormat:     providerSSH.OutputFormatType(maps.GetValueOrDefaultAsString(options.ProviderDeployConfig, "format", string(providerSSH.OUTPUT_FORMAT_PEM))),
				OutputCertPath:   maps.GetValueAsString(options.ProviderDeployConfig, "certPath"),
				OutputKeyPath:    maps.GetValueAsString(options.ProviderDeployConfig, "keyPath"),
				PfxPassword:      maps.GetValueAsString(options.ProviderDeployConfig, "pfxPassword"),
				JksAlias:         maps.GetValueAsString(options.ProviderDeployConfig, "jksAlias"),
				JksKeypass:       maps.GetValueAsString(options.ProviderDeployConfig, "jksKeypass"),
				JksStorepass:     maps.GetValueAsString(options.ProviderDeployConfig, "jksStorepass"),
			}, logger)
			return deployer, logger, err
		}

	case domain.DeployProviderTypeTencentCloudCDN, domain.DeployProviderTypeTencentCloudCLB, domain.DeployProviderTypeTencentCloudCOS, domain.DeployProviderTypeTencentCloudCSS, domain.DeployProviderTypeTencentCloudECDN, domain.DeployProviderTypeTencentCloudEO, domain.DeployProviderTypeTencentCloudSSLDeploy:
		{
			access := domain.AccessConfigForTencentCloud{}
			if err := maps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeployProviderTypeTencentCloudCDN:
				deployer, err := pTencentCloudCDN.NewWithLogger(&pTencentCloudCDN.TencentCloudCDNDeployerConfig{
					SecretId:  access.SecretId,
					SecretKey: access.SecretKey,
					Domain:    maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
				}, logger)
				return deployer, logger, err

			case domain.DeployProviderTypeTencentCloudCLB:
				deployer, err := pTencentCloudCLB.NewWithLogger(&pTencentCloudCLB.TencentCloudCLBDeployerConfig{
					SecretId:       access.SecretId,
					SecretKey:      access.SecretKey,
					Region:         maps.GetValueAsString(options.ProviderDeployConfig, "region"),
					ResourceType:   pTencentCloudCLB.DeployResourceType(maps.GetValueAsString(options.ProviderDeployConfig, "resourceType")),
					LoadbalancerId: maps.GetValueAsString(options.ProviderDeployConfig, "loadbalancerId"),
					ListenerId:     maps.GetValueAsString(options.ProviderDeployConfig, "listenerId"),
					Domain:         maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
				}, logger)
				return deployer, logger, err

			case domain.DeployProviderTypeTencentCloudCOS:
				deployer, err := pTencentCloudCOS.NewWithLogger(&pTencentCloudCOS.TencentCloudCOSDeployerConfig{
					SecretId:  access.SecretId,
					SecretKey: access.SecretKey,
					Region:    maps.GetValueAsString(options.ProviderDeployConfig, "region"),
					Bucket:    maps.GetValueAsString(options.ProviderDeployConfig, "bucket"),
					Domain:    maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
				}, logger)
				return deployer, logger, err

			case domain.DeployProviderTypeTencentCloudCSS:
				deployer, err := pTencentCloudCSS.NewWithLogger(&pTencentCloudCSS.TencentCloudCSSDeployerConfig{
					SecretId:  access.SecretId,
					SecretKey: access.SecretKey,
					Domain:    maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
				}, logger)
				return deployer, logger, err

			case domain.DeployProviderTypeTencentCloudECDN:
				deployer, err := pTencentCloudECDN.NewWithLogger(&pTencentCloudECDN.TencentCloudECDNDeployerConfig{
					SecretId:  access.SecretId,
					SecretKey: access.SecretKey,
					Domain:    maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
				}, logger)
				return deployer, logger, err

			case domain.DeployProviderTypeTencentCloudEO:
				deployer, err := pTencentCloudEO.NewWithLogger(&pTencentCloudEO.TencentCloudEODeployerConfig{
					SecretId:  access.SecretId,
					SecretKey: access.SecretKey,
					ZoneId:    maps.GetValueAsString(options.ProviderDeployConfig, "zoneId"),
					Domain:    maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
				}, logger)
				return deployer, logger, err

			case domain.DeployProviderTypeTencentCloudSSLDeploy:
				deployer, err := pTencentCloudSSLDeploy.NewWithLogger(&pTencentCloudSSLDeploy.TencentCloudSSLDeployDeployerConfig{
					SecretId:     access.SecretId,
					SecretKey:    access.SecretKey,
					Region:       maps.GetValueAsString(options.ProviderDeployConfig, "region"),
					ResourceType: maps.GetValueAsString(options.ProviderDeployConfig, "resourceType"),
					ResourceIds:  slices.Filter(strings.Split(maps.GetValueAsString(options.ProviderDeployConfig, "resourceIds"), ";"), func(s string) bool { return s != "" }),
				}, logger)
				return deployer, logger, err

			default:
				break
			}
		}

	case domain.DeployProviderTypeUCloudUCDN, domain.DeployProviderTypeUCloudUS3:
		{
			access := domain.AccessConfigForUCloud{}
			if err := maps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeployProviderTypeUCloudUCDN:
				deployer, err := pUCloudUCDN.NewWithLogger(&pUCloudUCDN.UCloudUCDNDeployerConfig{
					PrivateKey: access.PrivateKey,
					PublicKey:  access.PublicKey,
					ProjectId:  access.ProjectId,
					DomainId:   maps.GetValueAsString(options.ProviderDeployConfig, "domainId"),
				}, logger)
				return deployer, logger, err

			case domain.DeployProviderTypeUCloudUS3:
				deployer, err := pUCloudUS3.NewWithLogger(&pUCloudUS3.UCloudUS3DeployerConfig{
					PrivateKey: access.PrivateKey,
					PublicKey:  access.PublicKey,
					ProjectId:  access.ProjectId,
					Region:     maps.GetValueAsString(options.ProviderDeployConfig, "region"),
					Bucket:     maps.GetValueAsString(options.ProviderDeployConfig, "bucket"),
					Domain:     maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
				}, logger)
				return deployer, logger, err

			default:
				break
			}
		}

	case domain.DeployProviderTypeVolcEngineCDN, domain.DeployProviderTypeVolcEngineCLB, domain.DeployProviderTypeVolcEngineDCDN, domain.DeployProviderTypeVolcEngineImageX, domain.DeployProviderTypeVolcEngineLive, domain.DeployProviderTypeVolcEngineTOS:
		{
			access := domain.AccessConfigForVolcEngine{}
			if err := maps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeployProviderTypeVolcEngineCDN:
				deployer, err := pVolcEngineCDN.NewWithLogger(&pVolcEngineCDN.VolcEngineCDNDeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.SecretAccessKey,
					Domain:          maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
				}, logger)
				return deployer, logger, err

			case domain.DeployProviderTypeVolcEngineCLB:
				deployer, err := pVolcEngineCLB.NewWithLogger(&pVolcEngineCLB.VolcEngineCLBDeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.SecretAccessKey,
					Region:          maps.GetValueAsString(options.ProviderDeployConfig, "region"),
					ResourceType:    pVolcEngineCLB.DeployResourceType(maps.GetValueAsString(options.ProviderDeployConfig, "resourceType")),
					ListenerId:      maps.GetValueAsString(options.ProviderDeployConfig, "listenerId"),
				}, logger)
				return deployer, logger, err

			case domain.DeployProviderTypeVolcEngineDCDN:
				deployer, err := pVolcEngineDCDN.NewWithLogger(&pVolcEngineDCDN.VolcEngineDCDNDeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.SecretAccessKey,
					Domain:          maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
				}, logger)
				return deployer, logger, err

			case domain.DeployProviderTypeVolcEngineImageX:
				deployer, err := pVolcEngineImageX.NewWithLogger(&pVolcEngineImageX.VolcEngineImageXDeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.SecretAccessKey,
					Region:          maps.GetValueAsString(options.ProviderDeployConfig, "region"),
					ServiceId:       maps.GetValueAsString(options.ProviderDeployConfig, "serviceId"),
					Domain:          maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
				}, logger)
				return deployer, logger, err

			case domain.DeployProviderTypeVolcEngineLive:
				deployer, err := pVolcEngineLive.NewWithLogger(&pVolcEngineLive.VolcEngineLiveDeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.SecretAccessKey,
					Domain:          maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
				}, logger)
				return deployer, logger, err

			case domain.DeployProviderTypeVolcEngineTOS:
				deployer, err := pVolcEngineTOS.NewWithLogger(&pVolcEngineTOS.VolcEngineTOSDeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.SecretAccessKey,
					Region:          maps.GetValueAsString(options.ProviderDeployConfig, "region"),
					Bucket:          maps.GetValueAsString(options.ProviderDeployConfig, "bucket"),
					Domain:          maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
				}, logger)
				return deployer, logger, err

			default:
				break
			}
		}

	case domain.DeployProviderTypeWebhook:
		{
			access := domain.AccessConfigForWebhook{}
			if err := maps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			deployer, err := pWebhook.NewWithLogger(&pWebhook.WebhookDeployerConfig{
				WebhookUrl:  access.Url,
				WebhookData: maps.GetValueAsString(options.ProviderDeployConfig, "webhookData"),
			}, logger)
			return deployer, logger, err
		}
	}

	return nil, nil, fmt.Errorf("unsupported deployer provider: %s", string(options.Provider))
}
