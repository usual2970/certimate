package bunnycdn_test

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"

	provider "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/bunny-cdn"
)

var (
	fInputCertPath string
	fInputKeyPath  string
	fApiKey        string
	fPullZoneId    string
	fHostName      string
)

func init() {
	argsPrefix := "CERTIMATE_DEPLOYER_BUNNYCDN_"

	flag.StringVar(&fInputCertPath, argsPrefix+"INPUTCERTPATH", "", "")
	flag.StringVar(&fInputKeyPath, argsPrefix+"INPUTKEYPATH", "", "")
	flag.StringVar(&fApiKey, argsPrefix+"APIKEY", "", "")
	flag.StringVar(&fPullZoneId, argsPrefix+"PULLZONEID", "", "")
	flag.StringVar(&fHostName, argsPrefix+"HOSTNAME", "", "")
}

/*
Shell command to run this test:

	go test -v ./bunny_cdn_test.go -args \
	--CERTIMATE_DEPLOYER_BUNNYCDN_INPUTCERTPATH="/path/to/your-input-cert.pem" \
	--CERTIMATE_DEPLOYER_BUNNYCDN_INPUTKEYPATH="/path/to/your-input-key.pem" \
	--CERTIMATE_DEPLOYER_BUNNYCDN_APITOKEN="your-api-token" \
	--CERTIMATE_DEPLOYER_BUNNYCDN_PULLZONEID="your-pull-zone-id" \
	--CERTIMATE_DEPLOYER_BUNNYCDN_HOSTNAME="example.com"
*/
func TestDeploy(t *testing.T) {
	flag.Parse()

	t.Run("Deploy", func(t *testing.T) {
		t.Log(strings.Join([]string{
			"args:",
			fmt.Sprintf("INPUTCERTPATH: %v", fInputCertPath),
			fmt.Sprintf("INPUTKEYPATH: %v", fInputKeyPath),
			fmt.Sprintf("APIKEY: %v", fApiKey),
			fmt.Sprintf("PULLZONEID: %v", fPullZoneId),
			fmt.Sprintf("HOSTNAME: %v", fHostName),
		}, "\n"))

		deployer, err := provider.NewDeployer(&provider.DeployerConfig{
			ApiKey:     fApiKey,
			PullZoneId: fPullZoneId,
			Hostname:   fHostName,
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
