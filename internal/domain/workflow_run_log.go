package domain

type RunLogOutput struct {
	Time    string `json:"time"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Error   string `json:"error"`
}

type RunLog struct {
	NodeName string         `json:"nodeName"`
	Error    string         `json:"error"`
	Outputs  []RunLogOutput `json:"outputs"`
}

type WorkflowRunLog struct {
	Meta
	Workflow string   `json:"workflow"`
	Log      []RunLog `json:"log"`
	Succeed  bool     `json:"succeed"`
	Error    string   `json:"error"`
}
