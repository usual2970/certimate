package baotapanelwaf_test

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"

	provider "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/baotawaf-site"
)

var (
	fInputCertPath string
	fInputKeyPath  string
	fServerUrl     string
	fApiKey        string
	fSiteName      string
	fSitePort      int64
)

func init() {
	argsPrefix := "CERTIMATE_DEPLOYER_BAOTAWAFSITE_"

	flag.StringVar(&fInputCertPath, argsPrefix+"INPUTCERTPATH", "", "")
	flag.StringVar(&fInputKeyPath, argsPrefix+"INPUTKEYPATH", "", "")
	flag.StringVar(&fServerUrl, argsPrefix+"SERVERURL", "", "")
	flag.StringVar(&fApiKey, argsPrefix+"APIKEY", "", "")
	flag.StringVar(&fSiteName, argsPrefix+"SITENAME", "", "")
	flag.Int64Var(&fSitePort, argsPrefix+"SITEPORT", 0, "")
}

/*
Shell command to run this test:

	go test -v ./baotawaf_site_test.go -args \
	--CERTIMATE_DEPLOYER_BAOTAWAFSITE_INPUTCERTPATH="/path/to/your-input-cert.pem" \
	--CERTIMATE_DEPLOYER_BAOTAWAFSITE_INPUTKEYPATH="/path/to/your-input-key.pem" \
	--CERTIMATE_DEPLOYER_BAOTAWAFSITE_SERVERURL="http://127.0.0.1:8888" \
	--CERTIMATE_DEPLOYER_BAOTAWAFSITE_APIKEY="your-api-key" \
	--CERTIMATE_DEPLOYER_BAOTAWAFSITE_SITENAME="your-site-name"\
	--CERTIMATE_DEPLOYER_BAOTAWAFSITE_SITEPORT=443
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
			fmt.Sprintf("SITENAME: %v", fSiteName),
			fmt.Sprintf("SITEPORT: %v", fSitePort),
		}, "\n"))

		deployer, err := provider.NewDeployer(&provider.DeployerConfig{
			ServerUrl:                fServerUrl,
			ApiKey:                   fApiKey,
			AllowInsecureConnections: true,
			SiteName:                 fSiteName,
			SitePort:                 int32(fSitePort),
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
