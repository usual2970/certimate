package notify

import (
	"context"
	"fmt"
	"strconv"

	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/utils/app"

	notifyPackage "github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/dingding"
	"github.com/nikoksr/notify/service/http"
	"github.com/nikoksr/notify/service/lark"
	"github.com/nikoksr/notify/service/telegram"
)

func Send(title, content string) error {
	// 获取所有的推送渠道
	notifiers, err := getNotifiers()
	if err != nil {
		return err
	}
	if len(notifiers) == 0 {
		return nil
	}

	n := notifyPackage.New()
	// 添加推送渠道
	n.UseServices(notifiers...)

	// 发送消息
	return n.Send(context.Background(), title, content)
}

type sendTestParam struct {
	Title   string         `json:"title"`
	Content string         `json:"content"`
	Channel string         `json:"channel"`
	Conf    map[string]any `json:"conf"`
}

func SendTest(param *sendTestParam) error {
	notifier, err := getNotifier(param.Channel, param.Conf)
	if err != nil {
		return err
	}

	n := notifyPackage.New()

	// 添加推送渠道
	n.UseServices(notifier)

	// 发送消息
	return n.Send(context.Background(), param.Title, param.Content)
}

func getNotifiers() ([]notifyPackage.Notifier, error) {
	resp, err := app.GetApp().Dao().FindFirstRecordByFilter("settings", "name='notifyChannels'")
	if err != nil {
		return nil, fmt.Errorf("find notifyChannels error: %w", err)
	}

	notifiers := make([]notifyPackage.Notifier, 0)

	rs := make(map[string]map[string]any)

	if err := resp.UnmarshalJSONField("content", &rs); err != nil {
		return nil, fmt.Errorf("unmarshal notifyChannels error: %w", err)
	}

	for k, v := range rs {

		if !getBool(v, "enabled") {
			continue
		}

		notifier, err := getNotifier(k, v)
		if err != nil {
			continue
		}

		notifiers = append(notifiers, notifier)

	}

	return notifiers, nil
}

func getNotifier(channel string, conf map[string]any) (notifyPackage.Notifier, error) {
	switch channel {
	case domain.NotifyChannelTelegram:
		temp := getTelegramNotifier(conf)
		if temp == nil {
			return nil, fmt.Errorf("telegram notifier config error")
		}

		return temp, nil
	case domain.NotifyChannelDingtalk:
		return getDingTalkNotifier(conf), nil
	case domain.NotifyChannelLark:
		return getLarkNotifier(conf), nil
	case domain.NotifyChannelWebhook:
		return getWebhookNotifier(conf), nil
	}

	return nil, fmt.Errorf("notifier not found")
}

func getWebhookNotifier(conf map[string]any) notifyPackage.Notifier {
	rs := http.New()

	rs.AddReceiversURLs(getString(conf, "url"))

	return rs
}

func getTelegramNotifier(conf map[string]any) notifyPackage.Notifier {
	rs, err := telegram.New(getString(conf, "apiToken"))
	if err != nil {
		return nil
	}

	chatId := getString(conf, "chatId")

	id, err := strconv.ParseInt(chatId, 10, 64)
	if err != nil {
		return nil
	}

	rs.AddReceivers(id)
	return rs
}

func getDingTalkNotifier(conf map[string]any) notifyPackage.Notifier {
	return dingding.New(&dingding.Config{
		Token:  getString(conf, "accessToken"),
		Secret: getString(conf, "secret"),
	})
}

func getLarkNotifier(conf map[string]any) notifyPackage.Notifier {
	return lark.NewWebhookService(getString(conf, "webhookUrl"))
}

func getString(conf map[string]any, key string) string {
	if _, ok := conf[key]; !ok {
		return ""
	}

	return conf[key].(string)
}

func getBool(conf map[string]any, key string) bool {
	if _, ok := conf[key]; !ok {
		return false
	}

	return conf[key].(bool)
}
