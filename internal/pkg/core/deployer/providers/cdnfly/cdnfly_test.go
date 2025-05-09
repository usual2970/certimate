package cdnfly_test

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"

	provider "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/cdnfly"
)

var (
	fInputCertPath string
	fInputKeyPath  string
	fApiUrl        string
	fApiKey        string
	fApiSecret     string
	fCertificateId string
)

func init() {
	argsPrefix := "CERTIMATE_DEPLOYER_CDNFLY_"

	flag.StringVar(&fInputCertPath, argsPrefix+"INPUTCERTPATH", "", "")
	flag.StringVar(&fInputKeyPath, argsPrefix+"INPUTKEYPATH", "", "")
	flag.StringVar(&fApiUrl, argsPrefix+"APIURL", "", "")
	flag.StringVar(&fApiKey, argsPrefix+"APIKEY", "", "")
	flag.StringVar(&fApiSecret, argsPrefix+"APISECRET", "", "")
	flag.StringVar(&fCertificateId, argsPrefix+"CERTIFICATEID", "", "")
}

/*
Shell command to run this test:

	go test -v ./cdnfly_test.go -args \
	--CERTIMATE_DEPLOYER_CDNFLY_INPUTCERTPATH="/path/to/your-input-cert.pem" \
	--CERTIMATE_DEPLOYER_CDNFLY_INPUTKEYPATH="/path/to/your-input-key.pem" \
	--CERTIMATE_DEPLOYER_CDNFLY_APIURL="http://127.0.0.1:88" \
	--CERTIMATE_DEPLOYER_CDNFLY_APIKEY="your-api-key" \
	--CERTIMATE_DEPLOYER_CDNFLY_APISECRET="your-api-secret" \
	--CERTIMATE_DEPLOYER_CDNFLY_CERTIFICATEID="your-cert-id"
*/
func TestDeploy(t *testing.T) {
	flag.Parse()

	t.Run("Deploy", func(t *testing.T) {
		t.Log(strings.Join([]string{
			"args:",
			fmt.Sprintf("INPUTCERTPATH: %v", fInputCertPath),
			fmt.Sprintf("INPUTKEYPATH: %v", fInputKeyPath),
			fmt.Sprintf("APIURL: %v", fApiUrl),
			fmt.Sprintf("APIKEY: %v", fApiKey),
			fmt.Sprintf("APISECRET: %v", fApiSecret),
			fmt.Sprintf("CERTIFICATEID: %v", fCertificateId),
		}, "\n"))

		deployer, err := provider.NewDeployer(&provider.DeployerConfig{
			ApiUrl:                   fApiUrl,
			ApiKey:                   fApiKey,
			ApiSecret:                fApiSecret,
			AllowInsecureConnections: true,
			ResourceType:             provider.RESOURCE_TYPE_CERTIFICATE,
			CertificateId:            fCertificateId,
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
