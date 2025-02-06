package nodeprocessor

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/exp/maps"

	"github.com/usual2970/certimate/internal/applicant"
	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/pkg/utils/certs"
	"github.com/usual2970/certimate/internal/repository"
)

type applyNode struct {
	node       *domain.WorkflowNode
	certRepo   certificateRepository
	outputRepo workflowOutputRepository
	*nodeLogger
}

func NewApplyNode(node *domain.WorkflowNode) *applyNode {
	return &applyNode{
		node:       node,
		nodeLogger: NewNodeLogger(node),
		outputRepo: repository.NewWorkflowOutputRepository(),
		certRepo:   repository.NewCertificateRepository(),
	}
}

func (n *applyNode) Run(ctx context.Context) error {
	n.AddOutput(ctx, n.node.Name, "开始执行")

	// 查询上次执行结果
	lastOutput, err := n.outputRepo.GetByNodeId(ctx, n.node.Id)
	if err != nil && !domain.IsRecordNotFoundError(err) {
		n.AddOutput(ctx, n.node.Name, "查询申请记录失败", err.Error())
		return err
	}

	// 检测是否可以跳过本次执行
	if skippable, skipReason := n.checkCanSkip(ctx, lastOutput); skippable {
		n.AddOutput(ctx, n.node.Name, skipReason)
		return nil
	}

	// 初始化申请器
	applicant, err := applicant.NewWithApplyNode(n.node)
	if err != nil {
		n.AddOutput(ctx, n.node.Name, "获取申请对象失败", err.Error())
		return err
	}

	// 申请证书
	applyResult, err := applicant.Apply()
	if err != nil {
		n.AddOutput(ctx, n.node.Name, "申请失败", err.Error())
		return err
	}
	n.AddOutput(ctx, n.node.Name, "申请成功")

	// 解析证书并生成实体
	certX509, err := certs.ParseCertificateFromPEM(applyResult.CertificateFullChain)
	if err != nil {
		n.AddOutput(ctx, n.node.Name, "解析证书失败", err.Error())
		return err
	}
	certificate := &domain.Certificate{
		Source:            domain.CertificateSourceTypeWorkflow,
		Certificate:       applyResult.CertificateFullChain,
		PrivateKey:        applyResult.PrivateKey,
		IssuerCertificate: applyResult.IssuerCertificate,
		ACMEAccountUrl:    applyResult.ACMEAccountUrl,
		ACMECertUrl:       applyResult.ACMECertUrl,
		ACMECertStableUrl: applyResult.ACMECertStableUrl,
	}
	certificate.PopulateFromX509(certX509)

	// 保存执行结果
	// TODO: 先保持一个节点始终只有一个输出，后续增加版本控制
	currentOutput := &domain.WorkflowOutput{
		WorkflowId: getContextWorkflowId(ctx),
		NodeId:     n.node.Id,
		Node:       n.node,
		Succeeded:  true,
		Outputs:    n.node.Outputs,
	}
	if lastOutput != nil {
		currentOutput.Id = lastOutput.Id
	}
	if _, err := n.outputRepo.SaveWithCertificate(ctx, currentOutput, certificate); err != nil {
		n.AddOutput(ctx, n.node.Name, "保存申请记录失败", err.Error())
		return err
	}
	n.AddOutput(ctx, n.node.Name, "保存申请记录成功")

	return nil
}

func (n *applyNode) checkCanSkip(ctx context.Context, lastOutput *domain.WorkflowOutput) (skip bool, reason string) {
	if lastOutput != nil && lastOutput.Succeeded {
		// 比较和上次申请时的关键配置（即影响证书签发的）参数是否一致
		currentNodeConfig := n.node.GetConfigForApply()
		lastNodeConfig := lastOutput.Node.GetConfigForApply()
		if currentNodeConfig.Domains != lastNodeConfig.Domains {
			return false, "配置项变化：域名"
		}
		if currentNodeConfig.ContactEmail != lastNodeConfig.ContactEmail {
			return false, "配置项变化：联系邮箱"
		}
		if currentNodeConfig.ProviderAccessId != lastNodeConfig.ProviderAccessId {
			return false, "配置项变化：DNS 提供商授权"
		}
		if !maps.Equal(currentNodeConfig.ProviderConfig, lastNodeConfig.ProviderConfig) {
			return false, "配置项变化：DNS 提供商参数"
		}
		if currentNodeConfig.KeyAlgorithm != lastNodeConfig.KeyAlgorithm {
			return false, "配置项变化：数字签名算法"
		}

		lastCertificate, _ := n.certRepo.GetByWorkflowNodeId(ctx, n.node.Id)
		if lastCertificate != nil {
			renewalInterval := time.Duration(currentNodeConfig.SkipBeforeExpiryDays) * time.Hour * 24
			expirationTime := time.Until(lastCertificate.ExpireAt)
			if expirationTime > renewalInterval {
				return true, fmt.Sprintf("已申请过证书，且证书尚未临近过期（到期尚余 %d 天，预计距 %d 天时续期）", int(expirationTime.Hours()/24), currentNodeConfig.SkipBeforeExpiryDays)
			}
		}
	}

	return false, ""
}
