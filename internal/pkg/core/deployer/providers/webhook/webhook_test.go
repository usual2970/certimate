package webhook_test

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"

	dpWebhook "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/webhook"
)

var (
	fInputCertPath string
	fInputKeyPath  string
	fUrl           string
)

func init() {
	argsPrefix := "CERTIMATE_DEPLOYER_WEBHOOK_"

	flag.StringVar(&fInputCertPath, argsPrefix+"INPUTCERTPATH", "", "")
	flag.StringVar(&fInputKeyPath, argsPrefix+"INPUTKEYPATH", "", "")
	flag.StringVar(&fUrl, argsPrefix+"URL", "", "")
}

/*
Shell command to run this test:

	go test -v webhook_test.go -args \
	--CERTIMATE_DEPLOYER_WEBHOOK_INPUTCERTPATH="/path/to/your-input-cert.pem" \
	--CERTIMATE_DEPLOYER_WEBHOOK_INPUTKEYPATH="/path/to/your-input-key.pem" \
	--CERTIMATE_DEPLOYER_WEBHOOK_URL="https://example.com/your-webhook-url"
*/
func Test(t *testing.T) {
	flag.Parse()

	t.Run("Deploy", func(t *testing.T) {
		t.Log(strings.Join([]string{
			"args:",
			fmt.Sprintf("INPUTCERTPATH: %v", fInputCertPath),
			fmt.Sprintf("INPUTKEYPATH: %v", fInputKeyPath),
			fmt.Sprintf("URL: %v", fUrl),
		}, "\n"))

		deployer, err := dpWebhook.New(&dpWebhook.WebhookDeployerConfig{
			Url: fUrl,
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
