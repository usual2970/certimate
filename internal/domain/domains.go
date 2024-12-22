package domain

import (
	"encoding/json"
	"strings"

	"github.com/usual2970/certimate/internal/pkg/utils/maps"
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

// Deprecated: 以字符串形式获取配置项。
//
// 入参：
//   - key: 配置项的键。
//
// 出参：
//   - 配置项的值。如果配置项不存在或者类型不是字符串，则返回空字符串。
func (dc *DeployConfig) GetConfigAsString(key string) string {
	return maps.GetValueAsString(dc.Config, key)
}

// Deprecated: 以字符串形式获取配置项。
//
// 入参：
//   - key: 配置项的键。
//   - defaultValue: 默认值。
//
// 出参：
//   - 配置项的值。如果配置项不存在、类型不是字符串或者值为零值，则返回默认值。
func (dc *DeployConfig) GetConfigOrDefaultAsString(key string, defaultValue string) string {
	return maps.GetValueOrDefaultAsString(dc.Config, key, defaultValue)
}

// Deprecated: 以 32 位整数形式获取配置项。
//
// 入参：
//   - key: 配置项的键。
//
// 出参：
//   - 配置项的值。如果配置项不存在或者类型不是 32 位整数，则返回 0。
func (dc *DeployConfig) GetConfigAsInt32(key string) int32 {
	return maps.GetValueAsInt32(dc.Config, key)
}

// Deprecated: 以 32 位整数形式获取配置项。
//
// 入参：
//   - key: 配置项的键。
//   - defaultValue: 默认值。
//
// 出参：
//   - 配置项的值。如果配置项不存在、类型不是 32 位整数或者值为零值，则返回默认值。
func (dc *DeployConfig) GetConfigOrDefaultAsInt32(key string, defaultValue int32) int32 {
	return maps.GetValueOrDefaultAsInt32(dc.Config, key, defaultValue)
}

// Deprecated: 以布尔形式获取配置项。
//
// 入参：
//   - key: 配置项的键。
//
// 出参：
//   - 配置项的值。如果配置项不存在或者类型不是布尔，则返回 false。
func (dc *DeployConfig) GetConfigAsBool(key string) bool {
	return maps.GetValueAsBool(dc.Config, key)
}

// Deprecated: 以布尔形式获取配置项。
//
// 入参：
//   - key: 配置项的键。
//   - defaultValue: 默认值。
//
// 出参：
//   - 配置项的值。如果配置项不存在或者类型不是布尔，则返回默认值。
func (dc *DeployConfig) GetConfigOrDefaultAsBool(key string, defaultValue bool) bool {
	return maps.GetValueOrDefaultAsBool(dc.Config, key, defaultValue)
}

// Deprecated: 以变量字典形式获取配置项。
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

// Deprecated: GetDomain returns the domain from the deploy config,
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
