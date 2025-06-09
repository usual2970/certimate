package nodeprocessor

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/usual2970/certimate/internal/deployer"
	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/repository"
	"golang.org/x/exp/maps"
)

type deployNode struct {
	node *domain.WorkflowNode
	*nodeProcessor
	*nodeOutputer

	certRepo   certificateRepository
	outputRepo workflowOutputRepository
}

func NewDeployNode(node *domain.WorkflowNode) *deployNode {
	return &deployNode{
		node:          node,
		nodeProcessor: newNodeProcessor(node),
		nodeOutputer:  newNodeOutputer(),

		certRepo:   repository.NewCertificateRepository(),
		outputRepo: repository.NewWorkflowOutputRepository(),
	}
}

func (n *deployNode) Process(ctx context.Context) error {
	nodeCfg := n.node.GetConfigForDeploy()
	n.logger.Info("ready to deploy certificate ...", slog.Any("config", nodeCfg))

	// 查询上次执行结果
	lastOutput, err := n.outputRepo.GetByNodeId(ctx, n.node.Id)
	if err != nil && !domain.IsRecordNotFoundError(err) {
		return err
	}

	// 获取前序节点输出证书
	const DELIMITER = "#"
	previousNodeOutputCertificateSource := n.node.GetConfigForDeploy().Certificate
	previousNodeOutputCertificateSourceSlice := strings.Split(previousNodeOutputCertificateSource, DELIMITER)
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
		if skippable, reason := n.checkCanSkip(ctx, lastOutput); skippable {
			n.outputs[outputKeyForNodeSkipped] = strconv.FormatBool(true)
			n.logger.Info(fmt.Sprintf("skip this deployment, because %s", reason))
			return nil
		} else if reason != "" {
			n.logger.Info(fmt.Sprintf("re-deploy, because %s", reason))
		}
	}

	// 初始化部署器
	deployer, err := deployer.NewWithWorkflowNode(deployer.DeployerWithWorkflowNodeConfig{
		Node:           n.node,
		Logger:         n.logger,
		CertificatePEM: certificate.Certificate,
		PrivateKeyPEM:  certificate.PrivateKey,
	})
	if err != nil {
		n.logger.Warn("failed to create deployer provider")
		return err
	}

	// 部署证书
	if err := deployer.Deploy(ctx); err != nil {
		n.logger.Warn("failed to deploy certificate")
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

	// 记录中间结果
	n.outputs[outputKeyForNodeSkipped] = strconv.FormatBool(false)

	n.logger.Info("deployment completed")
	return nil
}

func (n *deployNode) checkCanSkip(ctx context.Context, lastOutput *domain.WorkflowOutput) (_skip bool, _reason string) {
	if lastOutput != nil && lastOutput.Succeeded {
		// 比较和上次部署时的关键配置（即影响证书部署的）参数是否一致
		thisNodeCfg := n.node.GetConfigForDeploy()
		lastNodeCfg := lastOutput.Node.GetConfigForDeploy()

		if thisNodeCfg.ProviderAccessId != lastNodeCfg.ProviderAccessId {
			return false, "the configuration item 'ProviderAccessId' changed"
		}
		if !maps.Equal(thisNodeCfg.ProviderConfig, lastNodeCfg.ProviderConfig) {
			return false, "the configuration item 'ProviderConfig' changed"
		}

		if thisNodeCfg.SkipOnLastSucceeded {
			return true, "the certificate has already been deployed"
		}
	}

	return false, ""
}
