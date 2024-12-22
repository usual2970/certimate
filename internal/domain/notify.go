package domain

/*
消息通知渠道常量值。

	注意：如果追加新的常量值，请保持以 ASCII 排序。
	NOTICE: If you add new constant, please keep ASCII order.
*/
const (
	NotifyChannelBark       = "bark"
	NotifyChannelDingtalk   = "dingtalk"
	NotifyChannelEmail      = "email"
	NotifyChannelLark       = "lark"
	NotifyChannelServerChan = "serverchan"
	NotifyChannelTelegram   = "telegram"
	NotifyChannelWebhook    = "webhook"
	NotifyChannelWeCom      = "wecom"
)

type NotifyTestPushReq struct {
	Channel string `json:"channel"`
}
