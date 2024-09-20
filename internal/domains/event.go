package domains

import (
	"certimate/internal/utils/app"

	"github.com/pocketbase/pocketbase/core"
)

const tableName = "domains"

func AddEvent() error {
	app := app.GetApp()

	app.OnRecordAfterCreateRequest(tableName).Add(func(e *core.RecordCreateEvent) error {
		return create(e.HttpContext.Request().Context(), e.Record)
	})

	app.OnRecordAfterUpdateRequest(tableName).Add(func(e *core.RecordUpdateEvent) error {
		return update(e.HttpContext.Request().Context(), e.Record)
	})

	app.OnRecordAfterDeleteRequest(tableName).Add(func(e *core.RecordDeleteEvent) error {
		return delete(e.HttpContext.Request().Context(), e.Record)
	})

	return nil
}
