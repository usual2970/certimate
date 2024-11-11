package serverchan

import (
	"context"
	"errors"
	"net/http"

	notifyHttp "github.com/nikoksr/notify/service/http"

	"github.com/usual2970/certimate/internal/pkg/core/notifier"
)

type ServerChanNotifierConfig struct {
	Url string `json:"url"`
}

type ServerChanNotifier struct {
	config *ServerChanNotifierConfig
}

var _ notifier.Notifier = (*ServerChanNotifier)(nil)

func New(config *ServerChanNotifierConfig) (*ServerChanNotifier, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}

	return &ServerChanNotifier{
		config: config,
	}, nil
}

func (n *ServerChanNotifier) Notify(ctx context.Context, subject string, message string) (res *notifier.NotifyResult, err error) {
	srv := notifyHttp.New()

	srv.AddReceivers(&notifyHttp.Webhook{
		URL:         n.config.Url,
		Header:      http.Header{},
		ContentType: "application/json",
		Method:      http.MethodPost,
		BuildPayload: func(subject, message string) (payload any) {
			return map[string]string{
				"text": subject,
				"desp": message,
			}
		},
	})

	err = srv.Send(ctx, subject, message)
	if err != nil {
		return nil, err
	}

	return &notifier.NotifyResult{}, nil
}
