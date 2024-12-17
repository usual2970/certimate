package deployer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/pocketbase/pocketbase/models"

	"github.com/usual2970/certimate/internal/applicant"
	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/repository"
)

/*
提供商部署目标常量值。

	注意：如果追加新的枚举值，请保持以 ASCII 排序。
	NOTICE: If you add new enum, please keep ASCII order.
*/
const (
	targetAliyunALB        = "aliyun-alb"
	targetAliyunCDN        = "aliyun-cdn"
	targetAliyunCLB        = "aliyun-clb"
	targetAliyunDCDN       = "aliyun-dcdn"
	targetAliyunNLB        = "aliyun-nlb"
	targetAliyunOSS        = "aliyun-oss"
	targetBaiduCloudCDN    = "baiducloud-cdn"
	targetBytePlusCDN      = "byteplus-cdn"
	targetDogeCloudCDN     = "dogecloud-cdn"
	targetHuaweiCloudCDN   = "huaweicloud-cdn"
	targetHuaweiCloudELB   = "huaweicloud-elb"
	targetK8sSecret        = "k8s-secret"
	targetLocal            = "local"
	targetQiniuCDN         = "qiniu-cdn"
	targetSSH              = "ssh"
	targetTencentCloudCDN  = "tencentcloud-cdn"
	targetTencentCloudCLB  = "tencentcloud-clb"
	targetTencentCloudCOS  = "tencentcloud-cos"
	targetTencentCloudECDN = "tencentcloud-ecdn"
	targetTencentCloudEO   = "tencentcloud-eo"
	targetVolcEngineCDN    = "volcengine-cdn"
	targetVolcEngineLive   = "volcengine-live"
	targetWebhook          = "webhook"
)

type DeployerOption struct {
	DomainId     string                `json:"domainId"`
	Domain       string                `json:"domain"`
	Access       string                `json:"access"`
	AccessRecord *domain.Access        `json:"-"`
	DeployConfig domain.DeployConfig   `json:"deployConfig"`
	Certificate  applicant.Certificate `json:"certificate"`
	Variables    map[string]string     `json:"variables"`
}

type Deployer interface {
	Deploy(ctx context.Context) error
	GetInfos() []string
	GetID() string
}

func Gets(record *models.Record, cert *applicant.Certificate) ([]Deployer, error) {
	rs := make([]Deployer, 0)
	if record.GetString("deployConfig") == "" {
		return rs, nil
	}

	deployConfigs := make([]domain.DeployConfig, 0)

	err := record.UnmarshalJSONField("deployConfig", &deployConfigs)
	if err != nil {
		return nil, fmt.Errorf("解析部署配置失败: %w", err)
	}

	if len(deployConfigs) == 0 {
		return rs, nil
	}

	for _, deployConfig := range deployConfigs {
		deployer, err := getWithDeployConfig(record, cert, deployConfig)
		if err != nil {
			return nil, err
		}

		rs = append(rs, deployer)
	}

	return rs, nil
}

func GetWithTypeAndOption(deployType string, option *DeployerOption) (Deployer, error) {
	return getWithTypeAndOption(deployType, option)
}

func getWithDeployConfig(record *models.Record, cert *applicant.Certificate, deployConfig domain.DeployConfig) (Deployer, error) {
	accessRepo := repository.NewAccessRepository()
	access, err := accessRepo.GetById(context.Background(), deployConfig.Access)
	if err != nil {
		return nil, fmt.Errorf("获取access失败:%w", err)
	}

	option := &DeployerOption{
		DomainId:     record.Id,
		Domain:       record.GetString("domain"),
		Access:       access.Config,
		AccessRecord: access,
		DeployConfig: deployConfig,
	}
	if cert != nil {
		option.Certificate = *cert
	} else {
		option.Certificate = applicant.Certificate{
			Certificate: record.GetString("certificate"),
			PrivateKey:  record.GetString("privateKey"),
		}
	}

	return getWithTypeAndOption(deployConfig.Type, option)
}

func getWithTypeAndOption(deployType string, option *DeployerOption) (Deployer, error) {
	switch deployType {
	case targetAliyunOSS:
		return NewAliyunOSSDeployer(option)
	case targetAliyunCDN:
		return NewAliyunCDNDeployer(option)
	case targetAliyunDCDN:
		return NewAliyunDCDNDeployer(option)
	case targetAliyunCLB:
		return NewAliyunCLBDeployer(option)
	case targetAliyunALB:
		return NewAliyunALBDeployer(option)
	case targetAliyunNLB:
		return NewAliyunNLBDeployer(option)
	case targetTencentCloudCDN:
		return NewTencentCDNDeployer(option)
	case targetTencentCloudECDN:
		return NewTencentECDNDeployer(option)
	case targetTencentCloudCLB:
		return NewTencentCLBDeployer(option)
	case targetTencentCloudCOS:
		return NewTencentCOSDeployer(option)
	case targetTencentCloudEO:
		return NewTencentTEODeployer(option)
	case targetHuaweiCloudCDN:
		return NewHuaweiCloudCDNDeployer(option)
	case targetHuaweiCloudELB:
		return NewHuaweiCloudELBDeployer(option)
	case targetBaiduCloudCDN:
		return NewBaiduCloudCDNDeployer(option)
	case targetQiniuCDN:
		return NewQiniuCDNDeployer(option)
	case targetDogeCloudCDN:
		return NewDogeCloudCDNDeployer(option)
	case targetLocal:
		return NewLocalDeployer(option)
	case targetSSH:
		return NewSSHDeployer(option)
	case targetWebhook:
		return NewWebhookDeployer(option)
	case targetK8sSecret:
		return NewK8sSecretDeployer(option)
	case targetVolcEngineLive:
		return NewVolcengineLiveDeployer(option)
	case targetVolcEngineCDN:
		return NewVolcengineCDNDeployer(option)
	case targetBytePlusCDN:
		return NewByteplusCDNDeployer(option)
	}
	return nil, errors.New("unsupported deploy target")
}

func toStr(tag string, data any) string {
	if data == nil {
		return tag
	}
	byts, _ := json.Marshal(data)
	return tag + "：" + string(byts)
}
