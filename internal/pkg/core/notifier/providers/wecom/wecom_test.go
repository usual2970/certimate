﻿package serverchan_test

import (
	"context"
	"flag"
	"fmt"
	"strings"
	"testing"

	provider "github.com/usual2970/certimate/internal/pkg/core/notifier/providers/wecom"
)

const (
	mockSubject = "test_subject"
	mockMessage = "test_message"
)

var fWebhookUrl string

func init() {
	argsPrefix := "CERTIMATE_NOTIFIER_WECOM_"

	flag.StringVar(&fWebhookUrl, argsPrefix+"WEBHOOKURL", "", "")
}

/*
Shell command to run this test:

	go test -v ./wecom_test.go -args \
	--CERTIMATE_NOTIFIER_WECOM_WEBHOOKURL="https://example.com/your-webhook-url" \
*/
func TestNotify(t *testing.T) {
	flag.Parse()

	t.Run("Notify", func(t *testing.T) {
		t.Log(strings.Join([]string{
			"args:",
			fmt.Sprintf("WEBHOOKURL: %v", fWebhookUrl),
		}, "\n"))

		notifier, err := provider.New(&provider.WeComNotifierConfig{
			WebhookUrl: fWebhookUrl,
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
