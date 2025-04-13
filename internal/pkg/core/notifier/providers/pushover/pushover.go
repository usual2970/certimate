package pushover

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
	Token string `json:"token"` // 应用 API Token
	User  string `json:"user"`  // 用户/分组 Key
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

// Notify 发送通知
// 参考文档：https://pushover.net/api
func (n *NotifierProvider) Notify(ctx context.Context, subject string, message string) (res *notifier.NotifyResult, err error) {
	// 请求体
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

	// Make request
	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, errors.Wrap(err, "encode message body")
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://api.pushover.net/1/messages.json",
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, errors.Wrap(err, "create new request")
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	// Send request to pushover service
	resp, err := n.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "send request to pushover server")
	}
	defer resp.Body.Close()

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read response")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("pushover returned status code %d: %s", resp.StatusCode, string(result))
	}

	return &notifier.NotifyResult{}, nil
}
