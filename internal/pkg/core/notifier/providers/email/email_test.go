package email_test

import (
	"context"
	"os"
	"strconv"
	"testing"

	npEmail "github.com/usual2970/certimate/internal/pkg/core/notifier/providers/email"
)

const (
	MockSubject = "test_subject"
	MockMessage = "test_message"
)

/*
Shell command to run this test:

	CERTIMATE_NOTIFIER_EMAIL_SMTPHOST="smtp.example.com" \
	CERTIMATE_NOTIFIER_EMAIL_SMTPPORT=465 \
	CERTIMATE_NOTIFIER_EMAIL_SMTPTLS=true \
	CERTIMATE_NOTIFIER_EMAIL_USERNAME="your-username" \
	CERTIMATE_NOTIFIER_EMAIL_PASSWORD="your-password" \
	CERTIMATE_NOTIFIER_EMAIL_SENDERADDRESS="sender@example.com" \
	CERTIMATE_NOTIFIER_EMAIL_RECEIVERADDRESS="receiver@example.com" \
	go test -v -run TestNotify email_test.go
*/
func TestNotify(t *testing.T) {
	envPrefix := "CERTIMATE_NOTIFIER_EMAIL_"
	tSmtpHost := os.Getenv(envPrefix + "SMTPHOST")
	tSmtpPort, _ := strconv.ParseInt(os.Getenv(envPrefix+"SMTPPORT"), 10, 32)
	tSmtpTLS, _ := strconv.ParseBool(os.Getenv(envPrefix + "SMTPTLS"))
	tSmtpUsername := os.Getenv(envPrefix + "USERNAME")
	tSmtpPassword := os.Getenv(envPrefix + "PASSWORD")
	tSenderAddress := os.Getenv(envPrefix + "SENDERADDRESS")
	tReceiverAddress := os.Getenv(envPrefix + "RECEIVERADDRESS")

	notifier, err := npEmail.New(&npEmail.EmailNotifierConfig{
		SmtpHost:        tSmtpHost,
		SmtpPort:        int32(tSmtpPort),
		SmtpTLS:         tSmtpTLS,
		Username:        tSmtpUsername,
		Password:        tSmtpPassword,
		SenderAddress:   tSenderAddress,
		ReceiverAddress: tReceiverAddress,
	})
	if err != nil {
		t.Errorf("err: %+v", err)
		panic(err)
	}

	res, err := notifier.Notify(context.Background(), MockSubject, MockMessage)
	if err != nil {
		t.Errorf("err: %+v", err)
		panic(err)
	}

	t.Logf("ok: %v", res)
}
