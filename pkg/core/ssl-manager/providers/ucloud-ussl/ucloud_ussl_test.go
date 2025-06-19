package ucloudussl_test

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"

	provider "github.com/certimate-go/certimate/pkg/core/ssl-manager/providers/ucloud-ussl"
)

var (
	fInputCertPath string
	fInputKeyPath  string
	fPrivateKey    string
	fPublicKey     string
)

func init() {
	argsPrefix := "CERTIMATE_SSLMANAGER_UCLOUDUSSL_"

	flag.StringVar(&fInputCertPath, argsPrefix+"INPUTCERTPATH", "", "")
	flag.StringVar(&fInputKeyPath, argsPrefix+"INPUTKEYPATH", "", "")
	flag.StringVar(&fPrivateKey, argsPrefix+"PRIVATEKEY", "", "")
	flag.StringVar(&fPublicKey, argsPrefix+"PUBLICKEY", "", "")
}

/*
Shell command to run this test:

	go test -v ./ucloud_ussl_test.go -args \
	--CERTIMATE_SSLMANAGER_UCLOUDUSSL_INPUTCERTPATH="/path/to/your-input-cert.pem" \
	--CERTIMATE_SSLMANAGER_UCLOUDUSSL_INPUTKEYPATH="/path/to/your-input-key.pem" \
	--CERTIMATE_SSLMANAGER_UCLOUDUSSL_PRIVATEKEY="your-private-key" \
	--CERTIMATE_SSLMANAGER_UCLOUDUSSL_PUBLICKEY="your-public-key"
*/
func TestDeploy(t *testing.T) {
	flag.Parse()

	t.Run("Deploy", func(t *testing.T) {
		t.Log(strings.Join([]string{
			"args:",
			fmt.Sprintf("INPUTCERTPATH: %v", fInputCertPath),
			fmt.Sprintf("INPUTKEYPATH: %v", fInputKeyPath),
			fmt.Sprintf("PRIVATEKEY: %v", fPrivateKey),
			fmt.Sprintf("PUBLICKEY: %v", fPublicKey),
		}, "\n"))

		sslmanager, err := provider.NewSSLManagerProvider(&provider.SSLManagerProviderConfig{
			PrivateKey: fPrivateKey,
			PublicKey:  fPublicKey,
		})
		if err != nil {
			t.Errorf("err: %+v", err)
			return
		}

		fInputCertData, _ := os.ReadFile(fInputCertPath)
		fInputKeyData, _ := os.ReadFile(fInputKeyPath)
		res, err := sslmanager.Upload(context.Background(), string(fInputCertData), string(fInputKeyData))
		if err != nil {
			t.Errorf("err: %+v", err)
			return
		}

		sres, _ := json.Marshal(res)
		t.Logf("ok: %s", string(sres))
	})
}
