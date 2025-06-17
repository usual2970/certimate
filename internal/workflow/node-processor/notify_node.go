package nodeprocessor

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/internal/notify"
	"github.com/certimate-go/certimate/internal/repository"
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
	nodeCfg := n.node.GetConfigForNotify()
	n.logger.Info("ready to send notification ...", slog.Any("config", nodeCfg))

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

	// 检测是否可以跳过本次执行
	if skippable := n.checkCanSkip(ctx); skippable {
		n.logger.Info(fmt.Sprintf("skip this notification, because all the previous nodes have been skipped"))
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

func (n *notifyNode) checkCanSkip(ctx context.Context) (_skip bool) {
	thisNodeCfg := n.node.GetConfigForNotify()
	if !thisNodeCfg.SkipOnAllPrevSkipped {
		return false
	}

	prevNodeOutputs := GetAllNodeOutputs(ctx)
	for _, nodeOutput := range prevNodeOutputs {
		if nodeOutput[outputKeyForNodeSkipped] != nil {
			if nodeOutput[outputKeyForNodeSkipped].(string) != strconv.FormatBool(true) {
				return false
			}
		}
	}

	return true
}
