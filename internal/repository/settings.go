package repository

import (
	"context"

	"github.com/pocketbase/dbx"
	"github.com/usual2970/certimate/internal/app"
	"github.com/usual2970/certimate/internal/domain"
)

type SettingsRepository struct{}

func NewSettingsRepository() *SettingsRepository {
	return &SettingsRepository{}
}

func (s *SettingsRepository) GetByName(ctx context.Context, name string) (*domain.Settings, error) {
	resp, err := app.GetApp().Dao().FindFirstRecordByFilter("settings", "name={:name}", dbx.Params{"name": name})
	if err != nil {
		return nil, err
	}

	rs := &domain.Settings{
		Meta: domain.Meta{
			Id:        resp.GetString("id"),
			CreatedAt: resp.GetTime("created"),
			UpdatedAt: resp.GetTime("updated"),
		},
		Name:    resp.GetString("name"),
		Content: resp.GetString("content"),
	}

	return rs, nil
}
