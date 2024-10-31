package domain

import (
	"encoding/json"
	"strings"
)

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

// 以字符串形式获取配置项。
//
// 入参：
//   - key: 配置项的键。
//
// 出参：
//   - 配置项的值。如果配置项不存在或者类型不是字符串，则返回空字符串。
func (dc *DeployConfig) GetConfigAsString(key string) string {
	return dc.GetConfigOrDefaultAsString(key, "")
}

// 以字符串形式获取配置项。
//
// 入参：
//   - key: 配置项的键。
//   - defaultValue: 默认值。
//
// 出参：
//   - 配置项的值。如果配置项不存在或者类型不是字符串，则返回默认值。
func (dc *DeployConfig) GetConfigOrDefaultAsString(key string, defaultValue string) string {
	if dc.Config == nil {
		return defaultValue
	}

	if value, ok := dc.Config[key]; ok {
		if result, ok := value.(string); ok {
			return result
		}
	}

	return defaultValue
}

// 以 32 位整数形式获取配置项。
//
// 入参：
//   - key: 配置项的键。
//
// 出参：
//   - 配置项的值。如果配置项不存在或者类型不是 32 位整数，则返回 0。
func (dc *DeployConfig) GetConfigAsInt32(key string) int32 {
	return dc.GetConfigOrDefaultAsInt32(key, 0)
}

// 以 32 位整数形式获取配置项。
//
// 入参：
//   - key: 配置项的键。
//   - defaultValue: 默认值。
//
// 出参：
//   - 配置项的值。如果配置项不存在或者类型不是 32 位整数，则返回默认值。
func (dc *DeployConfig) GetConfigOrDefaultAsInt32(key string, defaultValue int32) int32 {
	if dc.Config == nil {
		return defaultValue
	}

	if value, ok := dc.Config[key]; ok {
		if result, ok := value.(int32); ok {
			return result
		}
	}

	return defaultValue
}

// 以布尔形式获取配置项。
//
// 入参：
//   - key: 配置项的键。
//
// 出参：
//   - 配置项的值。如果配置项不存在或者类型不是布尔，则返回 false。
func (dc *DeployConfig) GetConfigAsBool(key string) bool {
	return dc.GetConfigOrDefaultAsBool(key, false)
}

// 以布尔形式获取配置项。
//
// 入参：
//   - key: 配置项的键。
//   - defaultValue: 默认值。
//
// 出参：
//   - 配置项的值。如果配置项不存在或者类型不是布尔，则返回默认值。
func (dc *DeployConfig) GetConfigOrDefaultAsBool(key string, defaultValue bool) bool {
	if dc.Config == nil {
		return defaultValue
	}

	if value, ok := dc.Config[key]; ok {
		if result, ok := value.(bool); ok {
			return result
		}
	}

	return defaultValue
}

// 以变量字典形式获取配置项。
//
// 出参：
//   - 变量字典。
func (dc *DeployConfig) GetConfigAsVariables() map[string]string {
	rs := make(map[string]string)

	if dc.Config != nil {
		value, ok := dc.Config["variables"]
		if !ok {
			return rs
		}

		kvs := make([]KV, 0)
		bts, _ := json.Marshal(value)
		if err := json.Unmarshal(bts, &kvs); err != nil {
			return rs
		}

		for _, kv := range kvs {
			rs[kv.Key] = kv.Value
		}
	}

	return rs
}

// GetDomain returns the domain from the deploy config
// if the domain is a wildcard domain, and wildcard is true, return the wildcard domain
func (dc *DeployConfig) GetDomain(wildcard ...bool) string {
	val := dc.GetConfigAsString("domain")
	if val == "" {
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
