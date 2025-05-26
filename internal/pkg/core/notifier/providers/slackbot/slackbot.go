package discordbot

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/go-resty/resty/v2"

	"github.com/usual2970/certimate/internal/pkg/core/notifier"
)

type NotifierConfig struct {
	// Slack Bot API Token。
	BotToken string `json:"botToken"`
	// Slack Channel ID。
	ChannelId string `json:"channelId"`
}

type NotifierProvider struct {
	config     *NotifierConfig
	logger     *slog.Logger
	httpClient *resty.Client
}

var _ notifier.Notifier = (*NotifierProvider)(nil)

func NewNotifier(config *NotifierConfig) (*NotifierProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	client := resty.New()

	return &NotifierProvider{
		config:     config,
		logger:     slog.Default(),
		httpClient: client,
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
	// REF: https://docs.slack.dev/messaging/sending-and-scheduling-messages#publishing
	req := n.httpClient.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+n.config.BotToken).
		SetBody(map[string]any{
			"token":   n.config.BotToken,
			"channel": n.config.ChannelId,
			"text":    subject + "\n" + message,
		})
	resp, err := req.Post("https://slack.com/api/chat.postMessage")
	if err != nil {
		return nil, fmt.Errorf("slack api error: failed to send request: %w", err)
	} else if resp.IsError() {
		return nil, fmt.Errorf("slack api error: unexpected status code: %d, resp: %s", resp.StatusCode(), resp.String())
	}

	return &notifier.NotifyResult{}, nil
}
