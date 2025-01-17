package workflow

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/usual2970/certimate/internal/app"
	"github.com/usual2970/certimate/internal/domain"
	nodeprocessor "github.com/usual2970/certimate/internal/workflow/node-processor"
)

const defaultRoutines = 10

type workflowRunData struct {
	Workflow *domain.Workflow
	Options  *domain.WorkflowRunReq
}

type workflowRepository interface {
	ListEnabledAuto(ctx context.Context) ([]*domain.Workflow, error)
	GetById(ctx context.Context, id string) (*domain.Workflow, error)
	Save(ctx context.Context, workflow *domain.Workflow) error
	SaveRun(ctx context.Context, workflowRun *domain.WorkflowRun) error
}

type WorkflowService struct {
	ch     chan *workflowRunData
	repo   workflowRepository
	wg     sync.WaitGroup
	cancel context.CancelFunc
}

func NewWorkflowService(repo workflowRepository) *WorkflowService {
	rs := &WorkflowService{
		repo: repo,
		ch:   make(chan *workflowRunData, 1),
	}

	ctx, cancel := context.WithCancel(context.Background())
	rs.cancel = cancel

	rs.wg.Add(defaultRoutines)
	for i := 0; i < defaultRoutines; i++ {
		go rs.process(ctx)
	}

	return rs
}

func (s *WorkflowService) process(ctx context.Context) {
	defer s.wg.Done()
	for {
		select {
		case data := <-s.ch:
			// 执行
			if err := s.run(ctx, data); err != nil {
				app.GetLogger().Error("failed to run workflow", "id", data.Workflow.Id, "err", err)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (s *WorkflowService) InitSchedule(ctx context.Context) error {
	workflows, err := s.repo.ListEnabledAuto(ctx)
	if err != nil {
		return err
	}

	scheduler := app.GetScheduler()
	for _, workflow := range workflows {
		err := scheduler.Add(fmt.Sprintf("workflow#%s", workflow.Id), workflow.TriggerCron, func() {
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

	return nil
}

func (s *WorkflowService) Run(ctx context.Context, req *domain.WorkflowRunReq) error {
	// 查询
	workflow, err := s.repo.GetById(ctx, req.WorkflowId)
	if err != nil {
		app.GetLogger().Error("failed to get workflow", "id", req.WorkflowId, "err", err)
		return err
	}

	if workflow.LastRunStatus == domain.WorkflowRunStatusTypeRunning {
		return errors.New("workflow is running")
	}

	// set last run
	workflow.LastRunTime = time.Now()
	workflow.LastRunStatus = domain.WorkflowRunStatusTypeRunning
	workflow.LastRunId = ""

	if err := s.repo.Save(ctx, workflow); err != nil {
		return err
	}

	s.ch <- &workflowRunData{
		Workflow: workflow,
		Options:  req,
	}

	return nil
}

func (s *WorkflowService) run(ctx context.Context, runData *workflowRunData) error {
	// 执行
	workflow := runData.Workflow
	options := runData.Options

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

	// 保存日志
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

func (s *WorkflowService) Stop(ctx context.Context) {
	s.cancel()
	s.wg.Wait()
}
