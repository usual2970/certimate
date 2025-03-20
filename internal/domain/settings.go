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

type NotifyTemplatesSettingsContent struct {
	NotifyTemplates []struct {
		Subject string `json:"subject"`
		Message string `json:"message"`
	} `json:"notifyTemplates"`
}

type NotifyChannelsSettingsContent map[string]map[string]any

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
