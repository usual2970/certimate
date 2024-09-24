package notify

import (
	"certimate/internal/utils/app"
	"context"
	"fmt"
	"strconv"

	notifyPackage "github.com/nikoksr/notify"

	"github.com/nikoksr/notify/service/dingding"

	"github.com/nikoksr/notify/service/telegram"

	"github.com/nikoksr/notify/service/http"
)

const (
	notifyChannelDingtalk = "dingtalk"
	notifyChannelWebhook  = "webhook"
	notifyChannelTelegram = "telegram"
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

	// 添加推送渠道
	notifyPackage.UseServices(notifiers...)

	// 发送消息
	return notifyPackage.Send(context.Background(), title, content)
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

		switch k {
		case notifyChannelTelegram:
			temp := getTelegramNotifier(v)
			if temp == nil {
				continue
			}

			notifiers = append(notifiers, temp)
		case notifyChannelDingtalk:
			notifiers = append(notifiers, getDingTalkNotifier(v))
		case notifyChannelWebhook:
			notifiers = append(notifiers, getWebhookNotifier(v))
		}

	}

	return notifiers, nil

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
