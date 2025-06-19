package email

import (
	"context"
	"crypto/tls"
	"errors"
	"log/slog"
	"net"
	"net/smtp"
	"strconv"

	"github.com/domodwyer/mailyak/v3"

	"github.com/certimate-go/certimate/pkg/core"
)

type NotifierProviderConfig struct {
	// SMTP 服务器地址。
	SmtpHost string `json:"smtpHost"`
	// SMTP 服务器端口。
	// 零值时根据是否启用 TLS 决定。
	SmtpPort int32 `json:"smtpPort"`
	// 是否启用 TLS。
	SmtpTls bool `json:"smtpTls"`
	// 用户名。
	Username string `json:"username"`
	// 密码。
	Password string `json:"password"`
	// 发件人邮箱。
	SenderAddress string `json:"senderAddress"`
	// 发件人显示名称。
	SenderName string `json:"senderName,omitempty"`
	// 收件人邮箱。
	ReceiverAddress string `json:"receiverAddress"`
}

type NotifierProvider struct {
	config *NotifierProviderConfig
	logger *slog.Logger
}

var _ core.Notifier = (*NotifierProvider)(nil)

func NewNotifierProvider(config *NotifierProviderConfig) (*NotifierProvider, error) {
	if config == nil {
		return nil, errors.New("the configuration of the notifier provider is nil")
	}

	return &NotifierProvider{
		config: config,
		logger: slog.Default(),
	}, nil
}

func (n *NotifierProvider) SetLogger(logger *slog.Logger) {
	if logger == nil {
		n.logger = slog.New(slog.DiscardHandler)
	} else {
		n.logger = logger
	}
}

func (n *NotifierProvider) Notify(ctx context.Context, subject string, message string) (*core.NotifyResult, error) {
	var smtpAuth smtp.Auth
	if n.config.Username != "" || n.config.Password != "" {
		smtpAuth = smtp.PlainAuth("", n.config.Username, n.config.Password, n.config.SmtpHost)
	}

	var smtpAddr string
	if n.config.SmtpPort == 0 {
		if n.config.SmtpTls {
			smtpAddr = net.JoinHostPort(n.config.SmtpHost, "465")
		} else {
			smtpAddr = net.JoinHostPort(n.config.SmtpHost, "25")
		}
	} else {
		smtpAddr = net.JoinHostPort(n.config.SmtpHost, strconv.Itoa(int(n.config.SmtpPort)))
	}

	var yak *mailyak.MailYak
	if n.config.SmtpTls {
		yakWithTls, err := mailyak.NewWithTLS(smtpAddr, smtpAuth, newTlsConfig())
		if err != nil {
			return nil, err
		}
		yak = yakWithTls
	} else {
		yak = mailyak.New(smtpAddr, smtpAuth)
	}

	yak.From(n.config.SenderAddress)
	yak.FromName(n.config.SenderName)
	yak.To(n.config.ReceiverAddress)
	yak.Subject(subject)
	yak.Plain().Set(message)

	if err := yak.Send(); err != nil {
		return nil, err
	}

	return &core.NotifyResult{}, nil
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
