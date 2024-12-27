package notify

import (
	"fmt"

	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/pkg/core/notifier"
	providerBark "github.com/usual2970/certimate/internal/pkg/core/notifier/providers/bark"
	providerDingTalk "github.com/usual2970/certimate/internal/pkg/core/notifier/providers/dingtalk"
	providerEmail "github.com/usual2970/certimate/internal/pkg/core/notifier/providers/email"
	providerLark "github.com/usual2970/certimate/internal/pkg/core/notifier/providers/lark"
	providerServerChan "github.com/usual2970/certimate/internal/pkg/core/notifier/providers/serverchan"
	providerTelegram "github.com/usual2970/certimate/internal/pkg/core/notifier/providers/telegram"
	providerWebhook "github.com/usual2970/certimate/internal/pkg/core/notifier/providers/webhook"
	providerWeCom "github.com/usual2970/certimate/internal/pkg/core/notifier/providers/wecom"
	"github.com/usual2970/certimate/internal/pkg/utils/maps"
)

func createNotifier(channel domain.NotifyChannelType, channelConfig map[string]any) (notifier.Notifier, error) {
	/*
	  注意：如果追加新的常量值，请保持以 ASCII 排序。
	  NOTICE: If you add new constant, please keep ASCII order.
	*/
	switch channel {
	case domain.NOTIFY_CHANNEL_BARK:
		return providerBark.New(&providerBark.BarkNotifierConfig{
			DeviceKey: maps.GetValueAsString(channelConfig, "deviceKey"),
			ServerUrl: maps.GetValueAsString(channelConfig, "serverUrl"),
		})

	case domain.NOTIFY_CHANNEL_DINGTALK:
		return providerDingTalk.New(&providerDingTalk.DingTalkNotifierConfig{
			AccessToken: maps.GetValueAsString(channelConfig, "accessToken"),
			Secret:      maps.GetValueAsString(channelConfig, "secret"),
		})

	case domain.NOTIFY_CHANNEL_EMAIL:
		return providerEmail.New(&providerEmail.EmailNotifierConfig{
			SmtpHost:        maps.GetValueAsString(channelConfig, "smtpHost"),
			SmtpPort:        maps.GetValueAsInt32(channelConfig, "smtpPort"),
			SmtpTLS:         maps.GetValueOrDefaultAsBool(channelConfig, "smtpTLS", true),
			Username:        maps.GetValueOrDefaultAsString(channelConfig, "username", maps.GetValueAsString(channelConfig, "senderAddress")),
			Password:        maps.GetValueAsString(channelConfig, "password"),
			SenderAddress:   maps.GetValueAsString(channelConfig, "senderAddress"),
			ReceiverAddress: maps.GetValueAsString(channelConfig, "receiverAddress"),
		})

	case domain.NOTIFY_CHANNEL_LARK:
		return providerLark.New(&providerLark.LarkNotifierConfig{
			WebhookUrl: maps.GetValueAsString(channelConfig, "webhookUrl"),
		})

	case domain.NOTIFY_CHANNEL_SERVERCHAN:
		return providerServerChan.New(&providerServerChan.ServerChanNotifierConfig{
			Url: maps.GetValueAsString(channelConfig, "url"),
		})

	case domain.NOTIFY_CHANNEL_TELEGRAM:
		return providerTelegram.New(&providerTelegram.TelegramNotifierConfig{
			ApiToken: maps.GetValueAsString(channelConfig, "apiToken"),
			ChatId:   maps.GetValueAsInt64(channelConfig, "chatId"),
		})

	case domain.NOTIFY_CHANNEL_WEBHOOK:
		return providerWebhook.New(&providerWebhook.WebhookNotifierConfig{
			Url: maps.GetValueAsString(channelConfig, "url"),
		})

	case domain.NOTIFY_CHANNEL_WECOM:
		return providerWeCom.New(&providerWeCom.WeComNotifierConfig{
			WebhookUrl: maps.GetValueAsString(channelConfig, "webhookUrl"),
		})
	}

	return nil, fmt.Errorf("unsupported notifier channel: %s", channelConfig)
}
