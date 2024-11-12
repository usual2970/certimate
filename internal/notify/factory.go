package notify

import (
	"errors"

	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/pkg/core/notifier"
	notifierBark "github.com/usual2970/certimate/internal/pkg/core/notifier/providers/bark"
	notifierDingTalk "github.com/usual2970/certimate/internal/pkg/core/notifier/providers/dingtalk"
	notifierEmail "github.com/usual2970/certimate/internal/pkg/core/notifier/providers/email"
	notifierLark "github.com/usual2970/certimate/internal/pkg/core/notifier/providers/lark"
	notifierServerChan "github.com/usual2970/certimate/internal/pkg/core/notifier/providers/serverchan"
	notifierTelegram "github.com/usual2970/certimate/internal/pkg/core/notifier/providers/telegram"
	notifierWebhook "github.com/usual2970/certimate/internal/pkg/core/notifier/providers/webhook"
	"github.com/usual2970/certimate/internal/pkg/utils/maps"
)

func createNotifier(channel string, channelConfig map[string]any) (notifier.Notifier, error) {
	switch channel {
	case domain.NotifyChannelEmail:
		return notifierEmail.New(&notifierEmail.EmailNotifierConfig{
			SmtpHost:        maps.GetValueAsString(channelConfig, "smtpHost"),
			SmtpPort:        maps.GetValueAsInt32(channelConfig, "smtpPort"),
			SmtpTLS:         maps.GetValueOrDefaultAsBool(channelConfig, "smtpTLS", true),
			Username:        maps.GetValueOrDefaultAsString(channelConfig, "username", maps.GetValueAsString(channelConfig, "senderAddress")),
			Password:        maps.GetValueAsString(channelConfig, "password"),
			SenderAddress:   maps.GetValueAsString(channelConfig, "senderAddress"),
			ReceiverAddress: maps.GetValueAsString(channelConfig, "receiverAddress"),
		})

	case domain.NotifyChannelWebhook:
		return notifierWebhook.New(&notifierWebhook.WebhookNotifierConfig{
			Url: maps.GetValueAsString(channelConfig, "url"),
		})

	case domain.NotifyChannelDingtalk:
		return notifierDingTalk.New(&notifierDingTalk.DingTalkNotifierConfig{
			AccessToken: maps.GetValueAsString(channelConfig, "accessToken"),
			Secret:      maps.GetValueAsString(channelConfig, "secret"),
		})

	case domain.NotifyChannelLark:
		return notifierLark.New(&notifierLark.LarkNotifierConfig{
			WebhookUrl: maps.GetValueAsString(channelConfig, "webhookUrl"),
		})

	case domain.NotifyChannelTelegram:
		return notifierTelegram.New(&notifierTelegram.TelegramNotifierConfig{
			ApiToken: maps.GetValueAsString(channelConfig, "apiToken"),
			ChatId:   maps.GetValueAsInt64(channelConfig, "chatId"),
		})

	case domain.NotifyChannelServerChan:
		return notifierServerChan.New(&notifierServerChan.ServerChanNotifierConfig{
			Url: maps.GetValueAsString(channelConfig, "url"),
		})

	case domain.NotifyChannelBark:
		return notifierBark.New(&notifierBark.BarkNotifierConfig{
			DeviceKey: maps.GetValueAsString(channelConfig, "deviceKey"),
			ServerUrl: maps.GetValueAsString(channelConfig, "serverUrl"),
		})
	}

	return nil, errors.New("unsupported notifier channel")
}
