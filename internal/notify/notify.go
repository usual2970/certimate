package notify

import (
	"context"
	"fmt"

	stdhttp "net/http"

	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/utils/app"

	notifyPackage "github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/bark"
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

		if !getConfigAsBool(v, "enabled") {
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
	case domain.NotifyChannelServerChan:
		return getServerChanNotifier(conf), nil
	case domain.NotifyChannelMail:
		return getMailNotifier(conf)
	case domain.NotifyChannelBark:
		return getBarkNotifier(conf), nil
	}

	return nil, fmt.Errorf("notifier not found")
}

func getWebhookNotifier(conf map[string]any) notifyPackage.Notifier {
	rs := http.New()

	rs.AddReceiversURLs(getConfigAsString(conf, "url"))

	return rs
}

func getTelegramNotifier(conf map[string]any) notifyPackage.Notifier {
	rs, err := telegram.New(getConfigAsString(conf, "apiToken"))
	if err != nil {
		return nil
	}

	rs.AddReceivers(getConfigAsInt64(conf, "chatId"))
	return rs
}

func getServerChanNotifier(conf map[string]any) notifyPackage.Notifier {
	rs := http.New()

	rs.AddReceivers(&http.Webhook{
		URL:         getConfigAsString(conf, "url"),
		Header:      stdhttp.Header{},
		ContentType: "application/json",
		Method:      stdhttp.MethodPost,
		BuildPayload: func(subject, message string) (payload any) {
			return map[string]string{
				"text": subject,
				"desp": message,
			}
		},
	})

	return rs
}

func getBarkNotifier(conf map[string]any) notifyPackage.Notifier {
	deviceKey := getConfigAsString(conf, "deviceKey")
	serverURL := getConfigAsString(conf, "serverUrl")
	if serverURL == "" {
		return bark.New(deviceKey)
	}
	return bark.NewWithServers(deviceKey, serverURL)
}

func getDingTalkNotifier(conf map[string]any) notifyPackage.Notifier {
	return dingding.New(&dingding.Config{
		Token:  getConfigAsString(conf, "accessToken"),
		Secret: getConfigAsString(conf, "secret"),
	})
}

func getLarkNotifier(conf map[string]any) notifyPackage.Notifier {
	return lark.NewWebhookService(getConfigAsString(conf, "webhookUrl"))
}

func getMailNotifier(conf map[string]any) (notifyPackage.Notifier, error) {
	rs, err := NewMail(getConfigAsString(conf, "senderAddress"),
		getConfigAsString(conf, "receiverAddresses"),
		getConfigAsString(conf, "smtpHostAddr"),
		getConfigAsString(conf, "smtpHostPort"),
		getConfigAsString(conf, "password"),
	)
	if err != nil {
		return nil, err
	}

	return rs, nil
}
