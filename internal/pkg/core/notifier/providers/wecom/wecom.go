package serverchan

import (
	"context"
	"errors"
	"net/http"

	notifyHttp "github.com/nikoksr/notify/service/http"

	"github.com/usual2970/certimate/internal/pkg/core/notifier"
)

type WeComNotifierConfig struct {
	// 企业微信机器人 Webhook 地址。
	WebhookUrl string `json:"webhookUrl"`
}

type WeComNotifier struct {
	config *WeComNotifierConfig
}

var _ notifier.Notifier = (*WeComNotifier)(nil)

func New(config *WeComNotifierConfig) (*WeComNotifier, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}

	return &WeComNotifier{
		config: config,
	}, nil
}

func (n *WeComNotifier) Notify(ctx context.Context, subject string, message string) (res *notifier.NotifyResult, err error) {
	srv := notifyHttp.New()

	srv.AddReceivers(&notifyHttp.Webhook{
		URL:         n.config.WebhookUrl,
		Header:      http.Header{},
		ContentType: "application/json",
		Method:      http.MethodPost,
		BuildPayload: func(subject, message string) (payload any) {
			return map[string]any{
				"msgtype": "text",
				"text": map[string]string{
					"content": subject + "\n\n" + message,
				},
			}
		},
	})

	err = srv.Send(ctx, subject, message)
	if err != nil {
		return nil, err
	}

	return &notifier.NotifyResult{}, nil
}
