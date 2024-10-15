package domains

import (
	"context"

	"certimate/internal/notify"
	"certimate/internal/utils/app"
)

func InitSchedule() {
	// 查询所有启用的域名
	records, err := app.GetApp().Dao().FindRecordsByFilter("domains", "enabled=true", "-id", 500, 0)
	if err != nil {
		app.GetApp().Logger().Error("查询所有启用的域名失败", "err", err)
		return
	}

	// 加入到定时任务
	for _, record := range records {
		if err := app.GetScheduler().Add(record.Id, record.GetString("crontab"), func() {
			if err := deploy(context.Background(), record); err != nil {
				app.GetApp().Logger().Error("部署失败", "err", err)
				return
			}
		}); err != nil {
			app.GetApp().Logger().Error("加入到定时任务失败", "err", err)
		}
	}

	// 过期提醒
	app.GetScheduler().Add("expire", "0 0 * * *", func() {
		notify.PushExpireMsg()
	})

	// 启动定时任务
	app.GetScheduler().Start()
	app.GetApp().Logger().Info("定时任务启动成功", "total", app.GetScheduler().Total())
}
