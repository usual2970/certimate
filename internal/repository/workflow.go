package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/certimate-go/certimate/internal/app"
	"github.com/certimate-go/certimate/internal/domain"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

type WorkflowRepository struct{}

func NewWorkflowRepository() *WorkflowRepository {
	return &WorkflowRepository{}
}

func (r *WorkflowRepository) ListEnabledAuto(ctx context.Context) ([]*domain.Workflow, error) {
	records, err := app.GetApp().FindRecordsByFilter(
		domain.CollectionNameWorkflow,
		"enabled={:enabled} && trigger={:trigger}",
		"-created",
		0, 0,
		dbx.Params{"enabled": true, "trigger": string(domain.WorkflowTriggerTypeAuto)},
	)
	if err != nil {
		return nil, err
	}

	workflows := make([]*domain.Workflow, 0)
	for _, record := range records {
		workflow, err := r.castRecordToModel(record)
		if err != nil {
			return nil, err
		}

		workflows = append(workflows, workflow)
	}

	return workflows, nil
}

func (r *WorkflowRepository) GetById(ctx context.Context, id string) (*domain.Workflow, error) {
	record, err := app.GetApp().FindRecordById(domain.CollectionNameWorkflow, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrRecordNotFound
		}
		return nil, err
	}

	return r.castRecordToModel(record)
}

func (r *WorkflowRepository) Save(ctx context.Context, workflow *domain.Workflow) (*domain.Workflow, error) {
	collection, err := app.GetApp().FindCollectionByNameOrId(domain.CollectionNameWorkflow)
	if err != nil {
		return workflow, err
	}

	var record *core.Record
	if workflow.Id == "" {
		record = core.NewRecord(collection)
	} else {
		record, err = app.GetApp().FindRecordById(collection, workflow.Id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return workflow, domain.ErrRecordNotFound
			}
			return workflow, err
		}
	}

	record.Set("name", workflow.Name)
	record.Set("description", workflow.Description)
	record.Set("trigger", string(workflow.Trigger))
	record.Set("triggerCron", workflow.TriggerCron)
	record.Set("enabled", workflow.Enabled)
	record.Set("content", workflow.Content)
	record.Set("draft", workflow.Draft)
	record.Set("hasDraft", workflow.HasDraft)
	record.Set("lastRunId", workflow.LastRunId)
	record.Set("lastRunStatus", string(workflow.LastRunStatus))
	record.Set("lastRunTime", workflow.LastRunTime)
	if err := app.GetApp().Save(record); err != nil {
		return workflow, err
	}

	workflow.Id = record.Id
	workflow.CreatedAt = record.GetDateTime("created").Time()
	workflow.UpdatedAt = record.GetDateTime("updated").Time()
	return workflow, nil
}

func (r *WorkflowRepository) castRecordToModel(record *core.Record) (*domain.Workflow, error) {
	if record == nil {
		return nil, fmt.Errorf("record is nil")
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
			Id:        record.Id,
			CreatedAt: record.GetDateTime("created").Time(),
			UpdatedAt: record.GetDateTime("updated").Time(),
		},
		Name:          record.GetString("name"),
		Description:   record.GetString("description"),
		Trigger:       domain.WorkflowTriggerType(record.GetString("trigger")),
		TriggerCron:   record.GetString("triggerCron"),
		Enabled:       record.GetBool("enabled"),
		Content:       content,
		Draft:         draft,
		HasDraft:      record.GetBool("hasDraft"),
		LastRunId:     record.GetString("lastRunId"),
		LastRunStatus: domain.WorkflowRunStatusType(record.GetString("lastRunStatus")),
		LastRunTime:   record.GetDateTime("lastRunTime").Time(),
	}
	return workflow, nil
}
