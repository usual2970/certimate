package certificate

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/pocketbase/dbx"

	"github.com/usual2970/certimate/internal/app"
	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/domain/dtos"
	"github.com/usual2970/certimate/internal/notify"
	certutil "github.com/usual2970/certimate/internal/pkg/utils/cert"
	"github.com/usual2970/certimate/internal/repository"
)

const (
	defaultExpireSubject = "有 ${COUNT} 张证书即将过期"
	defaultExpireMessage = "有 ${COUNT} 张证书即将过期，域名分别为 ${DOMAINS}，请保持关注！"
)

type certificateRepository interface {
	ListExpireSoon(ctx context.Context) ([]*domain.Certificate, error)
	GetById(ctx context.Context, id string) (*domain.Certificate, error)
	DeleteWhere(ctx context.Context, exprs ...dbx.Expression) (int, error)
}

type settingsRepository interface {
	GetByName(ctx context.Context, name string) (*domain.Settings, error)
}

type CertificateService struct {
	certificateRepo certificateRepository
	settingsRepo    settingsRepository
}

func NewCertificateService(certificateRepo certificateRepository, settingsRepo settingsRepository) *CertificateService {
	return &CertificateService{
		certificateRepo: certificateRepo,
		settingsRepo:    settingsRepo,
	}
}

func (s *CertificateService) InitSchedule(ctx context.Context) error {
	// 每日发送过期证书提醒
	app.GetScheduler().MustAdd("certificateExpireSoonNotify", "0 0 * * *", func() {
		certificates, err := s.certificateRepo.ListExpireSoon(context.Background())
		if err != nil {
			app.GetLogger().Error("failed to get certificates which expire soon", "err", err)
			return
		}

		notification := buildExpireSoonNotification(certificates)
		if notification == nil {
			return
		}

		if err := notify.SendToAllChannels(notification.Subject, notification.Message); err != nil {
			app.GetLogger().Error("failed to send notification", "err", err)
		}
	})

	// 每日清理过期证书
	app.GetScheduler().MustAdd("certificateExpiredCleanup", "0 0 * * *", func() {
		settings, err := s.settingsRepo.GetByName(ctx, "persistence")
		if err != nil {
			app.GetLogger().Error("failed to get persistence settings", "err", err)
			return
		}

		var settingsContent *domain.PersistenceSettingsContent
		json.Unmarshal([]byte(settings.Content), &settingsContent)
		if settingsContent != nil && settingsContent.ExpiredCertificatesMaxDaysRetention != 0 {
			ret, err := s.certificateRepo.DeleteWhere(
				context.Background(),
				dbx.NewExp(fmt.Sprintf("expireAt<DATETIME('now', '-%d days')", settingsContent.ExpiredCertificatesMaxDaysRetention)),
			)
			if err != nil {
				app.GetLogger().Error("failed to delete expired certificates", "err", err)
			}

			if ret > 0 {
				app.GetLogger().Info(fmt.Sprintf("cleanup %d expired certificates", ret))
			}
		}
	})

	return nil
}

func (s *CertificateService) ArchiveFile(ctx context.Context, req *dtos.CertificateArchiveFileReq) (*dtos.CertificateArchiveFileResp, error) {
	certificate, err := s.certificateRepo.GetById(ctx, req.CertificateId)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)
	defer zipWriter.Close()

	resp := &dtos.CertificateArchiveFileResp{
		FileFormat: "zip",
	}

	switch strings.ToUpper(req.Format) {
	case "", "PEM":
		{
			certWriter, err := zipWriter.Create("certbundle.pem")
			if err != nil {
				return nil, err
			}

			_, err = certWriter.Write([]byte(certificate.Certificate))
			if err != nil {
				return nil, err
			}

			keyWriter, err := zipWriter.Create("privkey.pem")
			if err != nil {
				return nil, err
			}

			_, err = keyWriter.Write([]byte(certificate.PrivateKey))
			if err != nil {
				return nil, err
			}

			err = zipWriter.Close()
			if err != nil {
				return nil, err
			}

			resp.FileBytes = buf.Bytes()
			return resp, nil
		}

	case "PFX":
		{
			const pfxPassword = "certimate"

			certPFX, err := certutil.TransformCertificateFromPEMToPFX(certificate.Certificate, certificate.PrivateKey, pfxPassword)
			if err != nil {
				return nil, err
			}

			certWriter, err := zipWriter.Create("cert.pfx")
			if err != nil {
				return nil, err
			}

			_, err = certWriter.Write(certPFX)
			if err != nil {
				return nil, err
			}

			keyWriter, err := zipWriter.Create("pfx-password.txt")
			if err != nil {
				return nil, err
			}

			_, err = keyWriter.Write([]byte(pfxPassword))
			if err != nil {
				return nil, err
			}

			err = zipWriter.Close()
			if err != nil {
				return nil, err
			}

			resp.FileBytes = buf.Bytes()
			return resp, nil
		}

	case "JKS":
		{
			const jksPassword = "certimate"

			certJKS, err := certutil.TransformCertificateFromPEMToJKS(certificate.Certificate, certificate.PrivateKey, jksPassword, jksPassword, jksPassword)
			if err != nil {
				return nil, err
			}

			certWriter, err := zipWriter.Create("cert.jks")
			if err != nil {
				return nil, err
			}

			_, err = certWriter.Write(certJKS)
			if err != nil {
				return nil, err
			}

			keyWriter, err := zipWriter.Create("jks-password.txt")
			if err != nil {
				return nil, err
			}

			_, err = keyWriter.Write([]byte(jksPassword))
			if err != nil {
				return nil, err
			}

			err = zipWriter.Close()
			if err != nil {
				return nil, err
			}

			resp.FileBytes = buf.Bytes()
			return resp, nil
		}

	default:
		return nil, domain.ErrInvalidParams
	}
}

func (s *CertificateService) ValidateCertificate(ctx context.Context, req *dtos.CertificateValidateCertificateReq) (*dtos.CertificateValidateCertificateResp, error) {
	certX509, err := certutil.ParseCertificateFromPEM(req.Certificate)
	if err != nil {
		return nil, err
	} else if time.Now().After(certX509.NotAfter) {
		return nil, fmt.Errorf("certificate has expired at %s", certX509.NotAfter.UTC().Format(time.RFC3339))
	}

	return &dtos.CertificateValidateCertificateResp{
		IsValid: true,
		Domains: strings.Join(certX509.DNSNames, ";"),
	}, nil
}

func (s *CertificateService) ValidatePrivateKey(ctx context.Context, req *dtos.CertificateValidatePrivateKeyReq) (*dtos.CertificateValidatePrivateKeyResp, error) {
	_, err := certcrypto.ParsePEMPrivateKey([]byte(req.PrivateKey))
	if err != nil {
		return nil, err
	}

	return &dtos.CertificateValidatePrivateKeyResp{
		IsValid: true,
	}, nil
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
