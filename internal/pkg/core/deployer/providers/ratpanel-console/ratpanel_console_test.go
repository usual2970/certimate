package ratpanelconsole_test

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"

	provider "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/ratpanel-console"
)

var (
	fInputCertPath string
	fInputKeyPath  string
	fApiUrl        string
	fTokenId       uint
	fToken         string
)

func init() {
	argsPrefix := "CERTIMATE_DEPLOYER_RATPANELCONSOLE_"

	flag.StringVar(&fInputCertPath, argsPrefix+"INPUTCERTPATH", "", "")
	flag.StringVar(&fInputKeyPath, argsPrefix+"INPUTKEYPATH", "", "")
	flag.StringVar(&fApiUrl, argsPrefix+"APIURL", "", "")
	flag.UintVar(&fTokenId, argsPrefix+"TOKENID", 0, "")
	flag.StringVar(&fToken, argsPrefix+"TOKEN", "", "")
}

/*
Shell command to run this test:

	go test -v ./ratpanel_console_test.go -args \
	--CERTIMATE_DEPLOYER_RATPANELCONSOLE_INPUTCERTPATH="/path/to/your-input-cert.pem" \
	--CERTIMATE_DEPLOYER_RATPANELCONSOLE_INPUTKEYPATH="/path/to/your-input-key.pem" \
	--CERTIMATE_DEPLOYER_RATPANELCONSOLE_APIURL="http://127.0.0.1:8888" \
	--CERTIMATE_DEPLOYER_RATPANELCONSOLE_TOKENID=your-access-token-id \
	--CERTIMATE_DEPLOYER_RATPANELCONSOLE_TOKEN="your-access-token"
*/
func TestDeploy(t *testing.T) {
	flag.Parse()

	t.Run("Deploy", func(t *testing.T) {
		t.Log(strings.Join([]string{
			"args:",
			fmt.Sprintf("INPUTCERTPATH: %v", fInputCertPath),
			fmt.Sprintf("INPUTKEYPATH: %v", fInputKeyPath),
			fmt.Sprintf("APIURL: %v", fApiUrl),
			fmt.Sprintf("TOKENID: %v", fTokenId),
			fmt.Sprintf("TOKEN: %v", fToken),
		}, "\n"))

		deployer, err := provider.NewDeployer(&provider.DeployerConfig{
			ApiUrl:                   fApiUrl,
			AccessTokenId:            fTokenId,
			AccessToken:              fToken,
			AllowInsecureConnections: true,
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
