package gotify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/pkg/errors"

	"github.com/usual2970/certimate/internal/pkg/core/notifier"
)

type NotifierConfig struct {
	// Gotify 服务地址
	// 示例：https://gotify.example.com
	Url string `json:"url"`
	// Gotify Token
	Token string `json:"token"`
	// Gotify 消息优先级
	Priority int64 `json:"priority"`
}

type NotifierProvider struct {
	config *NotifierConfig
	logger *slog.Logger
	// 未来将移除
	httpClient *http.Client
}

var _ notifier.Notifier = (*NotifierProvider)(nil)

func NewNotifier(config *NotifierConfig) (*NotifierProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	return &NotifierProvider{
		config:     config,
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
	// Gotify 原生实现, notify 库没有实现, 等待合并
	reqBody := &struct {
		Title    string `json:"title"`
		Message  string `json:"message"`
		Priority int64  `json:"priority"`
	}{
		Title:    subject,
		Message:  message,
		Priority: n.config.Priority,
	}

	// Make request
	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, errors.Wrap(err, "encode message body")
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/message", n.config.Url),
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, errors.Wrap(err, "create new request")
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", n.config.Token))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	// Send request to gotify service
	resp, err := n.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "send request to gotify server")
	}
	defer resp.Body.Close()

	// Read response and verify success
	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read response")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gotify returned status code %d: %s", resp.StatusCode, string(result))
	}
	return &notifier.NotifyResult{}, nil
}
