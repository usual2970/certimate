package email

import (
	"context"
	"errors"
	"fmt"
	"net/smtp"
	"os"

	"github.com/domodwyer/mailyak/v3"

	"github.com/usual2970/certimate/internal/pkg/core/notifier"
)

type EmailNotifierConfig struct {
	SmtpHost        string `json:"smtpHost"`
	SmtpPort        int32  `json:"smtpPort"`
	SmtpTLS         bool   `json:"smtpTLS"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	SenderAddress   string `json:"senderAddress"`
	ReceiverAddress string `json:"receiverAddress"`
}

type EmailNotifier struct {
	config *EmailNotifierConfig
}

var _ notifier.Notifier = (*EmailNotifier)(nil)

func New(config *EmailNotifierConfig) (*EmailNotifier, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}

	return &EmailNotifier{
		config: config,
	}, nil
}

func (n *EmailNotifier) Notify(ctx context.Context, subject string, message string) (res *notifier.NotifyResult, err error) {
	var smtpAuth smtp.Auth
	if n.config.Username != "" || n.config.Password != "" {
		smtpAuth = smtp.PlainAuth("", n.config.Username, n.config.Password, n.config.SmtpHost)
	}

	var yak *mailyak.MailYak
	if n.config.SmtpTLS {
		os.Setenv("GODEBUG", "tlsrsakex=1") // Fix for TLS handshake error
		yak, err = mailyak.NewWithTLS(fmt.Sprintf("%s:%d", n.config.SmtpHost, n.config.SmtpPort), smtpAuth, nil)
		if err != nil {
			return nil, err
		}
	} else {
		yak = mailyak.New(fmt.Sprintf("%s:%d", n.config.SmtpHost, n.config.SmtpPort), smtpAuth)
	}

	yak.From(n.config.SenderAddress)
	yak.To(n.config.ReceiverAddress)
	yak.Subject(subject)
	yak.Plain().Set(message)

	err = yak.Send()
	if err != nil {
		return nil, err
	}

	return &notifier.NotifyResult{}, nil
}
