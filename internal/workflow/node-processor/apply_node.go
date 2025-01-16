package nodeprocessor

import (
	"context"
	"strings"
	"time"

	"github.com/usual2970/certimate/internal/applicant"
	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/pkg/utils/certs"
	"github.com/usual2970/certimate/internal/repository"
)

type applyNode struct {
	node       *domain.WorkflowNode
	outputRepo WorkflowOutputRepository
	*Logger
}

var validityDuration = time.Hour * 24 * 10

func NewApplyNode(node *domain.WorkflowNode) *applyNode {
	return &applyNode{
		node:       node,
		Logger:     NewLogger(node),
		outputRepo: repository.NewWorkflowOutputRepository(),
	}
}

type WorkflowOutputRepository interface {
	// 查询节点输出
	GetByNodeId(ctx context.Context, nodeId string) (*domain.WorkflowOutput, error)

	// 查询申请节点的证书
	GetCertificateByNodeId(ctx context.Context, nodeId string) (*domain.Certificate, error)

	// 保存节点输出
	Save(ctx context.Context, output *domain.WorkflowOutput, certificate *domain.Certificate, cb func(id string) error) error
}

// 申请节点根据申请类型执行不同的操作
func (a *applyNode) Run(ctx context.Context) error {
	a.AddOutput(ctx, a.node.Name, "开始执行")
	// 查询是否申请过，已申请过则直接返回
	// TODO: 先保持和 v0.2 一致，后续增加是否强制申请的参数
	output, err := a.outputRepo.GetByNodeId(ctx, a.node.Id)
	if err != nil && !domain.IsRecordNotFound(err) {
		a.AddOutput(ctx, a.node.Name, "查询申请记录失败", err.Error())
		return err
	}

	if output != nil && output.Succeeded {
		lastCertificate, _ := a.outputRepo.GetCertificateByNodeId(ctx, a.node.Id)
		if lastCertificate != nil {
			if time.Until(lastCertificate.ExpireAt) > validityDuration {
				a.AddOutput(ctx, a.node.Name, "已申请过证书，且证书在有效期内")
				return nil
			}
		}
	}

	// 获取Applicant
	applicant, err := applicant.NewWithApplyNode(a.node)
	if err != nil {
		a.AddOutput(ctx, a.node.Name, "获取申请对象失败", err.Error())
		return err
	}

	// 申请
	applyResult, err := applicant.Apply()
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
		Meta:       domain.Meta{Id: outputId},
		WorkflowId: GetWorkflowId(ctx),
		NodeId:     a.node.Id,
		Node:       a.node,
		Succeeded:  true,
		Outputs:    a.node.Outputs,
	}

	certX509, err := certs.ParseCertificateFromPEM(applyResult.CertificateFullChain)
	if err != nil {
		a.AddOutput(ctx, a.node.Name, "解析证书失败", err.Error())
		return err
	}

	certificate := &domain.Certificate{
		Source:            domain.CertificateSourceTypeWorkflow,
		SubjectAltNames:   strings.Join(certX509.DNSNames, ";"),
		Certificate:       applyResult.CertificateFullChain,
		PrivateKey:        applyResult.PrivateKey,
		IssuerCertificate: applyResult.IssuerCertificate,
		ACMECertUrl:       applyResult.ACMECertUrl,
		ACMECertStableUrl: applyResult.ACMECertStableUrl,
		EffectAt:          certX509.NotBefore,
		ExpireAt:          certX509.NotAfter,
		WorkflowId:        GetWorkflowId(ctx),
		WorkflowNodeId:    a.node.Id,
	}

	if err := a.outputRepo.Save(ctx, output, certificate, func(id string) error {
		if certificate != nil {
			certificate.WorkflowOutputId = id
		}

		return nil
	}); err != nil {
		a.AddOutput(ctx, a.node.Name, "保存申请记录失败", err.Error())
		return err
	}

	a.AddOutput(ctx, a.node.Name, "保存申请记录成功")

	return nil
}
