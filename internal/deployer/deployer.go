package deployer

import (
	"context"

	"github.com/usual2970/certimate/internal/applicant"
	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/logger"
)

type DeployerOption struct {
	NodeId       string                `json:"nodeId"`
	Domains      string                `json:"domains"`
	AccessConfig string                `json:"accessConfig"`
	AccessRecord *domain.Access        `json:"-"`
	DeployConfig domain.DeployConfig   `json:"deployConfig"`
	Certificate  applicant.Certificate `json:"certificate"`
}

type Deployer interface {
	Deploy(ctx context.Context) error
}

func GetWithProviderAndOption(provider string, option *DeployerOption) (Deployer, error) {
	deployer, logger, err := createDeployer(domain.DeployProviderType(provider), option.AccessRecord.Config, option.DeployConfig.NodeConfig)
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

func (d *proxyDeployer) Deploy(ctx context.Context) error {
	_, err := d.deployer.Deploy(ctx, d.option.Certificate.Certificate, d.option.Certificate.PrivateKey)
	return err
}
