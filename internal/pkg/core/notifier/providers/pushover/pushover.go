package pushover

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
	// Pushover API Token。
	Token string `json:"token"`
	// 用户或分组标识。
	User string `json:"user"`
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
	// REF: https://pushover.net/api
	reqBody := &struct {
		Token   string `json:"token"`
		User    string `json:"user"`
		Title   string `json:"title"`
		Message string `json:"message"`
	}{
		Token:   n.config.Token,
		User:    n.config.User,
		Title:   subject,
		Message: message,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("pushover api error: failed to encode message body: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://api.pushover.net/1/messages.json",
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, fmt.Errorf("pushover api error: failed to create new request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := n.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("pushover api error: failed to send request: %w", err)
	}
	defer resp.Body.Close()

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("pushover api error: failed to read response: %w", err)
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("pushover api error: unexpected status code: %d, resp: %s", resp.StatusCode, string(result))
	}

	return &notifier.NotifyResult{}, nil
}
