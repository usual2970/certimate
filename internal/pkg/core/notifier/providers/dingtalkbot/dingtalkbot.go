package dingtalkbot

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"

	"github.com/blinkbean/dingtalk"

	"github.com/usual2970/certimate/internal/pkg/core/notifier"
)

type NotifierConfig struct {
	// 钉钉机器人的 Webhook 地址。
	WebhookUrl string `json:"webhookUrl"`
	// 钉钉机器人的 Secret。
	Secret string `json:"secret"`
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
		n.logger = slog.New(slog.DiscardHandler)
	} else {
		n.logger = logger
	}
	return n
}

func (n *NotifierProvider) Notify(ctx context.Context, subject string, message string) (*notifier.NotifyResult, error) {
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

	return &notifier.NotifyResult{}, nil
}
