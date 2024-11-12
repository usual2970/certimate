package domain

const (
	NotifyChannelEmail      = "email"
	NotifyChannelWebhook    = "webhook"
	NotifyChannelDingtalk   = "dingtalk"
	NotifyChannelLark       = "lark"
	NotifyChannelTelegram   = "telegram"
	NotifyChannelServerChan = "serverchan"
	NotifyChannelBark       = "bark"
)

type NotifyTestPushReq struct {
	Channel string `json:"channel"`
}
