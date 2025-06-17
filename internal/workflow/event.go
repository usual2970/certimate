package workflow

import (
	"context"
	"fmt"

	"github.com/pocketbase/pocketbase/core"

	"github.com/certimate-go/certimate/internal/app"
	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/internal/domain/dtos"
	"github.com/certimate-go/certimate/internal/repository"
)

func Register() {
	app := app.GetApp()
	app.OnRecordCreateRequest(domain.CollectionNameWorkflow).BindFunc(func(e *core.RecordRequestEvent) error {
		if err := e.Next(); err != nil {
			return err
		}

		if err := onWorkflowRecordCreateOrUpdate(e.Request.Context(), e.Record); err != nil {
			return err
		}

		return nil
	})
	app.OnRecordUpdateRequest(domain.CollectionNameWorkflow).BindFunc(func(e *core.RecordRequestEvent) error {
		if err := e.Next(); err != nil {
			return err
		}

		if err := onWorkflowRecordCreateOrUpdate(e.Request.Context(), e.Record); err != nil {
			return err
		}

		return nil
	})
	app.OnRecordDeleteRequest(domain.CollectionNameWorkflow).BindFunc(func(e *core.RecordRequestEvent) error {
		if err := e.Next(); err != nil {
			return err
		}

		if err := onWorkflowRecordDelete(e.Request.Context(), e.Record); err != nil {
			return err
		}

		return nil
	})
}

func onWorkflowRecordCreateOrUpdate(ctx context.Context, record *core.Record) error {
	scheduler := app.GetScheduler()

	// 向数据库插入/更新时，同时更新定时任务
	workflowId := record.Id
	enabled := record.GetBool("enabled")
	trigger := record.GetString("trigger")

	// 如果是手动触发或未启用，移除定时任务
	if !enabled || trigger == string(domain.WorkflowTriggerTypeManual) {
		scheduler.Remove(fmt.Sprintf("workflow#%s", workflowId))
		return nil
	}

	// 反之，重新添加定时任务
	err := scheduler.Add(fmt.Sprintf("workflow#%s", workflowId), record.GetString("triggerCron"), func() {
		workflowSrv := NewWorkflowService(repository.NewWorkflowRepository(), repository.NewWorkflowRunRepository(), repository.NewSettingsRepository())
		workflowSrv.StartRun(ctx, &dtos.WorkflowStartRunReq{
			WorkflowId: workflowId,
			RunTrigger: domain.WorkflowTriggerTypeAuto,
		})
	})
	if err != nil {
		return fmt.Errorf("add cron job failed: %w", err)
	}

	return nil
}

func onWorkflowRecordDelete(_ context.Context, record *core.Record) error {
	scheduler := app.GetScheduler()

	// 从数据库删除时，同时移除定时任务
	workflowId := record.Id
	scheduler.Remove(fmt.Sprintf("workflow#%s", workflowId))

	return nil
}
