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

const defaultRoutines = 10

type workflowRunData struct {
	Workflow   *domain.Workflow
	RunTrigger domain.WorkflowTriggerType
}

type workflowRepository interface {
	ListEnabledAuto(ctx context.Context) ([]*domain.Workflow, error)
	GetById(ctx context.Context, id string) (*domain.Workflow, error)
	Save(ctx context.Context, workflow *domain.Workflow) (*domain.Workflow, error)
	SaveRun(ctx context.Context, workflowRun *domain.WorkflowRun) (*domain.WorkflowRun, error)
}

type WorkflowService struct {
	ch     chan *workflowRunData
	repo   workflowRepository
	wg     sync.WaitGroup
	cancel context.CancelFunc
}

func NewWorkflowService(repo workflowRepository) *WorkflowService {
	srv := &WorkflowService{
		repo: repo,
		ch:   make(chan *workflowRunData, 1),
	}

	ctx, cancel := context.WithCancel(context.Background())
	srv.cancel = cancel

	srv.wg.Add(defaultRoutines)
	for i := 0; i < defaultRoutines; i++ {
		go srv.run(ctx)
	}

	return srv
}

func (s *WorkflowService) InitSchedule(ctx context.Context) error {
	workflows, err := s.repo.ListEnabledAuto(ctx)
	if err != nil {
		return err
	}

	scheduler := app.GetScheduler()
	for _, workflow := range workflows {
		err := scheduler.Add(fmt.Sprintf("workflow#%s", workflow.Id), workflow.TriggerCron, func() {
			s.StartRun(ctx, &dtos.WorkflowStartRunReq{
				WorkflowId: workflow.Id,
				Trigger:    domain.WorkflowTriggerTypeAuto,
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
	workflow, err := s.repo.GetById(ctx, req.WorkflowId)
	if err != nil {
		app.GetLogger().Error("failed to get workflow", "id", req.WorkflowId, "err", err)
		return err
	}

	if workflow.LastRunStatus == domain.WorkflowRunStatusTypeRunning {
		return errors.New("workflow is running")
	}

	workflow.LastRunTime = time.Now()
	workflow.LastRunStatus = domain.WorkflowRunStatusTypePending
	workflow.LastRunId = ""
	if resp, err := s.repo.Save(ctx, workflow); err != nil {
		return err
	} else {
		workflow = resp
	}

	s.ch <- &workflowRunData{
		Workflow:   workflow,
		RunTrigger: req.Trigger,
	}

	return nil
}

func (s *WorkflowService) CancelRun(ctx context.Context, req *dtos.WorkflowCancelRunReq) error {
	// TODO: 取消运行，防止因为某些原因意外挂起（如进程被杀死）导致工作流一直处于 running 状态无法重新运行

	return errors.New("TODO: 尚未实现")
}

func (s *WorkflowService) Stop(ctx context.Context) {
	s.cancel()
	s.wg.Wait()
}

func (s *WorkflowService) run(ctx context.Context) {
	defer s.wg.Done()
	for {
		select {
		case data := <-s.ch:
			if err := s.runWithData(ctx, data); err != nil {
				app.GetLogger().Error("failed to run workflow", "id", data.Workflow.Id, "err", err)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (s *WorkflowService) runWithData(ctx context.Context, runData *workflowRunData) error {
	workflow := runData.Workflow
	run := &domain.WorkflowRun{
		WorkflowId: workflow.Id,
		Status:     domain.WorkflowRunStatusTypeRunning,
		Trigger:    runData.RunTrigger,
		StartedAt:  time.Now(),
	}
	if resp, err := s.repo.SaveRun(ctx, run); err != nil {
		return err
	} else {
		run = resp
	}

	processor := processor.NewWorkflowProcessor(workflow)
	if runErr := processor.Run(ctx); runErr != nil {
		run.Status = domain.WorkflowRunStatusTypeFailed
		run.EndedAt = time.Now()
		run.Logs = processor.GetRunLogs()
		run.Error = runErr.Error()
		if _, err := s.repo.SaveRun(ctx, run); err != nil {
			app.GetLogger().Error("failed to save workflow run", "err", err)
		}

		return fmt.Errorf("failed to run workflow: %w", runErr)
	}

	run.EndedAt = time.Now()
	run.Logs = processor.GetRunLogs()
	run.Error = domain.WorkflowRunLogs(run.Logs).ErrorString()
	if run.Error == "" {
		run.Status = domain.WorkflowRunStatusTypeSucceeded
	} else {
		run.Status = domain.WorkflowRunStatusTypeFailed
	}
	if _, err := s.repo.SaveRun(ctx, run); err != nil {
		app.GetLogger().Error("failed to save workflow run", "err", err)
		return err
	}

	return nil
}
