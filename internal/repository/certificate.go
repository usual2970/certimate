package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/certimate-go/certimate/internal/app"
	"github.com/certimate-go/certimate/internal/domain"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
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
		return nil, err
	}

	if len(records) == 0 {
		return nil, domain.ErrRecordNotFound
	}

	return r.castRecordToModel(records[0])
}

func (r *CertificateRepository) GetByWorkflowRunIdAndNodeId(ctx context.Context, workflowRunId string, workflowNodeId string) (*domain.Certificate, error) {
	records, err := app.GetApp().FindRecordsByFilter(
		domain.CollectionNameCertificate,
		"workflowRunId={:workflowRunId} && workflowNodeId={:workflowNodeId} && deleted=null",
		"-created",
		1, 0,
		dbx.Params{"workflowRunId": workflowRunId},
		dbx.Params{"workflowNodeId": workflowNodeId},
	)
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, domain.ErrRecordNotFound
	}

	return r.castRecordToModel(records[0])
}

func (r *CertificateRepository) Save(ctx context.Context, certificate *domain.Certificate) (*domain.Certificate, error) {
	collection, err := app.GetApp().FindCollectionByNameOrId(domain.CollectionNameCertificate)
	if err != nil {
		return certificate, err
	}

	var record *core.Record
	if certificate.Id == "" {
		record = core.NewRecord(collection)
	} else {
		record, err = app.GetApp().FindRecordById(collection, certificate.Id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return certificate, domain.ErrRecordNotFound
			}
			return certificate, err
		}
	}

	record.Set("source", string(certificate.Source))
	record.Set("subjectAltNames", certificate.SubjectAltNames)
	record.Set("serialNumber", certificate.SerialNumber)
	record.Set("certificate", certificate.Certificate)
	record.Set("privateKey", certificate.PrivateKey)
	record.Set("issuerOrg", certificate.IssuerOrg)
	record.Set("issuerCertificate", certificate.IssuerCertificate)
	record.Set("keyAlgorithm", string(certificate.KeyAlgorithm))
	record.Set("effectAt", certificate.EffectAt)
	record.Set("expireAt", certificate.ExpireAt)
	record.Set("acmeAccountUrl", certificate.ACMEAccountUrl)
	record.Set("acmeCertUrl", certificate.ACMECertUrl)
	record.Set("acmeCertStableUrl", certificate.ACMECertStableUrl)
	record.Set("acmeRenewed", certificate.ACMERenewed)
	record.Set("workflowId", certificate.WorkflowId)
	record.Set("workflowRunId", certificate.WorkflowRunId)
	record.Set("workflowNodeId", certificate.WorkflowNodeId)
	record.Set("workflowOutputId", certificate.WorkflowOutputId)
	if err := app.GetApp().Save(record); err != nil {
		return certificate, err
	}

	certificate.Id = record.Id
	certificate.CreatedAt = record.GetDateTime("created").Time()
	certificate.UpdatedAt = record.GetDateTime("updated").Time()
	return certificate, nil
}

func (r *CertificateRepository) DeleteWhere(ctx context.Context, exprs ...dbx.Expression) (int, error) {
	records, err := app.GetApp().FindAllRecords(domain.CollectionNameCertificate, exprs...)
	if err != nil {
		return 0, nil
	}

	var ret int
	var errs []error
	for _, record := range records {
		if err := app.GetApp().Delete(record); err != nil {
			errs = append(errs, err)
		} else {
			ret++
		}
	}

	if len(errs) > 0 {
		return ret, errors.Join(errs...)
	}

	return ret, nil
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
		SerialNumber:      record.GetString("serialNumber"),
		Certificate:       record.GetString("certificate"),
		PrivateKey:        record.GetString("privateKey"),
		IssuerOrg:         record.GetString("issuerOrg"),
		IssuerCertificate: record.GetString("issuerCertificate"),
		KeyAlgorithm:      domain.CertificateKeyAlgorithmType(record.GetString("keyAlgorithm")),
		EffectAt:          record.GetDateTime("effectAt").Time(),
		ExpireAt:          record.GetDateTime("expireAt").Time(),
		ACMEAccountUrl:    record.GetString("acmeAccountUrl"),
		ACMECertUrl:       record.GetString("acmeCertUrl"),
		ACMECertStableUrl: record.GetString("acmeCertStableUrl"),
		ACMERenewed:       record.GetBool("acmeRenewed"),
		WorkflowId:        record.GetString("workflowId"),
		WorkflowRunId:     record.GetString("workflowRunId"),
		WorkflowNodeId:    record.GetString("workflowNodeId"),
		WorkflowOutputId:  record.GetString("workflowOutputId"),
	}
	return certificate, nil
}
