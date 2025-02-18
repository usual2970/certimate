package lark

import (
	"context"

	"github.com/nikoksr/notify/service/lark"

	"github.com/usual2970/certimate/internal/pkg/core/notifier"
)

type LarkNotifierConfig struct {
	// 飞书机器人 Webhook 地址。
	WebhookUrl string `json:"webhookUrl"`
}

type LarkNotifier struct {
	config *LarkNotifierConfig
}

var _ notifier.Notifier = (*LarkNotifier)(nil)

func New(config *LarkNotifierConfig) (*LarkNotifier, error) {
	if config == nil {
		panic("config is nil")
	}

	return &LarkNotifier{
		config: config,
	}, nil
}

func (n *LarkNotifier) Notify(ctx context.Context, subject string, message string) (res *notifier.NotifyResult, err error) {
	srv := lark.NewWebhookService(n.config.WebhookUrl)

	err = srv.Send(ctx, subject, message)
	if err != nil {
		return nil, err
	}

	return &notifier.NotifyResult{}, nil
}
