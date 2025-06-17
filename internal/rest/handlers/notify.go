package handlers

import (
	"context"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/router"

	"github.com/certimate-go/certimate/internal/domain/dtos"
	"github.com/certimate-go/certimate/internal/rest/resp"
)

type notifyService interface {
	Test(ctx context.Context, req *dtos.NotifyTestPushReq) error
}

type NotifyHandler struct {
	service notifyService
}

func NewNotifyHandler(router *router.RouterGroup[*core.RequestEvent], service notifyService) {
	handler := &NotifyHandler{
		service: service,
	}

	group := router.Group("/notify")
	group.POST("/test", handler.test)
}

func (handler *NotifyHandler) test(e *core.RequestEvent) error {
	req := &dtos.NotifyTestPushReq{}
	if err := e.BindBody(req); err != nil {
		return resp.Err(e, err)
	}

	if err := handler.service.Test(e.Request.Context(), req); err != nil {
		return resp.Err(e, err)
	}

	return resp.Ok(e, nil)
}
