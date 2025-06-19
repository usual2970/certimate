package rainyunrcdn_test

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"

	provider "github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/rainyun-rcdn"
)

var (
	fInputCertPath string
	fInputKeyPath  string
	fApiKey        string
	fInstanceId    int64
	fDomain        string
)

func init() {
	argsPrefix := "CERTIMATE_SSLDEPLOYER_RAINYUNRCDN_"

	flag.StringVar(&fInputCertPath, argsPrefix+"INPUTCERTPATH", "", "")
	flag.StringVar(&fInputKeyPath, argsPrefix+"INPUTKEYPATH", "", "")
	flag.StringVar(&fApiKey, argsPrefix+"APIKEY", "", "")
	flag.Int64Var(&fInstanceId, argsPrefix+"INSTANCEID", 0, "")
	flag.StringVar(&fDomain, argsPrefix+"DOMAIN", "", "")
}

/*
Shell command to run this test:

	go test -v ./ucloud_ucdn_test.go -args \
	--CERTIMATE_SSLDEPLOYER_RAINYUNRCDN_INPUTCERTPATH="/path/to/your-input-cert.pem" \
	--CERTIMATE_SSLDEPLOYER_RAINYUNRCDN_INPUTKEYPATH="/path/to/your-input-key.pem" \
	--CERTIMATE_SSLDEPLOYER_RAINYUNRCDN_APIKEY="your-api-key" \
	--CERTIMATE_SSLDEPLOYER_RAINYUNRCDN_INSTANCEID="your-rcdn-instance-id" \
	--CERTIMATE_SSLDEPLOYER_RAINYUNRCDN_DOMAIN="example.com"
*/
func TestDeploy(t *testing.T) {
	flag.Parse()

	t.Run("Deploy", func(t *testing.T) {
		t.Log(strings.Join([]string{
			"args:",
			fmt.Sprintf("INPUTCERTPATH: %v", fInputCertPath),
			fmt.Sprintf("INPUTKEYPATH: %v", fInputKeyPath),
			fmt.Sprintf("APIKEY: %v", fApiKey),
			fmt.Sprintf("INSTANCEID: %v", fInstanceId),
			fmt.Sprintf("DOMAIN: %v", fDomain),
		}, "\n"))

		deployer, err := provider.NewSSLDeployerProvider(&provider.SSLDeployerProviderConfig{
			ApiKey:     fApiKey,
			InstanceId: int32(fInstanceId),
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
