package safeline_test

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"

	provider "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/safeline"
)

var (
	fInputCertPath string
	fInputKeyPath  string
	fApiUrl        string
	fApiToken      string
	fCertificateId int
)

func init() {
	argsPrefix := "CERTIMATE_DEPLOYER_SAFELINE_"

	flag.StringVar(&fInputCertPath, argsPrefix+"INPUTCERTPATH", "", "")
	flag.StringVar(&fInputKeyPath, argsPrefix+"INPUTKEYPATH", "", "")
	flag.StringVar(&fApiUrl, argsPrefix+"APIURL", "", "")
	flag.StringVar(&fApiToken, argsPrefix+"APITOKEN", "", "")
	flag.IntVar(&fCertificateId, argsPrefix+"CERTIFICATEID", 0, "")
}

/*
Shell command to run this test:

	go test -v ./safeline_test.go -args \
	--CERTIMATE_DEPLOYER_SAFELINE_INPUTCERTPATH="/path/to/your-input-cert.pem" \
	--CERTIMATE_DEPLOYER_SAFELINE_INPUTKEYPATH="/path/to/your-input-key.pem" \
	--CERTIMATE_DEPLOYER_SAFELINE_APIURL="http://127.0.0.1:9443" \
	--CERTIMATE_DEPLOYER_SAFELINE_APITOKEN="your-api-token" \
	--CERTIMATE_DEPLOYER_SAFELINE_CERTIFICATEID="your-cerficiate-id"
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
			fmt.Sprintf("CERTIFICATEID: %v", fCertificateId),
		}, "\n"))

		deployer, err := provider.NewDeployer(&provider.DeployerConfig{
			ApiUrl:                   fApiUrl,
			ApiToken:                 fApiToken,
			AllowInsecureConnections: true,
			ResourceType:             provider.ResourceType("certificate"),
			CertificateId:            int32(fCertificateId),
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
