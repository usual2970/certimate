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

type WorkflowRunRepository struct{}

func NewWorkflowRunRepository() *WorkflowRunRepository {
	return &WorkflowRunRepository{}
}

func (r *WorkflowRunRepository) GetById(ctx context.Context, id string) (*domain.WorkflowRun, error) {
	record, err := app.GetApp().FindRecordById(domain.CollectionNameWorkflowRun, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrRecordNotFound
		}
		return nil, err
	}

	return r.castRecordToModel(record)
}

func (r *WorkflowRunRepository) Save(ctx context.Context, workflowRun *domain.WorkflowRun) (*domain.WorkflowRun, error) {
	collection, err := app.GetApp().FindCollectionByNameOrId(domain.CollectionNameWorkflowRun)
	if err != nil {
		return workflowRun, err
	}

	var record *core.Record
	if workflowRun.Id == "" {
		record = core.NewRecord(collection)
	} else {
		record, err = app.GetApp().FindRecordById(collection, workflowRun.Id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return workflowRun, err
			}
			record = core.NewRecord(collection)
		}
	}

	err = app.GetApp().RunInTransaction(func(txApp core.App) error {
		record.Set("workflowId", workflowRun.WorkflowId)
		record.Set("trigger", string(workflowRun.Trigger))
		record.Set("status", string(workflowRun.Status))
		record.Set("startedAt", workflowRun.StartedAt)
		record.Set("endedAt", workflowRun.EndedAt)
		record.Set("detail", workflowRun.Detail)
		record.Set("error", workflowRun.Error)
		err = txApp.Save(record)
		if err != nil {
			return err
		}

		workflowRun.Id = record.Id
		workflowRun.CreatedAt = record.GetDateTime("created").Time()
		workflowRun.UpdatedAt = record.GetDateTime("updated").Time()

		// 事务级联更新所属工作流的最后运行记录
		workflowRecord, err := txApp.FindRecordById(domain.CollectionNameWorkflow, workflowRun.WorkflowId)
		if err != nil {
			return err
		} else if workflowRun.Id == workflowRecord.GetString("lastRunId") {
			workflowRecord.IgnoreUnchangedFields(true)
			workflowRecord.Set("lastRunStatus", record.GetString("status"))
			err = txApp.Save(workflowRecord)
			if err != nil {
				return err
			}
		} else if workflowRecord.GetDateTime("lastRunTime").Time().IsZero() || workflowRun.StartedAt.After(workflowRecord.GetDateTime("lastRunTime").Time()) {
			workflowRecord.IgnoreUnchangedFields(true)
			workflowRecord.Set("lastRunId", record.Id)
			workflowRecord.Set("lastRunStatus", record.GetString("status"))
			workflowRecord.Set("lastRunTime", record.GetString("startedAt"))
			err = txApp.Save(workflowRecord)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return workflowRun, err
	}

	return workflowRun, nil
}

func (r *WorkflowRunRepository) DeleteWhere(ctx context.Context, exprs ...dbx.Expression) (int, error) {
	records, err := app.GetApp().FindAllRecords(domain.CollectionNameWorkflowRun, exprs...)
	if err != nil {
		return 0, nil
	}

	var ret int
	var errs []error
	for _, record := range records {
		if err := app.GetApp().Delete(record); err != nil {
			errs = append(errs, err)
		} else {
			ret++
		}
	}

	if len(errs) > 0 {
		return ret, errors.Join(errs...)
	}

	return ret, nil
}

func (r *WorkflowRunRepository) castRecordToModel(record *core.Record) (*domain.WorkflowRun, error) {
	if record == nil {
		return nil, fmt.Errorf("record is nil")
	}

	detail := &domain.WorkflowNode{}
	if err := record.UnmarshalJSONField("detail", &detail); err != nil {
		return nil, err
	}

	workflowRun := &domain.WorkflowRun{
		Meta: domain.Meta{
			Id:        record.Id,
			CreatedAt: record.GetDateTime("created").Time(),
			UpdatedAt: record.GetDateTime("updated").Time(),
		},
		WorkflowId: record.GetString("workflowId"),
		Status:     domain.WorkflowRunStatusType(record.GetString("status")),
		Trigger:    domain.WorkflowTriggerType(record.GetString("trigger")),
		StartedAt:  record.GetDateTime("startedAt").Time(),
		EndedAt:    record.GetDateTime("endedAt").Time(),
		Detail:     detail,
		Error:      record.GetString("error"),
	}
	return workflowRun, nil
}
