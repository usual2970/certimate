package gotify

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/go-resty/resty/v2"

	"github.com/usual2970/certimate/internal/pkg/core/notifier"
)

type NotifierConfig struct {
	// Gotify 服务地址。
	ServerUrl string `json:"serverUrl"`
	// Gotify Token。
	Token string `json:"token"`
	// Gotify 消息优先级。
	Priority int64 `json:"priority,omitempty"`
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
		n.logger = slog.New(slog.DiscardHandler)
	} else {
		n.logger = logger
	}
	return n
}

func (n *NotifierProvider) Notify(ctx context.Context, subject string, message string) (*notifier.NotifyResult, error) {
	serverUrl := strings.TrimRight(n.config.ServerUrl, "/")

	// REF: https://gotify.net/api-docs#/message/createMessage
	req := n.httpClient.R().
		SetContext(ctx).
		SetHeader("Authorization", "Bearer "+n.config.Token).
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", "certimate").
		SetBody(map[string]any{
			"title":    subject,
			"message":  message,
			"priority": n.config.Priority,
		})
	resp, err := req.Post(fmt.Sprintf("%s/message", serverUrl))
	if err != nil {
		return nil, fmt.Errorf("gotify api error: failed to send request: %w", err)
	} else if resp.IsError() {
		return nil, fmt.Errorf("gotify api error: unexpected status code: %d, resp: %s", resp.StatusCode(), resp.String())
	}

	return &notifier.NotifyResult{}, nil
}
