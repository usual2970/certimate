package email_test

import (
	"context"
	"flag"
	"fmt"
	"strings"
	"testing"

	provider "github.com/usual2970/certimate/internal/pkg/core/notifier/providers/email"
)

const (
	mockSubject = "test_subject"
	mockMessage = "test_message"
)

var (
	fSmtpHost        string
	fSmtpPort        int64
	fSmtpTLS         bool
	fUsername        string
	fPassword        string
	fSenderAddress   string
	fReceiverAddress string
)

func init() {
	argsPrefix := "CERTIMATE_NOTIFIER_EMAIL_"

	flag.StringVar(&fSmtpHost, argsPrefix+"SMTPHOST", "", "")
	flag.Int64Var(&fSmtpPort, argsPrefix+"SMTPPORT", 0, "")
	flag.BoolVar(&fSmtpTLS, argsPrefix+"SMTPTLS", false, "")
	flag.StringVar(&fUsername, argsPrefix+"USERNAME", "", "")
	flag.StringVar(&fPassword, argsPrefix+"PASSWORD", "", "")
	flag.StringVar(&fSenderAddress, argsPrefix+"SENDERADDRESS", "", "")
	flag.StringVar(&fReceiverAddress, argsPrefix+"RECEIVERADDRESS", "", "")
}

/*
Shell command to run this test:

	go test -v ./email_test.go -args \
	--CERTIMATE_NOTIFIER_EMAIL_SMTPHOST="smtp.example.com" \
	--CERTIMATE_NOTIFIER_EMAIL_SMTPPORT=465 \
	--CERTIMATE_NOTIFIER_EMAIL_SMTPTLS=true \
	--CERTIMATE_NOTIFIER_EMAIL_USERNAME="your-username" \
	--CERTIMATE_NOTIFIER_EMAIL_PASSWORD="your-password" \
	--CERTIMATE_NOTIFIER_EMAIL_SENDERADDRESS="sender@example.com" \
	--CERTIMATE_NOTIFIER_EMAIL_RECEIVERADDRESS="receiver@example.com"
*/
func TestNotify(t *testing.T) {
	flag.Parse()

	t.Run("Notify", func(t *testing.T) {
		t.Log(strings.Join([]string{
			"args:",
			fmt.Sprintf("SMTPHOST: %v", fSmtpHost),
			fmt.Sprintf("SMTPPORT: %v", fSmtpPort),
			fmt.Sprintf("SMTPTLS: %v", fSmtpTLS),
			fmt.Sprintf("USERNAME: %v", fUsername),
			fmt.Sprintf("PASSWORD: %v", fPassword),
			fmt.Sprintf("SENDERADDRESS: %v", fSenderAddress),
			fmt.Sprintf("RECEIVERADDRESS: %v", fReceiverAddress),
		}, "\n"))

		notifier, err := provider.NewNotifier(&provider.NotifierConfig{
			SmtpHost:        fSmtpHost,
			SmtpPort:        int32(fSmtpPort),
			SmtpTls:         fSmtpTLS,
			Username:        fUsername,
			Password:        fPassword,
			SenderAddress:   fSenderAddress,
			ReceiverAddress: fReceiverAddress,
		})
		if err != nil {
			t.Errorf("err: %+v", err)
			return
		}

		res, err := notifier.Notify(context.Background(), mockSubject, mockMessage)
		if err != nil {
			t.Errorf("err: %+v", err)
			return
		}

		t.Logf("ok: %v", res)
	})
}
