package routes

import (
	"context"

	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/router"

	"github.com/usual2970/certimate/internal/notify"
	"github.com/usual2970/certimate/internal/repository"
	"github.com/usual2970/certimate/internal/rest/handlers"
	"github.com/usual2970/certimate/internal/statistics"
	"github.com/usual2970/certimate/internal/workflow"
)

var (
	notifySvc     *notify.NotifyService
	workflowSvc   *workflow.WorkflowService
	statisticsSvc *statistics.StatisticsService
)

func Register(router *router.Router[*core.RequestEvent]) {
	notifyRepo := repository.NewSettingsRepository()
	notifySvc = notify.NewNotifyService(notifyRepo)

	workflowRepo := repository.NewWorkflowRepository()
	workflowSvc = workflow.NewWorkflowService(workflowRepo)

	statisticsRepo := repository.NewStatisticsRepository()
	statisticsSvc = statistics.NewStatisticsService(statisticsRepo)

	group := router.Group("/api")
	group.Bind(apis.RequireSuperuserAuth())
	handlers.NewWorkflowHandler(group, workflowSvc)
	handlers.NewNotifyHandler(group, notifySvc)
	handlers.NewStatisticsHandler(group, statisticsSvc)
}

func Unregister() {
	if workflowSvc != nil {
		workflowSvc.Stop(context.Background())
	}
}
