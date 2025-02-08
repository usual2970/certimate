package workflow

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/usual2970/certimate/internal/app"
	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/domain/dtos"
	"github.com/usual2970/certimate/internal/workflow/dispatcher"
)

type workflowRepository interface {
	ListEnabledAuto(ctx context.Context) ([]*domain.Workflow, error)
	GetById(ctx context.Context, id string) (*domain.Workflow, error)
	Save(ctx context.Context, workflow *domain.Workflow) (*domain.Workflow, error)
}

type workflowRunRepository interface {
	GetById(ctx context.Context, id string) (*domain.WorkflowRun, error)
	Save(ctx context.Context, workflowRun *domain.WorkflowRun) (*domain.WorkflowRun, error)
}

type WorkflowService struct {
	dispatcher *dispatcher.WorkflowDispatcher

	workflowRepo    workflowRepository
	workflowRunRepo workflowRunRepository
}

func NewWorkflowService(workflowRepo workflowRepository, workflowRunRepo workflowRunRepository) *WorkflowService {
	srv := &WorkflowService{
		dispatcher: dispatcher.GetSingletonDispatcher(workflowRepo, workflowRunRepo),

		workflowRepo:    workflowRepo,
		workflowRunRepo: workflowRunRepo,
	}
	return srv
}

func (s *WorkflowService) InitSchedule(ctx context.Context) error {
	workflows, err := s.workflowRepo.ListEnabledAuto(ctx)
	if err != nil {
		return err
	}

	scheduler := app.GetScheduler()
	for _, workflow := range workflows {
		err := scheduler.Add(fmt.Sprintf("workflow#%s", workflow.Id), workflow.TriggerCron, func() {
			s.StartRun(ctx, &dtos.WorkflowStartRunReq{
				WorkflowId: workflow.Id,
				RunTrigger: domain.WorkflowTriggerTypeAuto,
			})
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *WorkflowService) StartRun(ctx context.Context, req *dtos.WorkflowStartRunReq) error {
	workflow, err := s.workflowRepo.GetById(ctx, req.WorkflowId)
	if err != nil {
		return err
	}

	if workflow.LastRunStatus == domain.WorkflowRunStatusTypePending || workflow.LastRunStatus == domain.WorkflowRunStatusTypeRunning {
		return errors.New("workflow is already pending or running")
	}

	run := &domain.WorkflowRun{
		WorkflowId: workflow.Id,
		Status:     domain.WorkflowRunStatusTypePending,
		Trigger:    req.RunTrigger,
		StartedAt:  time.Now(),
	}
	if resp, err := s.workflowRunRepo.Save(ctx, run); err != nil {
		return err
	} else {
		run = resp
	}

	s.dispatcher.Dispatch(&dispatcher.WorkflowWorkerData{
		WorkflowId:      workflow.Id,
		WorkflowContent: workflow.Content,
		RunId:           run.Id,
	})

	return nil
}

func (s *WorkflowService) CancelRun(ctx context.Context, req *dtos.WorkflowCancelRunReq) error {
	workflow, err := s.workflowRepo.GetById(ctx, req.WorkflowId)
	if err != nil {
		return err
	}

	workflowRun, err := s.workflowRunRepo.GetById(ctx, req.RunId)
	if err != nil {
		return err
	} else if workflowRun.WorkflowId != workflow.Id {
		return errors.New("workflow run not found")
	} else if workflowRun.Status != domain.WorkflowRunStatusTypePending && workflowRun.Status != domain.WorkflowRunStatusTypeRunning {
		return errors.New("workflow run is not pending or running")
	}

	s.dispatcher.Cancel(workflowRun.Id)

	return nil
}

func (s *WorkflowService) Shutdown(ctx context.Context) {
	s.dispatcher.Shutdown()
}
