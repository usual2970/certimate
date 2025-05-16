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
	fAccessTokenId int64
	fAccessToken   string
)

func init() {
	argsPrefix := "CERTIMATE_DEPLOYER_RATPANELCONSOLE_"

	flag.StringVar(&fInputCertPath, argsPrefix+"INPUTCERTPATH", "", "")
	flag.StringVar(&fInputKeyPath, argsPrefix+"INPUTKEYPATH", "", "")
	flag.StringVar(&fApiUrl, argsPrefix+"APIURL", "", "")
	flag.Int64Var(&fAccessTokenId, argsPrefix+"ACCESSTOKENID", 0, "")
	flag.StringVar(&fAccessToken, argsPrefix+"ACCESSTOKEN", "", "")
}

/*
Shell command to run this test:

	go test -v ./ratpanel_console_test.go -args \
	--CERTIMATE_DEPLOYER_RATPANELCONSOLE_INPUTCERTPATH="/path/to/your-input-cert.pem" \
	--CERTIMATE_DEPLOYER_RATPANELCONSOLE_INPUTKEYPATH="/path/to/your-input-key.pem" \
	--CERTIMATE_DEPLOYER_RATPANELCONSOLE_APIURL="http://127.0.0.1:8888" \
	--CERTIMATE_DEPLOYER_RATPANELCONSOLE_ACCESSTOKENID="your-access-token-id" \
	--CERTIMATE_DEPLOYER_RATPANELCONSOLE_ACCESSTOKEN="your-access-token"
*/
func TestDeploy(t *testing.T) {
	flag.Parse()

	t.Run("Deploy", func(t *testing.T) {
		t.Log(strings.Join([]string{
			"args:",
			fmt.Sprintf("INPUTCERTPATH: %v", fInputCertPath),
			fmt.Sprintf("INPUTKEYPATH: %v", fInputKeyPath),
			fmt.Sprintf("APIURL: %v", fApiUrl),
			fmt.Sprintf("ACCESSTOKENID: %v", fAccessTokenId),
			fmt.Sprintf("ACCESSTOKEN: %v", fAccessToken),
		}, "\n"))

		deployer, err := provider.NewDeployer(&provider.DeployerConfig{
			ApiUrl:                   fApiUrl,
			AccessTokenId:            int32(fAccessTokenId),
			AccessToken:              fAccessToken,
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
