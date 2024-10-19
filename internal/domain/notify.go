package domain

const (
	NotifyChannelDingtalk = "dingtalk"
	NotifyChannelWebhook  = "webhook"
	NotifyChannelTelegram = "telegram"
	NotifyChannelLark     = "lark"
)

type NotifyTestPushReq struct {
	Channel string `json:"channel"`
}
