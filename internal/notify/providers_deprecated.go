package notify

import (
	"fmt"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core"
	pBark "github.com/certimate-go/certimate/pkg/core/notifier/providers/bark"
	pDingTalk "github.com/certimate-go/certimate/pkg/core/notifier/providers/dingtalkbot"
	pEmail "github.com/certimate-go/certimate/pkg/core/notifier/providers/email"
	pGotify "github.com/certimate-go/certimate/pkg/core/notifier/providers/gotify"
	pLark "github.com/certimate-go/certimate/pkg/core/notifier/providers/larkbot"
	pMattermost "github.com/certimate-go/certimate/pkg/core/notifier/providers/mattermost"
	pPushover "github.com/certimate-go/certimate/pkg/core/notifier/providers/pushover"
	pPushPlus "github.com/certimate-go/certimate/pkg/core/notifier/providers/pushplus"
	pServerChan "github.com/certimate-go/certimate/pkg/core/notifier/providers/serverchan"
	pTelegram "github.com/certimate-go/certimate/pkg/core/notifier/providers/telegrambot"
	pWebhook "github.com/certimate-go/certimate/pkg/core/notifier/providers/webhook"
	pWeCom "github.com/certimate-go/certimate/pkg/core/notifier/providers/wecombot"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

// Deprecated: v0.4.x 将废弃
func createNotifierProviderUseGlobalSettings(channel domain.NotifyChannelType, channelConfig map[string]any) (core.Notifier, error) {
	/*
	  注意：如果追加新的常量值，请保持以 ASCII 排序。
	  NOTICE: If you add new constant, please keep ASCII order.
	*/
	switch channel {
	case domain.NotifyChannelTypeBark:
		return pBark.NewNotifierProvider(&pBark.NotifierProviderConfig{
			DeviceKey: xmaps.GetString(channelConfig, "deviceKey"),
			ServerUrl: xmaps.GetString(channelConfig, "serverUrl"),
		})

	case domain.NotifyChannelTypeDingTalk:
		return pDingTalk.NewNotifierProvider(&pDingTalk.NotifierProviderConfig{
			WebhookUrl: "https://oapi.dingtalk.com/robot/send?access_token=" + xmaps.GetString(channelConfig, "accessToken"),
			Secret:     xmaps.GetString(channelConfig, "secret"),
		})

	case domain.NotifyChannelTypeEmail:
		return pEmail.NewNotifierProvider(&pEmail.NotifierProviderConfig{
			SmtpHost:        xmaps.GetString(channelConfig, "smtpHost"),
			SmtpPort:        xmaps.GetInt32(channelConfig, "smtpPort"),
			SmtpTls:         xmaps.GetOrDefaultBool(channelConfig, "smtpTLS", true),
			Username:        xmaps.GetOrDefaultString(channelConfig, "username", xmaps.GetString(channelConfig, "senderAddress")),
			Password:        xmaps.GetString(channelConfig, "password"),
			SenderAddress:   xmaps.GetString(channelConfig, "senderAddress"),
			ReceiverAddress: xmaps.GetString(channelConfig, "receiverAddress"),
		})

	case domain.NotifyChannelTypeGotify:
		return pGotify.NewNotifierProvider(&pGotify.NotifierProviderConfig{
			ServerUrl: xmaps.GetString(channelConfig, "url"),
			Token:     xmaps.GetString(channelConfig, "token"),
			Priority:  xmaps.GetOrDefaultInt64(channelConfig, "priority", 1),
		})

	case domain.NotifyChannelTypeLark:
		return pLark.NewNotifierProvider(&pLark.NotifierProviderConfig{
			WebhookUrl: xmaps.GetString(channelConfig, "webhookUrl"),
		})

	case domain.NotifyChannelTypeMattermost:
		return pMattermost.NewNotifierProvider(&pMattermost.NotifierProviderConfig{
			ServerUrl: xmaps.GetString(channelConfig, "serverUrl"),
			ChannelId: xmaps.GetString(channelConfig, "channelId"),
			Username:  xmaps.GetString(channelConfig, "username"),
			Password:  xmaps.GetString(channelConfig, "password"),
		})

	case domain.NotifyChannelTypePushover:
		return pPushover.NewNotifierProvider(&pPushover.NotifierProviderConfig{
			Token: xmaps.GetString(channelConfig, "token"),
			User:  xmaps.GetString(channelConfig, "user"),
		})

	case domain.NotifyChannelTypePushPlus:
		return pPushPlus.NewNotifierProvider(&pPushPlus.NotifierProviderConfig{
			Token: xmaps.GetString(channelConfig, "token"),
		})

	case domain.NotifyChannelTypeServerChan:
		return pServerChan.NewNotifierProvider(&pServerChan.NotifierProviderConfig{
			ServerUrl: xmaps.GetString(channelConfig, "url"),
		})

	case domain.NotifyChannelTypeTelegram:
		return pTelegram.NewNotifierProvider(&pTelegram.NotifierProviderConfig{
			BotToken: xmaps.GetString(channelConfig, "apiToken"),
			ChatId:   xmaps.GetInt64(channelConfig, "chatId"),
		})

	case domain.NotifyChannelTypeWebhook:
		return pWebhook.NewNotifierProvider(&pWebhook.NotifierProviderConfig{
			WebhookUrl:               xmaps.GetString(channelConfig, "url"),
			AllowInsecureConnections: xmaps.GetBool(channelConfig, "allowInsecureConnections"),
		})

	case domain.NotifyChannelTypeWeCom:
		return pWeCom.NewNotifierProvider(&pWeCom.NotifierProviderConfig{
			WebhookUrl: xmaps.GetString(channelConfig, "webhookUrl"),
		})
	}

	return nil, fmt.Errorf("unsupported notifier channel '%s'", channelConfig)
}
