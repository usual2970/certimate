package mattermost_test

import (
	"context"
	"flag"
	"fmt"
	"strings"
	"testing"

	provider "github.com/certimate-go/certimate/pkg/core/notifier/providers/mattermost"
)

const (
	mockSubject = "test_subject"
	mockMessage = "test_message"
)

var (
	fServerUrl string
	fChannelId string
	fUsername  string
	fPassword  string
)

func init() {
	argsPrefix := "CERTIMATE_NOTIFIER_MATTERMOST_"

	flag.StringVar(&fServerUrl, argsPrefix+"SERVERURL", "", "")
	flag.StringVar(&fChannelId, argsPrefix+"CHANNELID", "", "")
	flag.StringVar(&fUsername, argsPrefix+"USERNAME", "", "")
	flag.StringVar(&fPassword, argsPrefix+"PASSWORD", "", "")
}

/*
Shell command to run this test:

	go test -v ./mattermost_test.go -args \
	--CERTIMATE_NOTIFIER_MATTERMOST_SERVERURL="https://example.com/your-server-url" \
	--CERTIMATE_NOTIFIER_MATTERMOST_CHANNELID="your-chanel-id" \
	--CERTIMATE_NOTIFIER_MATTERMOST_USERNAME="your-username" \
	--CERTIMATE_NOTIFIER_MATTERMOST_PASSWORD="your-password"
*/
func TestNotify(t *testing.T) {
	flag.Parse()

	t.Run("Notify", func(t *testing.T) {
		t.Log(strings.Join([]string{
			"args:",
			fmt.Sprintf("SERVERURL: %v", fServerUrl),
			fmt.Sprintf("CHANNELID: %v", fChannelId),
			fmt.Sprintf("USERNAME: %v", fUsername),
			fmt.Sprintf("PASSWORD: %v", fPassword),
		}, "\n"))

		notifier, err := provider.NewNotifierProvider(&provider.NotifierProviderConfig{
			ServerUrl: fServerUrl,
			ChannelId: fChannelId,
			Username:  fUsername,
			Password:  fPassword,
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
