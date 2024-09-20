package variables

import "strings"

// Parse2Map 将变量赋值字符串解析为map
func Parse2Map(str string) map[string]string {

	m := make(map[string]string)

	lines := strings.Split(str, ";")

	for _, line := range lines {

		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		kv := strings.Split(line, "=")

		if len(kv) != 2 {
			continue
		}

		m[kv[0]] = kv[1]
	}

	return m
}
