package domain

// Deprecated: TODO: 即将废弃
type DeployConfig struct {
	NodeId           string         `json:"nodeId"`
	NodeConfig       map[string]any `json:"nodeConfig"`
	Provider         string         `json:"provider"`
	ProviderAccessId string         `json:"providerAccessId"`
}
