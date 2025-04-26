package gotify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/usual2970/certimate/internal/pkg/core/notifier"
)

type NotifierConfig struct {
	// Gotify 服务地址。
	Url string `json:"url"`
	// Gotify Token。
	Token string `json:"token"`
	// Gotify 消息优先级。
	Priority int64 `json:"priority,omitempty"`
}

type NotifierProvider struct {
	config     *NotifierConfig
	logger     *slog.Logger
	httpClient *http.Client
}

var _ notifier.Notifier = (*NotifierProvider)(nil)

func NewNotifier(config *NotifierConfig) (*NotifierProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	return &NotifierProvider{
		config:     config,
		logger:     slog.Default(),
		httpClient: http.DefaultClient,
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
	reqBody := &struct {
		Title    string `json:"title"`
		Message  string `json:"message"`
		Priority int64  `json:"priority"`
	}{
		Title:    subject,
		Message:  message,
		Priority: n.config.Priority,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("gotify api error: failed to encode message body: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/message", n.config.Url),
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, fmt.Errorf("gotify api error: failed to create new request: %w", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", n.config.Token))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := n.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("gotify api error: failed to send request: %w", err)
	}
	defer resp.Body.Close()

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("gotify api error: failed to read response: %w", err)
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gotify api error: unexpected status code: %d, resp: %s", resp.StatusCode, string(result))
	}

	return &notifier.NotifyResult{}, nil
}
