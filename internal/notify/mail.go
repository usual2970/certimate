package notify

import (
	"context"
	"fmt"
	"net/mail"
	"strconv"

	"github.com/pocketbase/pocketbase/tools/mailer"
)

const defaultSmtpHostPort = "25"

type Mail struct {
	username string
	to       string
	client   *mailer.SmtpClient
}

func NewMail(senderAddress, receiverAddresses, smtpHostAddr, smtpHostPort, password string) (*Mail, error) {
	if smtpHostPort == "" {
		smtpHostPort = defaultSmtpHostPort
	}

	port, err := strconv.Atoi(smtpHostPort)
	if err != nil {
		return nil, fmt.Errorf("invalid smtp port: %w", err)
	}

	client := mailer.SmtpClient{
		Host:     smtpHostAddr,
		Port:     port,
		Username: senderAddress,
		Password: password,
		Tls:      true,
	}

	return &Mail{
		username: senderAddress,
		client:   &client,
		to:       receiverAddresses,
	}, nil
}

func (m *Mail) Send(ctx context.Context, subject, content string) error {
	message := &mailer.Message{
		From: mail.Address{
			Address: m.username,
		},
		To:      []mail.Address{{Address: m.to}},
		Subject: subject,
		Text:    content,
	}

	return m.client.Send(message)
}
