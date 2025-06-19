package tencentcloudgaap_test

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"

	provider "github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/tencentcloud-gaap"
)

var (
	fInputCertPath string
	fInputKeyPath  string
	fSecretId      string
	fSecretKey     string
	fProxyId       string
	fListenerId    string
)

func init() {
	argsPrefix := "CERTIMATE_SSLDEPLOYER_TENCENTCLOUDCDN_"

	flag.StringVar(&fInputCertPath, argsPrefix+"INPUTCERTPATH", "", "")
	flag.StringVar(&fInputKeyPath, argsPrefix+"INPUTKEYPATH", "", "")
	flag.StringVar(&fSecretId, argsPrefix+"SECRETID", "", "")
	flag.StringVar(&fSecretKey, argsPrefix+"SECRETKEY", "", "")
	flag.StringVar(&fProxyId, argsPrefix+"PROXYID", "", "")
	flag.StringVar(&fListenerId, argsPrefix+"LISTENERID", "", "")
}

/*
Shell command to run this test:

	go test -v ./tencentcloud_gaap_test.go -args \
	--CERTIMATE_SSLDEPLOYER_TENCENTCLOUDGAAP_INPUTCERTPATH="/path/to/your-input-cert.pem" \
	--CERTIMATE_SSLDEPLOYER_TENCENTCLOUDGAAP_INPUTKEYPATH="/path/to/your-input-key.pem" \
	--CERTIMATE_SSLDEPLOYER_TENCENTCLOUDGAAP_SECRETID="your-secret-id" \
	--CERTIMATE_SSLDEPLOYER_TENCENTCLOUDGAAP_SECRETKEY="your-secret-key" \
	--CERTIMATE_SSLDEPLOYER_TENCENTCLOUDGAAP_PROXYID="your-gaap-group-id" \
	--CERTIMATE_SSLDEPLOYER_TENCENTCLOUDGAAP_LISTENERID="your-clb-listener-id"
*/
func TestDeploy(t *testing.T) {
	flag.Parse()

	t.Run("Deploy_ToListener", func(t *testing.T) {
		t.Log(strings.Join([]string{
			"args:",
			fmt.Sprintf("INPUTCERTPATH: %v", fInputCertPath),
			fmt.Sprintf("INPUTKEYPATH: %v", fInputKeyPath),
			fmt.Sprintf("SECRETID: %v", fSecretId),
			fmt.Sprintf("SECRETKEY: %v", fSecretKey),
			fmt.Sprintf("PROXYID: %v", fProxyId),
			fmt.Sprintf("LISTENERID: %v", fListenerId),
		}, "\n"))

		deployer, err := provider.NewSSLDeployerProvider(&provider.SSLDeployerProviderConfig{
			SecretId:     fSecretId,
			SecretKey:    fSecretKey,
			ResourceType: provider.RESOURCE_TYPE_LISTENER,
			ProxyId:      fProxyId,
			ListenerId:   fListenerId,
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
