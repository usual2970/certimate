package domain

import (
	"time"

	"github.com/usual2970/certimate/internal/pkg/utils/maps"
)

type WorkflowNodeType string

const (
	WorkflowNodeTypeStart     = WorkflowNodeType("start")
	WorkflowNodeTypeEnd       = WorkflowNodeType("end")
	WorkflowNodeTypeApply     = WorkflowNodeType("apply")
	WorkflowNodeTypeDeploy    = WorkflowNodeType("deploy")
	WorkflowNodeTypeNotify    = WorkflowNodeType("notify")
	WorkflowNodeTypeBranch    = WorkflowNodeType("branch")
	WorkflowNodeTypeCondition = WorkflowNodeType("condition")
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

type WorkflowNodeConfigForApply struct {
	Domains            string         `json:"domains"`
	ContactEmail       string         `json:"contactEmail"`
	Provider           string         `json:"provider"`
	ProviderAccessId   string         `json:"providerAccessId"`
	ProviderConfig     map[string]any `json:"providerConfig"`
	KeyAlgorithm       string         `json:"keyAlgorithm"`
	Nameservers        string         `json:"nameservers"`
	PropagationTimeout int32          `json:"propagationTimeout"`
	DisableFollowCNAME bool           `json:"disableFollowCNAME"`
}

type WorkflowNodeConfigForDeploy struct {
	Certificate      string         `json:"certificate"`
	Provider         string         `json:"provider"`
	ProviderAccessId string         `json:"providerAccessId"`
	ProviderConfig   map[string]any `json:"providerConfig"`
}

type WorkflowNodeConfigForNotify struct {
	Channel string `json:"channel"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

func (n *WorkflowNode) getConfigValueAsString(key string) string {
	return maps.GetValueAsString(n.Config, key)
}

func (n *WorkflowNode) getConfigValueAsBool(key string) bool {
	return maps.GetValueAsBool(n.Config, key)
}

func (n *WorkflowNode) getConfigValueAsInt32(key string) int32 {
	return maps.GetValueAsInt32(n.Config, key)
}

func (n *WorkflowNode) getConfigValueAsMap(key string) map[string]any {
	if val, ok := n.Config[key]; ok {
		if result, ok := val.(map[string]any); ok {
			return result
		}
	}

	return make(map[string]any)
}

func (n *WorkflowNode) GetConfigForApply() WorkflowNodeConfigForApply {
	return WorkflowNodeConfigForApply{
		Domains:            n.getConfigValueAsString("domains"),
		ContactEmail:       n.getConfigValueAsString("contactEmail"),
		Provider:           n.getConfigValueAsString("provider"),
		ProviderAccessId:   n.getConfigValueAsString("providerAccessId"),
		ProviderConfig:     n.getConfigValueAsMap("providerConfig"),
		KeyAlgorithm:       n.getConfigValueAsString("keyAlgorithm"),
		Nameservers:        n.getConfigValueAsString("nameservers"),
		PropagationTimeout: n.getConfigValueAsInt32("propagationTimeout"),
		DisableFollowCNAME: n.getConfigValueAsBool("disableFollowCNAME"),
	}
}

func (n *WorkflowNode) GetConfigForDeploy() WorkflowNodeConfigForDeploy {
	return WorkflowNodeConfigForDeploy{
		Certificate:      n.getConfigValueAsString("certificate"),
		Provider:         n.getConfigValueAsString("provider"),
		ProviderAccessId: n.getConfigValueAsString("providerAccessId"),
		ProviderConfig:   n.getConfigValueAsMap("providerConfig"),
	}
}

func (n *WorkflowNode) GetConfigForNotify() WorkflowNodeConfigForNotify {
	return WorkflowNodeConfigForNotify{
		Channel: n.getConfigValueAsString("channel"),
		Subject: n.getConfigValueAsString("subject"),
		Message: n.getConfigValueAsString("message"),
	}
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
	WorkflowId string              `json:"-"`
	Trigger    WorkflowTriggerType `json:"trigger"`
}
