package domain

type WorkflowRunLog struct {
	Meta
	WorkflowId string   `json:"workflowId" db:"workflowId"`
	Logs       []RunLog `json:"logs" db:"logs"`
	Succeeded  bool     `json:"succeeded" db:"succeeded"`
	Error      string   `json:"error" db:"error"`
}

type RunLogOutput struct {
	Time    string `json:"time"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Error   string `json:"error"`
}

type RunLog struct {
	NodeId   string         `json:"nodeId"`
	NodeName string         `json:"nodeName"`
	Error    string         `json:"error"`
	Outputs  []RunLogOutput `json:"outputs"`
}

type RunLogs []RunLog

func (r RunLogs) Error() string {
	for _, log := range r {
		if log.Error != "" {
			return log.Error
		}
	}
	return ""
}
