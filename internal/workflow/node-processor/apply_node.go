package nodeprocessor

import (
	"context"
	"strings"
	"time"

	"github.com/usual2970/certimate/internal/applicant"
	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/pkg/utils/x509"
	"github.com/usual2970/certimate/internal/repository"
)

type applyNode struct {
	node       *domain.WorkflowNode
	outputRepo WorkflowOutputRepository
	*Logger
}

func NewApplyNode(node *domain.WorkflowNode) *applyNode {
	return &applyNode{
		node:       node,
		Logger:     NewLogger(node),
		outputRepo: repository.NewWorkflowOutputRepository(),
	}
}

type WorkflowOutputRepository interface {
	// 查询节点输出
	Get(ctx context.Context, nodeId string) (*domain.WorkflowOutput, error)

	// 查询申请节点的证书
	GetCertificate(ctx context.Context, nodeId string) (*domain.Certificate, error)

	// 保存节点输出
	Save(ctx context.Context, output *domain.WorkflowOutput, certificate *domain.Certificate, cb func(id string) error) error
}

// 申请节点根据申请类型执行不同的操作
func (a *applyNode) Run(ctx context.Context) error {
	a.AddOutput(ctx, a.node.Name, "开始执行")
	// 查询是否申请过，已申请过则直接返回（先保持和 v0.2 一致）
	output, err := a.outputRepo.Get(ctx, a.node.Id)
	if err != nil && !domain.IsRecordNotFound(err) {
		a.AddOutput(ctx, a.node.Name, "查询申请记录失败", err.Error())
		return err
	}

	if output != nil && output.Succeed {
		cert, err := a.outputRepo.GetCertificate(ctx, a.node.Id)
		if err != nil {
			a.AddOutput(ctx, a.node.Name, "获取证书失败", err.Error())
			return err
		}

		if time.Until(cert.ExpireAt) > domain.ValidityDuration {
			a.AddOutput(ctx, a.node.Name, "已申请过证书，且证书在有效期内")
			return nil
		}
	}

	// 获取Applicant
	apply, err := applicant.GetWithApplyNode(a.node)
	if err != nil {
		a.AddOutput(ctx, a.node.Name, "获取申请对象失败", err.Error())
		return err
	}

	// 申请
	certificate, err := apply.Apply()
	if err != nil {
		a.AddOutput(ctx, a.node.Name, "申请失败", err.Error())
		return err
	}
	a.AddOutput(ctx, a.node.Name, "申请成功")

	// 记录申请结果
	// 保持一个节点只有一个输出
	outputId := ""
	if output != nil {
		outputId = output.Id
	}
	output = &domain.WorkflowOutput{
		Workflow: GetWorkflowId(ctx),
		NodeId:   a.node.Id,
		Node:     a.node,
		Succeed:  true,
		Output:   a.node.Output,
		Meta:     domain.Meta{Id: outputId},
	}

	cert, err := x509.ParseCertificateFromPEM(certificate.Certificate)
	if err != nil {
		a.AddOutput(ctx, a.node.Name, "解析证书失败", err.Error())
		return err
	}

	certificateRecord := &domain.Certificate{
		SAN:               strings.Join(cert.DNSNames, ";"),
		Certificate:       certificate.Certificate,
		PrivateKey:        certificate.PrivateKey,
		IssuerCertificate: certificate.IssuerCertificate,
		CertUrl:           certificate.CertUrl,
		CertStableUrl:     certificate.CertStableUrl,
		ExpireAt:          cert.NotAfter,
		WorkflowId:        GetWorkflowId(ctx),
		WorkflowNodeId:    a.node.Id,
	}

	if err := a.outputRepo.Save(ctx, output, certificateRecord, func(id string) error {
		if certificateRecord != nil {
			certificateRecord.WorkflowOutputId = id
		}

		return nil
	}); err != nil {
		a.AddOutput(ctx, a.node.Name, "保存申请记录失败", err.Error())
		return err
	}

	a.AddOutput(ctx, a.node.Name, "保存申请记录成功")

	return nil
}
