package domain

import (
	"strings"
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
	Logs       []WorkflowRunLog      `json:"logs" db:"logs"`
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

type WorkflowRunLog struct {
	NodeId   string                 `json:"nodeId"`
	NodeName string                 `json:"nodeName"`
	Records  []WorkflowRunLogRecord `json:"records"`
	Error    string                 `json:"error"`
}

type WorkflowRunLogRecord struct {
	Time    string              `json:"time"`
	Level   WorkflowRunLogLevel `json:"level"`
	Content string              `json:"content"`
	Error   string              `json:"error"`
}

type WorkflowRunLogLevel string

const (
	WorkflowRunLogLevelDebug WorkflowRunLogLevel = "DEBUG"
	WorkflowRunLogLevelInfo  WorkflowRunLogLevel = "INFO"
	WorkflowRunLogLevelWarn  WorkflowRunLogLevel = "WARN"
	WorkflowRunLogLevelError WorkflowRunLogLevel = "ERROR"
)

type WorkflowRunLogs []WorkflowRunLog

func (r WorkflowRunLogs) ErrorString() string {
	var builder strings.Builder
	for _, log := range r {
		if log.Error != "" {
			builder.WriteString(log.Error)
			builder.WriteString("\n")
		}
	}
	return builder.String()
}
