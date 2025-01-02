package domain

const WorkflowOutputCertificate = "certificate"

type WorkflowOutput struct {
	Meta
	Workflow string           `json:"workflow"`
	NodeId   string           `json:"nodeId"`
	Node     *WorkflowNode    `json:"node"`
	Output   []WorkflowNodeIO `json:"output"`
	Succeed  bool             `json:"succeed"`
}
