package nodeprocessor

import (
	"context"
	"fmt"
	"strings"

	"github.com/usual2970/certimate/internal/deployer"
	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/repository"
)

type deployNode struct {
	node       *domain.WorkflowNode
	outputRepo WorkflowOutputRepository
	*Logger
}

func NewDeployNode(node *domain.WorkflowNode) *deployNode {
	return &deployNode{
		node:       node,
		Logger:     NewLogger(node),
		outputRepo: repository.NewWorkflowOutputRepository(),
	}
}

func (d *deployNode) Run(ctx context.Context) error {
	d.AddOutput(ctx, d.node.Name, "开始执行")
	// 检查是否部署过（部署过则直接返回，和 v0.2 暂时保持一致）
	output, err := d.outputRepo.Get(ctx, d.node.Id)
	if err != nil && !domain.IsRecordNotFound(err) {
		d.AddOutput(ctx, d.node.Name, "查询部署记录失败", err.Error())
		return err
	}
	if output != nil && output.Succeed {
		d.AddOutput(ctx, d.node.Name, "已部署过")
		return nil
	}
	// 获取部署对象
	// 获取证书
	certSource := d.node.GetConfigString("certificate")

	certSourceSlice := strings.Split(certSource, "#")
	if len(certSourceSlice) != 2 {
		d.AddOutput(ctx, d.node.Name, "证书来源配置错误", certSource)
		return fmt.Errorf("证书来源配置错误: %s", certSource)
	}

	cert, err := d.outputRepo.GetCertificate(ctx, certSourceSlice[0])
	if err != nil {
		d.AddOutput(ctx, d.node.Name, "获取证书失败", err.Error())
		return err
	}

	accessRepo := repository.NewAccessRepository()
	access, err := accessRepo.GetById(context.Background(), d.node.GetConfigString("access"))
	if err != nil {
		d.AddOutput(ctx, d.node.Name, "获取授权配置失败", err.Error())
		return err
	}
	option := &deployer.DeployerOption{
		DomainId:     d.node.Id,
		Domain:       cert.SAN,
		Access:       access.Config,
		AccessRecord: access,
		DeployConfig: domain.DeployConfig{
			Id:     d.node.Id,
			Access: access.Id,
			Type:   d.node.GetConfigString("providerType"),
			Config: d.node.Config,
		},
	}

	deploy, err := deployer.GetWithTypeAndOption(d.node.GetConfigString("providerType"), option)
	if err != nil {
		d.AddOutput(ctx, d.node.Name, "获取部署对象失败", err.Error())
		return err
	}

	// 部署
	if err := deploy.Deploy(ctx); err != nil {
		d.AddOutput(ctx, d.node.Name, "部署失败", err.Error())
		return err
	}

	d.AddOutput(ctx, d.node.Name, "部署成功")

	// 记录部署结果
	output = &domain.WorkflowOutput{
		Workflow: GetWorkflowId(ctx),
		NodeId:   d.node.Id,
		Node:     d.node,
		Succeed:  true,
	}

	if err := d.outputRepo.Save(ctx, output, nil, nil); err != nil {
		d.AddOutput(ctx, d.node.Name, "保存部署记录失败", err.Error())
		return err
	}

	d.AddOutput(ctx, d.node.Name, "保存部署记录成功")

	return nil
}
