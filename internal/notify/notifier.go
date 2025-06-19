package notify

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/internal/repository"
	"github.com/certimate-go/certimate/pkg/core"
)

type Notifier interface {
	Notify(ctx context.Context) error
}

type NotifierWithWorkflowNodeConfig struct {
	Node    *domain.WorkflowNode
	Logger  *slog.Logger
	Subject string
	Message string
}

func NewWithWorkflowNode(config NotifierWithWorkflowNodeConfig) (Notifier, error) {
	if config.Node == nil {
		return nil, fmt.Errorf("node is nil")
	}
	if config.Node.Type != domain.WorkflowNodeTypeNotify {
		return nil, fmt.Errorf("node type is not '%s'", string(domain.WorkflowNodeTypeNotify))
	}

	nodeCfg := config.Node.GetConfigForNotify()
	options := &notifierProviderOptions{
		Provider:              domain.NotificationProviderType(nodeCfg.Provider),
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

	notifier, err := createNotifierProvider(options)
	if err != nil {
		return nil, err
	} else {
		notifier.SetLogger(config.Logger)
	}

	return &notifierImpl{
		provider: notifier,
		subject:  config.Subject,
		message:  config.Message,
	}, nil
}

type notifierImpl struct {
	provider core.Notifier
	subject  string
	message  string
}

var _ Notifier = (*notifierImpl)(nil)

func (n *notifierImpl) Notify(ctx context.Context) error {
	_, err := n.provider.Notify(ctx, n.subject, n.message)
	return err
}
