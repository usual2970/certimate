package deployer

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/repository"
)

type Deployer interface {
	SetLogger(*slog.Logger)

	Deploy(ctx context.Context) error
}

type deployerOptions struct {
	Provider             domain.DeployProviderType
	ProviderAccessConfig map[string]any
	ProviderDeployConfig map[string]any
}

func NewWithDeployNode(node *domain.WorkflowNode, certdata struct {
	Certificate string
	PrivateKey  string
},
) (Deployer, error) {
	if node.Type != domain.WorkflowNodeTypeDeploy {
		return nil, fmt.Errorf("node type is not deploy")
	}

	nodeConfig := node.GetConfigForDeploy()

	accessRepo := repository.NewAccessRepository()
	access, err := accessRepo.GetById(context.Background(), nodeConfig.ProviderAccessId)
	if err != nil {
		return nil, fmt.Errorf("failed to get access #%s record: %w", nodeConfig.ProviderAccessId, err)
	}

	deployer, err := createDeployer(&deployerOptions{
		Provider:             domain.DeployProviderType(nodeConfig.Provider),
		ProviderAccessConfig: access.Config,
		ProviderDeployConfig: nodeConfig.ProviderConfig,
	})
	if err != nil {
		return nil, err
	}

	return &proxyDeployer{
		deployer:          deployer,
		deployCertificate: certdata.Certificate,
		deployPrivateKey:  certdata.PrivateKey,
	}, nil
}

// TODO: 暂时使用代理模式以兼容之前版本代码，后续重新实现此处逻辑
type proxyDeployer struct {
	deployer          deployer.Deployer
	deployCertificate string
	deployPrivateKey  string
}

func (d *proxyDeployer) SetLogger(logger *slog.Logger) {
	if logger == nil {
		panic("logger is nil")
	}

	d.deployer.WithLogger(logger)
}

func (d *proxyDeployer) Deploy(ctx context.Context) error {
	_, err := d.deployer.Deploy(ctx, d.deployCertificate, d.deployPrivateKey)
	return err
}
