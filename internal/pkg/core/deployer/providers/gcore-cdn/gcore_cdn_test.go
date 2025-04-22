package gcorecdn_test

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"

	provider "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/gcore-cdn"
)

var (
	fInputCertPath string
	fInputKeyPath  string
	fApiToken      string
	fResourceId    int64
)

func init() {
	argsPrefix := "CERTIMATE_DEPLOYER_GCORECDN_"

	flag.StringVar(&fInputCertPath, argsPrefix+"INPUTCERTPATH", "", "")
	flag.StringVar(&fInputKeyPath, argsPrefix+"INPUTKEYPATH", "", "")
	flag.StringVar(&fApiToken, argsPrefix+"APITOKEN", "", "")
	flag.Int64Var(&fResourceId, argsPrefix+"RESOURCEID", 0, "")
}

/*
Shell command to run this test:

	go test -v ./gcore_cdn_test.go -args \
	--CERTIMATE_DEPLOYER_GCORECDN_INPUTCERTPATH="/path/to/your-input-cert.pem" \
	--CERTIMATE_DEPLOYER_GCORECDN_INPUTKEYPATH="/path/to/your-input-key.pem" \
	--CERTIMATE_DEPLOYER_GCORECDN_APITOKEN="your-api-token" \
	--CERTIMATE_DEPLOYER_GCORECDN_RESOURCEID="your-cdn-resource-id"
*/
func TestDeploy(t *testing.T) {
	flag.Parse()

	t.Run("Deploy", func(t *testing.T) {
		t.Log(strings.Join([]string{
			"args:",
			fmt.Sprintf("INPUTCERTPATH: %v", fInputCertPath),
			fmt.Sprintf("INPUTKEYPATH: %v", fInputKeyPath),
			fmt.Sprintf("APITOKEN: %v", fApiToken),
			fmt.Sprintf("RESOURCEID: %v", fResourceId),
		}, "\n"))

		deployer, err := provider.NewDeployer(&provider.DeployerConfig{
			ApiToken:   fApiToken,
			ResourceId: fResourceId,
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
