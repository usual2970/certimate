package statistics

import (
	"context"

	"github.com/usual2970/certimate/internal/domain"
)

type statisticsRepository interface {
	Get(ctx context.Context) (*domain.Statistics, error)
}

type StatisticsService struct {
	repo statisticsRepository
}

func NewStatisticsService(repo statisticsRepository) *StatisticsService {
	return &StatisticsService{
		repo: repo,
	}
}

func (s *StatisticsService) Get(ctx context.Context) (*domain.Statistics, error) {
	return s.repo.Get(ctx)
}
