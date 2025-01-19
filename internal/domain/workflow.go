package domain

import (
	"time"

	"github.com/usual2970/certimate/internal/pkg/utils/maps"
)

type WorkflowNodeType string

const (
	WorkflowNodeTypeStart               = WorkflowNodeType("start")
	WorkflowNodeTypeEnd                 = WorkflowNodeType("end")
	WorkflowNodeTypeApply               = WorkflowNodeType("apply")
	WorkflowNodeTypeDeploy              = WorkflowNodeType("deploy")
	WorkflowNodeTypeNotify              = WorkflowNodeType("notify")
	WorkflowNodeTypeBranch              = WorkflowNodeType("branch")
	WorkflowNodeTypeCondition           = WorkflowNodeType("condition")
	WorkflowNodeTypeExecuteResultBranch = WorkflowNodeType("execute_result_branch")
	WorkflowNodeTypeExecuteSuccess      = WorkflowNodeType("execute_success")
	WorkflowNodeTypeExecuteFailure      = WorkflowNodeType("execute_failure")
)

type WorkflowTriggerType string

const (
	WorkflowTriggerTypeAuto   = WorkflowTriggerType("auto")
	WorkflowTriggerTypeManual = WorkflowTriggerType("manual")
)

type Workflow struct {
	Meta
	Name          string                `json:"name" db:"name"`
	Description   string                `json:"description" db:"description"`
	Trigger       WorkflowTriggerType   `json:"trigger" db:"trigger"`
	TriggerCron   string                `json:"triggerCron" db:"triggerCron"`
	Enabled       bool                  `json:"enabled" db:"enabled"`
	Content       *WorkflowNode         `json:"content" db:"content"`
	Draft         *WorkflowNode         `json:"draft" db:"draft"`
	HasDraft      bool                  `json:"hasDraft" db:"hasDraft"`
	LastRunId     string                `json:"lastRunId" db:"lastRunId"`
	LastRunStatus WorkflowRunStatusType `json:"lastRunStatus" db:"lastRunStatus"`
	LastRunTime   time.Time             `json:"lastRunTime" db:"lastRunTime"`
}

func (w *Workflow) Table() string {
	return "workflow"
}

type WorkflowNode struct {
	Id   string           `json:"id"`
	Type WorkflowNodeType `json:"type"`
	Name string           `json:"name"`

	Config  map[string]any   `json:"config"`
	Inputs  []WorkflowNodeIO `json:"inputs"`
	Outputs []WorkflowNodeIO `json:"outputs"`

	Next     *WorkflowNode  `json:"next"`
	Branches []WorkflowNode `json:"branches"`

	Validated bool `json:"validated"`
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

func (n *WorkflowNode) GetConfigMap(key string) map[string]any {
	if val, ok := n.Config[key]; ok {
		if result, ok := val.(map[string]any); ok {
			return result
		}
	}

	return make(map[string]any)
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
	WorkflowId string              `json:"workflowId"`
	Trigger    WorkflowTriggerType `json:"trigger"`
}
