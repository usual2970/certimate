package handlers

import (
	"context"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/router"

	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/rest/resp"
)

type certificateService interface {
	ArchiveFile(ctx context.Context, req *domain.CertificateArchiveFileReq) ([]byte, error)
}

type CertificateHandler struct {
	service certificateService
}

func NewCertificateHandler(router *router.RouterGroup[*core.RequestEvent], service certificateService) {
	handler := &CertificateHandler{
		service: service,
	}

	group := router.Group("/certificates")
	group.POST("/{id}/archive", handler.run)
}

func (handler *CertificateHandler) run(e *core.RequestEvent) error {
	req := &domain.CertificateArchiveFileReq{}
	req.CertificateId = e.Request.PathValue("id")
	if err := e.BindBody(req); err != nil {
		return resp.Err(e, err)
	}

	if bt, err := handler.service.ArchiveFile(e.Request.Context(), req); err != nil {
		return resp.Err(e, err)
	} else {
		return resp.Ok(e, bt)
	}
}
