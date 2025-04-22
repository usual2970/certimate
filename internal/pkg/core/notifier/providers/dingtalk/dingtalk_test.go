package dingtalk_test

import (
	"context"
	"flag"
	"fmt"
	"strings"
	"testing"

	provider "github.com/usual2970/certimate/internal/pkg/core/notifier/providers/dingtalk"
)

const (
	mockSubject = "test_subject"
	mockMessage = "test_message"
)

var (
	fAccessToken string
	fSecret      string
)

func init() {
	argsPrefix := "CERTIMATE_NOTIFIER_DINGTALK_"

	flag.StringVar(&fAccessToken, argsPrefix+"ACCESSTOKEN", "", "")
	flag.StringVar(&fSecret, argsPrefix+"SECRET", "", "")
}

/*
Shell command to run this test:

	go test -v ./dingtalk_test.go -args \
	--CERTIMATE_NOTIFIER_DINGTALK_URL="https://example.com/your-webhook-url"
*/
func TestNotify(t *testing.T) {
	flag.Parse()

	t.Run("Notify", func(t *testing.T) {
		t.Log(strings.Join([]string{
			"args:",
			fmt.Sprintf("ACCESSTOKEN: %v", fAccessToken),
			fmt.Sprintf("SECRET: %v", fSecret),
		}, "\n"))

		notifier, err := provider.NewNotifier(&provider.NotifierConfig{
			AccessToken: fAccessToken,
			Secret:      fSecret,
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
