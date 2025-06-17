package bark

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/go-resty/resty/v2"

	"github.com/certimate-go/certimate/pkg/core"
)

type NotifierProviderConfig struct {
	// Bark 服务地址。
	// 零值时使用官方服务器。
	ServerUrl string `json:"serverUrl"`
	// Bark 设备密钥。
	DeviceKey string `json:"deviceKey"`
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
	const defaultServerURL = "https://api.day.app/"
	serverUrl := defaultServerURL
	if n.config.ServerUrl != "" {
		serverUrl = n.config.ServerUrl
	}

	// REF: https://bark.day.app/#/tutorial
	req := n.httpClient.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]any{
			"title":      subject,
			"body":       message,
			"device_key": n.config.DeviceKey,
		})
	resp, err := req.Post(serverUrl)
	if err != nil {
		return nil, fmt.Errorf("bark api error: failed to send request: %w", err)
	} else if resp.IsError() {
		return nil, fmt.Errorf("bark api error: unexpected status code: %d, resp: %s", resp.StatusCode(), resp.String())
	}

	return &core.NotifyResult{}, nil
}
