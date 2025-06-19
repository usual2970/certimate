package dingtalkbot_test

import (
	"context"
	"flag"
	"fmt"
	"strings"
	"testing"

	provider "github.com/certimate-go/certimate/pkg/core/notifier/providers/dingtalkbot"
)

const (
	mockSubject = "test_subject"
	mockMessage = "test_message"
)

var (
	fWebhookUrl string
	fSecret     string
)

func init() {
	argsPrefix := "CERTIMATE_NOTIFIER_DINGTALKBOT_"

	flag.StringVar(&fWebhookUrl, argsPrefix+"WEBHOOKURL", "", "")
	flag.StringVar(&fSecret, argsPrefix+"SECRET", "", "")
}

/*
Shell command to run this test:

	go test -v ./dingtalkbot_test.go -args \
	--CERTIMATE_NOTIFIER_DINGTALKBOT_WEBHOOKURL="https://example.com/your-webhook-url" \
	--CERTIMATE_NOTIFIER_DINGTALKBOT_SECRET="your-secret"
*/
func TestNotify(t *testing.T) {
	flag.Parse()

	t.Run("Notify", func(t *testing.T) {
		t.Log(strings.Join([]string{
			"args:",
			fmt.Sprintf("WEBHOOKURL: %v", fWebhookUrl),
			fmt.Sprintf("SECRET: %v", fSecret),
		}, "\n"))

		notifier, err := provider.NewNotifierProvider(&provider.NotifierProviderConfig{
			WebhookUrl: fWebhookUrl,
			Secret:     fSecret,
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
