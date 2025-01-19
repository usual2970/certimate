package nodeprocessor

import (
	"context"

	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/notify"
	"github.com/usual2970/certimate/internal/repository"
)

type notifyNode struct {
	node         *domain.WorkflowNode
	settingsRepo settingsRepository
	*nodeLogger
}

func NewNotifyNode(node *domain.WorkflowNode) *notifyNode {
	return &notifyNode{
		node:         node,
		nodeLogger:   NewNodeLogger(node),
		settingsRepo: repository.NewSettingsRepository(),
	}
}

func (n *notifyNode) Run(ctx context.Context) error {
	n.AddOutput(ctx, n.node.Name, "开始执行")

	nodeConfig := n.node.GetConfigForNotify()

	// 获取通知配置
	settings, err := n.settingsRepo.GetByName(ctx, "notifyChannels")
	if err != nil {
		n.AddOutput(ctx, n.node.Name, "获取通知配置失败", err.Error())
		return err
	}

	// 获取通知渠道
	channelConfig, err := settings.GetNotifyChannelConfig(nodeConfig.Channel)
	if err != nil {
		n.AddOutput(ctx, n.node.Name, "获取通知渠道配置失败", err.Error())
		return err
	}

	// 发送通知
	if err := notify.SendToChannel(nodeConfig.Subject, nodeConfig.Message, nodeConfig.Channel, channelConfig); err != nil {
		n.AddOutput(ctx, n.node.Name, "发送通知失败", err.Error())
		return err
	}
	n.AddOutput(ctx, n.node.Name, "发送通知成功")

	return nil
}
