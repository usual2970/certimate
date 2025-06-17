package notify

import (
	"fmt"
	"net/http"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/pkg/core"
	pDingTalkBot "github.com/certimate-go/certimate/pkg/core/notifier/providers/dingtalkbot"
	pDiscordBot "github.com/certimate-go/certimate/pkg/core/notifier/providers/discordbot"
	pEmail "github.com/certimate-go/certimate/pkg/core/notifier/providers/email"
	pLarkBot "github.com/certimate-go/certimate/pkg/core/notifier/providers/larkbot"
	pMattermost "github.com/certimate-go/certimate/pkg/core/notifier/providers/mattermost"
	pSlackBot "github.com/certimate-go/certimate/pkg/core/notifier/providers/slackbot"
	pTelegramBot "github.com/certimate-go/certimate/pkg/core/notifier/providers/telegrambot"
	pWebhook "github.com/certimate-go/certimate/pkg/core/notifier/providers/webhook"
	pWeComBot "github.com/certimate-go/certimate/pkg/core/notifier/providers/wecombot"
	xhttp "github.com/certimate-go/certimate/pkg/utils/http"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

type notifierProviderOptions struct {
	Provider              domain.NotificationProviderType
	ProviderAccessConfig  map[string]any
	ProviderServiceConfig map[string]any
}

func createNotifierProvider(options *notifierProviderOptions) (core.Notifier, error) {
	/*
	  注意：如果追加新的常量值，请保持以 ASCII 排序。
	  NOTICE: If you add new constant, please keep ASCII order.
	*/
	switch options.Provider {
	case domain.NotificationProviderTypeDingTalkBot:
		{
			access := domain.AccessConfigForDingTalkBot{}
			if err := xmaps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			return pDingTalkBot.NewNotifierProvider(&pDingTalkBot.NotifierProviderConfig{
				WebhookUrl: access.WebhookUrl,
				Secret:     access.Secret,
			})
		}

	case domain.NotificationProviderTypeDiscordBot:
		{
			access := domain.AccessConfigForDiscordBot{}
			if err := xmaps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			return pDiscordBot.NewNotifierProvider(&pDiscordBot.NotifierProviderConfig{
				BotToken:  access.BotToken,
				ChannelId: xmaps.GetOrDefaultString(options.ProviderServiceConfig, "channelId", access.DefaultChannelId),
			})
		}

	case domain.NotificationProviderTypeEmail:
		{
			access := domain.AccessConfigForEmail{}
			if err := xmaps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			return pEmail.NewNotifierProvider(&pEmail.NotifierProviderConfig{
				SmtpHost:        access.SmtpHost,
				SmtpPort:        access.SmtpPort,
				SmtpTls:         access.SmtpTls,
				Username:        access.Username,
				Password:        access.Password,
				SenderAddress:   xmaps.GetOrDefaultString(options.ProviderServiceConfig, "senderAddress", access.DefaultSenderAddress),
				SenderName:      xmaps.GetOrDefaultString(options.ProviderServiceConfig, "senderName", access.DefaultSenderName),
				ReceiverAddress: xmaps.GetOrDefaultString(options.ProviderServiceConfig, "receiverAddress", access.DefaultReceiverAddress),
			})
		}

	case domain.NotificationProviderTypeLarkBot:
		{
			access := domain.AccessConfigForLarkBot{}
			if err := xmaps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			return pLarkBot.NewNotifierProvider(&pLarkBot.NotifierProviderConfig{
				WebhookUrl: access.WebhookUrl,
			})
		}

	case domain.NotificationProviderTypeMattermost:
		{
			access := domain.AccessConfigForMattermost{}
			if err := xmaps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			return pMattermost.NewNotifierProvider(&pMattermost.NotifierProviderConfig{
				ServerUrl: access.ServerUrl,
				Username:  access.Username,
				Password:  access.Password,
				ChannelId: xmaps.GetOrDefaultString(options.ProviderServiceConfig, "channelId", access.DefaultChannelId),
			})
		}

	case domain.NotificationProviderTypeSlackBot:
		{
			access := domain.AccessConfigForSlackBot{}
			if err := xmaps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			return pSlackBot.NewNotifierProvider(&pSlackBot.NotifierProviderConfig{
				BotToken:  access.BotToken,
				ChannelId: xmaps.GetOrDefaultString(options.ProviderServiceConfig, "channelId", access.DefaultChannelId),
			})
		}

	case domain.NotificationProviderTypeTelegramBot:
		{
			access := domain.AccessConfigForTelegramBot{}
			if err := xmaps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			return pTelegramBot.NewNotifierProvider(&pTelegramBot.NotifierProviderConfig{
				BotToken: access.BotToken,
				ChatId:   xmaps.GetOrDefaultInt64(options.ProviderServiceConfig, "chatId", access.DefaultChatId),
			})
		}

	case domain.NotificationProviderTypeWebhook:
		{
			access := domain.AccessConfigForWebhook{}
			if err := xmaps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			mergedHeaders := make(map[string]string)
			if defaultHeadersString := access.HeadersString; defaultHeadersString != "" {
				h, err := xhttp.ParseHeaders(defaultHeadersString)
				if err != nil {
					return nil, fmt.Errorf("failed to parse webhook headers: %w", err)
				}
				for key := range h {
					mergedHeaders[http.CanonicalHeaderKey(key)] = h.Get(key)
				}
			}
			if extendedHeadersString := xmaps.GetString(options.ProviderServiceConfig, "headers"); extendedHeadersString != "" {
				h, err := xhttp.ParseHeaders(extendedHeadersString)
				if err != nil {
					return nil, fmt.Errorf("failed to parse webhook headers: %w", err)
				}
				for key := range h {
					mergedHeaders[http.CanonicalHeaderKey(key)] = h.Get(key)
				}
			}

			return pWebhook.NewNotifierProvider(&pWebhook.NotifierProviderConfig{
				WebhookUrl:               access.Url,
				WebhookData:              xmaps.GetOrDefaultString(options.ProviderServiceConfig, "webhookData", access.DefaultDataForNotification),
				Method:                   access.Method,
				Headers:                  mergedHeaders,
				AllowInsecureConnections: access.AllowInsecureConnections,
			})
		}

	case domain.NotificationProviderTypeWeComBot:
		{
			access := domain.AccessConfigForWeComBot{}
			if err := xmaps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			return pWeComBot.NewNotifierProvider(&pWeComBot.NotifierProviderConfig{
				WebhookUrl: access.WebhookUrl,
			})
		}
	}

	return nil, fmt.Errorf("unsupported notifier provider '%s'", options.Provider)
}
