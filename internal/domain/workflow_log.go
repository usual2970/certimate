package domain

import "strings"

const CollectionNameWorkflowLog = "workflow_logs"

type WorkflowLog struct {
	Meta
	WorkflowId string         `json:"workflowId" db:"workflowId"`
	RunId      string         `json:"workflorunIdwId" db:"runId"`
	NodeId     string         `json:"nodeId" db:"nodeId"`
	NodeName   string         `json:"nodeName" db:"nodeName"`
	Timestamp  int64          `json:"timestamp" db:"timestamp"` // 毫秒级时间戳
	Level      string         `json:"level" db:"level"`
	Message    string         `json:"message" db:"message"`
	Data       map[string]any `json:"data" db:"data"`
}

type WorkflowLogs []WorkflowLog

func (r WorkflowLogs) ErrorString() string {
	var builder strings.Builder
	for _, log := range r {
		if log.Level == "ERROR" {
			builder.WriteString(log.Message)
			builder.WriteString("\n")
		}
	}
	return strings.TrimSpace(builder.String())
}
