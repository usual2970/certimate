package domain

import (
	"encoding/json"
	"fmt"
	"time"
)

type Setting struct {
	ID      string    `json:"id"`
	Name    string    `json:"name"`
	Content string    `json:"content"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

type ChannelsConfig map[string]map[string]any

func (s *Setting) GetChannelContent(channel string) (map[string]any, error) {
	conf := &ChannelsConfig{}
	if err := json.Unmarshal([]byte(s.Content), conf); err != nil {
		return nil, err
	}

	v, ok := (*conf)[channel]
	if !ok {
		return nil, fmt.Errorf("channel \"%s\" not found", channel)
	}

	return v, nil
}
