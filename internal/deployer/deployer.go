package deployer

import (
	"context"
	"fmt"

	"github.com/usual2970/certimate/internal/applicant"
	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/logger"
)

/*
提供商部署目标常量值。
短横线前的部分始终等于提供商类型。

	注意：如果追加新的常量值，请保持以 ASCII 排序。
	NOTICE: If you add new constant, please keep ASCII order.
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
	NodeId       string                `json:"nodeId"`
	Domains      string                `json:"domains"`
	AccessConfig string                `json:"accessConfig"`
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

func GetWithTypeAndOption(deployType string, option *DeployerOption) (Deployer, error) {
	deployer, logger, err := createDeployer(deployType, option.AccessRecord.Config, option.DeployConfig.NodeConfig)
	if err != nil {
		return nil, err
	}

	return &proxyDeployer{
		option:   option,
		logger:   logger,
		deployer: deployer,
	}, nil
}

// TODO: 暂时使用代理模式以兼容之前版本代码，后续重新实现此处逻辑
type proxyDeployer struct {
	option   *DeployerOption
	logger   logger.Logger
	deployer deployer.Deployer
}

func (d *proxyDeployer) GetID() string {
	return fmt.Sprintf("%s-%s", d.option.AccessRecord.Name, d.option.AccessRecord.Id)
}

func (d *proxyDeployer) GetInfos() []string {
	return d.logger.GetRecords()
}

func (d *proxyDeployer) Deploy(ctx context.Context) error {
	_, err := d.deployer.Deploy(ctx, d.option.Certificate.Certificate, d.option.Certificate.PrivateKey)
	return err
}
