package workflow

import (
	"context"
	"fmt"

	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/utils/app"
	nodeprocessor "github.com/usual2970/certimate/internal/workflow/node-processor"
)

type WorkflowRepository interface {
	Get(ctx context.Context, id string) (*domain.Workflow, error)
}

type WorkflowService struct {
	repo WorkflowRepository
}

func NewWorkflowService(repo WorkflowRepository) *WorkflowService {
	return &WorkflowService{
		repo: repo,
	}
}

func (s *WorkflowService) Run(ctx context.Context, req *domain.WorkflowRunReq) error {
	// 查询
	if req.Id == "" {
		return domain.ErrInvalidParams
	}

	workflow, err := s.repo.Get(ctx, req.Id)
	if err != nil {
		app.GetApp().Logger().Error("failed to get workflow", "id", req.Id, "err", err)
		return err
	}

	// 执行
	if !workflow.Enabled {
		app.GetApp().Logger().Error("workflow is disabled", "id", req.Id)
		return fmt.Errorf("workflow is disabled")
	}

	processor := nodeprocessor.NewWorkflowProcessor(workflow)
	if err := processor.Run(ctx); err != nil {
		return fmt.Errorf("failed to run workflow: %w", err)
	}

	// 保存执行日志

	return nil
}
