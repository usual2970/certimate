package routes

import (
	"context"

	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/router"

	"github.com/usual2970/certimate/internal/certificate"
	"github.com/usual2970/certimate/internal/notify"
	"github.com/usual2970/certimate/internal/repository"
	"github.com/usual2970/certimate/internal/rest/handlers"
	"github.com/usual2970/certimate/internal/statistics"
	"github.com/usual2970/certimate/internal/workflow"
)

var (
	certificateSvc *certificate.CertificateService
	workflowSvc    *workflow.WorkflowService
	statisticsSvc  *statistics.StatisticsService
	notifySvc      *notify.NotifyService
)

func Register(router *router.Router[*core.RequestEvent]) {
	workflowRepo := repository.NewWorkflowRepository()
	workflowRunRepo := repository.NewWorkflowRunRepository()
	certificateRepo := repository.NewCertificateRepository()
	settingsRepo := repository.NewSettingsRepository()
	statisticsRepo := repository.NewStatisticsRepository()

	certificateSvc = certificate.NewCertificateService(certificateRepo, settingsRepo)
	workflowSvc = workflow.NewWorkflowService(workflowRepo, workflowRunRepo, settingsRepo)
	statisticsSvc = statistics.NewStatisticsService(statisticsRepo)
	notifySvc = notify.NewNotifyService(settingsRepo)

	group := router.Group("/api")
	group.Bind(apis.RequireSuperuserAuth())
	handlers.NewCertificateHandler(group, certificateSvc)
	handlers.NewWorkflowHandler(group, workflowSvc)
	handlers.NewStatisticsHandler(group, statisticsSvc)
	handlers.NewNotifyHandler(group, notifySvc)
}

func Unregister() {
	if workflowSvc != nil {
		workflowSvc.Shutdown(context.Background())
	}
}
