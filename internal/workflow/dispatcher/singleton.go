package dispatcher

import (
	"context"
	"sync"

	"github.com/usual2970/certimate/internal/domain"
)

type workflowRepository interface {
	GetById(ctx context.Context, id string) (*domain.Workflow, error)
	Save(ctx context.Context, workflow *domain.Workflow) (*domain.Workflow, error)
}

type workflowRunRepository interface {
	GetById(ctx context.Context, id string) (*domain.WorkflowRun, error)
	Save(ctx context.Context, workflowRun *domain.WorkflowRun) (*domain.WorkflowRun, error)
}

var (
	instance    *WorkflowDispatcher
	intanceOnce sync.Once
)

func GetSingletonDispatcher(workflowRepo workflowRepository, workflowRunRepo workflowRunRepository) *WorkflowDispatcher {
	intanceOnce.Do(func() {
		instance = newWorkflowDispatcher(workflowRepo, workflowRunRepo)
	})

	return instance
}
