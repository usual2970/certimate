package ratpanelconsole_test

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"

	provider "github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/ratpanel-console"
)

var (
	fInputCertPath string
	fInputKeyPath  string
	fServerUrl     string
	fAccessTokenId int64
	fAccessToken   string
)

func init() {
	argsPrefix := "CERTIMATE_SSLDEPLOYER_RATPANELCONSOLE_"

	flag.StringVar(&fInputCertPath, argsPrefix+"INPUTCERTPATH", "", "")
	flag.StringVar(&fInputKeyPath, argsPrefix+"INPUTKEYPATH", "", "")
	flag.StringVar(&fServerUrl, argsPrefix+"SERVERURL", "", "")
	flag.Int64Var(&fAccessTokenId, argsPrefix+"ACCESSTOKENID", 0, "")
	flag.StringVar(&fAccessToken, argsPrefix+"ACCESSTOKEN", "", "")
}

/*
Shell command to run this test:

	go test -v ./ratpanel_console_test.go -args \
	--CERTIMATE_SSLDEPLOYER_RATPANELCONSOLE_INPUTCERTPATH="/path/to/your-input-cert.pem" \
	--CERTIMATE_SSLDEPLOYER_RATPANELCONSOLE_INPUTKEYPATH="/path/to/your-input-key.pem" \
	--CERTIMATE_SSLDEPLOYER_RATPANELCONSOLE_SERVERURL="http://127.0.0.1:8888" \
	--CERTIMATE_SSLDEPLOYER_RATPANELCONSOLE_ACCESSTOKENID="your-access-token-id" \
	--CERTIMATE_SSLDEPLOYER_RATPANELCONSOLE_ACCESSTOKEN="your-access-token"
*/
func TestDeploy(t *testing.T) {
	flag.Parse()

	t.Run("Deploy", func(t *testing.T) {
		t.Log(strings.Join([]string{
			"args:",
			fmt.Sprintf("INPUTCERTPATH: %v", fInputCertPath),
			fmt.Sprintf("INPUTKEYPATH: %v", fInputKeyPath),
			fmt.Sprintf("SERVERURL: %v", fServerUrl),
			fmt.Sprintf("ACCESSTOKENID: %v", fAccessTokenId),
			fmt.Sprintf("ACCESSTOKEN: %v", fAccessToken),
		}, "\n"))

		deployer, err := provider.NewSSLDeployerProvider(&provider.SSLDeployerProviderConfig{
			ServerUrl:                fServerUrl,
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
