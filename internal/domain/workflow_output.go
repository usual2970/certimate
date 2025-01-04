package domain

type WorkflowOutput struct {
	Meta
	WorkflowId string           `json:"workflowId" db:"workflow"`
	NodeId     string           `json:"nodeId" db:"nodeId"`
	Node       *WorkflowNode    `json:"node" db:"node"`
	Outputs    []WorkflowNodeIO `json:"outputs" db:"outputs"`
	Succeeded  bool             `json:"succeeded" db:"succeeded"`
}

const WORKFLOW_OUTPUT_CERTIFICATE = "certificate"
