package handlers

import (
	"context"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/router"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/internal/rest/resp"
)

type statisticsService interface {
	Get(ctx context.Context) (*domain.Statistics, error)
}

type StatisticsHandler struct {
	service statisticsService
}

func NewStatisticsHandler(router *router.RouterGroup[*core.RequestEvent], service statisticsService) {
	handler := &StatisticsHandler{
		service: service,
	}

	group := router.Group("/statistics")
	group.GET("/get", handler.get)
}

func (handler *StatisticsHandler) get(e *core.RequestEvent) error {
	if statistics, err := handler.service.Get(e.Request.Context()); err != nil {
		return resp.Err(e, err)
	} else {
		return resp.Ok(e, statistics)
	}
}
