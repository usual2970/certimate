package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/usual2970/certimate/internal/app"
	"github.com/usual2970/certimate/internal/domain"
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
		dbx.Params{"enabled": true, "trigger": domain.WorkflowTriggerTypeAuto},
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
		record, err = app.GetApp().FindRecordById(domain.CollectionNameWorkflow, workflow.Id)
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

func (r *WorkflowRepository) SaveRun(ctx context.Context, workflowRun *domain.WorkflowRun) (*domain.WorkflowRun, error) {
	collection, err := app.GetApp().FindCollectionByNameOrId(domain.CollectionNameWorkflowRun)
	if err != nil {
		return workflowRun, err
	}

	var workflowRunRecord *core.Record
	if workflowRun.Id == "" {
		workflowRunRecord = core.NewRecord(collection)
	} else {
		workflowRunRecord, err = app.GetApp().FindRecordById(domain.CollectionNameWorkflowRun, workflowRun.Id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return workflowRun, err
			}
			workflowRunRecord = core.NewRecord(collection)
		}
	}

	err = app.GetApp().RunInTransaction(func(txApp core.App) error {
		workflowRunRecord.Set("workflowId", workflowRun.WorkflowId)
		workflowRunRecord.Set("trigger", string(workflowRun.Trigger))
		workflowRunRecord.Set("status", string(workflowRun.Status))
		workflowRunRecord.Set("startedAt", workflowRun.StartedAt)
		workflowRunRecord.Set("endedAt", workflowRun.EndedAt)
		workflowRunRecord.Set("logs", workflowRun.Logs)
		workflowRunRecord.Set("error", workflowRun.Error)
		err = txApp.Save(workflowRunRecord)
		if err != nil {
			return err
		}

		workflowRecord, err := txApp.FindRecordById(domain.CollectionNameWorkflow, workflowRun.WorkflowId)
		if err != nil {
			return err
		}

		workflowRecord.IgnoreUnchangedFields(true)
		workflowRecord.Set("lastRunId", workflowRunRecord.Id)
		workflowRecord.Set("lastRunStatus", workflowRunRecord.GetString("status"))
		workflowRecord.Set("lastRunTime", workflowRunRecord.GetString("startedAt"))
		err = txApp.Save(workflowRecord)
		if err != nil {
			return err
		}

		workflowRun.Id = workflowRunRecord.Id
		workflowRun.CreatedAt = workflowRunRecord.GetDateTime("created").Time()
		workflowRun.UpdatedAt = workflowRunRecord.GetDateTime("updated").Time()

		return nil
	})
	if err != nil {
		return workflowRun, err
	}

	return workflowRun, nil
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
