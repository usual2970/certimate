package notify

import (
	"github.com/usual2970/certimate/internal/pkg/utils/maps"
)

func getConfigAsString(conf map[string]any, key string) string {
	return maps.GetValueAsString(conf, key)
}

func getConfigAsInt32(conf map[string]any, key string) int32 {
	return maps.GetValueAsInt32(conf, key)
}

func getConfigAsInt64(conf map[string]any, key string) int64 {
	return maps.GetValueAsInt64(conf, key)
}

func getConfigAsBool(conf map[string]any, key string) bool {
	return maps.GetValueAsBool(conf, key)
}
