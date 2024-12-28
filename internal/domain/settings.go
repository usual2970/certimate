package domain

import (
	"encoding/json"
	"fmt"
)

type Settings struct {
	Meta
	Name    string `json:"name" db:"name"`
	Content string `json:"content" db:"content"`
}

type NotifyTemplatesSettingsContent struct {
	NotifyTemplates []NotifyTemplate `json:"notifyTemplates"`
}

type NotifyTemplate struct {
	Subject string `json:"subject"`
	Message string `json:"message"`
}

type NotifyChannelsSettingsContent map[string]map[string]any

type NotifyMessage struct {
	Subject string `json:"subject"`
	Message string `json:"message"`
}

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
