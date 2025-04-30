package bark

import (
	"context"
	"log/slog"

	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/bark"

	"github.com/usual2970/certimate/internal/pkg/core/notifier"
)

type NotifierConfig struct {
	// Bark 服务地址。
	// 零值时默认使用官方服务器。
	ServerUrl string `json:"serverUrl"`
	// Bark 设备密钥。
	DeviceKey string `json:"deviceKey"`
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
		logger: slog.Default(),
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
	var srv notify.Notifier
	if n.config.ServerUrl == "" {
		srv = bark.New(n.config.DeviceKey)
	} else {
		srv = bark.NewWithServers(n.config.DeviceKey, n.config.ServerUrl)
	}

	err = srv.Send(ctx, subject, message)
	if err != nil {
		return nil, err
	}

	return &notifier.NotifyResult{}, nil
}
