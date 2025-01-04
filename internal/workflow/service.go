package workflow

import (
	"context"
	"fmt"
	"time"

	"github.com/usual2970/certimate/internal/app"
	"github.com/usual2970/certimate/internal/domain"
	nodeprocessor "github.com/usual2970/certimate/internal/workflow/node-processor"
)

type WorkflowRepository interface {
	Get(ctx context.Context, id string) (*domain.Workflow, error)
	SaveRun(ctx context.Context, run *domain.WorkflowRun) error
	ListEnabledAuto(ctx context.Context) ([]domain.Workflow, error)
}

type WorkflowService struct {
	repo WorkflowRepository
}

func NewWorkflowService(repo WorkflowRepository) *WorkflowService {
	return &WorkflowService{
		repo: repo,
	}
}

func (s *WorkflowService) InitSchedule(ctx context.Context) error {
	// 查询所有的 enabled auto workflow
	workflows, err := s.repo.ListEnabledAuto(ctx)
	if err != nil {
		return err
	}

	scheduler := app.GetScheduler()
	for _, workflow := range workflows {
		err := scheduler.Add(workflow.Id, workflow.TriggerCron, func() {
			s.Run(ctx, &domain.WorkflowRunReq{
				WorkflowId: workflow.Id,
				Trigger:    domain.WorkflowTriggerTypeAuto,
			})
		})
		if err != nil {
			app.GetLogger().Error("failed to add schedule", "err", err)
			return err
		}
	}
	scheduler.Start()

	app.GetLogger().Info("workflow schedule started")

	return nil
}

func (s *WorkflowService) Run(ctx context.Context, options *domain.WorkflowRunReq) error {
	// 查询
	if options.WorkflowId == "" {
		return domain.ErrInvalidParams
	}

	workflow, err := s.repo.Get(ctx, options.WorkflowId)
	if err != nil {
		app.GetLogger().Error("failed to get workflow", "id", options.WorkflowId, "err", err)
		return err
	}

	if !workflow.Enabled {
		app.GetLogger().Error("workflow is disabled", "id", options.WorkflowId)
		return fmt.Errorf("workflow is disabled")
	}

	// 执行
	run := &domain.WorkflowRun{
		WorkflowId: workflow.Id,
		Status:     domain.WorkflowRunStatusTypeRunning,
		Trigger:    options.Trigger,
		StartedAt:  time.Now(),
		EndedAt:    time.Now(),
	}

	processor := nodeprocessor.NewWorkflowProcessor(workflow)
	if err := processor.Run(ctx); err != nil {
		run.Status = domain.WorkflowRunStatusTypeFailed
		run.EndedAt = time.Now()
		run.Logs = processor.Log(ctx)
		run.Error = err.Error()

		if err := s.repo.SaveRun(ctx, run); err != nil {
			app.GetLogger().Error("failed to save workflow run", "err", err)
		}

		return fmt.Errorf("failed to run workflow: %w", err)
	}

	// 保存执行日志
	logs := processor.Log(ctx)
	runStatus := domain.WorkflowRunStatusTypeSucceeded
	runError := domain.WorkflowRunLogs(logs).FirstError()
	if runError != "" {
		runStatus = domain.WorkflowRunStatusTypeFailed
	}
	run.Status = runStatus
	run.EndedAt = time.Now()
	run.Logs = processor.Log(ctx)
	run.Error = runError
	if err := s.repo.SaveRun(ctx, run); err != nil {
		app.GetLogger().Error("failed to save workflow run", "err", err)
		return err
	}

	return nil
}
