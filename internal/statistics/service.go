package statistics

import (
	"context"

	"github.com/usual2970/certimate/internal/domain"
)

type StatisticsRepository interface {
	Get(ctx context.Context) (*domain.Statistics, error)
}

type StatisticsService struct {
	repo StatisticsRepository
}

func NewStatisticsService(repo StatisticsRepository) *StatisticsService {
	return &StatisticsService{
		repo: repo,
	}
}

func (s *StatisticsService) Get(ctx context.Context) (*domain.Statistics, error) {
	return s.repo.Get(ctx)
}
