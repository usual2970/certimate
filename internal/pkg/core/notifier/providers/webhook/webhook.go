package webhook

import (
	"context"
	"crypto/tls"
	"log/slog"
	"net/http"

	webhook "github.com/nikoksr/notify/service/http"

	"github.com/usual2970/certimate/internal/pkg/core/notifier"
)

type NotifierConfig struct {
	// Webhook URL。
	Url string `json:"url"`
	// 是否允许不安全的连接。
	AllowInsecureConnections bool `json:"allowInsecureConnections,omitempty"`
}

type NotifierProvider struct {
	config *NotifierConfig
	logger *slog.Logger
}

var _ notifier.Notifier = (*NotifierProvider)(nil)

func NewNotifier(config *NotifierConfig) (*NotifierProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	return &NotifierProvider{
		config: config,
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
	srv := webhook.New()
	srv.AddReceiversURLs(n.config.Url)

	if n.config.AllowInsecureConnections {
		tlsConfig := &tls.Config{InsecureSkipVerify: true}
		transport := &http.Transport{TLSClientConfig: tlsConfig}
		client := &http.Client{Transport: transport}
		srv.WithClient(client)
	}

	err = srv.Send(ctx, subject, message)
	if err != nil {
		return nil, err
	}

	return &notifier.NotifyResult{}, nil
}
