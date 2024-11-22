package rest

import (
	"context"

	"github.com/labstack/echo/v5"
	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/utils/resp"
)

type StatisticsService interface {
	Get(ctx context.Context) (*domain.Statistics, error)
}

type statisticsHandler struct {
	service StatisticsService
}

func NewStatisticsHandler(route *echo.Group, service StatisticsService) {
	handler := &statisticsHandler{
		service: service,
	}

	group := route.Group("/statistics")

	group.GET("/get", handler.get)
}

func (handler *statisticsHandler) get(c echo.Context) error {
	if statistics, err := handler.service.Get(c.Request().Context()); err != nil {
		return resp.Err(c, err)
	} else {
		return resp.Succ(c, statistics)
	}
}
