package workflow

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/usual2970/certimate/internal/app"
	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/domain/dtos"
	processor "github.com/usual2970/certimate/internal/workflow/processor"
)

const defaultRoutines = 16

type workflowRunData struct {
	WorkflowId      string
	WorkflowContent *domain.WorkflowNode
	RunTrigger      domain.WorkflowTriggerType
}

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
	ch     chan *workflowRunData
	wg     sync.WaitGroup
	cancel context.CancelFunc

	workflowRepo    workflowRepository
	workflowRunRepo workflowRunRepository
}

func NewWorkflowService(workflowRepo workflowRepository, workflowRunRepo workflowRunRepository) *WorkflowService {
	ctx, cancel := context.WithCancel(context.Background())

	srv := &WorkflowService{
		ch:     make(chan *workflowRunData, 1),
		cancel: cancel,

		workflowRepo:    workflowRepo,
		workflowRunRepo: workflowRunRepo,
	}

	srv.wg.Add(defaultRoutines)
	for i := 0; i < defaultRoutines; i++ {
		go srv.startRun(ctx)
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
			app.GetLogger().Error("failed to add schedule", "err", err)
			return err
		}
	}

	return nil
}

func (s *WorkflowService) StartRun(ctx context.Context, req *dtos.WorkflowStartRunReq) error {
	workflow, err := s.workflowRepo.GetById(ctx, req.WorkflowId)
	if err != nil {
		app.GetLogger().Error("failed to get workflow", "id", req.WorkflowId, "err", err)
		return err
	}

	if workflow.LastRunStatus == domain.WorkflowRunStatusTypePending || workflow.LastRunStatus == domain.WorkflowRunStatusTypeRunning {
		return errors.New("workflow is already pending or running")
	}

	s.ch <- &workflowRunData{
		WorkflowId:      workflow.Id,
		WorkflowContent: workflow.Content,
		RunTrigger:      req.RunTrigger,
	}

	return nil
}

func (s *WorkflowService) CancelRun(ctx context.Context, req *dtos.WorkflowCancelRunReq) error {
	workflow, err := s.workflowRepo.GetById(ctx, req.WorkflowId)
	if err != nil {
		app.GetLogger().Error("failed to get workflow", "id", req.WorkflowId, "err", err)
		return err
	}

	workflowRun, err := s.workflowRunRepo.GetById(ctx, req.RunId)
	if err != nil {
		app.GetLogger().Error("failed to get workflow run", "id", req.RunId, "err", err)
		return err
	} else if workflowRun.WorkflowId != workflow.Id {
		return errors.New("workflow run not found")
	} else if workflowRun.Status != domain.WorkflowRunStatusTypePending && workflowRun.Status != domain.WorkflowRunStatusTypeRunning {
		return errors.New("workflow run is not pending or running")
	}

	// TODO: 取消运行，防止因为某些原因意外挂起（如进程被杀死）导致工作流一直处于 running 状态无法重新运行
	// workflowRun.Status = domain.WorkflowRunStatusTypeCanceled
	// workflowRun.EndedAt = time.Now()
	// if _, err := s.workflowRunRepo.Save(ctx, workflowRun); err != nil {
	// 	return err
	// }

	// return nil

	return errors.New("TODO: 尚未实现")
}

func (s *WorkflowService) Stop(ctx context.Context) {
	s.cancel()
	s.wg.Wait()
}

func (s *WorkflowService) startRun(ctx context.Context) {
	defer s.wg.Done()

	for {
		select {
		case data := <-s.ch:
			if err := s.startRunWithData(ctx, data); err != nil {
				app.GetLogger().Error("failed to run workflow", "id", data.WorkflowId, "err", err)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (s *WorkflowService) startRunWithData(ctx context.Context, data *workflowRunData) error {
	run := &domain.WorkflowRun{
		WorkflowId: data.WorkflowId,
		Status:     domain.WorkflowRunStatusTypeRunning,
		Trigger:    data.RunTrigger,
		StartedAt:  time.Now(),
	}
	if resp, err := s.workflowRunRepo.Save(ctx, run); err != nil {
		return err
	} else {
		run = resp
	}

	processor := processor.NewWorkflowProcessor(data.WorkflowId, data.WorkflowContent, run.Id)
	if runErr := processor.Process(ctx); runErr != nil {
		run.Status = domain.WorkflowRunStatusTypeFailed
		run.EndedAt = time.Now()
		run.Logs = processor.GetLogs()
		run.Error = runErr.Error()
		if _, err := s.workflowRunRepo.Save(ctx, run); err != nil {
			app.GetLogger().Error("failed to save workflow run", "err", err)
		}

		return fmt.Errorf("failed to run workflow: %w", runErr)
	}

	run.EndedAt = time.Now()
	run.Logs = processor.GetLogs()
	run.Error = domain.WorkflowRunLogs(run.Logs).ErrorString()
	if run.Error == "" {
		run.Status = domain.WorkflowRunStatusTypeSucceeded
	} else {
		run.Status = domain.WorkflowRunStatusTypeFailed
	}
	if _, err := s.workflowRunRepo.Save(ctx, run); err != nil {
		app.GetLogger().Error("failed to save workflow run", "err", err)
		return err
	}

	return nil
}
