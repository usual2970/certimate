package nodeprocessor

import (
	"context"
	"fmt"
	"strings"

	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/repository"
)

type uploadNode struct {
	node *domain.WorkflowNode
	*nodeProcessor

	certRepo   certificateRepository
	outputRepo workflowOutputRepository
}

func NewUploadNode(node *domain.WorkflowNode) *uploadNode {
	return &uploadNode{
		node:          node,
		nodeProcessor: newNodeProcessor(node),

		certRepo:   repository.NewCertificateRepository(),
		outputRepo: repository.NewWorkflowOutputRepository(),
	}
}

func (n *uploadNode) Process(ctx context.Context) error {
	n.logger.Info("ready to upload ...")

	nodeConfig := n.node.GetConfigForUpload()

	// 查询上次执行结果
	lastOutput, err := n.outputRepo.GetByNodeId(ctx, n.node.Id)
	if err != nil && !domain.IsRecordNotFoundError(err) {
		return err
	}

	// 检测是否可以跳过本次执行
	if skippable, skipReason := n.checkCanSkip(ctx, lastOutput); skippable {
		n.logger.Info(fmt.Sprintf("skip this upload, because %s", skipReason))
		return nil
	} else if skipReason != "" {
		n.logger.Info(fmt.Sprintf("re-upload, because %s", skipReason))
	}

	// 生成证书实体
	certificate := &domain.Certificate{
		Source: domain.CertificateSourceTypeUpload,
	}
	certificate.PopulateFromPEM(nodeConfig.Certificate, nodeConfig.PrivateKey)

	// 保存执行结果
	output := &domain.WorkflowOutput{
		WorkflowId: getContextWorkflowId(ctx),
		RunId:      getContextWorkflowRunId(ctx),
		NodeId:     n.node.Id,
		Node:       n.node,
		Succeeded:  true,
		Outputs:    n.node.Outputs,
	}
	if _, err := n.outputRepo.SaveWithCertificate(ctx, output, certificate); err != nil {
		n.logger.Warn("failed to save node output")
		return err
	}

	n.logger.Info("upload completed")

	return nil
}

func (n *uploadNode) checkCanSkip(ctx context.Context, lastOutput *domain.WorkflowOutput) (skip bool, reason string) {
	if lastOutput != nil && lastOutput.Succeeded {
		// 比较和上次上传时的关键配置（即影响证书上传的）参数是否一致
		currentNodeConfig := n.node.GetConfigForUpload()
		lastNodeConfig := lastOutput.Node.GetConfigForUpload()
		if strings.TrimSpace(currentNodeConfig.Certificate) != strings.TrimSpace(lastNodeConfig.Certificate) {
			return false, "the configuration item 'Certificate' changed"
		}
		if strings.TrimSpace(currentNodeConfig.PrivateKey) != strings.TrimSpace(lastNodeConfig.PrivateKey) {
			return false, "the configuration item 'PrivateKey' changed"
		}

		lastCertificate, _ := n.certRepo.GetByWorkflowNodeId(ctx, n.node.Id)
		if lastCertificate != nil {
			return true, "the certificate has already been uploaded"
		}
	}

	return false, ""
}
