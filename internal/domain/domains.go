package domain

import "strings"

type ApplyConfig struct {
	Email              string `json:"email"`
	Access             string `json:"access"`
	KeyAlgorithm       string `json:"keyAlgorithm"`
	Nameservers        string `json:"nameservers"`
	Timeout            int64  `json:"timeout"`
	DisableFollowCNAME bool   `json:"disableFollowCNAME"`
}

type DeployConfig struct {
	Id     string         `json:"id"`
	Access string         `json:"access"`
	Type   string         `json:"type"`
	Config map[string]any `json:"config"`
}

// GetDomain returns the domain from the deploy config
// if the domain is a wildcard domain, and wildcard is true, return the wildcard domain
func (d *DeployConfig) GetDomain(wildcard ...bool) string {
	if _, ok := d.Config["domain"]; !ok {
		return ""
	}

	val, ok := d.Config["domain"].(string)
	if !ok {
		return ""
	}

	if !strings.HasPrefix(val, "*") {
		return val
	}

	if len(wildcard) > 0 && wildcard[0] {
		return val
	}

	return strings.TrimPrefix(val, "*")
}

type KV struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
