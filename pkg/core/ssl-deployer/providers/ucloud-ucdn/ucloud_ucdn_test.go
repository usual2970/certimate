package uclouducdn_test

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"

	provider "github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/ucloud-ucdn"
)

var (
	fInputCertPath string
	fInputKeyPath  string
	fPrivateKey    string
	fPublicKey     string
	fDomainId      string
)

func init() {
	argsPrefix := "CERTIMATE_SSLDEPLOYER_UCLOUDUCDN_"

	flag.StringVar(&fInputCertPath, argsPrefix+"INPUTCERTPATH", "", "")
	flag.StringVar(&fInputKeyPath, argsPrefix+"INPUTKEYPATH", "", "")
	flag.StringVar(&fPrivateKey, argsPrefix+"PRIVATEKEY", "", "")
	flag.StringVar(&fPublicKey, argsPrefix+"PUBLICKEY", "", "")
	flag.StringVar(&fDomainId, argsPrefix+"DOMAINID", "", "")
}

/*
Shell command to run this test:

	go test -v ./ucloud_ucdn_test.go -args \
	--CERTIMATE_SSLDEPLOYER_UCLOUDUCDN_INPUTCERTPATH="/path/to/your-input-cert.pem" \
	--CERTIMATE_SSLDEPLOYER_UCLOUDUCDN_INPUTKEYPATH="/path/to/your-input-key.pem" \
	--CERTIMATE_SSLDEPLOYER_UCLOUDUCDN_PRIVATEKEY="your-private-key" \
	--CERTIMATE_SSLDEPLOYER_UCLOUDUCDN_PUBLICKEY="your-public-key" \
	--CERTIMATE_SSLDEPLOYER_UCLOUDUCDN_DOMAINID="your-domain-id"
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
			fmt.Sprintf("DOMAIN: %v", fDomainId),
		}, "\n"))

		deployer, err := provider.NewSSLDeployerProvider(&provider.SSLDeployerProviderConfig{
			PrivateKey: fPrivateKey,
			PublicKey:  fPublicKey,
			DomainId:   fDomainId,
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
