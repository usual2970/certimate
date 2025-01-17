package workflow

import (
	"context"
	"fmt"

	"github.com/pocketbase/pocketbase/core"

	"github.com/usual2970/certimate/internal/app"
	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/repository"
)

func Register() {
	const tableName = "workflow"

	app := app.GetApp()
	app.OnRecordCreateRequest(tableName).BindFunc(func(e *core.RecordRequestEvent) error {
		if err := e.Next(); err != nil {
			return err
		}

		if err := update(e.Request.Context(), e.Record); err != nil {
			return err
		}

		return nil
	})
	app.OnRecordUpdateRequest(tableName).BindFunc(func(e *core.RecordRequestEvent) error {
		if err := e.Next(); err != nil {
			return err
		}

		if err := update(e.Request.Context(), e.Record); err != nil {
			return err
		}

		return nil
	})
	app.OnRecordDeleteRequest(tableName).BindFunc(func(e *core.RecordRequestEvent) error {
		if err := e.Next(); err != nil {
			return err
		}

		if err := delete(e.Request.Context(), e.Record); err != nil {
			return err
		}

		return nil
	})
}

func update(ctx context.Context, record *core.Record) error {
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
		NewWorkflowService(repository.NewWorkflowRepository()).Run(ctx, &domain.WorkflowRunReq{
			WorkflowId: workflowId,
			Trigger:    domain.WorkflowTriggerTypeAuto,
		})
	})
	if err != nil {
		app.GetLogger().Error("add cron job failed", "err", err)
		return fmt.Errorf("add cron job failed: %w", err)
	}

	return nil
}

func delete(_ context.Context, record *core.Record) error {
	scheduler := app.GetScheduler()

	// 从数据库删除时，同时移除定时任务
	workflowId := record.Id
	scheduler.Remove(fmt.Sprintf("workflow#%s", workflowId))

	return nil
}
