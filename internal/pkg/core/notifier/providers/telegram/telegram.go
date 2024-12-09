package telegram

import (
	"context"
	"errors"

	"github.com/nikoksr/notify/service/telegram"

	"github.com/usual2970/certimate/internal/pkg/core/notifier"
)

type TelegramNotifierConfig struct {
	// Telegram API Token。
	ApiToken string `json:"apiToken"`
	// Telegram Chat ID。
	ChatId int64 `json:"chatId"`
}

type TelegramNotifier struct {
	config *TelegramNotifierConfig
}

var _ notifier.Notifier = (*TelegramNotifier)(nil)

func New(config *TelegramNotifierConfig) (*TelegramNotifier, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}

	return &TelegramNotifier{
		config: config,
	}, nil
}

func (n *TelegramNotifier) Notify(ctx context.Context, subject string, message string) (res *notifier.NotifyResult, err error) {
	srv, err := telegram.New(n.config.ApiToken)
	if err != nil {
		return nil, err
	}

	srv.AddReceivers(n.config.ChatId)

	err = srv.Send(ctx, subject, message)
	if err != nil {
		return nil, err
	}

	return &notifier.NotifyResult{}, nil
}
