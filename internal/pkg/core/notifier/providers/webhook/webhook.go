package webhook

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/usual2970/certimate/internal/pkg/core/notifier"
)

type NotifierConfig struct {
	// Webhook URL。
	WebhookUrl string `json:"webhookUrl"`
	// Webhook 回调数据（application/json 或 application/x-www-form-urlencoded 格式）。
	WebhookData string `json:"webhookData,omitempty"`
	// 请求谓词。
	// 零值时默认为 "POST"。
	Method string `json:"method,omitempty"`
	// 请求标头。
	Headers map[string]string `json:"headers,omitempty"`
	// 是否允许不安全的连接。
	AllowInsecureConnections bool `json:"allowInsecureConnections,omitempty"`
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

	client := resty.New().
		SetTimeout(30 * time.Second).
		SetRetryCount(3).
		SetRetryWaitTime(5 * time.Second)
	if config.AllowInsecureConnections {
		client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}

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
	var webhookData interface{}
	if n.config.WebhookData == "" {
		webhookData = map[string]any{
			"subject": subject,
			"message": message,
		}
	} else {
		err = json.Unmarshal([]byte(n.config.WebhookData), &webhookData)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshall webhook data: %w", err)
		}

		replaceJsonValueRecursively(webhookData, "${SUBJECT}", subject)
		replaceJsonValueRecursively(webhookData, "${MESSAGE}", message)
	}

	req := n.httpClient.R().
		SetContext(ctx).
		SetHeaders(n.config.Headers)
	req.URL = n.config.WebhookUrl
	req.Method = n.config.Method
	if req.Method == "" {
		req.Method = http.MethodPost
	}

	resp, err := req.
		SetHeader("Content-Type", "application/json").
		SetBody(webhookData).
		Send()
	if err != nil {
		return nil, fmt.Errorf("failed to send webhook request: %w", err)
	} else if resp.IsError() {
		return nil, fmt.Errorf("unexpected webhook response status code: %d", resp.StatusCode())
	}

	n.logger.Debug("webhook responded", slog.Any("response", resp.String()))

	return &notifier.NotifyResult{}, nil
}

func replaceJsonValueRecursively(data interface{}, oldStr, newStr string) interface{} {
	switch v := data.(type) {
	case map[string]any:
		for k, val := range v {
			v[k] = replaceJsonValueRecursively(val, oldStr, newStr)
		}
	case []any:
		for i, val := range v {
			v[i] = replaceJsonValueRecursively(val, oldStr, newStr)
		}
	case string:
		return strings.ReplaceAll(v, oldStr, newStr)
	}
	return data
}
