package notify

func getString(conf map[string]any, key string) string {
	if _, ok := conf[key]; !ok {
		return ""
	}

	return conf[key].(string)
}

func getBool(conf map[string]any, key string) bool {
	if _, ok := conf[key]; !ok {
		return false
	}

	return conf[key].(bool)
}
