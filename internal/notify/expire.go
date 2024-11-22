package notify

import (
	"strconv"
	"strings"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/models"

	"github.com/usual2970/certimate/internal/utils/app"
	"github.com/usual2970/certimate/internal/utils/xtime"
)

const (
	defaultExpireSubject = "您有 {COUNT} 张证书即将过期"
	defaultExpireMessage = "有 {COUNT} 张证书即将过期，域名分别为 {DOMAINS}，请保持关注！"
)

func PushExpireMsg() {
	// 查询即将过期的证书
	records, err := app.GetApp().Dao().FindRecordsByFilter("certificate", "expireAt<{:time}&&certUrl!=''", "-created", 500, 0,
		dbx.Params{"time": xtime.GetTimeAfter(24 * time.Hour * 20)})
	if err != nil {
		app.GetApp().Logger().Error("find expired domains by filter", "error", err)
		return
	}

	// 组装消息
	msg := buildMsg(records)
	if msg == nil {
		return
	}

	// 发送通知
	if err := SendToAllChannels(msg.Subject, msg.Message); err != nil {
		app.GetApp().Logger().Error("send expire msg", "error", err)
	}
}

type notifyTemplates struct {
	NotifyTemplates []notifyTemplate `json:"notifyTemplates"`
}

type notifyTemplate struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type notifyMessage struct {
	Subject string
	Message string
}

func buildMsg(records []*models.Record) *notifyMessage {
	if len(records) == 0 {
		return nil
	}

	// 查询模板信息
	templateRecord, err := app.GetApp().Dao().FindFirstRecordByFilter("settings", "name='templates'")
	subject := defaultExpireSubject
	message := defaultExpireMessage

	if err == nil {
		var templates *notifyTemplates
		templateRecord.UnmarshalJSONField("content", templates)
		if templates != nil && len(templates.NotifyTemplates) > 0 {
			subject = templates.NotifyTemplates[0].Title
			message = templates.NotifyTemplates[0].Content
		}
	}

	// 替换变量
	count := len(records)
	domains := make([]string, count)

	for i, record := range records {
		domains[i] = record.GetString("san")
	}

	countStr := strconv.Itoa(count)
	domainStr := strings.Join(domains, ";")

	subject = strings.ReplaceAll(subject, "{COUNT}", countStr)
	subject = strings.ReplaceAll(subject, "{DOMAINS}", domainStr)

	message = strings.ReplaceAll(message, "{COUNT}", countStr)
	message = strings.ReplaceAll(message, "{DOMAINS}", domainStr)

	// 返回消息
	return &notifyMessage{
		Subject: subject,
		Message: message,
	}
}
