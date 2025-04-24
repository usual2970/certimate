package domain

import (
	"encoding/json"
	"fmt"
)

const CollectionNameSettings = "settings"

type Settings struct {
	Meta
	Name    string `json:"name" db:"name"`
	Content string `json:"content" db:"content"`
}

// Deprecated: v0.4.x 将废弃
type NotifyTemplatesSettingsContent struct {
	NotifyTemplates []struct {
		Subject string `json:"subject"`
		Message string `json:"message"`
	} `json:"notifyTemplates"`
}

// Deprecated: v0.4.x 将废弃
type NotifyChannelsSettingsContent map[string]map[string]any

// Deprecated: v0.4.x 将废弃
func (s *Settings) GetNotifyChannelConfig(channel string) (map[string]any, error) {
	conf := &NotifyChannelsSettingsContent{}
	if err := json.Unmarshal([]byte(s.Content), conf); err != nil {
		return nil, err
	}

	v, ok := (*conf)[channel]
	if !ok {
		return nil, fmt.Errorf("channel \"%s\" not found", channel)
	}

	return v, nil
}

type PersistenceSettingsContent struct {
	WorkflowRunsMaxDaysRetention        int `json:"workflowRunsMaxDaysRetention"`
	ExpiredCertificatesMaxDaysRetention int `json:"expiredCertificatesMaxDaysRetention"`
}
