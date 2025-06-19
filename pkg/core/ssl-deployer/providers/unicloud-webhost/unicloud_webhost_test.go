package unicloudwebhost_test

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"

	provider "github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/unicloud-webhost"
)

var (
	fInputCertPath string
	fInputKeyPath  string
	fUsername      string
	fPassword      string
	fSpaceProvider string
	fSpaceId       string
	fDomain        string
)

func init() {
	argsPrefix := "CERTIMATE_SSLDEPLOYER_UNICLOUDWEBHOST_"

	flag.StringVar(&fInputCertPath, argsPrefix+"INPUTCERTPATH", "", "")
	flag.StringVar(&fInputKeyPath, argsPrefix+"INPUTKEYPATH", "", "")
	flag.StringVar(&fUsername, argsPrefix+"USERNAME", "", "")
	flag.StringVar(&fPassword, argsPrefix+"PASSWORD", "", "")
	flag.StringVar(&fSpaceProvider, argsPrefix+"SPACEPROVIDER", "", "")
	flag.StringVar(&fSpaceId, argsPrefix+"SPACEID", "", "")
	flag.StringVar(&fDomain, argsPrefix+"DOMAIN", "", "")
}

/*
Shell command to run this test:

	go test -v ./unicloud_webhost_test.go -args \
	--CERTIMATE_SSLDEPLOYER_UNICLOUDWEBHOST_INPUTCERTPATH="/path/to/your-input-cert.pem" \
	--CERTIMATE_SSLDEPLOYER_UNICLOUDWEBHOST_INPUTKEYPATH="/path/to/your-input-key.pem" \
	--CERTIMATE_SSLDEPLOYER_UNICLOUDWEBHOST_USERNAME="your-username" \
	--CERTIMATE_SSLDEPLOYER_UNICLOUDWEBHOST_PASSWORD="your-password" \
	--CERTIMATE_SSLDEPLOYER_UNICLOUDWEBHOST_SPACEPROVIDER="aliyun/tencent" \
	--CERTIMATE_SSLDEPLOYER_UNICLOUDWEBHOST_SPACEID="your-space-id" \
	--CERTIMATE_SSLDEPLOYER_UNICLOUDWEBHOST_DOMAIN="example.com"
*/
func TestDeploy(t *testing.T) {
	flag.Parse()

	t.Run("Deploy", func(t *testing.T) {
		t.Log(strings.Join([]string{
			"args:",
			fmt.Sprintf("INPUTCERTPATH: %v", fInputCertPath),
			fmt.Sprintf("INPUTKEYPATH: %v", fInputKeyPath),
			fmt.Sprintf("USERNAME: %v", fUsername),
			fmt.Sprintf("PASSWORD: %v", fPassword),
			fmt.Sprintf("SPACEPROVIDER: %v", fSpaceProvider),
			fmt.Sprintf("SPACEID: %v", fSpaceId),
			fmt.Sprintf("DOMAIN: %v", fDomain),
		}, "\n"))

		deployer, err := provider.NewSSLDeployerProvider(&provider.SSLDeployerProviderConfig{
			Username:      fUsername,
			Password:      fPassword,
			SpaceProvider: fSpaceProvider,
			SpaceId:       fSpaceId,
			Domain:        fDomain,
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
