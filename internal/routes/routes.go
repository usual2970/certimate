package routes

import (
	"github.com/usual2970/certimate/internal/notify"
	"github.com/usual2970/certimate/internal/repository"
	"github.com/usual2970/certimate/internal/rest"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
)

func Register(e *echo.Echo) {
	notifyRepo := repository.NewSettingRepository()
	notifySvc := notify.NewNotifyService(notifyRepo)

	group := e.Group("/api", apis.RequireAdminAuth())

	rest.NewNotifyHandler(group, notifySvc)
}
