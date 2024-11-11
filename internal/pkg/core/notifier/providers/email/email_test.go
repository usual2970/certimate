package email_test

import (
	"os"
	"strconv"
	"testing"

	notifierEmail "github.com/usual2970/certimate/internal/pkg/core/notifier/providers/email"
)

/*
Shell command to run this test:

	CERTIMATE_NOTIFIER_EMAIL_SMTPPORT=465 \
	CERTIMATE_NOTIFIER_EMAIL_SMTPTLS=true \
	CERTIMATE_NOTIFIER_EMAIL_SMTPHOST="smtp.example.com" \
	CERTIMATE_NOTIFIER_EMAIL_USERNAME="your-username" \
	CERTIMATE_NOTIFIER_EMAIL_PASSWORD="your-password" \
	CERTIMATE_NOTIFIER_EMAIL_SENDERADDRESS="sender@example.com" \
	CERTIMATE_NOTIFIER_EMAIL_RECEIVERADDRESS="receiver@example.com" \
	go test -v -run TestNotify email_test.go
*/
func TestNotify(t *testing.T) {
	smtpPort, err := strconv.ParseInt(os.Getenv("CERTIMATE_NOTIFIER_EMAIL_SMTPPORT"), 10, 32)
	if err != nil {
		t.Errorf("invalid envvar: %+v", err)
		panic(err)
	}

	smtpTLS, err := strconv.ParseBool(os.Getenv("CERTIMATE_NOTIFIER_EMAIL_SMTPTLS"))
	if err != nil {
		t.Errorf("invalid envvar: %+v", err)
		panic(err)
	}

	res, err := notifierEmail.New(&notifierEmail.EmailNotifierConfig{
		SmtpHost:        os.Getenv("CERTIMATE_NOTIFIER_EMAIL_SMTPHOST"),
		SmtpPort:        int32(smtpPort),
		SmtpTLS:         smtpTLS,
		Username:        os.Getenv("CERTIMATE_NOTIFIER_EMAIL_USERNAME"),
		Password:        os.Getenv("CERTIMATE_NOTIFIER_EMAIL_PASSWORD"),
		SenderAddress:   os.Getenv("CERTIMATE_NOTIFIER_EMAIL_SENDERADDRESS"),
		ReceiverAddress: os.Getenv("CERTIMATE_NOTIFIER_EMAIL_RECEIVERADDRESS"),
	})
	if err != nil {
		t.Errorf("invalid envvar: %+v", err)
		panic(err)
	}

	t.Logf("notify result: %v", res)
}
