package scheduler

import "context"

type CertificateService interface {
	InitSchedule(ctx context.Context) error
}

func NewCertificateScheduler(service CertificateService) error {
	return service.InitSchedule(context.Background())
}
