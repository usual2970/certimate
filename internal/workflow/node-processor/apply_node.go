package nodeprocessor

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/exp/maps"

	"github.com/usual2970/certimate/internal/applicant"
	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/pkg/utils/certutil"
	"github.com/usual2970/certimate/internal/repository"
)

type applyNode struct {
	node *domain.WorkflowNode
	*nodeProcessor

	certRepo   certificateRepository
	outputRepo workflowOutputRepository
}

func NewApplyNode(node *domain.WorkflowNode) *applyNode {
	return &applyNode{
		node:          node,
		nodeProcessor: newNodeProcessor(node),

		certRepo:   repository.NewCertificateRepository(),
		outputRepo: repository.NewWorkflowOutputRepository(),
	}
}

func (n *applyNode) Process(ctx context.Context) error {
	n.logger.Info("ready to apply ...")

	// 查询上次执行结果
	lastOutput, err := n.outputRepo.GetByNodeId(ctx, n.node.Id)
	if err != nil && !domain.IsRecordNotFoundError(err) {
		return err
	}

	// 检测是否可以跳过本次执行
	if skippable, skipReason := n.checkCanSkip(ctx, lastOutput); skippable {
		n.logger.Info(fmt.Sprintf("skip this application, because %s", skipReason))
		return nil
	} else if skipReason != "" {
		n.logger.Info(fmt.Sprintf("re-apply, because %s", skipReason))
	}

	// 初始化申请器
	applicant, err := applicant.NewWithApplyNode(n.node)
	if err != nil {
		n.logger.Warn("failed to create applicant provider")
		return err
	}

	// 申请证书
	applyResult, err := applicant.Apply()
	if err != nil {
		n.logger.Warn("failed to apply")
		return err
	}

	// 解析证书并生成实体
	certX509, err := certutil.ParseCertificateFromPEM(applyResult.CertificateFullChain)
	if err != nil {
		n.logger.Warn("failed to parse certificate, may be the CA responded error")
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

	n.logger.Info("apply completed")

	return nil
}

func (n *applyNode) checkCanSkip(ctx context.Context, lastOutput *domain.WorkflowOutput) (skip bool, reason string) {
	if lastOutput != nil && lastOutput.Succeeded {
		// 比较和上次申请时的关键配置（即影响证书签发的）参数是否一致
		currentNodeConfig := n.node.GetConfigForApply()
		lastNodeConfig := lastOutput.Node.GetConfigForApply()
		if currentNodeConfig.Domains != lastNodeConfig.Domains {
			return false, "the configuration item 'Domains' changed"
		}
		if currentNodeConfig.ContactEmail != lastNodeConfig.ContactEmail {
			return false, "the configuration item 'ContactEmail' changed"
		}
		if currentNodeConfig.Provider != lastNodeConfig.Provider {
			return false, "the configuration item 'Provider' changed"
		}
		if currentNodeConfig.ProviderAccessId != lastNodeConfig.ProviderAccessId {
			return false, "the configuration item 'ProviderAccessId' changed"
		}
		if !maps.Equal(currentNodeConfig.ProviderConfig, lastNodeConfig.ProviderConfig) {
			return false, "the configuration item 'ProviderConfig' changed"
		}
		if currentNodeConfig.CAProvider != lastNodeConfig.CAProvider {
			return false, "the configuration item 'CAProvider' changed"
		}
		if currentNodeConfig.CAProviderAccessId != lastNodeConfig.CAProviderAccessId {
			return false, "the configuration item 'CAProviderAccessId' changed"
		}
		if !maps.Equal(currentNodeConfig.CAProviderConfig, lastNodeConfig.CAProviderConfig) {
			return false, "the configuration item 'CAProviderConfig' changed"
		}
		if currentNodeConfig.KeyAlgorithm != lastNodeConfig.KeyAlgorithm {
			return false, "the configuration item 'KeyAlgorithm' changed"
		}

		lastCertificate, _ := n.certRepo.GetByWorkflowNodeId(ctx, n.node.Id)
		if lastCertificate != nil {
			renewalInterval := time.Duration(currentNodeConfig.SkipBeforeExpiryDays) * time.Hour * 24
			expirationTime := time.Until(lastCertificate.ExpireAt)
			if expirationTime > renewalInterval {
				return true, fmt.Sprintf("the certificate has already been issued (expires in %dd, next renewal in %dd)", int(expirationTime.Hours()/24), currentNodeConfig.SkipBeforeExpiryDays)
			}
		}
	}

	return false, ""
}
