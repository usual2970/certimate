package domain

import "time"

type Workflow struct {
	Id          string        `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Type        string        `json:"type"`
	Content     *WorkflowNode `json:"content"`
	Draft       *WorkflowNode `json:"draft"`
	Enabled     bool          `json:"enabled"`
	HasDraft    bool          `json:"hasDraft"`
	Created     time.Time     `json:"created"`
	Updated     time.Time     `json:"updated"`
}

type WorkflowNode struct {
	Id     string           `json:"id"`
	Name   string           `json:"name"`
	Next   *WorkflowNode    `json:"next"`
	Config map[string]any   `json:"config"`
	Input  []WorkflowNodeIo `json:"input"`
	Output []WorkflowNodeIo `json:"output"`

	Validated bool   `json:"validated"`
	Type      string `json:"type"`

	Branches []WorkflowNode `json:"branches"`
}

type WorkflowNodeIo struct {
	Label    string `json:"label"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Required bool   `json:"required"`
}

type WorkflowRunReq struct {
	Id string `json:"id"`
}
