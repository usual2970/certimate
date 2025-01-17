package app

import (
	"log/slog"
	"sync"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

var instance core.App

var intanceOnce sync.Once

func GetApp() core.App {
	intanceOnce.Do(func() {
		instance = pocketbase.NewWithConfig(pocketbase.Config{
			HideStartBanner: true,
		})
	})

	return instance
}

func GetDB() dbx.Builder {
	return GetApp().DB()
}

func GetLogger() *slog.Logger {
	return GetApp().Logger()
}
