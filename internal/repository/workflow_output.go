package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/models"
	"github.com/usual2970/certimate/internal/app"
	"github.com/usual2970/certimate/internal/domain"
)

type WorkflowOutputRepository struct{}

func NewWorkflowOutputRepository() *WorkflowOutputRepository {
	return &WorkflowOutputRepository{}
}

func (w *WorkflowOutputRepository) Get(ctx context.Context, nodeId string) (*domain.WorkflowOutput, error) {
	records, err := app.GetApp().Dao().FindRecordsByFilter("workflow_output", "nodeId={:nodeId}", "-created", 1, 0, dbx.Params{"nodeId": nodeId})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrRecordNotFound
		}
		return nil, err
	}
	if len(records) == 0 {
		return nil, domain.ErrRecordNotFound
	}
	record := records[0]

	node := &domain.WorkflowNode{}
	if err := record.UnmarshalJSONField("node", node); err != nil {
		return nil, errors.New("failed to unmarshal node")
	}

	outputs := make([]domain.WorkflowNodeIO, 0)
	if err := record.UnmarshalJSONField("outputs", &outputs); err != nil {
		return nil, errors.New("failed to unmarshal output")
	}

	rs := &domain.WorkflowOutput{
		Meta: domain.Meta{
			Id:        record.GetId(),
			CreatedAt: record.GetCreated().Time(),
			UpdatedAt: record.GetUpdated().Time(),
		},
		WorkflowId: record.GetString("workflowId"),
		NodeId:     record.GetString("nodeId"),
		Node:       node,
		Outputs:    outputs,
		Succeeded:  record.GetBool("succeeded"),
	}

	return rs, nil
}

func (w *WorkflowOutputRepository) GetCertificate(ctx context.Context, nodeId string) (*domain.Certificate, error) {
	records, err := app.GetApp().Dao().FindRecordsByFilter("certificate", "workflowNodeId={:workflowNodeId}", "-created", 1, 0, dbx.Params{"workflowNodeId": nodeId})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrRecordNotFound
		}
		return nil, err
	}
	if len(records) == 0 {
		return nil, domain.ErrRecordNotFound
	}

	record := records[0]

	rs := &domain.Certificate{
		Meta: domain.Meta{
			Id:        record.GetId(),
			CreatedAt: record.GetCreated().Time(),
			UpdatedAt: record.GetUpdated().Time(),
		},
		Source:            domain.CertificateSourceType(record.GetString("source")),
		SubjectAltNames:   record.GetString("subjectAltNames"),
		Certificate:       record.GetString("certificate"),
		PrivateKey:        record.GetString("privateKey"),
		IssuerCertificate: record.GetString("issuerCertificate"),
		EffectAt:          record.GetTime("effectAt"),
		ExpireAt:          record.GetTime("expireAt"),
		AcmeCertUrl:       record.GetString("acmeCertUrl"),
		AcmeCertStableUrl: record.GetString("acmeCertStableUrl"),
		WorkflowId:        record.GetString("workflowId"),
		WorkflowNodeId:    record.GetString("workflowNodeId"),
		WorkflowOutputId:  record.GetString("workflowOutputId"),
	}
	return rs, nil
}

// 保存节点输出
func (w *WorkflowOutputRepository) Save(ctx context.Context, output *domain.WorkflowOutput, certificate *domain.Certificate, cb func(id string) error) error {
	var record *models.Record
	var err error

	if output.Id == "" {
		collection, err := app.GetApp().Dao().FindCollectionByNameOrId("workflow_output")
		if err != nil {
			return err
		}
		record = models.NewRecord(collection)
	} else {
		record, err = app.GetApp().Dao().FindRecordById("workflow_output", output.Id)
		if err != nil {
			return err
		}
	}
	record.Set("workflowId", output.WorkflowId)
	record.Set("nodeId", output.NodeId)
	record.Set("node", output.Node)
	record.Set("outputs", output.Outputs)
	record.Set("succeeded", output.Succeeded)

	if err := app.GetApp().Dao().SaveRecord(record); err != nil {
		return err
	}

	if cb != nil && certificate != nil {
		if err := cb(record.GetId()); err != nil {
			return err
		}

		certCollection, err := app.GetApp().Dao().FindCollectionByNameOrId("certificate")
		if err != nil {
			return err
		}

		certRecord := models.NewRecord(certCollection)
		certRecord.Set("source", certificate.Source)
		certRecord.Set("subjectAltNames", certificate.SubjectAltNames)
		certRecord.Set("certificate", certificate.Certificate)
		certRecord.Set("privateKey", certificate.PrivateKey)
		certRecord.Set("issuerCertificate", certificate.IssuerCertificate)
		certRecord.Set("effectAt", certificate.EffectAt)
		certRecord.Set("expireAt", certificate.ExpireAt)
		certRecord.Set("acmeCertUrl", certificate.AcmeCertUrl)
		certRecord.Set("acmeCertStableUrl", certificate.AcmeCertStableUrl)
		certRecord.Set("workflowId", certificate.WorkflowId)
		certRecord.Set("workflowNodeId", certificate.WorkflowNodeId)
		certRecord.Set("workflowOutputId", certificate.WorkflowOutputId)

		if err := app.GetApp().Dao().SaveRecord(certRecord); err != nil {
			return err
		}

		// 更新 certificate
		for i, item := range output.Outputs {
			if item.Name == "certificate" {
				output.Outputs[i].Value = certRecord.GetId()
				break
			}
		}

		record.Set("outputs", output.Outputs)

		if err := app.GetApp().Dao().SaveRecord(record); err != nil {
			return err
		}

	}
	return nil
}
