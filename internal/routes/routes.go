package routes

import (
	"github.com/usual2970/certimate/internal/notify"
	"github.com/usual2970/certimate/internal/repository"
	"github.com/usual2970/certimate/internal/rest"
	"github.com/usual2970/certimate/internal/workflow"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
)

func Register(e *echo.Echo) {
	group := e.Group("/api", apis.RequireAdminAuth())

	notifyRepo := repository.NewSettingRepository()
	notifySvc := notify.NewNotifyService(notifyRepo)

	workflowRepo := repository.NewWorkflowRepository()
	workflowSvc := workflow.NewWorkflowService(workflowRepo)

	rest.NewWorkflowHandler(group, workflowSvc)

	rest.NewNotifyHandler(group, notifySvc)
}
