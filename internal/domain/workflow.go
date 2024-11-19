package domain

import (
	"fmt"
	"strconv"
)

const (
	WorkflowNodeTypeStart     = "start"
	WorkflowNodeTypeEnd       = "end"
	WorkflowNodeTypeApply     = "apply"
	WorkflowNodeTypeDeploy    = "deploy"
	WorkflowNodeTypeNotify    = "notify"
	WorkflowNodeTypeBranch    = "branch"
	WorkflowNodeTypeCondition = "condition"
)

type Workflow struct {
	Meta
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Type        string        `json:"type"`
	Content     *WorkflowNode `json:"content"`
	Draft       *WorkflowNode `json:"draft"`
	Enabled     bool          `json:"enabled"`
	HasDraft    bool          `json:"hasDraft"`
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

func (n *WorkflowNode) GetConfigString(key string) string {
	if v, ok := n.Config[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func (n *WorkflowNode) GetConfigBool(key string) bool {
	if v, ok := n.Config[key]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return false
}

func (n *WorkflowNode) GetConfigInt64(key string) int64 {
	// 先转成字符串，再转成 int64
	if v, ok := n.Config[key]; ok {
		temp := fmt.Sprintf("%v", v)
		if i, err := strconv.ParseInt(temp, 10, 64); err == nil {
			return i
		}
	}

	return 0
}

type WorkflowNodeIo struct {
	Label         string                      `json:"label"`
	Name          string                      `json:"name"`
	Type          string                      `json:"type"`
	Required      bool                        `json:"required"`
	Value         any                         `json:"value"`
	ValueSelector WorkflowNodeIoValueSelector `json:"valueSelector"`
}

type WorkflowNodeIoValueSelector struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type WorkflowRunReq struct {
	Id string `json:"id"`
}
