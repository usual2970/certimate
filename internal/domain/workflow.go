package domain

import (
	"time"

	"github.com/usual2970/certimate/internal/pkg/utils/maputil"
)

const CollectionNameWorkflow = "workflow"

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

type WorkflowNodeType string

const (
	WorkflowNodeTypeStart               = WorkflowNodeType("start")
	WorkflowNodeTypeEnd                 = WorkflowNodeType("end")
	WorkflowNodeTypeApply               = WorkflowNodeType("apply")
	WorkflowNodeTypeUpload              = WorkflowNodeType("upload")
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

type WorkflowNode struct {
	Id   string           `json:"id"`
	Type WorkflowNodeType `json:"type"`
	Name string           `json:"name"`

	Config  map[string]any   `json:"config"`
	Inputs  []WorkflowNodeIO `json:"inputs"`
	Outputs []WorkflowNodeIO `json:"outputs"`

	Next     *WorkflowNode  `json:"next,omitempty"`
	Branches []WorkflowNode `json:"branches,omitempty"`

	Validated bool `json:"validated"`
}

type WorkflowNodeConfigForApply struct {
	Domains               string         `json:"domains"`               // 域名列表，以半角分号分隔
	ContactEmail          string         `json:"contactEmail"`          // 联系邮箱
	ChallengeType         string         `json:"challengeType"`         // TODO: 验证方式。目前仅支持 dns-01
	Provider              string         `json:"provider"`              // DNS 提供商
	ProviderAccessId      string         `json:"providerAccessId"`      // DNS 提供商授权记录 ID
	ProviderConfig        map[string]any `json:"providerConfig"`        // DNS 提供商额外配置
	KeyAlgorithm          string         `json:"keyAlgorithm"`          // 密钥算法
	Nameservers           string         `json:"nameservers"`           // DNS 服务器列表，以半角分号分隔
	DnsPropagationTimeout int32          `json:"dnsPropagationTimeout"` // DNS 传播超时时间（零值取决于提供商的默认值）
	DnsTTL                int32          `json:"dnsTTL"`                // DNS TTL（零值取决于提供商的默认值）
	DisableFollowCNAME    bool           `json:"disableFollowCNAME"`    // 是否关闭 CNAME 跟随
	DisableARI            bool           `json:"disableARI"`            // 是否关闭 ARI
	SkipBeforeExpiryDays  int32          `json:"skipBeforeExpiryDays"`  // 证书到期前多少天前跳过续期（零值将使用默认值 30）
}

type WorkflowNodeConfigForUpload struct {
	Certificate string `json:"certificate"`
	PrivateKey  string `json:"privateKey"`
	Domains     string `json:"domains"`
}

type WorkflowNodeConfigForDeploy struct {
	Certificate         string         `json:"certificate"`         // 前序节点输出的证书，形如“${NodeId}#certificate”
	Provider            string         `json:"provider"`            // 主机提供商
	ProviderAccessId    string         `json:"providerAccessId"`    // 主机提供商授权记录 ID
	ProviderConfig      map[string]any `json:"providerConfig"`      // 主机提供商额外配置
	SkipOnLastSucceeded bool           `json:"skipOnLastSucceeded"` // 上次部署成功时是否跳过
}

type WorkflowNodeConfigForNotify struct {
	Channel string `json:"channel"` // 通知渠道
	Subject string `json:"subject"` // 通知主题
	Message string `json:"message"` // 通知内容
}

func (n *WorkflowNode) getConfigString(key string) string {
	return maputil.GetString(n.Config, key)
}

func (n *WorkflowNode) getConfigBool(key string) bool {
	return maputil.GetBool(n.Config, key)
}

func (n *WorkflowNode) getConfigInt32(key string) int32 {
	return maputil.GetInt32(n.Config, key)
}

func (n *WorkflowNode) getConfigMap(key string) map[string]any {
	if val, ok := n.Config[key]; ok {
		if result, ok := val.(map[string]any); ok {
			return result
		}
	}

	return make(map[string]any)
}

func (n *WorkflowNode) GetConfigForApply() WorkflowNodeConfigForApply {
	skipBeforeExpiryDays := n.getConfigInt32("skipBeforeExpiryDays")
	if skipBeforeExpiryDays == 0 {
		skipBeforeExpiryDays = 30
	}

	return WorkflowNodeConfigForApply{
		Domains:               n.getConfigString("domains"),
		ContactEmail:          n.getConfigString("contactEmail"),
		Provider:              n.getConfigString("provider"),
		ProviderAccessId:      n.getConfigString("providerAccessId"),
		ProviderConfig:        n.getConfigMap("providerConfig"),
		KeyAlgorithm:          n.getConfigString("keyAlgorithm"),
		Nameservers:           n.getConfigString("nameservers"),
		DnsPropagationTimeout: n.getConfigInt32("dnsPropagationTimeout"),
		DnsTTL:                n.getConfigInt32("dnsTTL"),
		DisableFollowCNAME:    n.getConfigBool("disableFollowCNAME"),
		DisableARI:            n.getConfigBool("disableARI"),
		SkipBeforeExpiryDays:  skipBeforeExpiryDays,
	}
}

func (n *WorkflowNode) GetConfigForUpload() WorkflowNodeConfigForUpload {
	return WorkflowNodeConfigForUpload{
		Certificate: n.getConfigString("certificate"),
		PrivateKey:  n.getConfigString("privateKey"),
		Domains:     n.getConfigString("domains"),
	}
}

func (n *WorkflowNode) GetConfigForDeploy() WorkflowNodeConfigForDeploy {
	return WorkflowNodeConfigForDeploy{
		Certificate:         n.getConfigString("certificate"),
		Provider:            n.getConfigString("provider"),
		ProviderAccessId:    n.getConfigString("providerAccessId"),
		ProviderConfig:      n.getConfigMap("providerConfig"),
		SkipOnLastSucceeded: n.getConfigBool("skipOnLastSucceeded"),
	}
}

func (n *WorkflowNode) GetConfigForNotify() WorkflowNodeConfigForNotify {
	return WorkflowNodeConfigForNotify{
		Channel: n.getConfigString("channel"),
		Subject: n.getConfigString("subject"),
		Message: n.getConfigString("message"),
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

const WorkflowNodeIONameCertificate string = "certificate"
