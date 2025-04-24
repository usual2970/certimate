package notify

import (
	"fmt"

	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/pkg/core/notifier"
	pWebhook "github.com/usual2970/certimate/internal/pkg/core/notifier/providers/webhook"
	maputil "github.com/usual2970/certimate/internal/pkg/utils/map"
)

type notifierProviderOptions struct {
	Provider             domain.NotifyProviderType
	ProviderAccessConfig map[string]any
	ProviderNotifyConfig map[string]any
}

func createNotifierProvider(options *notifierProviderOptions) (notifier.Notifier, error) {
	/*
	  注意：如果追加新的常量值，请保持以 ASCII 排序。
	  NOTICE: If you add new constant, please keep ASCII order.
	*/
	switch options.Provider {
	case domain.NotifyProviderTypeWebhook:
		return pWebhook.NewNotifier(&pWebhook.NotifierConfig{
			Url:                      maputil.GetString(options.ProviderAccessConfig, "url"),
			AllowInsecureConnections: maputil.GetBool(options.ProviderAccessConfig, "allowInsecureConnections"),
		})
	}

	return nil, fmt.Errorf("unsupported notifier provider '%s'", options.Provider)
}
