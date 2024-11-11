package maps

import "strconv"

// 以字符串形式从字典中获取指定键的值。
//
// 入参：
//   - dict: 字典。
//   - key: 键。
//
// 出参：
//   - 字典中键对应的值。如果指定键不存在或者值的类型不是字符串，则返回空字符串。
func GetValueAsString(dict map[string]any, key string) string {
	return GetValueOrDefaultAsString(dict, key, "")
}

// 以字符串形式从字典中获取指定键的值。
//
// 入参：
//   - dict: 字典。
//   - key: 键。
//   - defaultValue: 默认值。
//
// 出参：
//   - 字典中键对应的值。如果指定键不存在或者值的类型不是字符串，则返回默认值。
func GetValueOrDefaultAsString(dict map[string]any, key string, defaultValue string) string {
	if dict == nil {
		return defaultValue
	}

	if value, ok := dict[key]; ok {
		if result, ok := value.(string); ok {
			return result
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
func GetValueAsInt32(dict map[string]any, key string) int32 {
	return GetValueOrDefaultAsInt32(dict, key, 0)
}

// 以 32 位整数形式从字典中获取指定键的值。
//
// 入参：
//   - dict: 字典。
//   - key: 键。
//   - defaultValue: 默认值。
//
// 出参：
//   - 字典中键对应的值。如果指定键不存在或者值的类型不是 32 位整数，则返回默认值。
func GetValueOrDefaultAsInt32(dict map[string]any, key string, defaultValue int32) int32 {
	if dict == nil {
		return defaultValue
	}

	if value, ok := dict[key]; ok {
		if result, ok := value.(int32); ok {
			return result
		}

		// 兼容字符串类型的值
		if str, ok := value.(string); ok {
			if result, err := strconv.ParseInt(str, 10, 32); err == nil {
				return int32(result)
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
func GetValueAsInt64(dict map[string]any, key string) int64 {
	return GetValueOrDefaultAsInt64(dict, key, 0)
}

// 以 64 位整数形式从字典中获取指定键的值。
//
// 入参：
//   - dict: 字典。
//   - key: 键。
//   - defaultValue: 默认值。
//
// 出参：
//   - 字典中键对应的值。如果指定键不存在或者值的类型不是 64 位整数，则返回默认值。
func GetValueOrDefaultAsInt64(dict map[string]any, key string, defaultValue int64) int64 {
	if dict == nil {
		return defaultValue
	}

	if value, ok := dict[key]; ok {
		if result, ok := value.(int64); ok {
			return result
		}

		// 兼容字符串类型的值
		if str, ok := value.(string); ok {
			if result, err := strconv.ParseInt(str, 10, 64); err == nil {
				return result
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
func GetValueAsBool(dict map[string]any, key string) bool {
	return GetValueOrDefaultAsBool(dict, key, false)
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
func GetValueOrDefaultAsBool(dict map[string]any, key string, defaultValue bool) bool {
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
