package domain

type ApplyConfig struct {
	Email        string `json:"email"`
	Access       string `json:"access"`
	KeyAlgorithm string `json:"keyAlgorithm"`
	Nameservers  string `json:"nameservers"`
	Timeout      int64  `json:"timeout"`
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
	if dc.Config == nil {
		return ""
	}

	if value, ok := dc.Config[key]; ok {
		if result, ok := value.(string); ok {
			return result
		}
	}

	return ""
}

// 以布尔形式获取配置项。
//
// 入参：
//   - key: 配置项的键。
//
// 出参：
//   - 配置项的值。如果配置项不存在或者类型不是布尔，则返回 false。
func (dc *DeployConfig) GetConfigAsBool(key string) bool {
	if dc.Config == nil {
		return false
	}

	if value, ok := dc.Config[key]; ok {
		if result, ok := value.(bool); ok {
			return result
		}
	}

	return false
}

type KV struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
