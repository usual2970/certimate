package scheduler

import (
	"github.com/usual2970/certimate/internal/certificate"
	"github.com/usual2970/certimate/internal/repository"
	"github.com/usual2970/certimate/internal/workflow"
)

func Register() {
	workflowRepo := repository.NewWorkflowRepository()
	workflowSvc := workflow.NewWorkflowService(workflowRepo)

	certificateRepo := repository.NewCertificateRepository()
	certificateSvc := certificate.NewCertificateService(certificateRepo)

	NewCertificateScheduler(certificateSvc)

	NewWorkflowScheduler(workflowSvc)
}
