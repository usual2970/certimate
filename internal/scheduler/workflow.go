package scheduler

import "context"

type WorkflowService interface {
	InitSchedule(ctx context.Context) error
}

func NewWorkflowScheduler(service WorkflowService) error {
	return service.InitSchedule(context.Background())
}
