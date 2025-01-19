package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/usual2970/certimate/internal/app"
	"github.com/usual2970/certimate/internal/domain"
)

type WorkflowOutputRepository struct{}

func NewWorkflowOutputRepository() *WorkflowOutputRepository {
	return &WorkflowOutputRepository{}
}

func (r *WorkflowOutputRepository) GetByNodeId(ctx context.Context, nodeId string) (*domain.WorkflowOutput, error) {
	records, err := app.GetApp().FindRecordsByFilter(
		domain.CollectionNameWorkflowOutput,
		"nodeId={:nodeId}",
		"-created",
		1, 0,
		dbx.Params{"nodeId": nodeId},
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
			Id:        record.Id,
			CreatedAt: record.GetDateTime("created").Time(),
			UpdatedAt: record.GetDateTime("updated").Time(),
		},
		WorkflowId: record.GetString("workflowId"),
		NodeId:     record.GetString("nodeId"),
		Node:       node,
		Outputs:    outputs,
		Succeeded:  record.GetBool("succeeded"),
	}

	return rs, nil
}

// 保存节点输出
func (r *WorkflowOutputRepository) Save(ctx context.Context, output *domain.WorkflowOutput, certificate *domain.Certificate, cb func(id string) error) error {
	var record *core.Record
	var err error

	if output.Id == "" {
		collection, err := app.GetApp().FindCollectionByNameOrId(domain.CollectionNameWorkflowOutput)
		if err != nil {
			return err
		}
		record = core.NewRecord(collection)
	} else {
		record, err = app.GetApp().FindRecordById(domain.CollectionNameWorkflowOutput, output.Id)
		if err != nil {
			return err
		}
	}
	record.Set("workflowId", output.WorkflowId)
	record.Set("nodeId", output.NodeId)
	record.Set("node", output.Node)
	record.Set("outputs", output.Outputs)
	record.Set("succeeded", output.Succeeded)

	if err := app.GetApp().Save(record); err != nil {
		return err
	}

	if cb != nil && certificate != nil {
		if err := cb(record.Id); err != nil {
			return err
		}

		certCollection, err := app.GetApp().FindCollectionByNameOrId(domain.CollectionNameCertificate)
		if err != nil {
			return err
		}

		certRecord := core.NewRecord(certCollection)
		certRecord.Set("source", string(certificate.Source))
		certRecord.Set("subjectAltNames", certificate.SubjectAltNames)
		certRecord.Set("certificate", certificate.Certificate)
		certRecord.Set("privateKey", certificate.PrivateKey)
		certRecord.Set("issuerCertificate", certificate.IssuerCertificate)
		certRecord.Set("effectAt", certificate.EffectAt)
		certRecord.Set("expireAt", certificate.ExpireAt)
		certRecord.Set("acmeCertUrl", certificate.ACMECertUrl)
		certRecord.Set("acmeCertStableUrl", certificate.ACMECertStableUrl)
		certRecord.Set("workflowId", certificate.WorkflowId)
		certRecord.Set("workflowNodeId", certificate.WorkflowNodeId)
		certRecord.Set("workflowOutputId", certificate.WorkflowOutputId)

		if err := app.GetApp().Save(certRecord); err != nil {
			return err
		}

		// 更新 certificate
		for i, item := range output.Outputs {
			if item.Name == string(domain.WorkflowNodeIONameCertificate) {
				output.Outputs[i].Value = certRecord.Id
				break
			}
		}

		record.Set("outputs", output.Outputs)

		if err := app.GetApp().Save(record); err != nil {
			return err
		}

	}
	return nil
}
