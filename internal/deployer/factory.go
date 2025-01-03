package deployer

import (
	"encoding/json"
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
	providerVolcEngineLive "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/volcengine-live"
	providerWebhook "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/webhook"
	"github.com/usual2970/certimate/internal/pkg/core/logger"
	"github.com/usual2970/certimate/internal/pkg/utils/maps"
)

func createDeployer(target string, accessConfig string, deployConfig map[string]any) (deployer.Deployer, logger.Logger, error) {
	logger := logger.NewDefaultLogger()

	/*
	  注意：如果追加新的常量值，请保持以 ASCII 排序。
	  NOTICE: If you add new constant, please keep ASCII order.
	*/
	switch target {
	case targetAliyunALB, targetAliyunCDN, targetAliyunCLB, targetAliyunDCDN, targetAliyunNLB, targetAliyunOSS:
		{
			access := &domain.AliyunAccessConfig{}
			if err := json.Unmarshal([]byte(accessConfig), access); err != nil {
				return nil, nil, fmt.Errorf("failed to unmarshal access config: %w", err)
			}

			switch target {
			case targetAliyunALB:
				deployer, err := providerAliyunALB.NewWithLogger(&providerAliyunALB.AliyunALBDeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Region:          maps.GetValueAsString(deployConfig, "region"),
					ResourceType:    providerAliyunALB.DeployResourceType(maps.GetValueAsString(deployConfig, "resourceType")),
					LoadbalancerId:  maps.GetValueAsString(deployConfig, "loadbalancerId"),
					ListenerId:      maps.GetValueAsString(deployConfig, "listenerId"),
				}, logger)
				return deployer, logger, err

			case targetAliyunCDN:
				deployer, err := providerAliyunCDN.NewWithLogger(&providerAliyunCDN.AliyunCDNDeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Domain:          maps.GetValueAsString(deployConfig, "domain"),
				}, logger)
				return deployer, logger, err

			case targetAliyunCLB:
				deployer, err := providerAliyunCLB.NewWithLogger(&providerAliyunCLB.AliyunCLBDeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Region:          maps.GetValueAsString(deployConfig, "region"),
					ResourceType:    providerAliyunCLB.DeployResourceType(maps.GetValueAsString(deployConfig, "resourceType")),
					LoadbalancerId:  maps.GetValueAsString(deployConfig, "loadbalancerId"),
					ListenerPort:    maps.GetValueAsInt32(deployConfig, "listenerPort"),
				}, logger)
				return deployer, logger, err

			case targetAliyunDCDN:
				deployer, err := providerAliyunDCDN.NewWithLogger(&providerAliyunDCDN.AliyunDCDNDeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Domain:          maps.GetValueAsString(deployConfig, "domain"),
				}, logger)
				return deployer, logger, err

			case targetAliyunNLB:
				deployer, err := providerAliyunNLB.NewWithLogger(&providerAliyunNLB.AliyunNLBDeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Region:          maps.GetValueAsString(deployConfig, "region"),
					ResourceType:    providerAliyunNLB.DeployResourceType(maps.GetValueAsString(deployConfig, "resourceType")),
					LoadbalancerId:  maps.GetValueAsString(deployConfig, "loadbalancerId"),
					ListenerId:      maps.GetValueAsString(deployConfig, "listenerId"),
				}, logger)
				return deployer, logger, err

			case targetAliyunOSS:
				deployer, err := providerAliyunOSS.NewWithLogger(&providerAliyunOSS.AliyunOSSDeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Region:          maps.GetValueAsString(deployConfig, "region"),
					Bucket:          maps.GetValueAsString(deployConfig, "bucket"),
					Domain:          maps.GetValueAsString(deployConfig, "domain"),
				}, logger)
				return deployer, logger, err

			default:
				break
			}
		}

	case targetBaiduCloudCDN:
		{
			access := &domain.BaiduCloudAccessConfig{}
			if err := json.Unmarshal([]byte(accessConfig), access); err != nil {
				return nil, nil, fmt.Errorf("failed to unmarshal access config: %w", err)
			}

			deployer, err := providerBaiduCloudCDN.NewWithLogger(&providerBaiduCloudCDN.BaiduCloudCDNDeployerConfig{
				AccessKeyId:     access.AccessKeyId,
				SecretAccessKey: access.SecretAccessKey,
				Domain:          maps.GetValueAsString(deployConfig, "domain"),
			}, logger)
			return deployer, logger, err
		}

	case targetBytePlusCDN:
		{
			access := &domain.BytePlusAccessConfig{}
			if err := json.Unmarshal([]byte(accessConfig), access); err != nil {
				return nil, nil, fmt.Errorf("failed to unmarshal access config: %w", err)
			}

			deployer, err := providerBytePlusCDN.NewWithLogger(&providerBytePlusCDN.BytePlusCDNDeployerConfig{
				AccessKey: access.AccessKey,
				SecretKey: access.SecretKey,
				Domain:    maps.GetValueAsString(deployConfig, "domain"),
			}, logger)
			return deployer, logger, err
		}

	case targetDogeCloudCDN:
		{
			access := &domain.DogeCloudAccessConfig{}
			if err := json.Unmarshal([]byte(accessConfig), access); err != nil {
				return nil, nil, fmt.Errorf("failed to unmarshal access config: %w", err)
			}

			deployer, err := providerDogeCDN.NewWithLogger(&providerDogeCDN.DogeCloudCDNDeployerConfig{
				AccessKey: access.AccessKey,
				SecretKey: access.SecretKey,
				Domain:    maps.GetValueAsString(deployConfig, "domain"),
			}, logger)
			return deployer, logger, err
		}

	case targetHuaweiCloudCDN, targetHuaweiCloudELB:
		{
			access := &domain.HuaweiCloudAccessConfig{}
			if err := json.Unmarshal([]byte(accessConfig), access); err != nil {
				return nil, nil, fmt.Errorf("failed to unmarshal access config: %w", err)
			}

			switch target {
			case targetHuaweiCloudCDN:
				deployer, err := providerHuaweiCloudCDN.NewWithLogger(&providerHuaweiCloudCDN.HuaweiCloudCDNDeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					SecretAccessKey: access.SecretAccessKey,
					Region:          maps.GetValueAsString(deployConfig, "region"),
					Domain:          maps.GetValueAsString(deployConfig, "domain"),
				}, logger)
				return deployer, logger, err

			case targetHuaweiCloudELB:
				deployer, err := providerHuaweiCloudELB.NewWithLogger(&providerHuaweiCloudELB.HuaweiCloudELBDeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					SecretAccessKey: access.SecretAccessKey,
					Region:          maps.GetValueAsString(deployConfig, "region"),
					ResourceType:    providerHuaweiCloudELB.DeployResourceType(maps.GetValueAsString(deployConfig, "resourceType")),
					CertificateId:   maps.GetValueAsString(deployConfig, "certificateId"),
					LoadbalancerId:  maps.GetValueAsString(deployConfig, "loadbalancerId"),
					ListenerId:      maps.GetValueAsString(deployConfig, "listenerId"),
				}, logger)
				return deployer, logger, err

			default:
				break
			}
		}

	case targetLocal:
		{
			deployer, err := providerLocal.NewWithLogger(&providerLocal.LocalDeployerConfig{
				ShellEnv:       providerLocal.ShellEnvType(maps.GetValueAsString(deployConfig, "shellEnv")),
				PreCommand:     maps.GetValueAsString(deployConfig, "preCommand"),
				PostCommand:    maps.GetValueAsString(deployConfig, "postCommand"),
				OutputFormat:   providerLocal.OutputFormatType(maps.GetValueOrDefaultAsString(deployConfig, "format", string(providerLocal.OUTPUT_FORMAT_PEM))),
				OutputCertPath: maps.GetValueAsString(deployConfig, "certPath"),
				OutputKeyPath:  maps.GetValueAsString(deployConfig, "keyPath"),
				PfxPassword:    maps.GetValueAsString(deployConfig, "pfxPassword"),
				JksAlias:       maps.GetValueAsString(deployConfig, "jksAlias"),
				JksKeypass:     maps.GetValueAsString(deployConfig, "jksKeypass"),
				JksStorepass:   maps.GetValueAsString(deployConfig, "jksStorepass"),
			}, logger)
			return deployer, logger, err
		}

	case targetK8sSecret:
		{
			access := &domain.KubernetesAccessConfig{}
			if err := json.Unmarshal([]byte(accessConfig), access); err != nil {
				return nil, nil, fmt.Errorf("failed to unmarshal access config: %w", err)
			}

			deployer, err := providerK8sSecret.NewWithLogger(&providerK8sSecret.K8sSecretDeployerConfig{
				KubeConfig:          access.KubeConfig,
				Namespace:           maps.GetValueOrDefaultAsString(deployConfig, "namespace", "default"),
				SecretName:          maps.GetValueAsString(deployConfig, "secretName"),
				SecretType:          maps.GetValueOrDefaultAsString(deployConfig, "secretType", "kubernetes.io/tls"),
				SecretDataKeyForCrt: maps.GetValueOrDefaultAsString(deployConfig, "secretDataKeyForCrt", "tls.crt"),
				SecretDataKeyForKey: maps.GetValueOrDefaultAsString(deployConfig, "secretDataKeyForKey", "tls.key"),
			}, logger)
			return deployer, logger, err
		}

	case targetQiniuCDN:
		{
			access := &domain.QiniuAccessConfig{}
			if err := json.Unmarshal([]byte(accessConfig), access); err != nil {
				return nil, nil, fmt.Errorf("failed to unmarshal access config: %w", err)
			}

			deployer, err := providerQiniuCDN.NewWithLogger(&providerQiniuCDN.QiniuCDNDeployerConfig{
				AccessKey: access.AccessKey,
				SecretKey: access.SecretKey,
				Domain:    maps.GetValueAsString(deployConfig, "domain"),
			}, logger)
			return deployer, logger, err
		}

	case targetSSH:
		{
			access := &domain.SSHAccessConfig{}
			if err := json.Unmarshal([]byte(accessConfig), access); err != nil {
				return nil, nil, fmt.Errorf("failed to unmarshal access config: %w", err)
			}

			sshPort, _ := strconv.ParseInt(access.Port, 10, 32)
			deployer, err := providerSSH.NewWithLogger(&providerSSH.SshDeployerConfig{
				SshHost:          access.Host,
				SshPort:          int32(sshPort),
				SshUsername:      access.Username,
				SshPassword:      access.Password,
				SshKey:           access.Key,
				SshKeyPassphrase: access.KeyPassphrase,
				PreCommand:       maps.GetValueAsString(deployConfig, "preCommand"),
				PostCommand:      maps.GetValueAsString(deployConfig, "postCommand"),
				OutputFormat:     providerSSH.OutputFormatType(maps.GetValueOrDefaultAsString(deployConfig, "format", string(providerSSH.OUTPUT_FORMAT_PEM))),
				OutputCertPath:   maps.GetValueAsString(deployConfig, "certPath"),
				OutputKeyPath:    maps.GetValueAsString(deployConfig, "keyPath"),
				PfxPassword:      maps.GetValueAsString(deployConfig, "pfxPassword"),
				JksAlias:         maps.GetValueAsString(deployConfig, "jksAlias"),
				JksKeypass:       maps.GetValueAsString(deployConfig, "jksKeypass"),
				JksStorepass:     maps.GetValueAsString(deployConfig, "jksStorepass"),
			}, logger)
			return deployer, logger, err
		}

	case targetTencentCloudCDN, targetTencentCloudCLB, targetTencentCloudCOS, targetTencentCloudECDN, targetTencentCloudEO:
		{
			access := &domain.TencentCloudAccessConfig{}
			if err := json.Unmarshal([]byte(accessConfig), access); err != nil {
				return nil, nil, fmt.Errorf("failed to unmarshal access config: %w", err)
			}

			switch target {
			case targetTencentCloudCDN:
				deployer, err := providerTencentCloudCDN.NewWithLogger(&providerTencentCloudCDN.TencentCloudCDNDeployerConfig{
					SecretId:  access.SecretId,
					SecretKey: access.SecretKey,
					Domain:    maps.GetValueAsString(deployConfig, "domain"),
				}, logger)
				return deployer, logger, err

			case targetTencentCloudCLB:
				deployer, err := providerTencentCloudCLB.NewWithLogger(&providerTencentCloudCLB.TencentCloudCLBDeployerConfig{
					SecretId:       access.SecretId,
					SecretKey:      access.SecretKey,
					Region:         maps.GetValueAsString(deployConfig, "region"),
					ResourceType:   providerTencentCloudCLB.DeployResourceType(maps.GetValueAsString(deployConfig, "resourceType")),
					LoadbalancerId: maps.GetValueAsString(deployConfig, "loadbalancerId"),
					ListenerId:     maps.GetValueAsString(deployConfig, "listenerId"),
					Domain:         maps.GetValueAsString(deployConfig, "domain"),
				}, logger)
				return deployer, logger, err

			case targetTencentCloudCOS:
				deployer, err := providerTencentCloudCOD.NewWithLogger(&providerTencentCloudCOD.TencentCloudCOSDeployerConfig{
					SecretId:  access.SecretId,
					SecretKey: access.SecretKey,
					Region:    maps.GetValueAsString(deployConfig, "region"),
					Bucket:    maps.GetValueAsString(deployConfig, "bucket"),
					Domain:    maps.GetValueAsString(deployConfig, "domain"),
				}, logger)
				return deployer, logger, err

			case targetTencentCloudECDN:
				deployer, err := providerTencentCloudECDN.NewWithLogger(&providerTencentCloudECDN.TencentCloudECDNDeployerConfig{
					SecretId:  access.SecretId,
					SecretKey: access.SecretKey,
					Domain:    maps.GetValueAsString(deployConfig, "domain"),
				}, logger)
				return deployer, logger, err

			case targetTencentCloudEO:
				deployer, err := providerTencentCloudEO.NewWithLogger(&providerTencentCloudEO.TencentCloudEODeployerConfig{
					SecretId:  access.SecretId,
					SecretKey: access.SecretKey,
					ZoneId:    maps.GetValueAsString(deployConfig, "zoneId"),
					Domain:    maps.GetValueAsString(deployConfig, "domain"),
				}, logger)
				return deployer, logger, err

			default:
				break
			}
		}

	case targetVolcEngineCDN, targetVolcEngineLive:
		{
			access := &domain.VolcEngineAccessConfig{}
			if err := json.Unmarshal([]byte(accessConfig), access); err != nil {
				return nil, nil, fmt.Errorf("failed to unmarshal access config: %w", err)
			}

			switch target {
			case targetVolcEngineCDN:
				deployer, err := providerVolcEngineCDN.NewWithLogger(&providerVolcEngineCDN.VolcEngineCDNDeployerConfig{
					AccessKey: access.AccessKeyId,
					SecretKey: access.SecretAccessKey,
					Domain:    maps.GetValueAsString(deployConfig, "domain"),
				}, logger)
				return deployer, logger, err

			case targetVolcEngineLive:
				deployer, err := providerVolcEngineLive.NewWithLogger(&providerVolcEngineLive.VolcEngineLiveDeployerConfig{
					AccessKey: access.AccessKeyId,
					SecretKey: access.SecretAccessKey,
					Domain:    maps.GetValueAsString(deployConfig, "domain"),
				}, logger)
				return deployer, logger, err

			default:
				break
			}
		}

	case targetWebhook:
		{
			access := &domain.WebhookAccessConfig{}
			if err := json.Unmarshal([]byte(accessConfig), access); err != nil {
				return nil, nil, fmt.Errorf("failed to unmarshal access config: %w", err)
			}

			deployer, err := providerWebhook.NewWithLogger(&providerWebhook.WebhookDeployerConfig{
				WebhookUrl:  access.Url,
				WebhookData: maps.GetValueAsString(deployConfig, "webhookData"),
			}, logger)
			return deployer, logger, err
		}
	}

	return nil, nil, fmt.Errorf("unsupported deployer target: %s", target)
}
