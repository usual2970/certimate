package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/models"
	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/utils/app"
)

type WorkflowOutputRepository struct{}

func NewWorkflowOutputRepository() *WorkflowOutputRepository {
	return &WorkflowOutputRepository{}
}

func (w *WorkflowOutputRepository) Get(ctx context.Context, nodeId string) (*domain.WorkflowOutput, error) {
	record, err := app.GetApp().Dao().FindFirstRecordByFilter("workflow_output", "nodeId={:nodeId}", dbx.Params{"nodeId": nodeId})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrRecordNotFound
		}
		return nil, err
	}

	node := &domain.WorkflowNode{}
	if err := record.UnmarshalJSONField("node", node); err != nil {
		return nil, errors.New("failed to unmarshal node")
	}

	output := make([]domain.WorkflowNodeIo, 0)
	if err := record.UnmarshalJSONField("output", &output); err != nil {
		return nil, errors.New("failed to unmarshal output")
	}

	rs := &domain.WorkflowOutput{
		Meta: domain.Meta{
			Id:      record.GetId(),
			Created: record.GetTime("created"),
			Updated: record.GetTime("updated"),
		},
		Workflow: record.GetString("workflow"),
		NodeId:   record.GetString("nodeId"),
		Node:     node,
		Output:   output,
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
	record.Set("workflow", output.Workflow)
	record.Set("nodeId", output.NodeId)
	record.Set("node", output.Node)
	record.Set("output", output.Output)

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
		certRecord.Set("certificate", certificate.Certificate)
		certRecord.Set("privateKey", certificate.PrivateKey)
		certRecord.Set("issuerCertificate", certificate.IssuerCertificate)
		certRecord.Set("san", certificate.SAN)
		certRecord.Set("workflowOutput", certificate.Output)
		certRecord.Set("expireAt", certificate.ExpireAt)
		certRecord.Set("certUrl", certificate.CertUrl)
		certRecord.Set("certStableUrl", certificate.CertStableUrl)

		if err := app.GetApp().Dao().SaveRecord(certRecord); err != nil {
			return err
		}
	}
	return nil
}
