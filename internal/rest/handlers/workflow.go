package handlers

import (
	"context"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/router"

	"github.com/certimate-go/certimate/internal/domain/dtos"
	"github.com/certimate-go/certimate/internal/rest/resp"
)

type workflowService interface {
	StartRun(ctx context.Context, req *dtos.WorkflowStartRunReq) error
	CancelRun(ctx context.Context, req *dtos.WorkflowCancelRunReq) error
	Shutdown(ctx context.Context)
}

type WorkflowHandler struct {
	service workflowService
}

func NewWorkflowHandler(router *router.RouterGroup[*core.RequestEvent], service workflowService) {
	handler := &WorkflowHandler{
		service: service,
	}

	group := router.Group("/workflows")
	group.POST("/{workflowId}/runs", handler.run)
	group.POST("/{workflowId}/runs/{runId}/cancel", handler.cancel)
}

func (handler *WorkflowHandler) run(e *core.RequestEvent) error {
	req := &dtos.WorkflowStartRunReq{}
	req.WorkflowId = e.Request.PathValue("workflowId")
	if err := e.BindBody(req); err != nil {
		return resp.Err(e, err)
	}

	if err := handler.service.StartRun(e.Request.Context(), req); err != nil {
		return resp.Err(e, err)
	}

	return resp.Ok(e, nil)
}

func (handler *WorkflowHandler) cancel(e *core.RequestEvent) error {
	req := &dtos.WorkflowCancelRunReq{}
	req.WorkflowId = e.Request.PathValue("workflowId")
	req.RunId = e.Request.PathValue("runId")

	if err := handler.service.CancelRun(e.Request.Context(), req); err != nil {
		return resp.Err(e, err)
	}

	return resp.Ok(e, nil)
}
