package certificate

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/usual2970/certimate/internal/app"
	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/domain/dtos"
	"github.com/usual2970/certimate/internal/notify"
	"github.com/usual2970/certimate/internal/pkg/utils/certs"
	"github.com/usual2970/certimate/internal/repository"
)

const (
	defaultExpireSubject = "有 ${COUNT} 张证书即将过期"
	defaultExpireMessage = "有 ${COUNT} 张证书即将过期，域名分别为 ${DOMAINS}，请保持关注！"
)

type certificateRepository interface {
	ListExpireSoon(ctx context.Context) ([]*domain.Certificate, error)
	GetById(ctx context.Context, id string) (*domain.Certificate, error)
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
	app.GetScheduler().MustAdd("certificateExpireSoonNotify", "0 0 * * *", func() {
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
	return nil
}

func (s *CertificateService) ArchiveFile(ctx context.Context, req *dtos.CertificateArchiveFileReq) ([]byte, error) {
	certificate, err := s.repo.GetById(ctx, req.CertificateId)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)
	defer zipWriter.Close()

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

			return buf.Bytes(), nil
		}

	case "PFX":
		{
			const pfxPassword = "certimate"

			certPFX, err := certs.TransformCertificateFromPEMToPFX(certificate.Certificate, certificate.PrivateKey, pfxPassword)
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

			return buf.Bytes(), nil
		}

	case "JKS":
		{
			const jksPassword = "certimate"

			certJKS, err := certs.TransformCertificateFromPEMToJKS(certificate.Certificate, certificate.PrivateKey, jksPassword, jksPassword, jksPassword)
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

			return buf.Bytes(), nil
		}

	default:
		return nil, domain.ErrInvalidParams
	}
}

func (s *CertificateService) ValidateCertificate(ctx context.Context, req *dtos.CertificateValidateCertificateReq) (*dtos.CertificateValidateCertificateResp, error) {
	info, err := certs.ParseCertificateFromPEM(req.Certificate)
	if err != nil {
		return nil, err
	}

	if time.Now().After(info.NotAfter) {
		return nil, errors.New("证书已过期")
	}

	return &dtos.CertificateValidateCertificateResp{
		Domains: strings.Join(info.DNSNames, ";"),
	}, nil
}

func (s *CertificateService) ValidatePrivateKey(ctx context.Context, req *dtos.CertificateValidatePrivateKeyReq) error {
	_, err := certcrypto.ParsePEMPrivateKey([]byte(req.PrivateKey))
	return err
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
