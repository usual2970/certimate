package notify

import (
	"errors"

	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/pkg/core/notifier"
	npBark "github.com/usual2970/certimate/internal/pkg/core/notifier/providers/bark"
	npDingTalk "github.com/usual2970/certimate/internal/pkg/core/notifier/providers/dingtalk"
	npEmail "github.com/usual2970/certimate/internal/pkg/core/notifier/providers/email"
	npLark "github.com/usual2970/certimate/internal/pkg/core/notifier/providers/lark"
	npServerChan "github.com/usual2970/certimate/internal/pkg/core/notifier/providers/serverchan"
	npTelegram "github.com/usual2970/certimate/internal/pkg/core/notifier/providers/telegram"
	npWebhook "github.com/usual2970/certimate/internal/pkg/core/notifier/providers/webhook"
	"github.com/usual2970/certimate/internal/pkg/utils/maps"
)

func createNotifier(channel string, channelConfig map[string]any) (notifier.Notifier, error) {
	switch channel {
	case domain.NotifyChannelEmail:
		return npEmail.New(&npEmail.EmailNotifierConfig{
			SmtpHost:        maps.GetValueAsString(channelConfig, "smtpHost"),
			SmtpPort:        maps.GetValueAsInt32(channelConfig, "smtpPort"),
			SmtpTLS:         maps.GetValueOrDefaultAsBool(channelConfig, "smtpTLS", true),
			Username:        maps.GetValueOrDefaultAsString(channelConfig, "username", maps.GetValueAsString(channelConfig, "senderAddress")),
			Password:        maps.GetValueAsString(channelConfig, "password"),
			SenderAddress:   maps.GetValueAsString(channelConfig, "senderAddress"),
			ReceiverAddress: maps.GetValueAsString(channelConfig, "receiverAddress"),
		})

	case domain.NotifyChannelWebhook:
		return npWebhook.New(&npWebhook.WebhookNotifierConfig{
			Url: maps.GetValueAsString(channelConfig, "url"),
		})

	case domain.NotifyChannelDingtalk:
		return npDingTalk.New(&npDingTalk.DingTalkNotifierConfig{
			AccessToken: maps.GetValueAsString(channelConfig, "accessToken"),
			Secret:      maps.GetValueAsString(channelConfig, "secret"),
		})

	case domain.NotifyChannelLark:
		return npLark.New(&npLark.LarkNotifierConfig{
			WebhookUrl: maps.GetValueAsString(channelConfig, "webhookUrl"),
		})

	case domain.NotifyChannelTelegram:
		return npTelegram.New(&npTelegram.TelegramNotifierConfig{
			ApiToken: maps.GetValueAsString(channelConfig, "apiToken"),
			ChatId:   maps.GetValueAsInt64(channelConfig, "chatId"),
		})

	case domain.NotifyChannelServerChan:
		return npServerChan.New(&npServerChan.ServerChanNotifierConfig{
			Url: maps.GetValueAsString(channelConfig, "url"),
		})

	case domain.NotifyChannelBark:
		return npBark.New(&npBark.BarkNotifierConfig{
			DeviceKey: maps.GetValueAsString(channelConfig, "deviceKey"),
			ServerUrl: maps.GetValueAsString(channelConfig, "serverUrl"),
		})
	}

	return nil, errors.New("unsupported notifier channel")
}
