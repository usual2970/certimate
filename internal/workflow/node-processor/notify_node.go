package nodeprocessor

import (
	"context"
	"log/slog"

	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/notify"
	"github.com/usual2970/certimate/internal/repository"
)

type notifyNode struct {
	node *domain.WorkflowNode
	*nodeProcessor

	settingsRepo settingsRepository
}

func NewNotifyNode(node *domain.WorkflowNode) *notifyNode {
	return &notifyNode{
		node:          node,
		nodeProcessor: newNodeProcessor(node),

		settingsRepo: repository.NewSettingsRepository(),
	}
}

func (n *notifyNode) Process(ctx context.Context) error {
	n.logger.Info("ready to notify ...")

	nodeConfig := n.node.GetConfigForNotify()

	// 获取通知配置
	settings, err := n.settingsRepo.GetByName(ctx, "notifyChannels")
	if err != nil {
		return err
	}

	// 获取通知渠道
	channelConfig, err := settings.GetNotifyChannelConfig(nodeConfig.Channel)
	if err != nil {
		return err
	}

	// 发送通知
	if err := notify.SendToChannel(nodeConfig.Subject, nodeConfig.Message, nodeConfig.Channel, channelConfig); err != nil {
		n.logger.Warn("failed to notify", slog.String("channel", nodeConfig.Channel))
		return err
	}

	n.logger.Info("notify completed")

	return nil
}
