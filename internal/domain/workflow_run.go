package domain

import (
	"time"
)

const CollectionNameWorkflowRun = "workflow_run"

type WorkflowRun struct {
	Meta
	WorkflowId string                `json:"workflowId" db:"workflowId"`
	Status     WorkflowRunStatusType `json:"status" db:"status"`
	Trigger    WorkflowTriggerType   `json:"trigger" db:"trigger"`
	StartedAt  time.Time             `json:"startedAt" db:"startedAt"`
	EndedAt    time.Time             `json:"endedAt" db:"endedAt"`
	Detail     *WorkflowNode         `json:"detail" db:"detail"`
	Error      string                `json:"error" db:"error"`
}

type WorkflowRunStatusType string

const (
	WorkflowRunStatusTypePending   WorkflowRunStatusType = "pending"
	WorkflowRunStatusTypeRunning   WorkflowRunStatusType = "running"
	WorkflowRunStatusTypeSucceeded WorkflowRunStatusType = "succeeded"
	WorkflowRunStatusTypeFailed    WorkflowRunStatusType = "failed"
	WorkflowRunStatusTypeCanceled  WorkflowRunStatusType = "canceled"
)
