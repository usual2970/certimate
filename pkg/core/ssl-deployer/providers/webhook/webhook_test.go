package webhook_test

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"

	provider "github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/webhook"
)

var (
	fInputCertPath      string
	fInputKeyPath       string
	fWebhookUrl         string
	fWebhookContentType string
	fWebhookData        string
)

func init() {
	argsPrefix := "CERTIMATE_SSLDEPLOYER_WEBHOOK_"

	flag.StringVar(&fInputCertPath, argsPrefix+"INPUTCERTPATH", "", "")
	flag.StringVar(&fInputKeyPath, argsPrefix+"INPUTKEYPATH", "", "")
	flag.StringVar(&fWebhookUrl, argsPrefix+"URL", "", "")
	flag.StringVar(&fWebhookContentType, argsPrefix+"CONTENTTYPE", "application/json", "")
	flag.StringVar(&fWebhookData, argsPrefix+"DATA", "", "")
}

/*
Shell command to run this test:

	go test -v ./webhook_test.go -args \
	--CERTIMATE_SSLDEPLOYER_WEBHOOK_INPUTCERTPATH="/path/to/your-input-cert.pem" \
	--CERTIMATE_SSLDEPLOYER_WEBHOOK_INPUTKEYPATH="/path/to/your-input-key.pem" \
	--CERTIMATE_SSLDEPLOYER_WEBHOOK_URL="https://example.com/your-webhook-url" \
	--CERTIMATE_SSLDEPLOYER_WEBHOOK_CONTENTTYPE="application/json" \
	--CERTIMATE_SSLDEPLOYER_WEBHOOK_DATA="{\"certificate\":\"${CERTIFICATE}\",\"privateKey\":\"${PRIVATE_KEY}\"}"
*/
func TestDeploy(t *testing.T) {
	flag.Parse()

	t.Run("Deploy", func(t *testing.T) {
		t.Log(strings.Join([]string{
			"args:",
			fmt.Sprintf("INPUTCERTPATH: %v", fInputCertPath),
			fmt.Sprintf("INPUTKEYPATH: %v", fInputKeyPath),
			fmt.Sprintf("WEBHOOKURL: %v", fWebhookUrl),
			fmt.Sprintf("WEBHOOKCONTENTTYPE: %v", fWebhookContentType),
			fmt.Sprintf("WEBHOOKDATA: %v", fWebhookData),
		}, "\n"))

		deployer, err := provider.NewSSLDeployerProvider(&provider.SSLDeployerProviderConfig{
			WebhookUrl:  fWebhookUrl,
			WebhookData: fWebhookData,
			Method:      "POST",
			Headers: map[string]string{
				"Content-Type": fWebhookContentType,
			},
			AllowInsecureConnections: true,
		})
		if err != nil {
			t.Errorf("err: %+v", err)
			return
		}

		fInputCertData, _ := os.ReadFile(fInputCertPath)
		fInputKeyData, _ := os.ReadFile(fInputKeyPath)
		res, err := deployer.Deploy(context.Background(), string(fInputCertData), string(fInputKeyData))
		if err != nil {
			t.Errorf("err: %+v", err)
			return
		}

		t.Logf("ok: %v", res)
	})
}
