package maputil

import (
	"strconv"
)

// 以字符串形式从字典中获取指定键的值。
//
// 入参：
//   - dict: 字典。
//   - key: 键。
//
// 出参：
//   - 字典中键对应的值。如果指定键不存在或者值的类型不是字符串，则返回空字符串。
func GetString(dict map[string]any, key string) string {
	return GetOrDefaultString(dict, key, "")
}

// 以字符串形式从字典中获取指定键的值。
//
// 入参：
//   - dict: 字典。
//   - key: 键。
//   - defaultValue: 默认值。
//
// 出参：
//   - 字典中键对应的值。如果指定键不存在、值的类型不是字符串或者值为零值，则返回默认值。
func GetOrDefaultString(dict map[string]any, key string, defaultValue string) string {
	if dict == nil {
		return defaultValue
	}

	if value, ok := dict[key]; ok {
		if result, ok := value.(string); ok {
			if result != "" {
				return result
			}
		}
	}

	return defaultValue
}

// 以 32 位整数形式从字典中获取指定键的值。
//
// 入参：
//   - dict: 字典。
//   - key: 键。
//
// 出参：
//   - 字典中键对应的值。如果指定键不存在或者值的类型不是 32 位整数，则返回 0。
func GetInt32(dict map[string]any, key string) int32 {
	return GetOrDefaultInt32(dict, key, 0)
}

// 以 32 位整数形式从字典中获取指定键的值。
//
// 入参：
//   - dict: 字典。
//   - key: 键。
//   - defaultValue: 默认值。
//
// 出参：
//   - 字典中键对应的值。如果指定键不存在、值的类型不是 32 位整数或者值为零值，则返回默认值。
func GetOrDefaultInt32(dict map[string]any, key string, defaultValue int32) int32 {
	if dict == nil {
		return defaultValue
	}

	if value, ok := dict[key]; ok {
		if result, ok := value.(int32); ok {
			if result != 0 {
				return result
			}
		}

		// 兼容字符串类型的值
		if str, ok := value.(string); ok {
			if result, err := strconv.ParseInt(str, 10, 32); err == nil {
				if result != 0 {
					return int32(result)
				}
			}
		}
	}

	return defaultValue
}

// 以 64 位整数形式从字典中获取指定键的值。
//
// 入参：
//   - dict: 字典。
//   - key: 键。
//
// 出参：
//   - 字典中键对应的值。如果指定键不存在或者值的类型不是 64 位整数，则返回 0。
func GetInt64(dict map[string]any, key string) int64 {
	return GetOrDefaultInt64(dict, key, 0)
}

// 以 64 位整数形式从字典中获取指定键的值。
//
// 入参：
//   - dict: 字典。
//   - key: 键。
//   - defaultValue: 默认值。
//
// 出参：
//   - 字典中键对应的值。如果指定键不存在、值的类型不是 64 位整数或者值为零值，则返回默认值。
func GetOrDefaultInt64(dict map[string]any, key string, defaultValue int64) int64 {
	if dict == nil {
		return defaultValue
	}

	if value, ok := dict[key]; ok {
		if result, ok := value.(int64); ok {
			if result != 0 {
				return result
			}
		}

		if result, ok := value.(int32); ok {
			if result != 0 {
				return int64(result)
			}
		}

		// 兼容字符串类型的值
		if str, ok := value.(string); ok {
			if result, err := strconv.ParseInt(str, 10, 64); err == nil {
				if result != 0 {
					return result
				}
			}
		}
	}

	return defaultValue
}

// 以布尔形式从字典中获取指定键的值。
//
// 入参：
//   - dict: 字典。
//   - key: 键。
//
// 出参：
//   - 字典中键对应的值。如果指定键不存在或者值的类型不是布尔，则返回 false。
func GetBool(dict map[string]any, key string) bool {
	return GetOrDefaultBool(dict, key, false)
}

// 以布尔形式从字典中获取指定键的值。
//
// 入参：
//   - dict: 字典。
//   - key: 键。
//   - defaultValue: 默认值。
//
// 出参：
//   - 字典中键对应的值。如果指定键不存在或者值的类型不是布尔，则返回默认值。
func GetOrDefaultBool(dict map[string]any, key string, defaultValue bool) bool {
	if dict == nil {
		return defaultValue
	}

	if value, ok := dict[key]; ok {
		if result, ok := value.(bool); ok {
			return result
		}

		// 兼容字符串类型的值
		if str, ok := value.(string); ok {
			if result, err := strconv.ParseBool(str); err == nil {
				return result
			}
		}
	}

	return defaultValue
}

// 以 `map[string]any` 形式从字典中获取指定键的值。
//
// 入参：
//   - dict: 字典。
//   - key: 键。
//
// 出参：
//   - 字典中键对应的 `map[string]any` 对象。
func GetAnyMap(dict map[string]any, key string) map[string]any {
	if dict == nil {
		return make(map[string]any)
	}

	if val, ok := dict[key]; ok {
		if result, ok := val.(map[string]any); ok {
			return result
		}
	}

	return make(map[string]any)
}
