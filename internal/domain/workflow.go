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
	Domains               string         `json:"domains"`                         // 域名列表，以半角分号分隔
	ContactEmail          string         `json:"contactEmail"`                    // 联系邮箱
	ChallengeType         string         `json:"challengeType"`                   // TODO: 验证方式。目前仅支持 dns-01
	Provider              string         `json:"provider"`                        // DNS 提供商
	ProviderAccessId      string         `json:"providerAccessId"`                // DNS 提供商授权记录 ID
	ProviderConfig        map[string]any `json:"providerConfig"`                  // DNS 提供商额外配置
	CAProvider            string         `json:"caProvider,omitempty"`            // CA 提供商（零值将使用全局配置）
	CAProviderAccessId    string         `json:"caProviderAccessId,omitempty"`    // CA 提供商授权记录 ID
	CAProviderConfig      map[string]any `json:"caProviderConfig,omitempty"`      // CA 提供商额外配置
	KeyAlgorithm          string         `json:"keyAlgorithm"`                    // 密钥算法
	Nameservers           string         `json:"nameservers,omitempty"`           // DNS 服务器列表，以半角分号分隔
	DnsPropagationTimeout int32          `json:"dnsPropagationTimeout,omitempty"` // DNS 传播超时时间（零值取决于提供商的默认值）
	DnsTTL                int32          `json:"dnsTTL,omitempty"`                // DNS TTL（零值取决于提供商的默认值）
	DisableFollowCNAME    bool           `json:"disableFollowCNAME,omitempty"`    // 是否关闭 CNAME 跟随
	DisableARI            bool           `json:"disableARI,omitempty"`            // 是否关闭 ARI
	SkipBeforeExpiryDays  int32          `json:"skipBeforeExpiryDays,omitempty"`  // 证书到期前多少天前跳过续期（零值将使用默认值 30）
}

type WorkflowNodeConfigForUpload struct {
	Certificate string `json:"certificate"`
	PrivateKey  string `json:"privateKey"`
	Domains     string `json:"domains"`
}

type WorkflowNodeConfigForDeploy struct {
	Certificate         string         `json:"certificate"`                // 前序节点输出的证书，形如“${NodeId}#certificate”
	Provider            string         `json:"provider"`                   // 主机提供商
	ProviderAccessId    string         `json:"providerAccessId,omitempty"` // 主机提供商授权记录 ID
	ProviderConfig      map[string]any `json:"providerConfig,omitempty"`   // 主机提供商额外配置
	SkipOnLastSucceeded bool           `json:"skipOnLastSucceeded"`        // 上次部署成功时是否跳过
}

type WorkflowNodeConfigForNotify struct {
	Channel string `json:"channel"` // 通知渠道
	Subject string `json:"subject"` // 通知主题
	Message string `json:"message"` // 通知内容
}

func (n *WorkflowNode) GetConfigForApply() WorkflowNodeConfigForApply {
	skipBeforeExpiryDays := maputil.GetInt32(n.Config, "skipBeforeExpiryDays")
	if skipBeforeExpiryDays == 0 {
		skipBeforeExpiryDays = 30
	}

	return WorkflowNodeConfigForApply{
		Domains:               maputil.GetString(n.Config, "domains"),
		ContactEmail:          maputil.GetString(n.Config, "contactEmail"),
		Provider:              maputil.GetString(n.Config, "provider"),
		ProviderAccessId:      maputil.GetString(n.Config, "providerAccessId"),
		ProviderConfig:        maputil.GetAnyMap(n.Config, "providerConfig"),
		CAProvider:            maputil.GetString(n.Config, "caProvider"),
		CAProviderAccessId:    maputil.GetString(n.Config, "caProviderAccessId"),
		CAProviderConfig:      maputil.GetAnyMap(n.Config, "caProviderConfig"),
		KeyAlgorithm:          maputil.GetString(n.Config, "keyAlgorithm"),
		Nameservers:           maputil.GetString(n.Config, "nameservers"),
		DnsPropagationTimeout: maputil.GetInt32(n.Config, "dnsPropagationTimeout"),
		DnsTTL:                maputil.GetInt32(n.Config, "dnsTTL"),
		DisableFollowCNAME:    maputil.GetBool(n.Config, "disableFollowCNAME"),
		DisableARI:            maputil.GetBool(n.Config, "disableARI"),
		SkipBeforeExpiryDays:  skipBeforeExpiryDays,
	}
}

func (n *WorkflowNode) GetConfigForUpload() WorkflowNodeConfigForUpload {
	return WorkflowNodeConfigForUpload{
		Certificate: maputil.GetString(n.Config, "certificate"),
		PrivateKey:  maputil.GetString(n.Config, "privateKey"),
		Domains:     maputil.GetString(n.Config, "domains"),
	}
}

func (n *WorkflowNode) GetConfigForDeploy() WorkflowNodeConfigForDeploy {
	return WorkflowNodeConfigForDeploy{
		Certificate:         maputil.GetString(n.Config, "certificate"),
		Provider:            maputil.GetString(n.Config, "provider"),
		ProviderAccessId:    maputil.GetString(n.Config, "providerAccessId"),
		ProviderConfig:      maputil.GetAnyMap(n.Config, "providerConfig"),
		SkipOnLastSucceeded: maputil.GetBool(n.Config, "skipOnLastSucceeded"),
	}
}

func (n *WorkflowNode) GetConfigForNotify() WorkflowNodeConfigForNotify {
	return WorkflowNodeConfigForNotify{
		Channel: maputil.GetString(n.Config, "channel"),
		Subject: maputil.GetString(n.Config, "subject"),
		Message: maputil.GetString(n.Config, "message"),
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
