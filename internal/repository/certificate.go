package repository

import (
	"context"

	"github.com/usual2970/certimate/internal/app"
	"github.com/usual2970/certimate/internal/domain"
)

type CertificateRepository struct{}

func NewCertificateRepository() *CertificateRepository {
	return &CertificateRepository{}
}

func (c *CertificateRepository) ListExpireSoon(ctx context.Context) ([]domain.Certificate, error) {
	rs := []domain.Certificate{}
	if err := app.GetApp().Dao().DB().
		NewQuery("SELECT * FROM certificate WHERE expireAt > DATETIME('now') AND expireAt < DATETIME('now', '+20 days')").
		All(&rs); err != nil {
		return nil, err
	}
	return rs, nil
}
