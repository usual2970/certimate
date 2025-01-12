package deployer

import (
	"fmt"
	"strconv"

	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	providerAliyunALB "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/aliyun-alb"
	providerAliyunCDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/aliyun-cdn"
	providerAliyunCLB "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/aliyun-clb"
	providerAliyunDCDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/aliyun-dcdn"
	providerAliyunNLB "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/aliyun-nlb"
	providerAliyunOSS "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/aliyun-oss"
	providerBaiduCloudCDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/baiducloud-cdn"
	providerBytePlusCDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/byteplus-cdn"
	providerDogeCDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/dogecloud-cdn"
	providerHuaweiCloudCDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/huaweicloud-cdn"
	providerHuaweiCloudELB "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/huaweicloud-elb"
	providerK8sSecret "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/k8s-secret"
	providerLocal "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/local"
	providerQiniuCDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/qiniu-cdn"
	providerSSH "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/ssh"
	providerTencentCloudCDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/tencentcloud-cdn"
	providerTencentCloudCLB "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/tencentcloud-clb"
	providerTencentCloudCOD "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/tencentcloud-cos"
	providerTencentCloudECDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/tencentcloud-ecdn"
	providerTencentCloudEO "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/tencentcloud-eo"
	providerVolcEngineCDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/volcengine-cdn"
	providerVolcEngineDCDN "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/volcengine-dcdn"
	providerVolcEngineLive "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/volcengine-live"
	providerWebhook "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/webhook"
	"github.com/usual2970/certimate/internal/pkg/core/logger"
	"github.com/usual2970/certimate/internal/pkg/utils/maps"
)

func createDeployer(options *deployerOptions) (deployer.Deployer, logger.Logger, error) {
	logger := logger.NewDefaultLogger()

	/*
	  注意：如果追加新的常量值，请保持以 ASCII 排序。
	  NOTICE: If you add new constant, please keep ASCII order.
	*/
	switch options.Provider {
	case domain.DeployProviderTypeAliyunALB, domain.DeployProviderTypeAliyunCDN, domain.DeployProviderTypeAliyunCLB, domain.DeployProviderTypeAliyunDCDN, domain.DeployProviderTypeAliyunNLB, domain.DeployProviderTypeAliyunOSS:
		{
			access := domain.AccessConfigForAliyun{}
			if err := maps.Decode(options.ProviderAccessConfig, &access); err != nil {
				return nil, nil, fmt.Errorf("failed to decode provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeployProviderTypeAliyunALB:
				deployer, err := providerAliyunALB.NewWithLogger(&providerAliyunALB.AliyunALBDeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Region:          maps.GetValueAsString(options.ProviderDeployConfig, "region"),
					ResourceType:    providerAliyunALB.DeployResourceType(maps.GetValueAsString(options.ProviderDeployConfig, "resourceType")),
					LoadbalancerId:  maps.GetValueAsString(options.ProviderDeployConfig, "loadbalancerId"),
					ListenerId:      maps.GetValueAsString(options.ProviderDeployConfig, "listenerId"),
				}, logger)
				return deployer, logger, err

			case domain.DeployProviderTypeAliyunCDN:
				deployer, err := providerAliyunCDN.NewWithLogger(&providerAliyunCDN.AliyunCDNDeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Domain:          maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
				}, logger)
				return deployer, logger, err

			case domain.DeployProviderTypeAliyunCLB:
				deployer, err := providerAliyunCLB.NewWithLogger(&providerAliyunCLB.AliyunCLBDeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Region:          maps.GetValueAsString(options.ProviderDeployConfig, "region"),
					ResourceType:    providerAliyunCLB.DeployResourceType(maps.GetValueAsString(options.ProviderDeployConfig, "resourceType")),
					LoadbalancerId:  maps.GetValueAsString(options.ProviderDeployConfig, "loadbalancerId"),
					ListenerPort:    maps.GetValueAsInt32(options.ProviderDeployConfig, "listenerPort"),
				}, logger)
				return deployer, logger, err

			case domain.DeployProviderTypeAliyunDCDN:
				deployer, err := providerAliyunDCDN.NewWithLogger(&providerAliyunDCDN.AliyunDCDNDeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Domain:          maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
				}, logger)
				return deployer, logger, err

			case domain.DeployProviderTypeAliyunNLB:
				deployer, err := providerAliyunNLB.NewWithLogger(&providerAliyunNLB.AliyunNLBDeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Region:          maps.GetValueAsString(options.ProviderDeployConfig, "region"),
					ResourceType:    providerAliyunNLB.DeployResourceType(maps.GetValueAsString(options.ProviderDeployConfig, "resourceType")),
					LoadbalancerId:  maps.GetValueAsString(options.ProviderDeployConfig, "loadbalancerId"),
					ListenerId:      maps.GetValueAsString(options.ProviderDeployConfig, "listenerId"),
				}, logger)
				return deployer, logger, err

			case domain.DeployProviderTypeAliyunOSS:
				deployer, err := providerAliyunOSS.NewWithLogger(&providerAliyunOSS.AliyunOSSDeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Region:          maps.GetValueAsString(options.ProviderDeployConfig, "region"),
					Bucket:          maps.GetValueAsString(options.ProviderDeployConfig, "bucket"),
					Domain:          maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
				}, logger)
				return deployer, logger, err

			default:
				break
			}
		}

	case domain.DeployProviderTypeBaiduCloudCDN:
		{
			access := domain.AccessConfigForBaiduCloud{}
			if err := maps.Decode(options.ProviderAccessConfig, &access); err != nil {
				return nil, nil, fmt.Errorf("failed to decode provider access config: %w", err)
			}

			deployer, err := providerBaiduCloudCDN.NewWithLogger(&providerBaiduCloudCDN.BaiduCloudCDNDeployerConfig{
				AccessKeyId:     access.AccessKeyId,
				SecretAccessKey: access.SecretAccessKey,
				Domain:          maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
			}, logger)
			return deployer, logger, err
		}

	case domain.DeployProviderTypeBytePlusCDN:
		{
			access := domain.AccessConfigForBytePlus{}
			if err := maps.Decode(options.ProviderAccessConfig, &access); err != nil {
				return nil, nil, fmt.Errorf("failed to decode provider access config: %w", err)
			}

			deployer, err := providerBytePlusCDN.NewWithLogger(&providerBytePlusCDN.BytePlusCDNDeployerConfig{
				AccessKey: access.AccessKey,
				SecretKey: access.SecretKey,
				Domain:    maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
			}, logger)
			return deployer, logger, err
		}

	case domain.DeployProviderTypeDogeCloudCDN:
		{
			access := domain.AccessConfigForDogeCloud{}
			if err := maps.Decode(options.ProviderAccessConfig, &access); err != nil {
				return nil, nil, fmt.Errorf("failed to decode provider access config: %w", err)
			}

			deployer, err := providerDogeCDN.NewWithLogger(&providerDogeCDN.DogeCloudCDNDeployerConfig{
				AccessKey: access.AccessKey,
				SecretKey: access.SecretKey,
				Domain:    maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
			}, logger)
			return deployer, logger, err
		}

	case domain.DeployProviderTypeHuaweiCloudCDN, domain.DeployProviderTypeHuaweiCloudELB:
		{
			access := domain.AccessConfigForHuaweiCloud{}
			if err := maps.Decode(options.ProviderAccessConfig, &access); err != nil {
				return nil, nil, fmt.Errorf("failed to decode provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeployProviderTypeHuaweiCloudCDN:
				deployer, err := providerHuaweiCloudCDN.NewWithLogger(&providerHuaweiCloudCDN.HuaweiCloudCDNDeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					SecretAccessKey: access.SecretAccessKey,
					Region:          maps.GetValueAsString(options.ProviderDeployConfig, "region"),
					Domain:          maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
				}, logger)
				return deployer, logger, err

			case domain.DeployProviderTypeHuaweiCloudELB:
				deployer, err := providerHuaweiCloudELB.NewWithLogger(&providerHuaweiCloudELB.HuaweiCloudELBDeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					SecretAccessKey: access.SecretAccessKey,
					Region:          maps.GetValueAsString(options.ProviderDeployConfig, "region"),
					ResourceType:    providerHuaweiCloudELB.DeployResourceType(maps.GetValueAsString(options.ProviderDeployConfig, "resourceType")),
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
			deployer, err := providerLocal.NewWithLogger(&providerLocal.LocalDeployerConfig{
				ShellEnv:       providerLocal.ShellEnvType(maps.GetValueAsString(options.ProviderDeployConfig, "shellEnv")),
				PreCommand:     maps.GetValueAsString(options.ProviderDeployConfig, "preCommand"),
				PostCommand:    maps.GetValueAsString(options.ProviderDeployConfig, "postCommand"),
				OutputFormat:   providerLocal.OutputFormatType(maps.GetValueOrDefaultAsString(options.ProviderDeployConfig, "format", string(providerLocal.OUTPUT_FORMAT_PEM))),
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
			if err := maps.Decode(options.ProviderAccessConfig, &access); err != nil {
				return nil, nil, fmt.Errorf("failed to decode provider access config: %w", err)
			}

			deployer, err := providerK8sSecret.NewWithLogger(&providerK8sSecret.K8sSecretDeployerConfig{
				KubeConfig:          access.KubeConfig,
				Namespace:           maps.GetValueOrDefaultAsString(options.ProviderDeployConfig, "namespace", "default"),
				SecretName:          maps.GetValueAsString(options.ProviderDeployConfig, "secretName"),
				SecretType:          maps.GetValueOrDefaultAsString(options.ProviderDeployConfig, "secretType", "kubernetes.io/tls"),
				SecretDataKeyForCrt: maps.GetValueOrDefaultAsString(options.ProviderDeployConfig, "secretDataKeyForCrt", "tls.crt"),
				SecretDataKeyForKey: maps.GetValueOrDefaultAsString(options.ProviderDeployConfig, "secretDataKeyForKey", "tls.key"),
			}, logger)
			return deployer, logger, err
		}

	case domain.DeployProviderTypeQiniuCDN:
		{
			access := domain.AccessConfigForQiniu{}
			if err := maps.Decode(options.ProviderAccessConfig, &access); err != nil {
				return nil, nil, fmt.Errorf("failed to decode provider access config: %w", err)
			}

			deployer, err := providerQiniuCDN.NewWithLogger(&providerQiniuCDN.QiniuCDNDeployerConfig{
				AccessKey: access.AccessKey,
				SecretKey: access.SecretKey,
				Domain:    maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
			}, logger)
			return deployer, logger, err
		}

	case domain.DeployProviderTypeSSH:
		{
			access := domain.AccessConfigForSSH{}
			if err := maps.Decode(options.ProviderAccessConfig, &access); err != nil {
				return nil, nil, fmt.Errorf("failed to decode provider access config: %w", err)
			}

			sshPort, _ := strconv.ParseInt(access.Port, 10, 32)
			deployer, err := providerSSH.NewWithLogger(&providerSSH.SshDeployerConfig{
				SshHost:          access.Host,
				SshPort:          int32(sshPort),
				SshUsername:      access.Username,
				SshPassword:      access.Password,
				SshKey:           access.Key,
				SshKeyPassphrase: access.KeyPassphrase,
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

	case domain.DeployProviderTypeTencentCloudCDN, domain.DeployProviderTypeTencentCloudCLB, domain.DeployProviderTypeTencentCloudCOS, domain.DeployProviderTypeTencentCloudECDN, domain.DeployProviderTypeTencentCloudEO:
		{
			access := domain.AccessConfigForTencentCloud{}
			if err := maps.Decode(options.ProviderAccessConfig, &access); err != nil {
				return nil, nil, fmt.Errorf("failed to decode provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeployProviderTypeTencentCloudCDN:
				deployer, err := providerTencentCloudCDN.NewWithLogger(&providerTencentCloudCDN.TencentCloudCDNDeployerConfig{
					SecretId:  access.SecretId,
					SecretKey: access.SecretKey,
					Domain:    maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
				}, logger)
				return deployer, logger, err

			case domain.DeployProviderTypeTencentCloudCLB:
				deployer, err := providerTencentCloudCLB.NewWithLogger(&providerTencentCloudCLB.TencentCloudCLBDeployerConfig{
					SecretId:       access.SecretId,
					SecretKey:      access.SecretKey,
					Region:         maps.GetValueAsString(options.ProviderDeployConfig, "region"),
					ResourceType:   providerTencentCloudCLB.DeployResourceType(maps.GetValueAsString(options.ProviderDeployConfig, "resourceType")),
					LoadbalancerId: maps.GetValueAsString(options.ProviderDeployConfig, "loadbalancerId"),
					ListenerId:     maps.GetValueAsString(options.ProviderDeployConfig, "listenerId"),
					Domain:         maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
				}, logger)
				return deployer, logger, err

			case domain.DeployProviderTypeTencentCloudCOS:
				deployer, err := providerTencentCloudCOD.NewWithLogger(&providerTencentCloudCOD.TencentCloudCOSDeployerConfig{
					SecretId:  access.SecretId,
					SecretKey: access.SecretKey,
					Region:    maps.GetValueAsString(options.ProviderDeployConfig, "region"),
					Bucket:    maps.GetValueAsString(options.ProviderDeployConfig, "bucket"),
					Domain:    maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
				}, logger)
				return deployer, logger, err

			case domain.DeployProviderTypeTencentCloudECDN:
				deployer, err := providerTencentCloudECDN.NewWithLogger(&providerTencentCloudECDN.TencentCloudECDNDeployerConfig{
					SecretId:  access.SecretId,
					SecretKey: access.SecretKey,
					Domain:    maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
				}, logger)
				return deployer, logger, err

			case domain.DeployProviderTypeTencentCloudEO:
				deployer, err := providerTencentCloudEO.NewWithLogger(&providerTencentCloudEO.TencentCloudEODeployerConfig{
					SecretId:  access.SecretId,
					SecretKey: access.SecretKey,
					ZoneId:    maps.GetValueAsString(options.ProviderDeployConfig, "zoneId"),
					Domain:    maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
				}, logger)
				return deployer, logger, err

			default:
				break
			}
		}

	case domain.DeployProviderTypeVolcEngineCDN, domain.DeployProviderTypeVolcEngineDCDN, domain.DeployProviderTypeVolcEngineLive:
		{
			access := domain.AccessConfigForVolcEngine{}
			if err := maps.Decode(options.ProviderAccessConfig, &access); err != nil {
				return nil, nil, fmt.Errorf("failed to decode provider access config: %w", err)
			}

			switch options.Provider {
			case domain.DeployProviderTypeVolcEngineCDN:
				deployer, err := providerVolcEngineCDN.NewWithLogger(&providerVolcEngineCDN.VolcEngineCDNDeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.SecretAccessKey,
					Domain:          maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
				}, logger)
				return deployer, logger, err

			case domain.DeployProviderTypeVolcEngineDCDN:
				deployer, err := providerVolcEngineDCDN.NewWithLogger(&providerVolcEngineDCDN.VolcEngineDCDNDeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.SecretAccessKey,
					Domain:          maps.GetValueAsString(options.ProviderDeployConfig, "domain"),
				}, logger)
				return deployer, logger, err

			case domain.DeployProviderTypeVolcEngineLive:
				deployer, err := providerVolcEngineLive.NewWithLogger(&providerVolcEngineLive.VolcEngineLiveDeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.SecretAccessKey,
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
			if err := maps.Decode(options.ProviderAccessConfig, &access); err != nil {
				return nil, nil, fmt.Errorf("failed to decode provider access config: %w", err)
			}

			deployer, err := providerWebhook.NewWithLogger(&providerWebhook.WebhookDeployerConfig{
				WebhookUrl:  access.Url,
				WebhookData: maps.GetValueAsString(options.ProviderDeployConfig, "webhookData"),
			}, logger)
			return deployer, logger, err
		}
	}

	return nil, nil, fmt.Errorf("unsupported deployer provider: %s", string(options.Provider))
}
