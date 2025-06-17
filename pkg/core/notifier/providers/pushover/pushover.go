package pushover

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/go-resty/resty/v2"

	"github.com/certimate-go/certimate/pkg/core"
)

type NotifierProviderConfig struct {
	// Pushover API Token。
	Token string `json:"token"`
	// 用户或分组标识。
	User string `json:"user"`
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
	// REF: https://pushover.net/api
	req := n.httpClient.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", "certimate").
		SetBody(map[string]any{
			"title":   subject,
			"message": message,
			"token":   n.config.Token,
			"user":    n.config.User,
		})
	resp, err := req.Post("https://api.pushover.net/1/messages.json")
	if err != nil {
		return nil, fmt.Errorf("pushover api error: failed to send request: %w", err)
	} else if resp.IsError() {
		return nil, fmt.Errorf("pushover api error: unexpected status code: %d, resp: %s", resp.StatusCode(), resp.String())
	}

	return &core.NotifyResult{}, nil
}
