package email

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/smtp"

	"github.com/domodwyer/mailyak/v3"

	"github.com/usual2970/certimate/internal/pkg/core/notifier"
)

type EmailNotifierConfig struct {
	// SMTP 服务器地址。
	SmtpHost string `json:"smtpHost"`
	// SMTP 服务器端口。
	// 零值时根据是否启用 TLS 决定。
	SmtpPort int32 `json:"smtpPort"`
	// 是否启用 TLS。
	SmtpTLS bool `json:"smtpTLS"`
	// 用户名。
	Username string `json:"username"`
	// 密码。
	Password string `json:"password"`
	// 发件人邮箱。
	SenderAddress string `json:"senderAddress"`
	// 收件人邮箱。
	ReceiverAddress string `json:"receiverAddress"`
}

type EmailNotifier struct {
	config *EmailNotifierConfig
}

var _ notifier.Notifier = (*EmailNotifier)(nil)

func New(config *EmailNotifierConfig) (*EmailNotifier, error) {
	if config == nil {
		panic("config is nil")
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

	var smtpAddr string
	if n.config.SmtpPort == 0 {
		if n.config.SmtpTLS {
			smtpAddr = fmt.Sprintf("%s:465", n.config.SmtpHost)
		} else {
			smtpAddr = fmt.Sprintf("%s:25", n.config.SmtpHost)
		}
	} else {
		smtpAddr = fmt.Sprintf("%s:%d", n.config.SmtpHost, n.config.SmtpPort)
	}

	var yak *mailyak.MailYak
	if n.config.SmtpTLS {
		yak, err = mailyak.NewWithTLS(smtpAddr, smtpAuth, newTlsConfig())
		if err != nil {
			return nil, err
		}
	} else {
		yak = mailyak.New(smtpAddr, smtpAuth)
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

func newTlsConfig() *tls.Config {
	var suiteIds []uint16
	for _, suite := range tls.CipherSuites() {
		suiteIds = append(suiteIds, suite.ID)
	}
	for _, suite := range tls.InsecureCipherSuites() {
		suiteIds = append(suiteIds, suite.ID)
	}

	// 为兼容国内部分低版本 TLS 的 SMTP 服务商
	return &tls.Config{
		MinVersion:   tls.VersionTLS10,
		CipherSuites: suiteIds,
	}
}
