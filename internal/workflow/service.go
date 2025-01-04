package workflow

import (
	"context"
	"fmt"

	"github.com/usual2970/certimate/internal/app"
	"github.com/usual2970/certimate/internal/domain"
	nodeprocessor "github.com/usual2970/certimate/internal/workflow/node-processor"
)

type WorkflowRepository interface {
	Get(ctx context.Context, id string) (*domain.Workflow, error)
	SaveRunLog(ctx context.Context, log *domain.WorkflowRun) error
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
				Id: workflow.Id,
			})
		})
		if err != nil {
			app.GetApp().Logger().Error("failed to add schedule", "err", err)
			return err
		}
	}
	scheduler.Start()
	app.GetApp().Logger().Info("workflow schedule started")
	return nil
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
		log := &domain.WorkflowRun{
			WorkflowId: workflow.Id,
			Logs:       processor.Log(ctx),
			Succeeded:  false,
			Error:      err.Error(),
		}
		if err := s.repo.SaveRunLog(ctx, log); err != nil {
			app.GetApp().Logger().Error("failed to save run log", "err", err)
		}
		return fmt.Errorf("failed to run workflow: %w", err)
	}

	// 保存执行日志
	logs := processor.Log(ctx)
	runLogs := domain.WorkflowRunLogs(logs)
	runErr := runLogs.FirstError()
	succeed := true
	if runErr != "" {
		succeed = false
	}
	log := &domain.WorkflowRun{
		WorkflowId: workflow.Id,
		Logs:       processor.Log(ctx),
		Error:      runErr,
		Succeeded:  succeed,
	}
	if err := s.repo.SaveRunLog(ctx, log); err != nil {
		app.GetApp().Logger().Error("failed to save run log", "err", err)
		return err
	}

	return nil
}
