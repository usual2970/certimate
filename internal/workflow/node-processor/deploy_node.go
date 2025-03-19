package nodeprocessor

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/usual2970/certimate/internal/deployer"
	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/repository"
	"golang.org/x/exp/maps"
)

type deployNode struct {
	node *domain.WorkflowNode
	*nodeProcessor

	certRepo   certificateRepository
	outputRepo workflowOutputRepository
}

func NewDeployNode(node *domain.WorkflowNode) *deployNode {
	return &deployNode{
		node:          node,
		nodeProcessor: newNodeProcessor(node),

		certRepo:   repository.NewCertificateRepository(),
		outputRepo: repository.NewWorkflowOutputRepository(),
	}
}

func (n *deployNode) Process(ctx context.Context) error {
	n.logger.Info("ready to deploy ...")

	// 查询上次执行结果
	lastOutput, err := n.outputRepo.GetByNodeId(ctx, n.node.Id)
	if err != nil && !domain.IsRecordNotFoundError(err) {
		return err
	}

	// 获取前序节点输出证书
	previousNodeOutputCertificateSource := n.node.GetConfigForDeploy().Certificate
	previousNodeOutputCertificateSourceSlice := strings.Split(previousNodeOutputCertificateSource, "#")
	if len(previousNodeOutputCertificateSourceSlice) != 2 {
		n.logger.Warn("invalid certificate source", slog.String("certificate.source", previousNodeOutputCertificateSource))
		return fmt.Errorf("invalid certificate source: %s", previousNodeOutputCertificateSource)
	}
	certificate, err := n.certRepo.GetByWorkflowNodeId(ctx, previousNodeOutputCertificateSourceSlice[0])
	if err != nil {
		n.logger.Warn("invalid certificate source", slog.String("certificate.source", previousNodeOutputCertificateSource))
		return err
	}

	// 检测是否可以跳过本次执行
	if lastOutput != nil && certificate.CreatedAt.Before(lastOutput.UpdatedAt) {
		if skippable, skipReason := n.checkCanSkip(ctx, lastOutput); skippable {
			n.logger.Info(fmt.Sprintf("skip this deployment, because %s", skipReason))
			return nil
		} else if skipReason != "" {
			n.logger.Info(fmt.Sprintf("continue to deploy, because %s", skipReason))
		}
	}

	// 初始化部署器
	deployer, err := deployer.NewWithDeployNode(n.node, struct {
		Certificate string
		PrivateKey  string
	}{Certificate: certificate.Certificate, PrivateKey: certificate.PrivateKey})
	if err != nil {
		n.logger.Warn("failed to create deployer provider")
		return err
	}

	// 部署证书
	deployer.SetLogger(n.logger)
	if err := deployer.Deploy(ctx); err != nil {
		n.logger.Warn("failed to deploy")
		return err
	}

	// 保存执行结果
	output := &domain.WorkflowOutput{
		WorkflowId: getContextWorkflowId(ctx),
		RunId:      getContextWorkflowRunId(ctx),
		NodeId:     n.node.Id,
		Node:       n.node,
		Succeeded:  true,
	}
	if _, err := n.outputRepo.Save(ctx, output); err != nil {
		n.logger.Warn("failed to save node output")
		return err
	}

	n.logger.Info("deploy completed")

	return nil
}

func (n *deployNode) checkCanSkip(ctx context.Context, lastOutput *domain.WorkflowOutput) (skip bool, reason string) {
	if lastOutput != nil && lastOutput.Succeeded {
		// 比较和上次部署时的关键配置（即影响证书部署的）参数是否一致
		currentNodeConfig := n.node.GetConfigForDeploy()
		lastNodeConfig := lastOutput.Node.GetConfigForDeploy()
		if currentNodeConfig.ProviderAccessId != lastNodeConfig.ProviderAccessId {
			return false, "the configuration item 'ProviderAccessId' changed"
		}
		if !maps.Equal(currentNodeConfig.ProviderConfig, lastNodeConfig.ProviderConfig) {
			return false, "the configuration item 'ProviderConfig' changed"
		}

		if currentNodeConfig.SkipOnLastSucceeded {
			return true, "the certificate has already been deployed"
		}
	}

	return false, ""
}
