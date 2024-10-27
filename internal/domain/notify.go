package domain

const (
	NotifyChannelDingtalk = "dingtalk"
	NotifyChannelWebhook  = "webhook"
	NotifyChannelTelegram = "telegram"
	NotifyChannelLark     = "lark"
	NotifyChannelServerChan = "serverchan"
)

type NotifyTestPushReq struct {
	Channel string `json:"channel"`
}
