package domain

import "time"

type WorkflowRun struct {
	Meta
	WorkflowId  string           `json:"workflowId" db:"workflowId"`
	Trigger     string           `json:"trigger" db:"trigger"`
	StartedAt   time.Time        `json:"startedAt" db:"startedAt"`
	CompletedAt time.Time        `json:"completedAt" db:"completedAt"`
	Logs        []WorkflowRunLog `json:"logs" db:"logs"`
	Succeeded   bool             `json:"succeeded" db:"succeeded"`
	Error       string           `json:"error" db:"error"`
}

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
