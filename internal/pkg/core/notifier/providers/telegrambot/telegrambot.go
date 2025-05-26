package telegrambot

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/go-resty/resty/v2"

	"github.com/usual2970/certimate/internal/pkg/core/notifier"
)

type NotifierConfig struct {
	// Telegram Bot API Token。
	BotToken string `json:"botToken"`
	// Telegram Chat ID。
	ChatId int64 `json:"chatId"`
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
	// REF: https://core.telegram.org/bots/api#sendmessage
	req := n.httpClient.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]any{
			"chat_id": n.config.ChatId,
			"text":    subject + "\n" + message,
		})
	resp, err := req.Post(fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", n.config.BotToken))
	if err != nil {
		return nil, fmt.Errorf("telegram api error: failed to send request: %w", err)
	} else if resp.IsError() {
		return nil, fmt.Errorf("telegram api error: unexpected status code: %d, resp: %s", resp.StatusCode(), resp.String())
	}

	return &notifier.NotifyResult{}, nil
}
