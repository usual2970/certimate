package webhook_test

import (
	"context"
	"os"
	"testing"

	npWebhook "github.com/usual2970/certimate/internal/pkg/core/notifier/providers/webhook"
)

const (
	MockSubject = "test_subject"
	MockMessage = "test_message"
)

/*
Shell command to run this test:

	CERTIMATE_NOTIFIER_WEBHOOK_URL="https://example.com/your-webhook-url" \
	go test -v -run TestNotify webhook_test.go
*/
func TestNotify(t *testing.T) {
	envPrefix := "CERTIMATE_NOTIFIER_WEBHOOK_"
	tUrl := os.Getenv(envPrefix + "URL")

	notifier, err := npWebhook.New(&npWebhook.WebhookNotifierConfig{
		Url: tUrl,
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
