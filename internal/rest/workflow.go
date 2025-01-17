package rest

import (
	"context"

	"github.com/labstack/echo/v5"
	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/rest/resp"
)

type workflowService interface {
	Run(ctx context.Context, req *domain.WorkflowRunReq) error
	Stop(ctx context.Context)
}

type WorkflowHandler struct {
	service workflowService
}

func NewWorkflowHandler(route *echo.Group, service workflowService) {
	handler := &WorkflowHandler{
		service: service,
	}

	group := route.Group("/workflow")
	group.POST("/run", handler.run)
}

func (handler *WorkflowHandler) run(c echo.Context) error {
	req := &domain.WorkflowRunReq{}
	if err := c.Bind(req); err != nil {
		return resp.Err(c, err)
	}

	if err := handler.service.Run(c.Request().Context(), req); err != nil {
		return resp.Err(c, err)
	}

	return resp.Ok(c, nil)
}
