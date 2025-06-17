package larkbot

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/go-lark/lark"

	"github.com/certimate-go/certimate/pkg/core"
)

type NotifierProviderConfig struct {
	// 飞书机器人 Webhook 地址。
	WebhookUrl string `json:"webhookUrl"`
}

type NotifierProvider struct {
	config *NotifierProviderConfig
	logger *slog.Logger
}

var _ core.Notifier = (*NotifierProvider)(nil)

func NewNotifierProvider(config *NotifierProviderConfig) (*NotifierProvider, error) {
	if config == nil {
		return nil, errors.New("the configuration of the notifier provider is nil")
	}

	return &NotifierProvider{
		config: config,
		logger: slog.Default(),
	}, nil
}

func (n *NotifierProvider) SetLogger(logger *slog.Logger) {
	if logger == nil {
		n.logger = slog.New(slog.DiscardHandler)
	} else {
		n.logger = logger
	}
}

func (n *NotifierProvider) Notify(ctx context.Context, subject string, message string) (*core.NotifyResult, error) {
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

	return &core.NotifyResult{}, nil
}
