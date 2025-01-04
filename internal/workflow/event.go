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

func AddEvent() error {
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

func delete(_ context.Context, record *models.Record) error {
	id := record.Id
	scheduler := app.GetScheduler()
	scheduler.Remove(id)
	scheduler.Start()
	return nil
}

func update(ctx context.Context, record *models.Record) error {
	// 是不是自动
	// 是不是 enabled

	id := record.Id
	enabled := record.GetBool("enabled")
	trigger := record.GetString("trigger")

	scheduler := app.GetScheduler()
	if !enabled || trigger == domain.WorkflowTriggerManual {
		scheduler.Remove(id)
		scheduler.Start()
		return nil
	}

	err := scheduler.Add(id, record.GetString("triggerCron"), func() {
		NewWorkflowService(repository.NewWorkflowRepository()).Run(ctx, &domain.WorkflowRunReq{
			Id: id,
		})
	})
	if err != nil {
		app.GetApp().Logger().Error("add cron job failed", "err", err)
		return fmt.Errorf("add cron job failed: %w", err)
	}
	app.GetApp().Logger().Error("add cron job failed", "subjectAltNames", record.GetString("subjectAltNames"))

	scheduler.Start()
	return nil
}
