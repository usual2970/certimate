package webhook_test

import (
	"context"
	"os"
	"testing"

	dpWebhook "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/webhook"
)

/*
Shell command to run this test:

	CERTIMATE_DEPLOYER_WEBHOOK_INPUTCERTPATH="/path/to/your-input-cert.pem" \
	CERTIMATE_DEPLOYER_WEBHOOK_INPUTKEYPATH="/path/to/your-input-key.pem" \
	CERTIMATE_DEPLOYER_WEBHOOK_URL="https://example.com/your-webhook-url" \
	go test -v -run TestDeploy webhook_test.go
*/
func TestDeploy(t *testing.T) {
	envPrefix := "CERTIMATE_DEPLOYER_WEBHOOK_"
	tInputCertData, _ := os.ReadFile(os.Getenv(envPrefix + "INPUTCERTPATH"))
	tInputKeyData, _ := os.ReadFile(os.Getenv(envPrefix + "INPUTKEYPATH"))
	tUrl := os.Getenv(envPrefix + "URL")

	deployer, err := dpWebhook.New(&dpWebhook.WebhookDeployerConfig{
		Url: tUrl,
	})
	if err != nil {
		t.Errorf("err: %+v", err)
		panic(err)
	}

	res, err := deployer.Deploy(context.Background(), string(tInputCertData), string(tInputKeyData))
	if err != nil {
		t.Errorf("err: %+v", err)
		panic(err)
	}

	t.Logf("ok: %v", res)
}
