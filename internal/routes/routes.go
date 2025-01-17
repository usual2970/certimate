package routes

import (
	"context"

	"github.com/usual2970/certimate/internal/notify"
	"github.com/usual2970/certimate/internal/repository"
	"github.com/usual2970/certimate/internal/rest"
	"github.com/usual2970/certimate/internal/statistics"
	"github.com/usual2970/certimate/internal/workflow"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
)

var (
	notifySvc     *notify.NotifyService
	workflowSvc   *workflow.WorkflowService
	statisticsSvc *statistics.StatisticsService
)

func Register(e *echo.Echo) {
	notifyRepo := repository.NewSettingsRepository()
	notifySvc = notify.NewNotifyService(notifyRepo)

	workflowRepo := repository.NewWorkflowRepository()
	workflowSvc = workflow.NewWorkflowService(workflowRepo)

	statisticsRepo := repository.NewStatisticsRepository()
	statisticsSvc = statistics.NewStatisticsService(statisticsRepo)

	group := e.Group("/api", apis.RequireAdminAuth())
	rest.NewWorkflowHandler(group, workflowSvc)
	rest.NewNotifyHandler(group, notifySvc)
	rest.NewStatisticsHandler(group, statisticsSvc)
}

func Unregister() {
	if workflowSvc != nil {
		workflowSvc.Stop(context.Background())
	}
}
