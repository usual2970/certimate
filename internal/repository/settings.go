package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/pocketbase/dbx"
	"github.com/usual2970/certimate/internal/app"
	"github.com/usual2970/certimate/internal/domain"
)

type SettingsRepository struct{}

func NewSettingsRepository() *SettingsRepository {
	return &SettingsRepository{}
}

func (r *SettingsRepository) GetByName(ctx context.Context, name string) (*domain.Settings, error) {
	record, err := app.GetApp().Dao().FindFirstRecordByFilter(
		"settings",
		"name={:name}",
		dbx.Params{"name": name},
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrRecordNotFound
		}
		return nil, err
	}

	settings := &domain.Settings{
		Meta: domain.Meta{
			Id:        record.GetId(),
			CreatedAt: record.GetCreated().Time(),
			UpdatedAt: record.GetUpdated().Time(),
		},
		Name:    record.GetString("name"),
		Content: record.GetString("content"),
	}
	return settings, nil
}
