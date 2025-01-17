package rest

import (
	"context"

	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/rest/resp"

	"github.com/labstack/echo/v5"
)

type notifyService interface {
	Test(ctx context.Context, req *domain.NotifyTestPushReq) error
}

type NotifyHandler struct {
	service notifyService
}

func NewNotifyHandler(route *echo.Group, service notifyService) {
	handler := &NotifyHandler{
		service: service,
	}

	group := route.Group("/notify")
	group.POST("/test", handler.test)
}

func (handler *NotifyHandler) test(c echo.Context) error {
	req := &domain.NotifyTestPushReq{}
	if err := c.Bind(req); err != nil {
		return resp.Err(c, err)
	}

	if err := handler.service.Test(c.Request().Context(), req); err != nil {
		return resp.Err(c, err)
	}

	return resp.Ok(c, nil)
}
