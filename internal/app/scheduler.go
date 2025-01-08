package app

import (
	"sync"
	"time"

	"github.com/pocketbase/pocketbase/tools/cron"
)

var scheduler *cron.Cron

var schedulerOnce sync.Once

func GetScheduler() *cron.Cron {
	schedulerOnce.Do(func() {
		scheduler = cron.New()
		location, err := time.LoadLocation("Local")
		if err == nil {
			scheduler.SetTimezone(location)
		}
	})

	return scheduler
}
