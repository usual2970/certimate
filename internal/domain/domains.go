package domain

type ApplyConfig struct {
	Email       string `json:"email"`
	Access      string `json:"access"`
	Timeout     int64  `json:"timeout"`
	Nameservers string `json:"nameservers"`
}

type DeployConfig struct {
	Id     string         `json:"id"`
	Access string         `json:"access"`
	Type   string         `json:"type"`
	Config map[string]any `json:"config"`
}

type KV struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
