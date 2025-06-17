package telegrambot

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/go-resty/resty/v2"

	"github.com/certimate-go/certimate/pkg/core"
)

type NotifierProviderConfig struct {
	// Telegram Bot API Token。
	BotToken string `json:"botToken"`
	// Telegram Chat ID。
	ChatId int64 `json:"chatId"`
}

type NotifierProvider struct {
	config     *NotifierProviderConfig
	logger     *slog.Logger
	httpClient *resty.Client
}

var _ core.Notifier = (*NotifierProvider)(nil)

func NewNotifierProvider(config *NotifierProviderConfig) (*NotifierProvider, error) {
	if config == nil {
		return nil, errors.New("the configuration of the notifier provider is nil")
	}

	client := resty.New()

	return &NotifierProvider{
		config:     config,
		logger:     slog.Default(),
		httpClient: client,
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
	// REF: https://core.telegram.org/bots/api#sendmessage
	req := n.httpClient.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", "certimate").
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

	return &core.NotifyResult{}, nil
}
