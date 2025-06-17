package edgioapplications_test

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"

	provider "github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/edgio-applications"
)

var (
	fInputCertPath string
	fInputKeyPath  string
	fClientId      string
	fClientSecret  string
	fEnvironmentId string
)

func init() {
	argsPrefix := "CERTIMATE_SSLDEPLOYER_EDGIOAPPLICATIONS_"

	flag.StringVar(&fInputCertPath, argsPrefix+"INPUTCERTPATH", "", "")
	flag.StringVar(&fInputKeyPath, argsPrefix+"INPUTKEYPATH", "", "")
	flag.StringVar(&fClientId, argsPrefix+"CLIENTID", "", "")
	flag.StringVar(&fClientSecret, argsPrefix+"CLIENTSECRET", "", "")
	flag.StringVar(&fEnvironmentId, argsPrefix+"ENVIRONMENTID", "", "")
}

/*
Shell command to run this test:

	go test -v ./edgio_applications_test.go -args \
	--CERTIMATE_SSLDEPLOYER_EDGIOAPPLICATIONS_INPUTCERTPATH="/path/to/your-input-cert.pem" \
	--CERTIMATE_SSLDEPLOYER_EDGIOAPPLICATIONS_INPUTKEYPATH="/path/to/your-input-key.pem" \
	--CERTIMATE_SSLDEPLOYER_EDGIOAPPLICATIONS_CLIENTID="your-client-id" \
	--CERTIMATE_SSLDEPLOYER_EDGIOAPPLICATIONS_CLIENTSECRET="your-client-secret" \
	--CERTIMATE_SSLDEPLOYER_EDGIOAPPLICATIONS_ENVIRONMENTID="your-enviroment-id"
*/
func TestDeploy(t *testing.T) {
	flag.Parse()

	t.Run("Deploy", func(t *testing.T) {
		t.Log(strings.Join([]string{
			"args:",
			fmt.Sprintf("INPUTCERTPATH: %v", fInputCertPath),
			fmt.Sprintf("INPUTKEYPATH: %v", fInputKeyPath),
			fmt.Sprintf("CLIENTID: %v", fClientId),
			fmt.Sprintf("CLIENTSECRET: %v", fClientSecret),
			fmt.Sprintf("ENVIRONMENTID: %v", fEnvironmentId),
		}, "\n"))

		deployer, err := provider.NewSSLDeployerProvider(&provider.SSLDeployerProviderConfig{
			ClientId:      fClientId,
			ClientSecret:  fClientSecret,
			EnvironmentId: fEnvironmentId,
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
