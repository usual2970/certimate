package dingtalk

import (
	"context"
	"errors"

	"github.com/nikoksr/notify/service/dingding"

	"github.com/usual2970/certimate/internal/pkg/core/notifier"
)

type DingTalkNotifierConfig struct {
	// 钉钉机器人的 Token。
	AccessToken string `json:"accessToken"`
	// 钉钉机器人的 Secret。
	Secret string `json:"secret"`
}

type DingTalkNotifier struct {
	config *DingTalkNotifierConfig
}

var _ notifier.Notifier = (*DingTalkNotifier)(nil)

func New(config *DingTalkNotifierConfig) (*DingTalkNotifier, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}

	return &DingTalkNotifier{
		config: config,
	}, nil
}

func (n *DingTalkNotifier) Notify(ctx context.Context, subject string, message string) (res *notifier.NotifyResult, err error) {
	srv := dingding.New(&dingding.Config{
		Token:  n.config.AccessToken,
		Secret: n.config.Secret,
	})

	err = srv.Send(ctx, subject, message)
	if err != nil {
		return nil, err
	}

	return &notifier.NotifyResult{}, nil
}
