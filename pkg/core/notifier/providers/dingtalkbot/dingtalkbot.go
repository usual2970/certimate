package dingtalkbot

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/url"

	"github.com/blinkbean/dingtalk"

	"github.com/certimate-go/certimate/pkg/core"
)

type NotifierProviderConfig struct {
	// 钉钉机器人的 Webhook 地址。
	WebhookUrl string `json:"webhookUrl"`
	// 钉钉机器人的 Secret。
	Secret string `json:"secret"`
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
	webhookUrl, err := url.Parse(n.config.WebhookUrl)
	if err != nil {
		return nil, fmt.Errorf("dingtalk api error: invalid webhook url: %w", err)
	}

	var bot *dingtalk.DingTalk
	if n.config.Secret == "" {
		bot = dingtalk.InitDingTalk([]string{webhookUrl.Query().Get("access_token")}, "")
	} else {
		bot = dingtalk.InitDingTalkWithSecret(webhookUrl.Query().Get("access_token"), n.config.Secret)
	}

	if err := bot.SendTextMessage(subject + "\n" + message); err != nil {
		return nil, fmt.Errorf("dingtalk api error: %w", err)
	}

	return &core.NotifyResult{}, nil
}
