package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/utils/app"
)

type WorkflowRepository struct{}

func NewWorkflowRepository() *WorkflowRepository {
	return &WorkflowRepository{}
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
