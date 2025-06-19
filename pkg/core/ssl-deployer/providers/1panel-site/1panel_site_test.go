package onepanelsite_test

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"

	provider "github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/1panel-site"
)

var (
	fInputCertPath string
	fInputKeyPath  string
	fServerUrl     string
	fApiVersion    string
	fApiKey        string
	fWebsiteId     int64
)

func init() {
	argsPrefix := "CERTIMATE_SSLDEPLOYER_1PANELSITE_"

	flag.StringVar(&fInputCertPath, argsPrefix+"INPUTCERTPATH", "", "")
	flag.StringVar(&fInputKeyPath, argsPrefix+"INPUTKEYPATH", "", "")
	flag.StringVar(&fServerUrl, argsPrefix+"SERVERURL", "", "")
	flag.StringVar(&fApiVersion, argsPrefix+"APIVERSION", "v1", "")
	flag.StringVar(&fApiKey, argsPrefix+"APIKEY", "", "")
	flag.Int64Var(&fWebsiteId, argsPrefix+"WEBSITEID", 0, "")
}

/*
Shell command to run this test:

	go test -v ./1panel_site_test.go -args \
	--CERTIMATE_SSLDEPLOYER_1PANELSITE_INPUTCERTPATH="/path/to/your-input-cert.pem" \
	--CERTIMATE_SSLDEPLOYER_1PANELSITE_INPUTKEYPATH="/path/to/your-input-key.pem" \
	--CERTIMATE_SSLDEPLOYER_1PANELSITE_SERVERURL="http://127.0.0.1:20410" \
	--CERTIMATE_SSLDEPLOYER_1PANELSITE_APIVERSION="v1" \
	--CERTIMATE_SSLDEPLOYER_1PANELSITE_APIKEY="your-api-key" \
	--CERTIMATE_SSLDEPLOYER_1PANELSITE_WEBSITEID="your-website-id"
*/
func TestDeploy(t *testing.T) {
	flag.Parse()

	t.Run("Deploy", func(t *testing.T) {
		t.Log(strings.Join([]string{
			"args:",
			fmt.Sprintf("INPUTCERTPATH: %v", fInputCertPath),
			fmt.Sprintf("INPUTKEYPATH: %v", fInputKeyPath),
			fmt.Sprintf("SERVERURL: %v", fServerUrl),
			fmt.Sprintf("APIVERSION: %v", fApiVersion),
			fmt.Sprintf("APIKEY: %v", fApiKey),
			fmt.Sprintf("WEBSITEID: %v", fWebsiteId),
		}, "\n"))

		deployer, err := provider.NewSSLDeployerProvider(&provider.SSLDeployerProviderConfig{
			ServerUrl:                fServerUrl,
			ApiVersion:               fApiVersion,
			ApiKey:                   fApiKey,
			AllowInsecureConnections: true,
			ResourceType:             provider.RESOURCE_TYPE_WEBSITE,
			WebsiteId:                fWebsiteId,
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
