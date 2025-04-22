package ucloudus3_test

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"

	provider "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/ucloud-us3"
)

var (
	fInputCertPath string
	fInputKeyPath  string
	fPrivateKey    string
	fPublicKey     string
	fRegion        string
	fBucket        string
	fDomain        string
)

func init() {
	argsPrefix := "CERTIMATE_DEPLOYER_UCLOUDUS3_"

	flag.StringVar(&fInputCertPath, argsPrefix+"INPUTCERTPATH", "", "")
	flag.StringVar(&fInputKeyPath, argsPrefix+"INPUTKEYPATH", "", "")
	flag.StringVar(&fPrivateKey, argsPrefix+"PRIVATEKEY", "", "")
	flag.StringVar(&fPublicKey, argsPrefix+"PUBLICKEY", "", "")
	flag.StringVar(&fRegion, argsPrefix+"REGION", "", "")
	flag.StringVar(&fBucket, argsPrefix+"BUCKET", "", "")
	flag.StringVar(&fDomain, argsPrefix+"DOMAIN", "", "")
}

/*
Shell command to run this test:

	go test -v ./ucloud_us3_test.go -args \
	--CERTIMATE_DEPLOYER_UCLOUDUS3_INPUTCERTPATH="/path/to/your-input-cert.pem" \
	--CERTIMATE_DEPLOYER_UCLOUDUS3_INPUTKEYPATH="/path/to/your-input-key.pem" \
	--CERTIMATE_DEPLOYER_UCLOUDUS3_PRIVATEKEY="your-private-key" \
	--CERTIMATE_DEPLOYER_UCLOUDUS3_PUBLICKEY="your-public-key" \
	--CERTIMATE_DEPLOYER_UCLOUDUS3_REGION="cn-bj2" \
	--CERTIMATE_DEPLOYER_UCLOUDUS3_BUCKET="your-us3-bucket" \
	--CERTIMATE_DEPLOYER_UCLOUDUS3_DOMAIN="example.com"
*/
func TestDeploy(t *testing.T) {
	flag.Parse()

	t.Run("Deploy", func(t *testing.T) {
		t.Log(strings.Join([]string{
			"args:",
			fmt.Sprintf("INPUTCERTPATH: %v", fInputCertPath),
			fmt.Sprintf("INPUTKEYPATH: %v", fInputKeyPath),
			fmt.Sprintf("PRIVATEKEY: %v", fPrivateKey),
			fmt.Sprintf("PUBLICKEY: %v", fPublicKey),
			fmt.Sprintf("REGION: %v", fRegion),
			fmt.Sprintf("BUCKET: %v", fBucket),
			fmt.Sprintf("DOMAIN: %v", fDomain),
		}, "\n"))

		deployer, err := provider.NewDeployer(&provider.DeployerConfig{
			PrivateKey: fPrivateKey,
			PublicKey:  fPublicKey,
			Region:     fRegion,
			Bucket:     fBucket,
			Domain:     fDomain,
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
