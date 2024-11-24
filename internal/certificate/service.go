package certificate

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/notify"
	"github.com/usual2970/certimate/internal/repository"
	"github.com/usual2970/certimate/internal/utils/app"
)

const (
	defaultExpireSubject = "您有 {COUNT} 张证书即将过期"
	defaultExpireMessage = "有 {COUNT} 张证书即将过期，域名分别为 {DOMAINS}，请保持关注！"
)

type CertificateRepository interface {
	GetExpireSoon(ctx context.Context) ([]domain.Certificate, error)
}

type certificateService struct {
	repo CertificateRepository
}

func NewCertificateService(repo CertificateRepository) *certificateService {
	return &certificateService{
		repo: repo,
	}
}

func (s *certificateService) InitSchedule(ctx context.Context) error {
	scheduler := app.GetScheduler()

	err := scheduler.Add("certificate", "0 0 * * *", func() {
		certs, err := s.repo.GetExpireSoon(context.Background())
		if err != nil {
			app.GetApp().Logger().Error("failed to get expire soon certificate", "err", err)
			return
		}
		msg := buildMsg(certs)
		if err := notify.SendToAllChannels(msg.Subject, msg.Message); err != nil {
			app.GetApp().Logger().Error("failed to send expire soon certificate", "err", err)
		}
	})
	if err != nil {
		app.GetApp().Logger().Error("failed to add schedule", "err", err)
		return err
	}
	scheduler.Start()
	app.GetApp().Logger().Info("certificate schedule started")
	return nil
}

func buildMsg(records []domain.Certificate) *domain.NotifyMessage {
	if len(records) == 0 {
		return nil
	}

	// 查询模板信息
	settingRepo := repository.NewSettingRepository()
	setting, err := settingRepo.GetByName(context.Background(), "templates")

	subject := defaultExpireSubject
	message := defaultExpireMessage

	if err == nil {
		var templates *domain.NotifyTemplates

		json.Unmarshal([]byte(setting.Content), &templates)

		if templates != nil && len(templates.NotifyTemplates) > 0 {
			subject = templates.NotifyTemplates[0].Title
			message = templates.NotifyTemplates[0].Content
		}
	}

	// 替换变量
	count := len(records)
	domains := make([]string, count)

	for i, record := range records {
		domains[i] = record.SAN
	}

	countStr := strconv.Itoa(count)
	domainStr := strings.Join(domains, ";")

	subject = strings.ReplaceAll(subject, "{COUNT}", countStr)
	subject = strings.ReplaceAll(subject, "{DOMAINS}", domainStr)

	message = strings.ReplaceAll(message, "{COUNT}", countStr)
	message = strings.ReplaceAll(message, "{DOMAINS}", domainStr)

	// 返回消息
	return &domain.NotifyMessage{
		Subject: subject,
		Message: message,
	}
}
