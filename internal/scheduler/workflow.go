package scheduler

import "context"

type workflowService interface {
	InitSchedule(ctx context.Context) error
}

func NewWorkflowScheduler(service workflowService) error {
	return service.InitSchedule(context.Background())
}
