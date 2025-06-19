package cdnfly_test

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"

	provider "github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/cdnfly"
)

var (
	fInputCertPath string
	fInputKeyPath  string
	fServerUrl     string
	fApiKey        string
	fApiSecret     string
	fCertificateId string
)

func init() {
	argsPrefix := "CERTIMATE_SSLDEPLOYER_CDNFLY_"

	flag.StringVar(&fInputCertPath, argsPrefix+"INPUTCERTPATH", "", "")
	flag.StringVar(&fInputKeyPath, argsPrefix+"INPUTKEYPATH", "", "")
	flag.StringVar(&fServerUrl, argsPrefix+"SERVERURL", "", "")
	flag.StringVar(&fApiKey, argsPrefix+"APIKEY", "", "")
	flag.StringVar(&fApiSecret, argsPrefix+"APISECRET", "", "")
	flag.StringVar(&fCertificateId, argsPrefix+"CERTIFICATEID", "", "")
}

/*
Shell command to run this test:

	go test -v ./cdnfly_test.go -args \
	--CERTIMATE_SSLDEPLOYER_CDNFLY_INPUTCERTPATH="/path/to/your-input-cert.pem" \
	--CERTIMATE_SSLDEPLOYER_CDNFLY_INPUTKEYPATH="/path/to/your-input-key.pem" \
	--CERTIMATE_SSLDEPLOYER_CDNFLY_SERVERURL="http://127.0.0.1:88" \
	--CERTIMATE_SSLDEPLOYER_CDNFLY_APIKEY="your-api-key" \
	--CERTIMATE_SSLDEPLOYER_CDNFLY_APISECRET="your-api-secret" \
	--CERTIMATE_SSLDEPLOYER_CDNFLY_CERTIFICATEID="your-cert-id"
*/
func TestDeploy(t *testing.T) {
	flag.Parse()

	t.Run("Deploy", func(t *testing.T) {
		t.Log(strings.Join([]string{
			"args:",
			fmt.Sprintf("INPUTCERTPATH: %v", fInputCertPath),
			fmt.Sprintf("INPUTKEYPATH: %v", fInputKeyPath),
			fmt.Sprintf("SERVERURL: %v", fServerUrl),
			fmt.Sprintf("APIKEY: %v", fApiKey),
			fmt.Sprintf("APISECRET: %v", fApiSecret),
			fmt.Sprintf("CERTIFICATEID: %v", fCertificateId),
		}, "\n"))

		deployer, err := provider.NewSSLDeployerProvider(&provider.SSLDeployerProviderConfig{
			ServerUrl:                fServerUrl,
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
