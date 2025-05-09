package proxmoxve_test

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"

	provider "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/proxmoxve"
)

var (
	fInputCertPath  string
	fInputKeyPath   string
	fApiUrl         string
	fApiToken       string
	fApiTokenSecret string
	fNodeName       string
)

func init() {
	argsPrefix := "CERTIMATE_DEPLOYER_PROXMOXVE_"

	flag.StringVar(&fInputCertPath, argsPrefix+"INPUTCERTPATH", "", "")
	flag.StringVar(&fInputKeyPath, argsPrefix+"INPUTKEYPATH", "", "")
	flag.StringVar(&fApiUrl, argsPrefix+"APIURL", "", "")
	flag.StringVar(&fApiToken, argsPrefix+"APITOKEN", "", "")
	flag.StringVar(&fApiTokenSecret, argsPrefix+"APITOKENSECRET", "", "")
	flag.StringVar(&fNodeName, argsPrefix+"NODENAME", "", "")
}

/*
Shell command to run this test:

	go test -v ./proxmoxve_test.go -args \
	--CERTIMATE_DEPLOYER_PROXMOXVE_INPUTCERTPATH="/path/to/your-input-cert.pem" \
	--CERTIMATE_DEPLOYER_PROXMOXVE_INPUTKEYPATH="/path/to/your-input-key.pem" \
	--CERTIMATE_DEPLOYER_PROXMOXVE_APIURL="http://127.0.0.1:8006" \
	--CERTIMATE_DEPLOYER_PROXMOXVE_APITOKEN="your-api-token" \
	--CERTIMATE_DEPLOYER_PROXMOXVE_APITOKENSECRET="your-api-token-secret" \
	--CERTIMATE_DEPLOYER_PROXMOXVE_NODENAME="your-cluster-node-name"
*/
func TestDeploy(t *testing.T) {
	flag.Parse()

	t.Run("Deploy", func(t *testing.T) {
		t.Log(strings.Join([]string{
			"args:",
			fmt.Sprintf("INPUTCERTPATH: %v", fInputCertPath),
			fmt.Sprintf("INPUTKEYPATH: %v", fInputKeyPath),
			fmt.Sprintf("APIURL: %v", fApiUrl),
			fmt.Sprintf("APITOKEN: %v", fApiToken),
			fmt.Sprintf("APITOKENSECRET: %v", fApiTokenSecret),
			fmt.Sprintf("NODENAME: %v", fNodeName),
		}, "\n"))

		deployer, err := provider.NewDeployer(&provider.DeployerConfig{
			ApiUrl:                   fApiUrl,
			ApiToken:                 fApiToken,
			ApiTokenSecret:           fApiTokenSecret,
			AllowInsecureConnections: true,
			NodeName:                 fNodeName,
			AutoRestart:              true,
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
