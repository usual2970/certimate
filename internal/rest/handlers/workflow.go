package handlers

import (
	"context"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/router"

	"github.com/usual2970/certimate/internal/domain/dtos"
	"github.com/usual2970/certimate/internal/rest/resp"
)

type workflowService interface {
	Run(ctx context.Context, req *dtos.WorkflowRunReq) error
	Stop(ctx context.Context)
}

type WorkflowHandler struct {
	service workflowService
}

func NewWorkflowHandler(router *router.RouterGroup[*core.RequestEvent], service workflowService) {
	handler := &WorkflowHandler{
		service: service,
	}

	group := router.Group("/workflows")
	group.POST("/{id}/run", handler.run)
}

func (handler *WorkflowHandler) run(e *core.RequestEvent) error {
	req := &dtos.WorkflowRunReq{}
	req.WorkflowId = e.Request.PathValue("id")
	if err := e.BindBody(req); err != nil {
		return resp.Err(e, err)
	}

	if err := handler.service.Run(e.Request.Context(), req); err != nil {
		return resp.Err(e, err)
	}

	return resp.Ok(e, nil)
}
