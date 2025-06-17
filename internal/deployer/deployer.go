package deployer

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/internal/repository"
	"github.com/certimate-go/certimate/pkg/core"
)

type Deployer interface {
	Deploy(ctx context.Context) error
}

type DeployerWithWorkflowNodeConfig struct {
	Node           *domain.WorkflowNode
	Logger         *slog.Logger
	CertificatePEM string
	PrivateKeyPEM  string
}

func NewWithWorkflowNode(config DeployerWithWorkflowNodeConfig) (Deployer, error) {
	if config.Node == nil {
		return nil, fmt.Errorf("node is nil")
	}
	if config.Node.Type != domain.WorkflowNodeTypeDeploy {
		return nil, fmt.Errorf("node type is not '%s'", string(domain.WorkflowNodeTypeDeploy))
	}

	nodeCfg := config.Node.GetConfigForDeploy()
	options := &deployerProviderOptions{
		Provider:              domain.DeploymentProviderType(nodeCfg.Provider),
		ProviderAccessConfig:  make(map[string]any),
		ProviderServiceConfig: nodeCfg.ProviderConfig,
	}

	accessRepo := repository.NewAccessRepository()
	if nodeCfg.ProviderAccessId != "" {
		access, err := accessRepo.GetById(context.Background(), nodeCfg.ProviderAccessId)
		if err != nil {
			return nil, fmt.Errorf("failed to get access #%s record: %w", nodeCfg.ProviderAccessId, err)
		} else {
			options.ProviderAccessConfig = access.Config
		}
	}

	deployer, err := createSSLDeployerProvider(options)
	if err != nil {
		return nil, err
	} else {
		deployer.SetLogger(config.Logger)
	}

	return &deployerImpl{
		provider:   deployer,
		certPEM:    config.CertificatePEM,
		privkeyPEM: config.PrivateKeyPEM,
	}, nil
}

type deployerImpl struct {
	provider   core.SSLDeployer
	certPEM    string
	privkeyPEM string
}

var _ Deployer = (*deployerImpl)(nil)

func (d *deployerImpl) Deploy(ctx context.Context) error {
	_, err := d.provider.Deploy(ctx, d.certPEM, d.privkeyPEM)
	return err
}
