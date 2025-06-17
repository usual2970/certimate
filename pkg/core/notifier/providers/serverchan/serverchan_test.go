package serverchan_test

import (
	"context"
	"flag"
	"fmt"
	"strings"
	"testing"

	provider "github.com/certimate-go/certimate/pkg/core/notifier/providers/serverchan"
)

const (
	mockSubject = "test_subject"
	mockMessage = "test_message"
)

var fUrl string

func init() {
	argsPrefix := "CERTIMATE_NOTIFIER_SERVERCHAN_"

	flag.StringVar(&fUrl, argsPrefix+"URL", "", "")
}

/*
Shell command to run this test:

	go test -v ./serverchan_test.go -args \
	--CERTIMATE_NOTIFIER_SERVERCHAN_URL="https://example.com/your-webhook-url" \
*/
func TestNotify(t *testing.T) {
	flag.Parse()

	t.Run("Notify", func(t *testing.T) {
		t.Log(strings.Join([]string{
			"args:",
			fmt.Sprintf("URL: %v", fUrl),
		}, "\n"))

		notifier, err := provider.NewNotifierProvider(&provider.NotifierProviderConfig{
			ServerUrl: fUrl,
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
