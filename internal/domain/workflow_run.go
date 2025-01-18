package domain

import "time"

const CollectionNameWorkflowRun = "workflow_run"

type WorkflowRun struct {
	Meta
	WorkflowId string                `json:"workflowId" db:"workflowId"`
	Status     WorkflowRunStatusType `json:"status" db:"status"`
	Trigger    WorkflowTriggerType   `json:"trigger" db:"trigger"`
	StartedAt  time.Time             `json:"startedAt" db:"startedAt"`
	EndedAt    time.Time             `json:"endedAt" db:"endedAt"`
	Logs       []WorkflowRunLog      `json:"logs" db:"logs"`
	Error      string                `json:"error" db:"error"`
}

type WorkflowRunStatusType string

const (
	WorkflowRunStatusTypePending   WorkflowRunStatusType = "pending"
	WorkflowRunStatusTypeRunning   WorkflowRunStatusType = "running"
	WorkflowRunStatusTypeSucceeded WorkflowRunStatusType = "succeeded"
	WorkflowRunStatusTypeFailed    WorkflowRunStatusType = "failed"
)

type WorkflowRunLog struct {
	NodeId   string                 `json:"nodeId"`
	NodeName string                 `json:"nodeName"`
	Error    string                 `json:"error"`
	Outputs  []WorkflowRunLogOutput `json:"outputs"`
}

type WorkflowRunLogOutput struct {
	Time    string `json:"time"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Error   string `json:"error"`
}

type WorkflowRunLogs []WorkflowRunLog

func (r WorkflowRunLogs) FirstError() string {
	for _, log := range r {
		if log.Error != "" {
			return log.Error
		}
	}

	return ""
}
