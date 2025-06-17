package discordbot

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/go-resty/resty/v2"

	"github.com/certimate-go/certimate/pkg/core"
)

type NotifierProviderConfig struct {
	// Discord Bot API Token。
	BotToken string `json:"botToken"`
	// Discord Channel ID。
	ChannelId string `json:"channelId"`
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
	// REF: https://discord.com/developers/docs/resources/message#create-message
	req := n.httpClient.R().
		SetContext(ctx).
		SetHeader("Authorization", "Bot "+n.config.BotToken).
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", "certimate").
		SetBody(map[string]any{
			"content": subject + "\n" + message,
		})
	resp, err := req.Post(fmt.Sprintf("https://discord.com/api/v9/channels/%s/messages", n.config.ChannelId))
	if err != nil {
		return nil, fmt.Errorf("discord api error: failed to send request: %w", err)
	} else if resp.IsError() {
		return nil, fmt.Errorf("discord api error: unexpected status code: %d, resp: %s", resp.StatusCode(), resp.String())
	}

	return &core.NotifyResult{}, nil
}
