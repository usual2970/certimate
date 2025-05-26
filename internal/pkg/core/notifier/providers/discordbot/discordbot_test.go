package discordbot_test

import (
	"context"
	"flag"
	"fmt"
	"strings"
	"testing"

	provider "github.com/usual2970/certimate/internal/pkg/core/notifier/providers/discordbot"
)

const (
	mockSubject = "test_subject"
	mockMessage = "test_message"
)

var (
	fApiToken  string
	fChannelId string
)

func init() {
	argsPrefix := "CERTIMATE_NOTIFIER_DISCORDBOT_"

	flag.StringVar(&fApiToken, argsPrefix+"APITOKEN", "", "")
	flag.StringVar(&fChannelId, argsPrefix+"CHANNELID", 0, "")
}

/*
Shell command to run this test:

	go test -v ./discordbot_test.go -args \
	--CERTIMATE_NOTIFIER_DISCORDBOT_APITOKEN="your-bot-token" \
	--CERTIMATE_NOTIFIER_DISCORDBOT_CHANNELID="your-channel-id"
*/
func TestNotify(t *testing.T) {
	flag.Parse()

	t.Run("Notify", func(t *testing.T) {
		t.Log(strings.Join([]string{
			"args:",
			fmt.Sprintf("APITOKEN: %v", fApiToken),
			fmt.Sprintf("CHANNELID: %v", fChannelId),
		}, "\n"))

		notifier, err := provider.NewNotifier(&provider.NotifierConfig{
			BotToken:  fApiToken,
			ChannelId: fChannelId,
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
