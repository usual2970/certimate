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

func (c *CertificateRepository) GetExpireSoon(ctx context.Context) ([]domain.Certificate, error) {
	rs := []domain.Certificate{}
	if err := app.GetApp().Dao().DB().
		NewQuery("select * from certificate where expireAt > datetime('now') and expireAt < datetime('now', '+20 days')").
		All(&rs); err != nil {
		return nil, err
	}
	return rs, nil
}
