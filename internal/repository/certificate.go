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
	certificates := []domain.Certificate{}
	err := app.GetApp().Dao().DB().
		NewQuery("SELECT * FROM certificate WHERE expireAt > DATETIME('now') AND expireAt < DATETIME('now', '+20 days')").
		All(&certificates)
	if err != nil {
		return nil, err
	}

	return certificates, nil
}
