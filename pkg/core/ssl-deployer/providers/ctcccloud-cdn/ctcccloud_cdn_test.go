package ctcccloudcdn_test

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"

	provider "github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/ctcccloud-cdn"
)

var (
	fInputCertPath   string
	fInputKeyPath    string
	fAccessKeyId     string
	fSecretAccessKey string
	fDomain          string
)

func init() {
	argsPrefix := "CERTIMATE_SSLDEPLOYER_CTCCCLOUDCDN_"

	flag.StringVar(&fInputCertPath, argsPrefix+"INPUTCERTPATH", "", "")
	flag.StringVar(&fInputKeyPath, argsPrefix+"INPUTKEYPATH", "", "")
	flag.StringVar(&fAccessKeyId, argsPrefix+"ACCESSKEYID", "", "")
	flag.StringVar(&fSecretAccessKey, argsPrefix+"SECRETACCESSKEY", "", "")
	flag.StringVar(&fDomain, argsPrefix+"DOMAIN", "", "")
}

/*
Shell command to run this test:

	go test -v ./ctcccloud_cdn_test.go -args \
	--CERTIMATE_SSLDEPLOYER_CTCCCLOUDCDN_INPUTCERTPATH="/path/to/your-input-cert.pem" \
	--CERTIMATE_SSLDEPLOYER_CTCCCLOUDCDN_INPUTKEYPATH="/path/to/your-input-key.pem" \
	--CERTIMATE_SSLDEPLOYER_CTCCCLOUDCDN_ACCESSKEYID="your-access-key-id" \
	--CERTIMATE_SSLDEPLOYER_CTCCCLOUDCDN_SECRETACCESSKEY="your-secret-access-key" \
	--CERTIMATE_SSLDEPLOYER_CTCCCLOUDCDN_DOMAIN="example.com"
*/
func TestDeploy(t *testing.T) {
	flag.Parse()

	t.Run("Deploy", func(t *testing.T) {
		t.Log(strings.Join([]string{
			"args:",
			fmt.Sprintf("INPUTCERTPATH: %v", fInputCertPath),
			fmt.Sprintf("INPUTKEYPATH: %v", fInputKeyPath),
			fmt.Sprintf("ACCESSKEYID: %v", fAccessKeyId),
			fmt.Sprintf("SECRETACCESSKEY: %v", fSecretAccessKey),
			fmt.Sprintf("DOMAIN: %v", fDomain),
		}, "\n"))

		deployer, err := provider.NewSSLDeployerProvider(&provider.SSLDeployerProviderConfig{
			AccessKeyId:     fAccessKeyId,
			SecretAccessKey: fSecretAccessKey,
			Domain:          fDomain,
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
