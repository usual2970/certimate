package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
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
		"-created", 1000, 0, dbx.Params{"enabled": true, "trigger": domain.WorkflowTriggerTypeAuto},
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

func (w *WorkflowRepository) SaveRun(ctx context.Context, run *domain.WorkflowRun) error {
	collection, err := app.GetApp().Dao().FindCollectionByNameOrId("workflow_run")
	if err != nil {
		return err
	}

	err = app.GetApp().Dao().RunInTransaction(func(txDao *daos.Dao) error {
		record := models.NewRecord(collection)
		record.Set("workflowId", run.WorkflowId)
		record.Set("trigger", string(run.Trigger))
		record.Set("status", string(run.Status))
		record.Set("startedAt", run.StartedAt)
		record.Set("endedAt", run.EndedAt)
		record.Set("logs", run.Logs)
		record.Set("error", run.Error)
		err = txDao.SaveRecord(record)
		if err != nil {
			return err
		}

		_, err = txDao.DB().Update(
			"workflow",
			dbx.Params{
				"lastRunId":     record.GetId(),
				"lastRunStatus": record.GetString("status"),
				"lastRunTime":   record.GetString("startedAt"),
			},
			dbx.NewExp("id={:id}", dbx.Params{"id": run.WorkflowId}),
		).Execute()
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (w *WorkflowRepository) GetById(ctx context.Context, id string) (*domain.Workflow, error) {
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
			CreatedAt: record.GetCreated().Time(),
			UpdatedAt: record.GetUpdated().Time(),
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
