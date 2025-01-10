package rest

import (
	"context"

	"github.com/labstack/echo/v5"
	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/rest/resp"
)

type WorkflowService interface {
	Run(ctx context.Context, req *domain.WorkflowRunReq) error
	Stop()
}

type workflowHandler struct {
	service WorkflowService
}

func NewWorkflowHandler(route *echo.Group, service WorkflowService) {
	handler := &workflowHandler{
		service: service,
	}

	group := route.Group("/workflow")

	group.POST("/run", handler.run)
}

func (handler *workflowHandler) run(c echo.Context) error {
	req := &domain.WorkflowRunReq{}
	if err := c.Bind(req); err != nil {
		return resp.Err(c, err)
	}

	if err := handler.service.Run(c.Request().Context(), req); err != nil {
		return resp.Err(c, err)
	}
	return resp.Succ(c, nil)
}
