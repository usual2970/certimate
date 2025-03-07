package notify

import (
	"fmt"

	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/pkg/core/notifier"
	pBark "github.com/usual2970/certimate/internal/pkg/core/notifier/providers/bark"
	pDingTalk "github.com/usual2970/certimate/internal/pkg/core/notifier/providers/dingtalk"
	pEmail "github.com/usual2970/certimate/internal/pkg/core/notifier/providers/email"
	pLark "github.com/usual2970/certimate/internal/pkg/core/notifier/providers/lark"
	pServerChan "github.com/usual2970/certimate/internal/pkg/core/notifier/providers/serverchan"
	pTelegram "github.com/usual2970/certimate/internal/pkg/core/notifier/providers/telegram"
	pWebhook "github.com/usual2970/certimate/internal/pkg/core/notifier/providers/webhook"
	pWeCom "github.com/usual2970/certimate/internal/pkg/core/notifier/providers/wecom"
	"github.com/usual2970/certimate/internal/pkg/utils/maps"
)

func createNotifier(channel domain.NotifyChannelType, channelConfig map[string]any) (notifier.Notifier, error) {
	/*
	  注意：如果追加新的常量值，请保持以 ASCII 排序。
	  NOTICE: If you add new constant, please keep ASCII order.
	*/
	switch channel {
	case domain.NotifyChannelTypeBark:
		return pBark.NewNotifier(&pBark.NotifierConfig{
			DeviceKey: maps.GetValueAsString(channelConfig, "deviceKey"),
			ServerUrl: maps.GetValueAsString(channelConfig, "serverUrl"),
		})

	case domain.NotifyChannelTypeDingTalk:
		return pDingTalk.NewNotifier(&pDingTalk.NotifierConfig{
			AccessToken: maps.GetValueAsString(channelConfig, "accessToken"),
			Secret:      maps.GetValueAsString(channelConfig, "secret"),
		})

	case domain.NotifyChannelTypeEmail:
		return pEmail.NewNotifier(&pEmail.NotifierConfig{
			SmtpHost:        maps.GetValueAsString(channelConfig, "smtpHost"),
			SmtpPort:        maps.GetValueAsInt32(channelConfig, "smtpPort"),
			SmtpTLS:         maps.GetValueOrDefaultAsBool(channelConfig, "smtpTLS", true),
			Username:        maps.GetValueOrDefaultAsString(channelConfig, "username", maps.GetValueAsString(channelConfig, "senderAddress")),
			Password:        maps.GetValueAsString(channelConfig, "password"),
			SenderAddress:   maps.GetValueAsString(channelConfig, "senderAddress"),
			ReceiverAddress: maps.GetValueAsString(channelConfig, "receiverAddress"),
		})

	case domain.NotifyChannelTypeLark:
		return pLark.NewNotifier(&pLark.NotifierConfig{
			WebhookUrl: maps.GetValueAsString(channelConfig, "webhookUrl"),
		})

	case domain.NotifyChannelTypeServerChan:
		return pServerChan.NewNotifier(&pServerChan.NotifierConfig{
			Url: maps.GetValueAsString(channelConfig, "url"),
		})

	case domain.NotifyChannelTypeTelegram:
		return pTelegram.NewNotifier(&pTelegram.NotifierConfig{
			ApiToken: maps.GetValueAsString(channelConfig, "apiToken"),
			ChatId:   maps.GetValueAsInt64(channelConfig, "chatId"),
		})

	case domain.NotifyChannelTypeWebhook:
		return pWebhook.NewNotifier(&pWebhook.NotifierConfig{
			Url:                      maps.GetValueAsString(channelConfig, "url"),
			AllowInsecureConnections: maps.GetValueAsBool(channelConfig, "allowInsecureConnections"),
		})

	case domain.NotifyChannelTypeWeCom:
		return pWeCom.NewNotifier(&pWeCom.NotifierConfig{
			WebhookUrl: maps.GetValueAsString(channelConfig, "webhookUrl"),
		})
	}

	return nil, fmt.Errorf("unsupported notifier channel: %s", channelConfig)
}
