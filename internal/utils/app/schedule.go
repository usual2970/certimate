package app

import (
	"sync"

	"github.com/pocketbase/pocketbase/tools/cron"
)

var schedulerOnce sync.Once

var scheduler *cron.Cron

func GetScheduler() *cron.Cron {
	schedulerOnce.Do(func() {
		scheduler = cron.New()
	})

	return scheduler
}
