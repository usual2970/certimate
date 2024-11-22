package deployer

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	providerAliyunAlb "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/aliyun-alb"
	providerAliyunCdn "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/aliyun-cdn"
	providerAliyunClb "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/aliyun-clb"
	providerAliyunDcdn "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/aliyun-dcdn"
	providerAliyunNlb "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/aliyun-nlb"
	providerAliyunOss "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/aliyun-oss"
	providerBaiduCloudCdn "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/baiducloud-cdn"
	providerBytePlusCdn "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/byteplus-cdn"
	providerDogeCdn "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/dogecloud-cdn"
	providerHuaweiCloudCdn "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/huaweicloud-cdn"
	providerHuaweiCloudElb "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/huaweicloud-elb"
	providerK8sSecret "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/k8s-secret"
	providerLocal "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/local"
	providerQiniuCdn "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/qiniu-cdn"
	providerSSH "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/ssh"
	providerTencentCloudCdn "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/tencentcloud-cdn"
	providerTencentCloudClb "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/tencentcloud-clb"
	providerTencentCloudCos "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/tencentcloud-cos"
	providerTencentCloudEcdn "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/tencentcloud-ecdn"
	providerTencentCloudTeo "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/tencentcloud-teo"
	providerVolcEngineCdn "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/volcengine-cdn"
	providerVolcEngineLive "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/volcengine-live"
	providerWebhook "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/webhook"
	"github.com/usual2970/certimate/internal/pkg/utils/maps"
)

// TODO: 该方法目前未实际使用，将在后续迭代中替换
func createDeployer(target string, accessConfig string, deployConfig map[string]any) (deployer.Deployer, deployer.Logger, error) {
	logger := deployer.NewDefaultLogger()

	switch target {
	case targetAliyunALB, targetAliyunCDN, targetAliyunCLB, targetAliyunDCDN, targetAliyunNLB, targetAliyunOSS:
		{
			access := &domain.AliyunAccess{}
			if err := json.Unmarshal([]byte(accessConfig), access); err != nil {
				return nil, nil, fmt.Errorf("failed to unmarshal access config: %w", err)
			}

			switch target {
			case targetAliyunALB:
				deployer, err := providerAliyunAlb.NewWithLogger(&providerAliyunAlb.AliyunALBDeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Region:          maps.GetValueAsString(deployConfig, "region"),
					ResourceType:    providerAliyunAlb.DeployResourceType(maps.GetValueAsString(deployConfig, "resourceType")),
					LoadbalancerId:  maps.GetValueAsString(deployConfig, "loadbalancerId"),
					ListenerId:      maps.GetValueAsString(deployConfig, "listenerId"),
				}, logger)
				return deployer, logger, err

			case targetAliyunCDN:
				deployer, err := providerAliyunCdn.NewWithLogger(&providerAliyunCdn.AliyunCDNDeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Domain:          maps.GetValueAsString(deployConfig, "domain"),
				}, logger)
				return deployer, logger, err

			case targetAliyunCLB:
				deployer, err := providerAliyunClb.NewWithLogger(&providerAliyunClb.AliyunCLBDeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Region:          maps.GetValueAsString(deployConfig, "region"),
					ResourceType:    providerAliyunClb.DeployResourceType(maps.GetValueAsString(deployConfig, "resourceType")),
					LoadbalancerId:  maps.GetValueAsString(deployConfig, "loadbalancerId"),
					ListenerPort:    maps.GetValueAsInt32(deployConfig, "listenerPort"),
				}, logger)
				return deployer, logger, err

			case targetAliyunDCDN:
				deployer, err := providerAliyunDcdn.NewWithLogger(&providerAliyunDcdn.AliyunDCDNDeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Domain:          maps.GetValueAsString(deployConfig, "domain"),
				}, logger)
				return deployer, logger, err

			case targetAliyunNLB:
				deployer, err := providerAliyunNlb.NewWithLogger(&providerAliyunNlb.AliyunNLBDeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					AccessKeySecret: access.AccessKeySecret,
					Region:          maps.GetValueAsString(deployConfig, "region"),
					ResourceType:    providerAliyunNlb.DeployResourceType(maps.GetValueAsString(deployConfig, "resourceType")),
					LoadbalancerId:  maps.GetValueAsString(deployConfig, "loadbalancerId"),
					ListenerId:      maps.GetValueAsString(deployConfig, "listenerId"),
				}, logger)
				return deployer, logger, err

			case targetAliyunOSS:
				deployer, err := providerAliyunOss.NewWithLogger(&providerAliyunOss.AliyunOSSDeployerConfig{
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
			access := &domain.BaiduCloudAccess{}
			if err := json.Unmarshal([]byte(accessConfig), access); err != nil {
				return nil, nil, fmt.Errorf("failed to unmarshal access config: %w", err)
			}

			deployer, err := providerBaiduCloudCdn.NewWithLogger(&providerBaiduCloudCdn.BaiduCloudCDNDeployerConfig{
				AccessKeyId:     access.AccessKeyId,
				SecretAccessKey: access.SecretAccessKey,
				Domain:          maps.GetValueAsString(deployConfig, "domain"),
			}, logger)
			return deployer, logger, err
		}

	case targetBytePlusCDN:
		{
			access := &domain.ByteplusAccess{}
			if err := json.Unmarshal([]byte(accessConfig), access); err != nil {
				return nil, nil, fmt.Errorf("failed to unmarshal access config: %w", err)
			}

			deployer, err := providerBytePlusCdn.NewWithLogger(&providerBytePlusCdn.BytePlusCDNDeployerConfig{
				AccessKey: access.AccessKey,
				SecretKey: access.SecretKey,
				Domain:    maps.GetValueAsString(deployConfig, "domain"),
			}, logger)
			return deployer, logger, err
		}

	case targetDogeCloudCdn:
		{
			access := &domain.DogeCloudAccess{}
			if err := json.Unmarshal([]byte(accessConfig), access); err != nil {
				return nil, nil, fmt.Errorf("failed to unmarshal access config: %w", err)
			}

			deployer, err := providerDogeCdn.NewWithLogger(&providerDogeCdn.DogeCloudCDNDeployerConfig{
				AccessKey: access.AccessKey,
				SecretKey: access.SecretKey,
				Domain:    maps.GetValueAsString(deployConfig, "domain"),
			}, logger)
			return deployer, logger, err
		}

	case targetHuaweiCloudCDN, targetHuaweiCloudELB:
		{
			access := &domain.HuaweiCloudAccess{}
			if err := json.Unmarshal([]byte(accessConfig), access); err != nil {
				return nil, nil, fmt.Errorf("failed to unmarshal access config: %w", err)
			}

			switch target {
			case targetHuaweiCloudCDN:
				deployer, err := providerHuaweiCloudCdn.NewWithLogger(&providerHuaweiCloudCdn.HuaweiCloudCDNDeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					SecretAccessKey: access.SecretAccessKey,
					Region:          maps.GetValueAsString(deployConfig, "region"),
					Domain:          maps.GetValueAsString(deployConfig, "domain"),
				}, logger)
				return deployer, logger, err

			case targetHuaweiCloudELB:
				deployer, err := providerHuaweiCloudElb.NewWithLogger(&providerHuaweiCloudElb.HuaweiCloudELBDeployerConfig{
					AccessKeyId:     access.AccessKeyId,
					SecretAccessKey: access.SecretAccessKey,
					Region:          maps.GetValueAsString(deployConfig, "region"),
					ResourceType:    providerHuaweiCloudElb.DeployResourceType(maps.GetValueAsString(deployConfig, "resourceType")),
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
				OutputFormat:   providerLocal.OutputFormatType(maps.GetValueOrDefaultAsString(deployConfig, "format", "PEM")),
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
			access := &domain.KubernetesAccess{}
			if err := json.Unmarshal([]byte(accessConfig), access); err != nil {
				return nil, nil, fmt.Errorf("failed to unmarshal access config: %w", err)
			}

			deployer, err := providerK8sSecret.NewWithLogger(&providerK8sSecret.K8sSecretDeployerConfig{
				KubeConfig:          access.KubeConfig,
				Namespace:           maps.GetValueOrDefaultAsString(deployConfig, "namespace", "default"),
				SecretName:          maps.GetValueAsString(deployConfig, "secretName"),
				SecretDataKeyForCrt: maps.GetValueOrDefaultAsString(deployConfig, "secretDataKeyForCrt", "tls.crt"),
				SecretDataKeyForKey: maps.GetValueOrDefaultAsString(deployConfig, "secretDataKeyForKey", "tls.key"),
			}, logger)
			return deployer, logger, err
		}

	case targetQiniuCdn:
		{
			access := &domain.QiniuAccess{}
			if err := json.Unmarshal([]byte(accessConfig), access); err != nil {
				return nil, nil, fmt.Errorf("failed to unmarshal access config: %w", err)
			}

			deployer, err := providerQiniuCdn.NewWithLogger(&providerQiniuCdn.QiniuCDNDeployerConfig{
				AccessKey: access.AccessKey,
				SecretKey: access.SecretKey,
				Domain:    maps.GetValueAsString(deployConfig, "domain"),
			}, logger)
			return deployer, logger, err
		}

	case targetSSH:
		{
			access := &domain.SSHAccess{}
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
				OutputFormat:     providerSSH.OutputFormatType(maps.GetValueOrDefaultAsString(deployConfig, "format", "PEM")),
				OutputCertPath:   maps.GetValueAsString(deployConfig, "certPath"),
				OutputKeyPath:    maps.GetValueAsString(deployConfig, "keyPath"),
				PfxPassword:      maps.GetValueAsString(deployConfig, "pfxPassword"),
				JksAlias:         maps.GetValueAsString(deployConfig, "jksAlias"),
				JksKeypass:       maps.GetValueAsString(deployConfig, "jksKeypass"),
				JksStorepass:     maps.GetValueAsString(deployConfig, "jksStorepass"),
			}, logger)
			return deployer, logger, err
		}

	case targetTencentCDN, targetTencentCLB, targetTencentCOS, targetTencentECDN, targetTencentTEO:
		{
			access := &domain.TencentAccess{}
			if err := json.Unmarshal([]byte(accessConfig), access); err != nil {
				return nil, nil, fmt.Errorf("failed to unmarshal access config: %w", err)
			}

			switch target {
			case targetTencentCDN:
				deployer, err := providerTencentCloudCdn.NewWithLogger(&providerTencentCloudCdn.TencentCloudCDNDeployerConfig{
					SecretId:  access.SecretId,
					SecretKey: access.SecretKey,
					Domain:    maps.GetValueAsString(deployConfig, "domain"),
				}, logger)
				return deployer, logger, err

			case targetTencentCLB:
				deployer, err := providerTencentCloudClb.NewWithLogger(&providerTencentCloudClb.TencentCloudCLBDeployerConfig{
					SecretId:       access.SecretId,
					SecretKey:      access.SecretKey,
					Region:         maps.GetValueAsString(deployConfig, "region"),
					ResourceType:   providerTencentCloudClb.DeployResourceType(maps.GetValueAsString(deployConfig, "resourceType")),
					LoadbalancerId: maps.GetValueAsString(deployConfig, "loadbalancerId"),
					ListenerId:     maps.GetValueAsString(deployConfig, "listenerId"),
					Domain:         maps.GetValueAsString(deployConfig, "domain"),
				}, logger)
				return deployer, logger, err

			case targetTencentCOS:
				deployer, err := providerTencentCloudCos.NewWithLogger(&providerTencentCloudCos.TencentCloudCOSDeployerConfig{
					SecretId:  access.SecretId,
					SecretKey: access.SecretKey,
					Region:    maps.GetValueAsString(deployConfig, "region"),
					Bucket:    maps.GetValueAsString(deployConfig, "bucket"),
					Domain:    maps.GetValueAsString(deployConfig, "domain"),
				}, logger)
				return deployer, logger, err

			case targetTencentECDN:
				deployer, err := providerTencentCloudEcdn.NewWithLogger(&providerTencentCloudEcdn.TencentCloudECDNDeployerConfig{
					SecretId:  access.SecretId,
					SecretKey: access.SecretKey,
					Domain:    maps.GetValueAsString(deployConfig, "domain"),
				}, logger)
				return deployer, logger, err

			case targetTencentTEO:
				deployer, err := providerTencentCloudTeo.NewWithLogger(&providerTencentCloudTeo.TencentCloudTEODeployerConfig{
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
			access := &domain.VolcEngineAccess{}
			if err := json.Unmarshal([]byte(accessConfig), access); err != nil {
				return nil, nil, fmt.Errorf("failed to unmarshal access config: %w", err)
			}

			switch target {
			case targetVolcEngineCDN:
				deployer, err := providerVolcEngineCdn.NewWithLogger(&providerVolcEngineCdn.VolcEngineCDNDeployerConfig{
					AccessKey: access.AccessKey,
					SecretKey: access.SecretKey,
					Domain:    maps.GetValueAsString(deployConfig, "domain"),
				}, logger)
				return deployer, logger, err

			case targetVolcEngineLive:
				deployer, err := providerVolcEngineLive.NewWithLogger(&providerVolcEngineLive.VolcEngineLiveDeployerConfig{
					AccessKey: access.AccessKey,
					SecretKey: access.SecretKey,
					Domain:    maps.GetValueAsString(deployConfig, "domain"),
				}, logger)
				return deployer, logger, err

			default:
				break
			}
		}

	case targetWebhook:
		{
			access := &domain.WebhookAccess{}
			if err := json.Unmarshal([]byte(accessConfig), access); err != nil {
				return nil, nil, fmt.Errorf("failed to unmarshal access config: %w", err)
			}

			deployer, err := providerWebhook.NewWithLogger(&providerWebhook.WebhookDeployerConfig{
				Url:       access.Url,
				Variables: nil, // TODO: 尚未实现
			}, logger)
			return deployer, logger, err
		}
	}

	return nil, nil, fmt.Errorf("unsupported deployer target: %s", target)
}
