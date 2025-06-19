package nodeprocessor

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/internal/repository"
)

type uploadNode struct {
	node *domain.WorkflowNode
	*nodeProcessor
	*nodeOutputer

	certRepo   certificateRepository
	outputRepo workflowOutputRepository
}

func NewUploadNode(node *domain.WorkflowNode) *uploadNode {
	return &uploadNode{
		node:          node,
		nodeProcessor: newNodeProcessor(node),
		nodeOutputer:  newNodeOutputer(),

		certRepo:   repository.NewCertificateRepository(),
		outputRepo: repository.NewWorkflowOutputRepository(),
	}
}

func (n *uploadNode) Process(ctx context.Context) error {
	nodeCfg := n.node.GetConfigForUpload()
	n.logger.Info("ready to upload certiticate ...", slog.Any("config", nodeCfg))

	// 查询上次执行结果
	lastOutput, err := n.outputRepo.GetByNodeId(ctx, n.node.Id)
	if err != nil && !domain.IsRecordNotFoundError(err) {
		return err
	}

	// 检测是否可以跳过本次执行
	if skippable, reason := n.checkCanSkip(ctx, lastOutput); skippable {
		n.outputs[outputKeyForNodeSkipped] = strconv.FormatBool(true)
		n.logger.Info(fmt.Sprintf("skip this uploading, because %s", reason))
		return nil
	} else if reason != "" {
		n.logger.Info(fmt.Sprintf("re-upload, because %s", reason))
	}

	// 生成证书实体
	certificate := &domain.Certificate{
		Source: domain.CertificateSourceTypeUpload,
	}
	certificate.PopulateFromPEM(nodeCfg.Certificate, nodeCfg.PrivateKey)

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

	// 记录中间结果
	n.outputs[outputKeyForNodeSkipped] = strconv.FormatBool(false)
	n.outputs[outputKeyForCertificateValidity] = strconv.FormatBool(true)
	n.outputs[outputKeyForCertificateDaysLeft] = strconv.FormatInt(int64(time.Until(certificate.ExpireAt).Hours()/24), 10)

	n.logger.Info("uploading completed")
	return nil
}

func (n *uploadNode) checkCanSkip(ctx context.Context, lastOutput *domain.WorkflowOutput) (_skip bool, _reason string) {
	if lastOutput != nil && lastOutput.Succeeded {
		// 比较和上次上传时的关键配置（即影响证书上传的）参数是否一致
		thisNodeCfg := n.node.GetConfigForUpload()
		lastNodeCfg := lastOutput.Node.GetConfigForUpload()

		if strings.TrimSpace(thisNodeCfg.Certificate) != strings.TrimSpace(lastNodeCfg.Certificate) {
			return false, "the configuration item 'Certificate' changed"
		}
		if strings.TrimSpace(thisNodeCfg.PrivateKey) != strings.TrimSpace(lastNodeCfg.PrivateKey) {
			return false, "the configuration item 'PrivateKey' changed"
		}

		lastCertificate, _ := n.certRepo.GetByWorkflowRunIdAndNodeId(ctx, lastOutput.RunId, lastOutput.NodeId)
		if lastCertificate != nil {
			daysLeft := int(time.Until(lastCertificate.ExpireAt).Hours() / 24)
			n.outputs[outputKeyForCertificateValidity] = strconv.FormatBool(daysLeft > 0)
			n.outputs[outputKeyForCertificateDaysLeft] = strconv.FormatInt(int64(daysLeft), 10)

			return true, "the certificate has already been uploaded"
		}
	}

	return false, ""
}
