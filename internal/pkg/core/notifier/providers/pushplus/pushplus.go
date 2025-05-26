package pushplus

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/go-resty/resty/v2"

	"github.com/usual2970/certimate/internal/pkg/core/notifier"
)

type NotifierConfig struct {
	// PushPlus Token。
	Token string `json:"token"`
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
	// REF: https://pushplus.plus/doc/guide/api.html#%E4%B8%80%E3%80%81%E5%8F%91%E9%80%81%E6%B6%88%E6%81%AF%E6%8E%A5%E5%8F%A3
	req := n.httpClient.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]any{
			"title":   subject,
			"content": message,
			"token":   n.config.Token,
		})
	resp, err := req.Post("https://www.pushplus.plus/send")
	if err != nil {
		return nil, fmt.Errorf("pushplus api error: failed to send request: %w", err)
	} else if resp.IsError() {
		return nil, fmt.Errorf("pushplus api error: unexpected status code: %d, resp: %s", resp.StatusCode(), resp.String())
	}

	var errorResponse struct {
		Code    int    `json:"code"`
		Message string `json:"msg"`
	}
	if err := json.Unmarshal(resp.Body(), &errorResponse); err != nil {
		return nil, fmt.Errorf("pushplus api error: failed to unmarshal response: %w", err)
	} else if errorResponse.Code != 200 {
		return nil, fmt.Errorf("pushplus api error: code='%d', message='%s'", errorResponse.Code, errorResponse.Message)
	}

	return &notifier.NotifyResult{}, nil
}
