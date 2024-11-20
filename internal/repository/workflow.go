package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/pocketbase/pocketbase/models"
	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/utils/app"
)

type WorkflowRepository struct{}

func NewWorkflowRepository() *WorkflowRepository {
	return &WorkflowRepository{}
}

func (w *WorkflowRepository) SaveRunLog(ctx context.Context, log *domain.WorkflowRunLog) error {
	collection, err := app.GetApp().Dao().FindCollectionByNameOrId("workflow_run_log")
	if err != nil {
		return err
	}
	record := models.NewRecord(collection)

	record.Set("workflow", log.Workflow)
	record.Set("log", log.Log)
	record.Set("succeed", log.Succeed)
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
			Id:      record.GetId(),
			Created: record.GetTime("created"),
			Updated: record.GetTime("updated"),
		},
		Name:        record.GetString("name"),
		Description: record.GetString("description"),
		Type:        record.GetString("type"),
		Enabled:     record.GetBool("enabled"),
		HasDraft:    record.GetBool("hasDraft"),

		Content: content,
		Draft:   draft,
	}

	return workflow, nil
}
