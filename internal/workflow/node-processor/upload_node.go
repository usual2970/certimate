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
	node *domain.WorkflowNode
	*nodeLogger

	certRepo   certificateRepository
	outputRepo workflowOutputRepository
}

func NewUploadNode(node *domain.WorkflowNode) *uploadNode {
	return &uploadNode{
		node:       node,
		nodeLogger: newNodeLogger(node),

		certRepo:   repository.NewCertificateRepository(),
		outputRepo: repository.NewWorkflowOutputRepository(),
	}
}

func (n *uploadNode) Process(ctx context.Context) error {
	n.AppendLogRecord(ctx, domain.WorkflowRunLogLevelInfo, "进入上传证书节点")

	nodeConfig := n.node.GetConfigForUpload()

	// 查询上次执行结果
	lastOutput, err := n.outputRepo.GetByNodeId(ctx, n.node.Id)
	if err != nil && !domain.IsRecordNotFoundError(err) {
		n.AppendLogRecord(ctx, domain.WorkflowRunLogLevelError, "查询申请记录失败", err.Error())
		return err
	}

	// 检测是否可以跳过本次执行
	if skippable, skipReason := n.checkCanSkip(ctx, lastOutput); skippable {
		n.AppendLogRecord(ctx, domain.WorkflowRunLogLevelInfo, skipReason)
		return nil
	}

	// 检查证书是否过期
	// 如果证书过期，则直接返回错误
	certX509, err := certs.ParseCertificateFromPEM(nodeConfig.Certificate)
	if err != nil {
		n.AppendLogRecord(ctx, domain.WorkflowRunLogLevelError, "解析证书失败")
		return err
	}
	if time.Now().After(certX509.NotAfter) {
		n.AppendLogRecord(ctx, domain.WorkflowRunLogLevelWarn, "证书已过期")
		return errors.New("certificate is expired")
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
		n.AppendLogRecord(ctx, domain.WorkflowRunLogLevelError, "保存上传记录失败", err.Error())
		return err
	}
	n.AppendLogRecord(ctx, domain.WorkflowRunLogLevelInfo, "保存上传记录成功")

	return nil
}

func (n *uploadNode) checkCanSkip(ctx context.Context, lastOutput *domain.WorkflowOutput) (skip bool, reason string) {
	if lastOutput != nil && lastOutput.Succeeded {
		// 比较和上次上传时的关键配置（即影响证书上传的）参数是否一致
		currentNodeConfig := n.node.GetConfigForUpload()
		lastNodeConfig := lastOutput.Node.GetConfigForUpload()
		if strings.TrimSpace(currentNodeConfig.Certificate) != strings.TrimSpace(lastNodeConfig.Certificate) {
			return false, "配置项变化：证书"
		}
		if strings.TrimSpace(currentNodeConfig.PrivateKey) != strings.TrimSpace(lastNodeConfig.PrivateKey) {
			return false, "配置项变化：私钥"
		}

		lastCertificate, _ := n.certRepo.GetByWorkflowNodeId(ctx, n.node.Id)
		if lastCertificate != nil {
			return true, "已上传过证书"
		}
	}

	return false, ""
}
