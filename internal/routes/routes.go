package routes

import (
	"certimate/internal/notify"
	"certimate/internal/repository"
	"certimate/internal/rest"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
)

func Register(e *echo.Echo) {
	notifyRepo := repository.NewSettingRepository()
	notifySvc := notify.NewNotifyService(notifyRepo)

	group := e.Group("/api", apis.RequireAdminAuth())

	rest.NewNotifyHandler(group, notifySvc)
}
