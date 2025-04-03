package pushplus_test

import (
	"context"
	"flag"
	"fmt"
	"strings"
	"testing"

	provider "github.com/usual2970/certimate/internal/pkg/core/notifier/providers/pushplus"
)

const (
	mockSubject = "test_subject"
	mockMessage = "test_message"
)

var fToken string

func init() {
	argsPrefix := "CERTIMATE_NOTIFIER_PUSHPLUS_"
	flag.StringVar(&fToken, argsPrefix+"TOKEN", "", "")
}

/*
Shell command to run this test:

	go test -v ./pushplus_test.go -args \
	--CERTIMATE_NOTIFIER_PUSHPLUS_TOKEN="your-pushplus-token" \
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
