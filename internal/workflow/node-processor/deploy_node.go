package nodeprocessor

import (
	"context"
	"fmt"
	"strings"

	"github.com/usual2970/certimate/internal/deployer"
	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/repository"
	"golang.org/x/exp/maps"
)

type deployNode struct {
	node *domain.WorkflowNode
	*nodeLogger

	certRepo   certificateRepository
	outputRepo workflowOutputRepository
}

func NewDeployNode(node *domain.WorkflowNode) *deployNode {
	return &deployNode{
		node:       node,
		nodeLogger: NewNodeLogger(node),

		certRepo:   repository.NewCertificateRepository(),
		outputRepo: repository.NewWorkflowOutputRepository(),
	}
}

func (n *deployNode) Run(ctx context.Context) error {
	n.AddOutput(ctx, n.node.Name, "开始执行")

	// 查询上次执行结果
	lastOutput, err := n.outputRepo.GetByNodeId(ctx, n.node.Id)
	if err != nil && !domain.IsRecordNotFoundError(err) {
		n.AddOutput(ctx, n.node.Name, "查询部署记录失败", err.Error())
		return err
	}

	// 获取前序节点输出证书
	previousNodeOutputCertificateSource := n.node.GetConfigForDeploy().Certificate
	previousNodeOutputCertificateSourceSlice := strings.Split(previousNodeOutputCertificateSource, "#")
	if len(previousNodeOutputCertificateSourceSlice) != 2 {
		n.AddOutput(ctx, n.node.Name, "证书来源配置错误", previousNodeOutputCertificateSource)
		return fmt.Errorf("证书来源配置错误: %s", previousNodeOutputCertificateSource)
	}
	certificate, err := n.certRepo.GetByWorkflowNodeId(ctx, previousNodeOutputCertificateSourceSlice[0])
	if err != nil {
		n.AddOutput(ctx, n.node.Name, "获取证书失败", err.Error())
		return err
	}

	// 检测是否可以跳过本次执行
	if certificate.CreatedAt.Before(lastOutput.UpdatedAt) {
		if skippable, skipReason := n.checkCanSkip(ctx, lastOutput); skippable {
			n.AddOutput(ctx, n.node.Name, skipReason)
			return nil
		}
	}

	// 初始化部署器
	deployer, err := deployer.NewWithDeployNode(n.node, struct {
		Certificate string
		PrivateKey  string
	}{Certificate: certificate.Certificate, PrivateKey: certificate.PrivateKey})
	if err != nil {
		n.AddOutput(ctx, n.node.Name, "获取部署对象失败", err.Error())
		return err
	}

	// 部署证书
	if err := deployer.Deploy(ctx); err != nil {
		n.AddOutput(ctx, n.node.Name, "部署失败", err.Error())
		return err
	}
	n.AddOutput(ctx, n.node.Name, "部署成功")

	// 保存执行结果
	output := &domain.WorkflowOutput{
		WorkflowId: getContextWorkflowId(ctx),
		RunId:      getContextWorkflowRunId(ctx),
		NodeId:     n.node.Id,
		Node:       n.node,
		Succeeded:  true,
	}
	if _, err := n.outputRepo.Save(ctx, output); err != nil {
		n.AddOutput(ctx, n.node.Name, "保存部署记录失败", err.Error())
		return err
	}
	n.AddOutput(ctx, n.node.Name, "保存部署记录成功")

	return nil
}

func (n *deployNode) checkCanSkip(ctx context.Context, lastOutput *domain.WorkflowOutput) (skip bool, reason string) {
	if lastOutput != nil && lastOutput.Succeeded {
		// 比较和上次部署时的关键配置（即影响证书部署的）参数是否一致
		currentNodeConfig := n.node.GetConfigForDeploy()
		lastNodeConfig := lastOutput.Node.GetConfigForDeploy()
		if currentNodeConfig.ProviderAccessId != lastNodeConfig.ProviderAccessId {
			return false, "配置项变化：主机提供商授权"
		}
		if !maps.Equal(currentNodeConfig.ProviderConfig, lastNodeConfig.ProviderConfig) {
			return false, "配置项变化：主机提供商参数"
		}

		if currentNodeConfig.SkipOnLastSucceeded {
			return true, "已部署过证书"
		}
	}

	return false, ""
}
