package gotify_test

import (
	"context"
	"flag"
	"fmt"
	"strings"
	"testing"

	provider "github.com/usual2970/certimate/internal/pkg/core/notifier/providers/gotify"
)

const (
	mockSubject = "test_subject"
	mockMessage = "test_message"
)

var (
	fUrl      string
	fToken    string
	fPriority int64
)

func init() {
	argsPrefix := "CERTIMATE_NOTIFIER_GOTIFY_"
	flag.StringVar(&fUrl, argsPrefix+"URL", "", "")
	flag.StringVar(&fToken, argsPrefix+"TOKEN", "", "")
	flag.Int64Var(&fPriority, argsPrefix+"PRIORITY", 0, "")
}

/*
Shell command to run this test:

	go test -v ./gotify_test.go -args \
	--CERTIMATE_NOTIFIER_GOTIFY_URL="https://example.com" \
	--CERTIMATE_NOTIFIER_GOTIFY_TOKEN="your-gotify-application-token" \
	--CERTIMATE_NOTIFIER_GOTIFY_PRIORITY="your-message-priority"
*/
func TestNotify(t *testing.T) {
	flag.Parse()

	t.Run("Notify", func(t *testing.T) {
		t.Log(strings.Join([]string{
			"args:",
			fmt.Sprintf("URL: %v", fUrl),
			fmt.Sprintf("TOKEN: %v", fToken),
			fmt.Sprintf("PRIORITY: %d", fPriority),
		}, "\n"))

		notifier, err := provider.NewNotifier(&provider.NotifierConfig{
			Url:      fUrl,
			Token:    fToken,
			Priority: fPriority,
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
