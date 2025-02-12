package statistics

import (
	"context"

	"github.com/usual2970/certimate/internal/domain"
)

type statisticsRepository interface {
	Get(ctx context.Context) (*domain.Statistics, error)
}

type StatisticsService struct {
	statRepo statisticsRepository
}

func NewStatisticsService(statRepo statisticsRepository) *StatisticsService {
	return &StatisticsService{
		statRepo: statRepo,
	}
}

func (s *StatisticsService) Get(ctx context.Context) (*domain.Statistics, error) {
	return s.statRepo.Get(ctx)
}
