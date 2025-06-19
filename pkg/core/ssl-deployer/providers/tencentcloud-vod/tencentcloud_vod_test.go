package tencentcloudvod_test

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"

	provider "github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/tencentcloud-vod"
)

var (
	fInputCertPath string
	fInputKeyPath  string
	fSecretId      string
	fSecretKey     string
	fDomain        string
	fSubAppId      int64
	fInstanceId    string
)

func init() {
	argsPrefix := "CERTIMATE_SSLDEPLOYER_TENCENTCLOUDVOD_"

	flag.StringVar(&fInputCertPath, argsPrefix+"INPUTCERTPATH", "", "")
	flag.StringVar(&fInputKeyPath, argsPrefix+"INPUTKEYPATH", "", "")
	flag.StringVar(&fSecretId, argsPrefix+"SECRETID", "", "")
	flag.StringVar(&fSecretKey, argsPrefix+"SECRETKEY", "", "")
	flag.StringVar(&fDomain, argsPrefix+"DOMAIN", "", "")
	flag.Int64Var(&fSubAppId, argsPrefix+"SUBAPPID", 0, "")
	flag.StringVar(&fInstanceId, argsPrefix+"INSTANCEID", "", "")
}

/*
Shell command to run this test:

	go test -v ./tencentcloud_vod_test.go -args \
	--CERTIMATE_SSLDEPLOYER_TENCENTCLOUDVOD_INPUTCERTPATH="/path/to/your-input-cert.pem" \
	--CERTIMATE_SSLDEPLOYER_TENCENTCLOUDVOD_INPUTKEYPATH="/path/to/your-input-key.pem" \
	--CERTIMATE_SSLDEPLOYER_TENCENTCLOUDVOD_SECRETID="your-secret-id" \
	--CERTIMATE_SSLDEPLOYER_TENCENTCLOUDVOD_SECRETKEY="your-secret-key" \
	--CERTIMATE_SSLDEPLOYER_TENCENTCLOUDVOD_SUBAPPID="your-app-id" \
	--CERTIMATE_SSLDEPLOYER_TENCENTCLOUDVOD_DOMAIN="example.com"
*/
func TestDeploy(t *testing.T) {
	flag.Parse()

	t.Run("Deploy", func(t *testing.T) {
		t.Log(strings.Join([]string{
			"args:",
			fmt.Sprintf("INPUTCERTPATH: %v", fInputCertPath),
			fmt.Sprintf("INPUTKEYPATH: %v", fInputKeyPath),
			fmt.Sprintf("SECRETID: %v", fSecretId),
			fmt.Sprintf("SECRETKEY: %v", fSecretKey),
			fmt.Sprintf("DOMAIN: %v", fDomain),
			fmt.Sprintf("INSTANCEID: %v", fInstanceId),
		}, "\n"))

		deployer, err := provider.NewSSLDeployerProvider(&provider.SSLDeployerProviderConfig{
			SecretId:  fSecretId,
			SecretKey: fSecretKey,
			SubAppId:  fSubAppId,
			Domain:    fDomain,
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
