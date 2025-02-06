package domain

const CollectionNameWorkflowOutput = "workflow_output"

type WorkflowOutput struct {
	Meta
	WorkflowId string           `json:"workflowId" db:"workflow"`
	RunId      string           `json:"runId" db:"runId"`
	NodeId     string           `json:"nodeId" db:"nodeId"`
	Node       *WorkflowNode    `json:"node" db:"node"`
	Outputs    []WorkflowNodeIO `json:"outputs" db:"outputs"`
	Succeeded  bool             `json:"succeeded" db:"succeeded"`
}
