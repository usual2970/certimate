package scheduler

import "context"

type workflowService interface {
	InitSchedule(ctx context.Context) error
}

func InitWorkflowScheduler(service workflowService) error {
	return service.InitSchedule(context.Background())
}
