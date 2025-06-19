package wecombot

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/go-resty/resty/v2"

	"github.com/certimate-go/certimate/pkg/core"
)

type NotifierProviderConfig struct {
	// 企业微信机器人 Webhook 地址。
	WebhookUrl string `json:"webhookUrl"`
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
	// REF: https://developer.work.weixin.qq.com/document/path/91770
	req := n.httpClient.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", "certimate").
		SetBody(map[string]any{
			"msgtype": "text",
			"text": map[string]string{
				"content": subject + "\n\n" + message,
			},
		})
	resp, err := req.Post(n.config.WebhookUrl)
	if err != nil {
		return nil, fmt.Errorf("wecom api error: failed to send request: %w", err)
	} else if resp.IsError() {
		return nil, fmt.Errorf("wecom api error: unexpected status code: %d, resp: %s", resp.StatusCode(), resp.String())
	}

	return &core.NotifyResult{}, nil
}
