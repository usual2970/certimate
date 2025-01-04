package domain

// Deprecated: TODO: 即将废弃
type ApplyConfig struct {
	ContactEmail       string `json:"contactEmail"`
	ProviderAccessId   string `json:"providerAccessId"`
	KeyAlgorithm       string `json:"keyAlgorithm"`
	Nameservers        string `json:"nameservers"`
	PropagationTimeout int32  `json:"propagationTimeout"`
	DisableFollowCNAME bool   `json:"disableFollowCNAME"`
}

// Deprecated: TODO: 即将废弃
type DeployConfig struct {
	NodeId           string         `json:"nodeId"`
	NodeConfig       map[string]any `json:"nodeConfig"`
	Provider         string         `json:"provider"`
	ProviderAccessId string         `json:"providerAccessId"`
}
