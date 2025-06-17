package dispatcher

import (
	"context"
	"sync"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/internal/repository"
)

type workflowRepository interface {
	GetById(ctx context.Context, id string) (*domain.Workflow, error)
	Save(ctx context.Context, workflow *domain.Workflow) (*domain.Workflow, error)
}

type workflowRunRepository interface {
	GetById(ctx context.Context, id string) (*domain.WorkflowRun, error)
	Save(ctx context.Context, workflowRun *domain.WorkflowRun) (*domain.WorkflowRun, error)
}

type workflowLogRepository interface {
	Save(ctx context.Context, workflowLog *domain.WorkflowLog) (*domain.WorkflowLog, error)
}

var (
	instance    *WorkflowDispatcher
	intanceOnce sync.Once
)

func GetSingletonDispatcher() *WorkflowDispatcher {
	intanceOnce.Do(func() {
		instance = newWorkflowDispatcher(repository.NewWorkflowRepository(), repository.NewWorkflowRunRepository(), repository.NewWorkflowLogRepository())
	})

	return instance
}
