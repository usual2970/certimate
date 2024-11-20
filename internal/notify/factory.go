package notify

import (
	"errors"

	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/pkg/core/notifier"
	providerBark "github.com/usual2970/certimate/internal/pkg/core/notifier/providers/bark"
	providerDingTalk "github.com/usual2970/certimate/internal/pkg/core/notifier/providers/dingtalk"
	providerEmail "github.com/usual2970/certimate/internal/pkg/core/notifier/providers/email"
	providerLark "github.com/usual2970/certimate/internal/pkg/core/notifier/providers/lark"
	providerServerChan "github.com/usual2970/certimate/internal/pkg/core/notifier/providers/serverchan"
	providerTelegram "github.com/usual2970/certimate/internal/pkg/core/notifier/providers/telegram"
	providerWebhook "github.com/usual2970/certimate/internal/pkg/core/notifier/providers/webhook"
	"github.com/usual2970/certimate/internal/pkg/utils/maps"
)

func createNotifier(channel string, channelConfig map[string]any) (notifier.Notifier, error) {
	switch channel {
	case domain.NotifyChannelEmail:
		return providerEmail.New(&providerEmail.EmailNotifierConfig{
			SmtpHost:        maps.GetValueAsString(channelConfig, "smtpHost"),
			SmtpPort:        maps.GetValueAsInt32(channelConfig, "smtpPort"),
			SmtpTLS:         maps.GetValueOrDefaultAsBool(channelConfig, "smtpTLS", true),
			Username:        maps.GetValueOrDefaultAsString(channelConfig, "username", maps.GetValueAsString(channelConfig, "senderAddress")),
			Password:        maps.GetValueAsString(channelConfig, "password"),
			SenderAddress:   maps.GetValueAsString(channelConfig, "senderAddress"),
			ReceiverAddress: maps.GetValueAsString(channelConfig, "receiverAddress"),
		})

	case domain.NotifyChannelWebhook:
		return providerWebhook.New(&providerWebhook.WebhookNotifierConfig{
			Url: maps.GetValueAsString(channelConfig, "url"),
		})

	case domain.NotifyChannelDingtalk:
		return providerDingTalk.New(&providerDingTalk.DingTalkNotifierConfig{
			AccessToken: maps.GetValueAsString(channelConfig, "accessToken"),
			Secret:      maps.GetValueAsString(channelConfig, "secret"),
		})

	case domain.NotifyChannelLark:
		return providerLark.New(&providerLark.LarkNotifierConfig{
			WebhookUrl: maps.GetValueAsString(channelConfig, "webhookUrl"),
		})

	case domain.NotifyChannelTelegram:
		return providerTelegram.New(&providerTelegram.TelegramNotifierConfig{
			ApiToken: maps.GetValueAsString(channelConfig, "apiToken"),
			ChatId:   maps.GetValueAsInt64(channelConfig, "chatId"),
		})

	case domain.NotifyChannelServerChan:
		return providerServerChan.New(&providerServerChan.ServerChanNotifierConfig{
			Url: maps.GetValueAsString(channelConfig, "url"),
		})

	case domain.NotifyChannelBark:
		return providerBark.New(&providerBark.BarkNotifierConfig{
			DeviceKey: maps.GetValueAsString(channelConfig, "deviceKey"),
			ServerUrl: maps.GetValueAsString(channelConfig, "serverUrl"),
		})
	}

	return nil, errors.New("unsupported notifier channel")
}
