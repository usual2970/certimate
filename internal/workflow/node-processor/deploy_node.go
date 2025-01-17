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
	node       *domain.WorkflowNode
	certRepo   certificateRepository
	outputRepo workflowOutputRepository
	*nodeLogger
}

func NewDeployNode(node *domain.WorkflowNode) *deployNode {
	return &deployNode{
		node:       node,
		nodeLogger: NewNodeLogger(node),
		outputRepo: repository.NewWorkflowOutputRepository(),
		certRepo:   repository.NewCertificateRepository(),
	}
}

func (d *deployNode) Run(ctx context.Context) error {
	d.AddOutput(ctx, d.node.Name, "开始执行")

	// 查询上次执行结果
	lastOutput, err := d.outputRepo.GetByNodeId(ctx, d.node.Id)
	if err != nil && !domain.IsRecordNotFoundError(err) {
		d.AddOutput(ctx, d.node.Name, "查询部署记录失败", err.Error())
		return err
	}

	// 获取前序节点输出证书
	previousNodeOutputCertificateSource := d.node.GetConfigForDeploy().Certificate
	previousNodeOutputCertificateSourceSlice := strings.Split(previousNodeOutputCertificateSource, "#")
	if len(previousNodeOutputCertificateSourceSlice) != 2 {
		d.AddOutput(ctx, d.node.Name, "证书来源配置错误", previousNodeOutputCertificateSource)
		return fmt.Errorf("证书来源配置错误: %s", previousNodeOutputCertificateSource)
	}
	certificate, err := d.certRepo.GetByWorkflowNodeId(ctx, previousNodeOutputCertificateSourceSlice[0])
	if err != nil {
		d.AddOutput(ctx, d.node.Name, "获取证书失败", err.Error())
		return err
	}

	// 检测是否可以跳过本次执行
	if skippable, skipReason := d.checkCanSkip(ctx, lastOutput); skippable {
		if certificate.CreatedAt.Before(lastOutput.UpdatedAt) {
			d.AddOutput(ctx, d.node.Name, "已部署过且证书未更新")
		} else {
			d.AddOutput(ctx, d.node.Name, skipReason)
		}
		return nil
	}

	// 初始化部署器
	deploy, err := deployer.NewWithDeployNode(d.node, struct {
		Certificate string
		PrivateKey  string
	}{Certificate: certificate.Certificate, PrivateKey: certificate.PrivateKey})
	if err != nil {
		d.AddOutput(ctx, d.node.Name, "获取部署对象失败", err.Error())
		return err
	}

	// 部署证书
	if err := deploy.Deploy(ctx); err != nil {
		d.AddOutput(ctx, d.node.Name, "部署失败", err.Error())
		return err
	}
	d.AddOutput(ctx, d.node.Name, "部署成功")

	// 保存执行结果
	// TODO: 先保持一个节点始终只有一个输出，后续增加版本控制
	currentOutput := &domain.WorkflowOutput{
		Meta:       domain.Meta{},
		WorkflowId: GetWorkflowId(ctx),
		NodeId:     d.node.Id,
		Node:       d.node,
		Succeeded:  true,
	}
	if lastOutput != nil {
		currentOutput.Id = lastOutput.Id
	}
	if err := d.outputRepo.Save(ctx, currentOutput, nil, nil); err != nil {
		d.AddOutput(ctx, d.node.Name, "保存部署记录失败", err.Error())
		return err
	}
	d.AddOutput(ctx, d.node.Name, "保存部署记录成功")

	return nil
}

func (d *deployNode) checkCanSkip(ctx context.Context, lastOutput *domain.WorkflowOutput) (skip bool, reason string) {
	// TODO: 可控制是否强制部署
	if lastOutput != nil && lastOutput.Succeeded {
		// 比较和上次部署时的关键配置（即影响证书部署的）参数是否一致
		currentNodeConfig := d.node.GetConfigForDeploy()
		lastNodeConfig := lastOutput.Node.GetConfigForDeploy()
		if currentNodeConfig.ProviderAccessId != lastNodeConfig.ProviderAccessId {
			return false, "配置项变化：主机提供商授权"
		}
		if !maps.Equal(currentNodeConfig.ProviderConfig, lastNodeConfig.ProviderConfig) {
			return false, "配置项变化：主机提供商参数"
		}

		return true, "已部署过证书"
	}

	return false, "无历史部署记录"
}
