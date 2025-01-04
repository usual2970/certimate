package workflow

import (
	"context"
	"fmt"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"

	"github.com/usual2970/certimate/internal/app"
	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/repository"
)

const tableName = "workflow"

func RegisterEvents() error {
	app := app.GetApp()

	app.OnRecordAfterCreateRequest(tableName).Add(func(e *core.RecordCreateEvent) error {
		return update(e.HttpContext.Request().Context(), e.Record)
	})

	app.OnRecordAfterUpdateRequest(tableName).Add(func(e *core.RecordUpdateEvent) error {
		return update(e.HttpContext.Request().Context(), e.Record)
	})

	app.OnRecordAfterDeleteRequest(tableName).Add(func(e *core.RecordDeleteEvent) error {
		return delete(e.HttpContext.Request().Context(), e.Record)
	})

	return nil
}

func update(ctx context.Context, record *models.Record) error {
	scheduler := app.GetScheduler()

	// 向数据库插入/更新时，同时更新定时任务
	workflowId := record.GetId()
	enabled := record.GetBool("enabled")
	trigger := record.GetString("trigger")

	// 如果是手动触发或未启用，移除定时任务
	if !enabled || trigger == string(domain.WorkflowTriggerTypeManual) {
		scheduler.Remove(fmt.Sprintf("workflow#%s", workflowId))
		scheduler.Start()
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

	scheduler.Start()

	return nil
}

func delete(_ context.Context, record *models.Record) error {
	scheduler := app.GetScheduler()

	// 从数据库删除时，同时移除定时任务
	workflowId := record.GetId()
	scheduler.Remove(fmt.Sprintf("workflow#%s", workflowId))
	scheduler.Start()

	return nil
}
