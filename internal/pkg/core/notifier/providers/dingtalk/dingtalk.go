package dingtalk

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"

	"github.com/nikoksr/notify/service/dingding"

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
		n.logger = slog.Default()
	} else {
		n.logger = logger
	}
	return n
}

func (n *NotifierProvider) Notify(ctx context.Context, subject string, message string) (res *notifier.NotifyResult, err error) {
	webhookUrl, err := url.Parse(n.config.WebhookUrl)
	if err != nil {
		return nil, fmt.Errorf("invalid webhook url: %w", err)
	}

	srv := dingding.New(&dingding.Config{
		Token:  webhookUrl.Query().Get("access_token"),
		Secret: n.config.Secret,
	})

	err = srv.Send(ctx, subject, message)
	if err != nil {
		return nil, err
	}

	return &notifier.NotifyResult{}, nil
}
