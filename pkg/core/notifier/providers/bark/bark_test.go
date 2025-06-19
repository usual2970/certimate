package bark_test

import (
	"context"
	"flag"
	"fmt"
	"strings"
	"testing"

	provider "github.com/certimate-go/certimate/pkg/core/notifier/providers/bark"
)

const (
	mockSubject = "test_subject"
	mockMessage = "test_message"
)

var (
	fServerUrl string
	fDeviceKey string
)

func init() {
	argsPrefix := "CERTIMATE_NOTIFIER_BARK_"

	flag.StringVar(&fServerUrl, argsPrefix+"SERVERURL", "", "")
	flag.StringVar(&fDeviceKey, argsPrefix+"DEVICEKEY", "", "")
}

/*
Shell command to run this test:

	go test -v ./bark_test.go -args \
	--CERTIMATE_NOTIFIER_BARK_SERVERURL="https://example.com/your-server-url" \
	--CERTIMATE_NOTIFIER_BARK_DEVICEKEY="your-device-key"
*/
func TestNotify(t *testing.T) {
	flag.Parse()

	t.Run("Notify", func(t *testing.T) {
		t.Log(strings.Join([]string{
			"args:",
			fmt.Sprintf("SERVERURL: %v", fServerUrl),
			fmt.Sprintf("DEVICEKEY: %v", fDeviceKey),
		}, "\n"))

		notifier, err := provider.NewNotifierProvider(&provider.NotifierProviderConfig{
			ServerUrl: fServerUrl,
			DeviceKey: fDeviceKey,
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
