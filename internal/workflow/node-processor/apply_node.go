package nodeprocessor

import (
	"context"
	"fmt"
	"strings"
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

// 申请节点根据申请类型执行不同的操作
func (a *applyNode) Run(ctx context.Context) error {
	a.AddOutput(ctx, a.node.Name, "开始执行")

	// 查询上次执行结果
	lastOutput, err := a.outputRepo.GetByNodeId(ctx, a.node.Id)
	if err != nil && !domain.IsRecordNotFoundError(err) {
		a.AddOutput(ctx, a.node.Name, "查询申请记录失败", err.Error())
		return err
	}

	// 检测是否可以跳过本次执行
	if skippable, skipReason := a.checkCanSkip(ctx, lastOutput); skippable {
		a.AddOutput(ctx, a.node.Name, skipReason)
		return nil
	}

	// 初始化申请器
	applicant, err := applicant.NewWithApplyNode(a.node)
	if err != nil {
		a.AddOutput(ctx, a.node.Name, "获取申请对象失败", err.Error())
		return err
	}

	// 申请证书
	applyResult, err := applicant.Apply()
	if err != nil {
		a.AddOutput(ctx, a.node.Name, "申请失败", err.Error())
		return err
	}
	a.AddOutput(ctx, a.node.Name, "申请成功")

	// 解析证书并生成实体
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
		WorkflowId:        getContextWorkflowId(ctx),
		WorkflowNodeId:    a.node.Id,
	}

	// 保存执行结果
	// TODO: 先保持一个节点始终只有一个输出，后续增加版本控制
	currentOutput := &domain.WorkflowOutput{
		WorkflowId: getContextWorkflowId(ctx),
		NodeId:     a.node.Id,
		Node:       a.node,
		Succeeded:  true,
		Outputs:    a.node.Outputs,
	}
	if lastOutput != nil {
		currentOutput.Id = lastOutput.Id
	}
	if err := a.outputRepo.Save(ctx, currentOutput, certificate, func(id string) error {
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

func (a *applyNode) checkCanSkip(ctx context.Context, lastOutput *domain.WorkflowOutput) (skip bool, reason string) {
	if lastOutput != nil && lastOutput.Succeeded {
		// 比较和上次申请时的关键配置（即影响证书签发的）参数是否一致
		currentNodeConfig := a.node.GetConfigForApply()
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

		lastCertificate, _ := a.certRepo.GetByWorkflowNodeId(ctx, a.node.Id)
		renewalInterval := time.Duration(currentNodeConfig.SkipBeforeExpiryDays) * time.Hour * 24
		expirationTime := time.Until(lastCertificate.ExpireAt)
		if lastCertificate != nil && expirationTime > renewalInterval {
			return true, fmt.Sprintf("已申请过证书，且证书尚未临近过期（到期尚余 %d 天，预计距 %d 天时续期）", int(expirationTime.Hours()/24), currentNodeConfig.SkipBeforeExpiryDays)
		}
	}

	return false, ""
}
