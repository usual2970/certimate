package baotapanelsite_test

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"

	provider "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/baotapanel-site"
)

var (
	fInputCertPath string
	fInputKeyPath  string
	fApiUrl        string
	fApiKey        string
	fSiteType      string
	fSiteName      string
)

func init() {
	argsPrefix := "CERTIMATE_DEPLOYER_BAOTAPANELSITE_"

	flag.StringVar(&fInputCertPath, argsPrefix+"INPUTCERTPATH", "", "")
	flag.StringVar(&fInputKeyPath, argsPrefix+"INPUTKEYPATH", "", "")
	flag.StringVar(&fApiUrl, argsPrefix+"APIURL", "", "")
	flag.StringVar(&fApiKey, argsPrefix+"APIKEY", "", "")
	flag.StringVar(&fSiteType, argsPrefix+"SITETYPE", "", "")
	flag.StringVar(&fSiteName, argsPrefix+"SITENAME", "", "")
}

/*
Shell command to run this test:

	go test -v ./baotapanel_site_test.go -args \
	--CERTIMATE_DEPLOYER_BAOTAPANELSITE_INPUTCERTPATH="/path/to/your-input-cert.pem" \
	--CERTIMATE_DEPLOYER_BAOTAPANELSITE_INPUTKEYPATH="/path/to/your-input-key.pem" \
	--CERTIMATE_DEPLOYER_BAOTAPANELSITE_APIURL="http://127.0.0.1:8888" \
	--CERTIMATE_DEPLOYER_BAOTAPANELSITE_APIKEY="your-api-key" \
	--CERTIMATE_DEPLOYER_BAOTAPANELSITE_SITETYPE="php" \
	--CERTIMATE_DEPLOYER_BAOTAPANELSITE_SITENAME="your-site-name"
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
			fmt.Sprintf("SITETYPE: %v", fSiteType),
			fmt.Sprintf("SITENAME: %v", fSiteName),
		}, "\n"))

		deployer, err := provider.NewDeployer(&provider.DeployerConfig{
			ApiUrl:                   fApiUrl,
			ApiKey:                   fApiKey,
			AllowInsecureConnections: true,
			SiteType:                 fSiteType,
			SiteName:                 fSiteName,
			SiteNames:                []string{fSiteName},
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
