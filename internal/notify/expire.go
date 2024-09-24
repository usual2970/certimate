package notify

import (
	"certimate/internal/utils/app"
	"certimate/internal/utils/xtime"
	"strconv"
	"strings"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/models"
)

type msg struct {
	subject string
	message string
}

const (
	defaultExpireSubject = "您有{COUNT}张证书即将过期"
	defaultExpireMsg     = "有{COUNT}张证书即将过期,域名分别为{DOMAINS},请保持关注！"
)

func PushExpireMsg() {
	// 查询即将过期的证书

	records, err := app.GetApp().Dao().FindRecordsByFilter("domains", "expiredAt<{:time}&&certUrl!=''", "-created", 500, 0,
		dbx.Params{"time": xtime.GetTimeAfter(24 * time.Hour * 15)})
	if err != nil {
		app.GetApp().Logger().Error("find expired domains by filter", "error", err)
		return
	}

	// 组装消息
	msg := buildMsg(records)

	if msg == nil {
		return
	}

	if err := Send(msg.subject, msg.message); err != nil {
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

func buildMsg(records []*models.Record) *msg {
	if len(records) == 0 {
		return nil
	}

	// 查询模板信息
	templateRecord, err := app.GetApp().Dao().FindFirstRecordByFilter("settings", "name='templates'")
	title := defaultExpireSubject
	content := defaultExpireMsg

	if err == nil {
		var templates *notifyTemplates
		templateRecord.UnmarshalJSONField("content", templates)
		if templates != nil && len(templates.NotifyTemplates) > 0 {
			title = templates.NotifyTemplates[0].Title
			content = templates.NotifyTemplates[0].Content
		}
	}

	// 替换变量
	count := len(records)
	domains := make([]string, count)

	for i, record := range records {
		domains[i] = record.GetString("domain")
	}

	countStr := strconv.Itoa(count)
	domainStr := strings.Join(domains, ",")

	title = strings.ReplaceAll(title, "{COUNT}", countStr)
	title = strings.ReplaceAll(title, "{DOMAINS}", domainStr)

	content = strings.ReplaceAll(content, "{COUNT}", countStr)
	content = strings.ReplaceAll(content, "{DOMAINS}", domainStr)

	// 返回消息
	return &msg{
		subject: title,
		message: content,
	}

}
