package scheduler

import (
	"github.com/usual2970/certimate/internal/app"
	"github.com/usual2970/certimate/internal/certificate"
	"github.com/usual2970/certimate/internal/repository"
	"github.com/usual2970/certimate/internal/workflow"
)

func Register() {
	workflowRepo := repository.NewWorkflowRepository()
	workflowRunRepo := repository.NewWorkflowRunRepository()
	certificateRepo := repository.NewCertificateRepository()
	settingsRepo := repository.NewSettingsRepository()

	workflowSvc := workflow.NewWorkflowService(workflowRepo, workflowRunRepo, settingsRepo)
	certificateSvc := certificate.NewCertificateService(certificateRepo, settingsRepo)

	if err := InitWorkflowScheduler(workflowSvc); err != nil {
		app.GetLogger().Error("failed to init workflow scheduler", "err", err)
	}

	if err := InitCertificateScheduler(certificateSvc); err != nil {
		app.GetLogger().Error("failed to init certificate scheduler", "err", err)
	}
}
