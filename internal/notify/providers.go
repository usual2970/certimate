package notify

import (
	"fmt"

	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/pkg/core/notifier"
	pEmail "github.com/usual2970/certimate/internal/pkg/core/notifier/providers/email"
	pMattermost "github.com/usual2970/certimate/internal/pkg/core/notifier/providers/mattermost"
	pTelegram "github.com/usual2970/certimate/internal/pkg/core/notifier/providers/telegram"
	pWebhook "github.com/usual2970/certimate/internal/pkg/core/notifier/providers/webhook"
	maputil "github.com/usual2970/certimate/internal/pkg/utils/map"
)

type notifierProviderOptions struct {
	Provider             domain.NotificationProviderType
	ProviderAccessConfig map[string]any
	ProviderNotifyConfig map[string]any
}

func createNotifierProvider(options *notifierProviderOptions) (notifier.Notifier, error) {
	/*
	  注意：如果追加新的常量值，请保持以 ASCII 排序。
	  NOTICE: If you add new constant, please keep ASCII order.
	*/
	switch options.Provider {
	case domain.NotificationProviderTypeEmail:
		{
			access := domain.AccessConfigForEmail{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			return pEmail.NewNotifier(&pEmail.NotifierConfig{
				SmtpHost:        access.SmtpHost,
				SmtpPort:        access.SmtpPort,
				SmtpTls:         access.SmtpTls,
				Username:        access.Username,
				Password:        access.Password,
				SenderAddress:   maputil.GetOrDefaultString(options.ProviderNotifyConfig, "senderAddress", access.DefaultSenderAddress),
				ReceiverAddress: maputil.GetOrDefaultString(options.ProviderNotifyConfig, "receiverAddress", access.DefaultReceiverAddress),
			})
		}

	case domain.NotificationProviderTypeMattermost:
		{
			access := domain.AccessConfigForMattermost{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			return pMattermost.NewNotifier(&pMattermost.NotifierConfig{
				ServerUrl: access.ServerUrl,
				Username:  access.Username,
				Password:  access.Password,
				ChannelId: maputil.GetOrDefaultString(options.ProviderNotifyConfig, "channelId", access.DefaultChannelId),
			})
		}

	case domain.NotificationProviderTypeTelegram:
		{
			access := domain.AccessConfigForTelegram{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			return pTelegram.NewNotifier(&pTelegram.NotifierConfig{
				BotToken: access.BotToken,
				ChatId:   maputil.GetOrDefaultInt64(options.ProviderNotifyConfig, "chatId", access.DefaultChatId),
			})
		}

	case domain.NotificationProviderTypeWebhook:
		{
			access := domain.AccessConfigForWebhook{}
			if err := maputil.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			return pWebhook.NewNotifier(&pWebhook.NotifierConfig{
				Url:                      access.Url,
				AllowInsecureConnections: access.AllowInsecureConnections,
			})
		}
	}

	return nil, fmt.Errorf("unsupported notifier provider '%s'", options.Provider)
}
