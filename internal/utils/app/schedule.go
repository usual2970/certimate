package app

import (
	"sync"
	"time"

	"github.com/pocketbase/pocketbase/tools/cron"
)

var schedulerOnce sync.Once

var scheduler *cron.Cron

func GetScheduler() *cron.Cron {
	schedulerOnce.Do(func() {
		scheduler = cron.New()
		location, err := time.LoadLocation("Asia/Shanghai")
		if err == nil {
			scheduler.SetTimezone(location)
		}
	})

	return scheduler
}
