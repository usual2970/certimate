package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/usual2970/certimate/internal/app"
	"github.com/usual2970/certimate/internal/domain"
)

type CertificateRepository struct{}

func NewCertificateRepository() *CertificateRepository {
	return &CertificateRepository{}
}

func (r *CertificateRepository) ListExpireSoon(ctx context.Context) ([]*domain.Certificate, error) {
	records, err := app.GetApp().FindAllRecords(
		domain.CollectionNameCertificate,
		dbx.NewExp("expireAt>DATETIME('now')"),
		dbx.NewExp("expireAt<DATETIME('now', '+20 days')"),
		dbx.NewExp("deleted=null"),
	)
	if err != nil {
		return nil, err
	}

	certificates := make([]*domain.Certificate, 0)
	for _, record := range records {
		certificate, err := r.castRecordToModel(record)
		if err != nil {
			return nil, err
		}

		certificates = append(certificates, certificate)
	}

	return certificates, nil
}

func (r *CertificateRepository) GetById(ctx context.Context, id string) (*domain.Certificate, error) {
	record, err := app.GetApp().FindRecordById(domain.CollectionNameCertificate, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrRecordNotFound
		}
		return nil, err
	}

	if !record.GetDateTime("deleted").Time().IsZero() {
		return nil, domain.ErrRecordNotFound
	}

	return r.castRecordToModel(record)
}

func (r *CertificateRepository) GetByWorkflowNodeId(ctx context.Context, workflowNodeId string) (*domain.Certificate, error) {
	records, err := app.GetApp().FindRecordsByFilter(
		domain.CollectionNameCertificate,
		"workflowNodeId={:workflowNodeId} && deleted=null",
		"-created",
		1, 0,
		dbx.Params{"workflowNodeId": workflowNodeId},
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrRecordNotFound
		}
		return nil, err
	}
	if len(records) == 0 {
		return nil, domain.ErrRecordNotFound
	}

	return r.castRecordToModel(records[0])
}

func (r *CertificateRepository) castRecordToModel(record *core.Record) (*domain.Certificate, error) {
	if record == nil {
		return nil, fmt.Errorf("record is nil")
	}

	certificate := &domain.Certificate{
		Meta: domain.Meta{
			Id:        record.Id,
			CreatedAt: record.GetDateTime("created").Time(),
			UpdatedAt: record.GetDateTime("updated").Time(),
		},
		Source:            domain.CertificateSourceType(record.GetString("source")),
		SubjectAltNames:   record.GetString("subjectAltNames"),
		Certificate:       record.GetString("certificate"),
		PrivateKey:        record.GetString("privateKey"),
		IssuerCertificate: record.GetString("issuerCertificate"),
		EffectAt:          record.GetDateTime("effectAt").Time(),
		ExpireAt:          record.GetDateTime("expireAt").Time(),
		ACMECertUrl:       record.GetString("acmeCertUrl"),
		ACMECertStableUrl: record.GetString("acmeCertStableUrl"),
		WorkflowId:        record.GetString("workflowId"),
		WorkflowNodeId:    record.GetString("workflowNodeId"),
		WorkflowOutputId:  record.GetString("workflowOutputId"),
	}
	return certificate, nil
}
