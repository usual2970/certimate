package repository

import (
	"context"

	"certimate/internal/domain"
	"certimate/internal/utils/app"
)

type SettingRepository struct{}

func NewSettingRepository() *SettingRepository {
	return &SettingRepository{}
}

func (s *SettingRepository) GetByName(ctx context.Context, name string) (*domain.Setting, error) {
	resp, err := app.GetApp().Dao().FindFirstRecordByFilter("settings", "name='"+name+"'")
	if err != nil {
		return nil, err
	}

	rs := &domain.Setting{
		ID:      resp.GetString("id"),
		Name:    resp.GetString("name"),
		Content: resp.GetString("content"),
		Created: resp.GetTime("created"),
		Updated: resp.GetTime("updated"),
	}

	return rs, nil
}
