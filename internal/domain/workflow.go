package domain

import (
	"github.com/usual2970/certimate/internal/pkg/utils/maps"
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

const (
	WorkflowTypeAuto   = "auto"
	WorkflowTypeManual = "manual"
)

type Workflow struct {
	Meta
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Type        string        `json:"type"`
	Crontab     string        `json:"crontab"`
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
	Input  []WorkflowNodeIO `json:"input"`
	Output []WorkflowNodeIO `json:"output"`

	Validated bool   `json:"validated"`
	Type      string `json:"type"`

	Branches []WorkflowNode `json:"branches"`
}

func (n *WorkflowNode) GetConfigString(key string) string {
	return maps.GetValueAsString(n.Config, key)
}

func (n *WorkflowNode) GetConfigBool(key string) bool {
	return maps.GetValueAsBool(n.Config, key)
}

func (n *WorkflowNode) GetConfigInt32(key string) int32 {
	return maps.GetValueAsInt32(n.Config, key)
}

func (n *WorkflowNode) GetConfigInt64(key string) int64 {
	return maps.GetValueAsInt64(n.Config, key)
}

type WorkflowNodeIO struct {
	Label         string                      `json:"label"`
	Name          string                      `json:"name"`
	Type          string                      `json:"type"`
	Required      bool                        `json:"required"`
	Value         any                         `json:"value"`
	ValueSelector WorkflowNodeIOValueSelector `json:"valueSelector"`
}

type WorkflowNodeIOValueSelector struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type WorkflowRunReq struct {
	Id string `json:"id"`
}
