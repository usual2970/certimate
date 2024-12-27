package domain

type NotifyChannelType string

/*
消息通知渠道常量值。

	注意：如果追加新的常量值，请保持以 ASCII 排序。
	NOTICE: If you add new constant, please keep ASCII order.
*/
const (
	NOTIFY_CHANNEL_BARK       = NotifyChannelType("bark")
	NOTIFY_CHANNEL_DINGTALK   = NotifyChannelType("dingtalk")
	NOTIFY_CHANNEL_EMAIL      = NotifyChannelType("email")
	NOTIFY_CHANNEL_LARK       = NotifyChannelType("lark")
	NOTIFY_CHANNEL_SERVERCHAN = NotifyChannelType("serverchan")
	NOTIFY_CHANNEL_TELEGRAM   = NotifyChannelType("telegram")
	NOTIFY_CHANNEL_WEBHOOK    = NotifyChannelType("webhook")
	NOTIFY_CHANNEL_WECOM      = NotifyChannelType("wecom")
)

type NotifyTestPushReq struct {
	Channel string `json:"channel"`
}
