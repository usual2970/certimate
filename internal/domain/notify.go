package domain

type NotifyChannelType string

/*
消息通知渠道常量值。

	注意：如果追加新的常量值，请保持以 ASCII 排序。
	NOTICE: If you add new constant, please keep ASCII order.
*/
const (
	NotifyChannelTypeBark       = NotifyChannelType("bark")
	NotifyChannelTypeDingTalk   = NotifyChannelType("dingtalk")
	NotifyChannelTypeEmail      = NotifyChannelType("email")
	NotifyChannelTypeGotify     = NotifyChannelType("gotify")
	NotifyChannelTypeLark       = NotifyChannelType("lark")
	NotifyChannelTypeMattermost = NotifyChannelType("mattermost")
	NotifyChannelTypePushover   = NotifyChannelType("pushover")
	NotifyChannelTypePushPlus   = NotifyChannelType("pushplus")
	NotifyChannelTypeServerChan = NotifyChannelType("serverchan")
	NotifyChannelTypeTelegram   = NotifyChannelType("telegram")
	NotifyChannelTypeWebhook    = NotifyChannelType("webhook")
	NotifyChannelTypeWeCom      = NotifyChannelType("wecom")
)
