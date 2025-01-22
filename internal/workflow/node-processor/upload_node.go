package nodeprocessor

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/pkg/utils/certs"
	"github.com/usual2970/certimate/internal/repository"
)

type uploadNode struct {
	node       *domain.WorkflowNode
	outputRepo workflowOutputRepository
	*nodeLogger
}

func NewUploadNode(node *domain.WorkflowNode) *uploadNode {
	return &uploadNode{
		node:       node,
		nodeLogger: NewNodeLogger(node),
		outputRepo: repository.NewWorkflowOutputRepository(),
	}
}

// Run 上传证书节点执行
// 包含上传证书的工作流，理论上应该手动执行，如果每天定时执行，也只是重新保存一下
func (n *uploadNode) Run(ctx context.Context) error {
	n.AddOutput(ctx,
		n.node.Name,
		"进入上传证书节点",
	)

	config := n.node.GetConfigForUpload()

	// 检查证书是否过期
	// 如果证书过期，则直接返回错误
	certX509, err := certs.ParseCertificateFromPEM(config.Certificate)
	if err != nil {
		n.AddOutput(ctx,
			n.node.Name,
			"解析证书失败",
		)
		return err
	}

	if time.Now().After(certX509.NotAfter) {
		n.AddOutput(ctx,
			n.node.Name,
			"证书已过期",
		)
		return errors.New("certificate is expired")
	}

	certificate := &domain.Certificate{
		Source:          domain.CertificateSourceTypeUpload,
		SubjectAltNames: strings.Join(certX509.DNSNames, ";"),
		Certificate:     config.Certificate,
		PrivateKey:      config.PrivateKey,

		EffectAt:       certX509.NotBefore,
		ExpireAt:       certX509.NotAfter,
		WorkflowId:     getContextWorkflowId(ctx),
		WorkflowNodeId: n.node.Id,
	}

	// 保存执行结果
	// TODO: 先保持一个节点始终只有一个输出，后续增加版本控制
	currentOutput := &domain.WorkflowOutput{
		WorkflowId: getContextWorkflowId(ctx),
		NodeId:     n.node.Id,
		Node:       n.node,
		Succeeded:  true,
		Outputs:    n.node.Outputs,
	}

	// 查询上次执行结果
	lastOutput, err := n.outputRepo.GetByNodeId(ctx, n.node.Id)
	if err != nil && !domain.IsRecordNotFoundError(err) {
		n.AddOutput(ctx, n.node.Name, "查询上传记录失败", err.Error())
		return err
	}
	if lastOutput != nil {
		currentOutput.Id = lastOutput.Id
	}
	if err := n.outputRepo.Save(ctx, currentOutput, certificate, func(id string) error {
		if certificate != nil {
			certificate.WorkflowOutputId = id
		}

		return nil
	}); err != nil {
		n.AddOutput(ctx, n.node.Name, "保存上传记录失败", err.Error())
		return err
	}
	n.AddOutput(ctx, n.node.Name, "保存上传记录成功")

	return nil
}
