package domain

// Deprecated: TODO: 即将废弃
type ApplyConfig struct {
	Email              string `json:"email"`
	Access             string `json:"access"`
	KeyAlgorithm       string `json:"keyAlgorithm"`
	Nameservers        string `json:"nameservers"`
	PropagationTimeout int64  `json:"propagationTimeout"`
	DisableFollowCNAME bool   `json:"disableFollowCNAME"`
}

// Deprecated: TODO: 即将废弃
type DeployConfig struct {
	Id     string         `json:"id"`
	Access string         `json:"access"`
	Type   string         `json:"type"`
	Config map[string]any `json:"config"`
}

// Deprecated: TODO: 即将废弃
type KV struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
