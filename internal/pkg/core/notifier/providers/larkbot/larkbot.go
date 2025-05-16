package larkbot

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/go-lark/lark"

	"github.com/usual2970/certimate/internal/pkg/core/notifier"
)

type NotifierConfig struct {
	// 飞书机器人 Webhook 地址。
	WebhookUrl string `json:"webhookUrl"`
}

type NotifierProvider struct {
	config *NotifierConfig
	logger *slog.Logger
}

var _ notifier.Notifier = (*NotifierProvider)(nil)

func NewNotifier(config *NotifierConfig) (*NotifierProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	return &NotifierProvider{
		config: config,
		logger: slog.Default(),
	}, nil
}

func (n *NotifierProvider) WithLogger(logger *slog.Logger) notifier.Notifier {
	if logger == nil {
		n.logger = slog.Default()
	} else {
		n.logger = logger
	}
	return n
}

func (n *NotifierProvider) Notify(ctx context.Context, subject string, message string) (res *notifier.NotifyResult, err error) {
	bot := lark.NewNotificationBot(n.config.WebhookUrl)
	content := lark.NewPostBuilder().
		Title(subject).
		TextTag(message, 1, false).
		Render()
	msg := lark.NewMsgBuffer(lark.MsgPost).Post(content)
	resp, err := bot.PostNotificationV2(msg.Build())
	if err != nil {
		return nil, fmt.Errorf("lark api error: %w", err)
	} else if resp.Code != 0 {
		return nil, fmt.Errorf("lark api error: code='%d', message='%s'", resp.Code, resp.Msg)
	}

	return &notifier.NotifyResult{}, nil
}
