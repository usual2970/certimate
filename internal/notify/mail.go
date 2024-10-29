package notify

import (
	"context"
	"net/smtp"
)

type Mail struct {
	senderAddress     string
	smtpHostAddr      string
	smtpHostPort	  string
	smtpAuth          smtp.Auth
	receiverAddresses string
}

func NewMail(senderAddress, receiverAddresses, smtpHostAddr, smtpHostPort string) *Mail {
	if(smtpHostPort == "") {
		smtpHostPort = "25"
	}

	return &Mail{
		senderAddress:     senderAddress,
		smtpHostAddr:      smtpHostAddr,
		smtpHostPort:      smtpHostPort,
		receiverAddresses: receiverAddresses,
	}
}

func (m *Mail) SetAuth(username, password string) {
	m.smtpAuth = smtp.PlainAuth("", username, password, m.smtpHostAddr)
}

func (m *Mail) Send(ctx context.Context, subject, message string) error {
	// 构建邮件
    from := m.senderAddress
    to := []string{m.receiverAddresses}
    msg := []byte(
		"From: " + from + "\r\n" +
		"To: " + m.receiverAddresses + "\r\n" +
        "Subject: " + subject + "\r\n" +
        "\r\n" +
        message + "\r\n")
	
	var smtpAddress string
	// 组装邮箱服务器地址
	if(m.smtpHostPort == "25"){
		smtpAddress = m.smtpHostAddr
	}else{
		smtpAddress = m.smtpHostAddr + ":" + m.smtpHostPort
	}

    err := smtp.SendMail(smtpAddress, m.smtpAuth, from, to, msg)
    if err != nil {
        return err
    }

	return nil
}