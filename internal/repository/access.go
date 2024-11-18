package repository

import (
	"context"

	"github.com/usual2970/certimate/internal/domain"
)

type AccessRepository struct{}

func NewAccessRepository() *AccessRepository {
	return &AccessRepository{}
}

func (a *AccessRepository) GetById(ctx context.Context, id string) (*domain.Access, error) {
	return nil, nil
}
