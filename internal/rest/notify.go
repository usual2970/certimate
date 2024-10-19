package rest

import (
	"context"

	"certimate/internal/domain"
	"certimate/internal/utils/resp"

	"github.com/labstack/echo/v5"
)

type NotifyService interface {
	Test(ctx context.Context, req *domain.NotifyTestPushReq) error
}

type notifyHandler struct {
	service NotifyService
}

func NewNotifyHandler(route *echo.Group, service NotifyService) {
	handler := &notifyHandler{
		service: service,
	}

	group := route.Group("/notify")

	group.POST("/test", handler.test)
}

func (handler *notifyHandler) test(c echo.Context) error {
	req := &domain.NotifyTestPushReq{}
	if err := c.Bind(req); err != nil {
		return err
	}

	if err := handler.service.Test(c.Request().Context(), req); err != nil {
		return resp.Err(c, err)
	}

	return resp.Succ(c, nil)
}
