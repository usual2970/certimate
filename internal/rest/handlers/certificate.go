package handlers

import (
	"context"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/router"

	"github.com/certimate-go/certimate/internal/domain/dtos"
	"github.com/certimate-go/certimate/internal/rest/resp"
)

type certificateService interface {
	ArchiveFile(ctx context.Context, req *dtos.CertificateArchiveFileReq) (*dtos.CertificateArchiveFileResp, error)
	ValidateCertificate(ctx context.Context, req *dtos.CertificateValidateCertificateReq) (*dtos.CertificateValidateCertificateResp, error)
	ValidatePrivateKey(ctx context.Context, req *dtos.CertificateValidatePrivateKeyReq) (*dtos.CertificateValidatePrivateKeyResp, error)
}

type CertificateHandler struct {
	service certificateService
}

func NewCertificateHandler(router *router.RouterGroup[*core.RequestEvent], service certificateService) {
	handler := &CertificateHandler{
		service: service,
	}

	group := router.Group("/certificates")
	group.POST("/{certificateId}/archive", handler.archiveFile)
	group.POST("/validate/certificate", handler.validateCertificate)
	group.POST("/validate/private-key", handler.validatePrivateKey)
}

func (handler *CertificateHandler) archiveFile(e *core.RequestEvent) error {
	req := &dtos.CertificateArchiveFileReq{}
	req.CertificateId = e.Request.PathValue("certificateId")
	if err := e.BindBody(req); err != nil {
		return resp.Err(e, err)
	}

	if res, err := handler.service.ArchiveFile(e.Request.Context(), req); err != nil {
		return resp.Err(e, err)
	} else {
		return resp.Ok(e, res)
	}
}

func (handler *CertificateHandler) validateCertificate(e *core.RequestEvent) error {
	req := &dtos.CertificateValidateCertificateReq{}
	if err := e.BindBody(req); err != nil {
		return resp.Err(e, err)
	}

	if res, err := handler.service.ValidateCertificate(e.Request.Context(), req); err != nil {
		return resp.Err(e, err)
	} else {
		return resp.Ok(e, res)
	}
}

func (handler *CertificateHandler) validatePrivateKey(e *core.RequestEvent) error {
	req := &dtos.CertificateValidatePrivateKeyReq{}
	if err := e.BindBody(req); err != nil {
		return resp.Err(e, err)
	}

	if res, err := handler.service.ValidatePrivateKey(e.Request.Context(), req); err != nil {
		return resp.Err(e, err)
	} else {
		return resp.Ok(e, res)
	}
}
