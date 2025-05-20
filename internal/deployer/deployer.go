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

	nodeConfig := config.Node.GetConfigForDeploy()
	options := &deployerProviderOptions{
		Provider:              domain.DeploymentProviderType(nodeConfig.Provider),
		ProviderAccessConfig:  make(map[string]any),
		ProviderServiceConfig: nodeConfig.ProviderConfig,
	}

	accessRepo := repository.NewAccessRepository()
	if nodeConfig.ProviderAccessId != "" {
		access, err := accessRepo.GetById(context.Background(), nodeConfig.ProviderAccessId)
		if err != nil {
			return nil, fmt.Errorf("failed to get access #%s record: %w", nodeConfig.ProviderAccessId, err)
		} else {
			options.ProviderAccessConfig = access.Config
		}
	}

	deployerProvider, err := createDeployerProvider(options)
	if err != nil {
		return nil, err
	}

	return &deployerImpl{
		provider:   deployerProvider.WithLogger(config.Logger),
		certPEM:    config.CertificatePEM,
		privkeyPEM: config.PrivateKeyPEM,
	}, nil
}

type deployerImpl struct {
	provider   deployer.Deployer
	certPEM    string
	privkeyPEM string
}

var _ Deployer = (*deployerImpl)(nil)

func (d *deployerImpl) Deploy(ctx context.Context) error {
	_, err := d.provider.Deploy(ctx, d.certPEM, d.privkeyPEM)
	return err
}
