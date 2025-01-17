package certificate

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/usual2970/certimate/internal/app"
	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/notify"
	"github.com/usual2970/certimate/internal/repository"
)

const (
	defaultExpireSubject = "有 ${COUNT} 张证书即将过期"
	defaultExpireMessage = "有 ${COUNT} 张证书即将过期，域名分别为 ${DOMAINS}，请保持关注！"
)

type certificateRepository interface {
	ListExpireSoon(ctx context.Context) ([]*domain.Certificate, error)
}

type CertificateService struct {
	repo certificateRepository
}

func NewCertificateService(repo certificateRepository) *CertificateService {
	return &CertificateService{
		repo: repo,
	}
}

func (s *CertificateService) InitSchedule(ctx context.Context) error {
	scheduler := app.GetScheduler()
	err := scheduler.Add("certificate", "0 0 * * *", func() {
		certs, err := s.repo.ListExpireSoon(context.Background())
		if err != nil {
			app.GetLogger().Error("failed to get certificates which expire soon", "err", err)
			return
		}

		notification := buildExpireSoonNotification(certs)
		if notification == nil {
			return
		}

		if err := notify.SendToAllChannels(notification.Subject, notification.Message); err != nil {
			app.GetLogger().Error("failed to send notification", "err", err)
		}
	})
	if err != nil {
		app.GetLogger().Error("failed to add schedule", "err", err)
		return err
	}
	scheduler.Start()
	app.GetLogger().Info("certificate schedule started")
	return nil
}

func buildExpireSoonNotification(certificates []*domain.Certificate) *struct {
	Subject string
	Message string
} {
	if len(certificates) == 0 {
		return nil
	}

	subject := defaultExpireSubject
	message := defaultExpireMessage

	// 查询模板信息
	settingsRepo := repository.NewSettingsRepository()
	settings, err := settingsRepo.GetByName(context.Background(), "notifyTemplates")
	if err == nil {
		var templates *domain.NotifyTemplatesSettingsContent
		json.Unmarshal([]byte(settings.Content), &templates)

		if templates != nil && len(templates.NotifyTemplates) > 0 {
			subject = templates.NotifyTemplates[0].Subject
			message = templates.NotifyTemplates[0].Message
		}
	}

	// 替换变量
	count := len(certificates)
	domains := make([]string, count)
	for i, record := range certificates {
		domains[i] = record.SubjectAltNames
	}
	countStr := strconv.Itoa(count)
	domainStr := strings.Join(domains, ";")
	subject = strings.ReplaceAll(subject, "${COUNT}", countStr)
	subject = strings.ReplaceAll(subject, "${DOMAINS}", domainStr)
	message = strings.ReplaceAll(message, "${COUNT}", countStr)
	message = strings.ReplaceAll(message, "${DOMAINS}", domainStr)

	// 返回消息
	return &struct {
		Subject string
		Message string
	}{Subject: subject, Message: message}
}
