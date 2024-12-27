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

type NotifyChannelsConfig map[string]map[string]any

func (s *Settings) GetChannelContent(channel string) (map[string]any, error) {
	conf := &NotifyChannelsConfig{}
	if err := json.Unmarshal([]byte(s.Content), conf); err != nil {
		return nil, err
	}

	v, ok := (*conf)[channel]
	if !ok {
		return nil, fmt.Errorf("channel \"%s\" not found", channel)
	}

	return v, nil
}

type NotifyTemplates struct {
	NotifyTemplates []NotifyTemplate `json:"notifyTemplates"`
}

type NotifyTemplate struct {
	Subject string `json:"subject"`
	Message string `json:"message"`
}

type NotifyMessage struct {
	Subject string `json:"subject"`
	Message string `json:"message"`
}
