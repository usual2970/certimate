package rest

import (
	"context"

	"github.com/labstack/echo/v5"
	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/rest/resp"
)

type statisticsService interface {
	Get(ctx context.Context) (*domain.Statistics, error)
}

type StatisticsHandler struct {
	service statisticsService
}

func NewStatisticsHandler(route *echo.Group, service statisticsService) {
	handler := &StatisticsHandler{
		service: service,
	}

	group := route.Group("/statistics")
	group.GET("/get", handler.get)
}

func (handler *StatisticsHandler) get(c echo.Context) error {
	if statistics, err := handler.service.Get(c.Request().Context()); err != nil {
		return resp.Err(c, err)
	} else {
		return resp.Ok(c, statistics)
	}
}
