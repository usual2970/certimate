package pushover_test

import (
	"context"
	"flag"
	"fmt"
	"strings"
	"testing"

	provider "github.com/usual2970/certimate/internal/pkg/core/notifier/providers/pushover"
)

const (
	mockSubject = "test_subject"
	mockMessage = "test_message"
)

var (
	fToken string
	fUser  string
)

func init() {
	argsPrefix := "CERTIMATE_NOTIFIER_PUSHOVER_"
	flag.StringVar(&fToken, argsPrefix+"TOKEN", "", "")
	flag.StringVar(&fUser, argsPrefix+"USER", "", "")
}

/*
Shell command to run this test:

	go test -v ./pushover_test.go -args \
	--CERTIMATE_NOTIFIER_PUSHOVER_TOKEN="your-pushover-token" \
	--CERTIMATE_NOTIFIER_PUSHOVER_USER="your-pushover-user" \
*/
func TestNotify(t *testing.T) {
	flag.Parse()

	t.Run("Notify", func(t *testing.T) {
		t.Log(strings.Join([]string{
			"args:",
			fmt.Sprintf("TOKEN: %v", fToken),
		}, "\n"))

		notifier, err := provider.NewNotifier(&provider.NotifierConfig{
			Token: fToken,
			User:  fUser,
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
