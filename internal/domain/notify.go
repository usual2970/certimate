package domain

const (
	NotifyChannelDingtalk = "dingtalk"
	NotifyChannelWebhook  = "webhook"
	NotifyChannelTelegram = "telegram"
	NotifyChannelLark     = "lark"
	NotifyChannelServerChan = "serverchan"
	NotifyChannelMail = "mail"
)

type NotifyTestPushReq struct {
	Channel string `json:"channel"`
}
