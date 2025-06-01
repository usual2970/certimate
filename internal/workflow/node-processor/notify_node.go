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
	*nodeOutputer

	settingsRepo settingsRepository
}

func NewNotifyNode(node *domain.WorkflowNode) *notifyNode {
	return &notifyNode{
		node:          node,
		nodeProcessor: newNodeProcessor(node),
		nodeOutputer:  newNodeOutputer(),

		settingsRepo: repository.NewSettingsRepository(),
	}
}

func (n *notifyNode) Process(ctx context.Context) error {
	n.logger.Info("ready to send notification ...")

	nodeCfg := n.node.GetConfigForNotify()

	if nodeCfg.Provider == "" {
		// Deprecated: v0.4.x 将废弃
		// 兼容旧版本的通知渠道
		n.logger.Warn("WARNING! you are using the notification channel from global settings, which will be deprecated in the future")

		// 获取通知配置
		settings, err := n.settingsRepo.GetByName(ctx, "notifyChannels")
		if err != nil {
			return err
		}

		// 获取通知渠道
		channelConfig, err := settings.GetNotifyChannelConfig(nodeCfg.Channel)
		if err != nil {
			return err
		}

		// 发送通知
		if err := notify.SendToChannel(nodeCfg.Subject, nodeCfg.Message, nodeCfg.Channel, channelConfig); err != nil {
			n.logger.Warn("failed to send notification", slog.String("channel", nodeCfg.Channel))
			return err
		}

		n.logger.Info("notification completed")
		return nil
	}

	// 初始化通知器
	deployer, err := notify.NewWithWorkflowNode(notify.NotifierWithWorkflowNodeConfig{
		Node:    n.node,
		Logger:  n.logger,
		Subject: nodeCfg.Subject,
		Message: nodeCfg.Message,
	})
	if err != nil {
		n.logger.Warn("failed to create notifier provider")
		return err
	}

	// 推送通知
	if err := deployer.Notify(ctx); err != nil {
		n.logger.Warn("failed to send notification")
		return err
	}

	n.logger.Info("notification completed")
	return nil
}
