package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/models"
	"github.com/usual2970/certimate/internal/app"
	"github.com/usual2970/certimate/internal/domain"
)

type WorkflowRepository struct{}

func NewWorkflowRepository() *WorkflowRepository {
	return &WorkflowRepository{}
}

func (w *WorkflowRepository) ListEnabledAuto(ctx context.Context) ([]domain.Workflow, error) {
	records, err := app.GetApp().Dao().FindRecordsByFilter(
		"workflow",
		"enabled={:enabled} && trigger={:trigger}",
		"-created", 1000, 0, dbx.Params{"enabled": true, "trigger": domain.WorkflowTriggerAuto},
	)
	if err != nil {
		return nil, err
	}
	rs := make([]domain.Workflow, 0)
	for _, record := range records {
		workflow, err := record2Workflow(record)
		if err != nil {
			return nil, err
		}
		rs = append(rs, *workflow)
	}
	return rs, nil
}

func (w *WorkflowRepository) SaveRunLog(ctx context.Context, log *domain.WorkflowRun) error {
	collection, err := app.GetApp().Dao().FindCollectionByNameOrId("workflow_run")
	if err != nil {
		return err
	}
	record := models.NewRecord(collection)

	record.Set("workflowId", log.WorkflowId)
	record.Set("trigger", log.Trigger)
	record.Set("startedAt", log.StartedAt)
	record.Set("completedAt", log.CompletedAt)
	record.Set("logs", log.Logs)
	record.Set("succeeded", log.Succeeded)
	record.Set("error", log.Error)

	return app.GetApp().Dao().SaveRecord(record)
}

func (w *WorkflowRepository) Get(ctx context.Context, id string) (*domain.Workflow, error) {
	record, err := app.GetApp().Dao().FindRecordById("workflow", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrRecordNotFound
		}
		return nil, err
	}

	return record2Workflow(record)
}

func record2Workflow(record *models.Record) (*domain.Workflow, error) {
	content := &domain.WorkflowNode{}
	if err := record.UnmarshalJSONField("content", content); err != nil {
		return nil, err
	}

	draft := &domain.WorkflowNode{}
	if err := record.UnmarshalJSONField("draft", draft); err != nil {
		return nil, err
	}

	workflow := &domain.Workflow{
		Meta: domain.Meta{
			Id:        record.GetId(),
			CreatedAt: record.GetTime("created"),
			UpdatedAt: record.GetTime("updated"),
		},
		Name:        record.GetString("name"),
		Description: record.GetString("description"),
		Trigger:     record.GetString("trigger"),
		TriggerCron: record.GetString("triggerCron"),
		Enabled:     record.GetBool("enabled"),
		Content:     content,
		Draft:       draft,
		HasDraft:    record.GetBool("hasDraft"),
	}

	return workflow, nil
}
