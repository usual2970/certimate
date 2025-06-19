package volcengineclb_test

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"

	provider "github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/volcengine-clb"
)

var (
	fInputCertPath   string
	fInputKeyPath    string
	fAccessKeyId     string
	fAccessKeySecret string
	fRegion          string
	fListenerId      string
)

func init() {
	argsPrefix := "CERTIMATE_SSLDEPLOYER_VOLCENGINECLB_"

	flag.StringVar(&fInputCertPath, argsPrefix+"INPUTCERTPATH", "", "")
	flag.StringVar(&fInputKeyPath, argsPrefix+"INPUTKEYPATH", "", "")
	flag.StringVar(&fAccessKeyId, argsPrefix+"ACCESSKEYID", "", "")
	flag.StringVar(&fAccessKeySecret, argsPrefix+"ACCESSKEYSECRET", "", "")
	flag.StringVar(&fRegion, argsPrefix+"REGION", "", "")
	flag.StringVar(&fListenerId, argsPrefix+"LISTENERID", "", "")
}

/*
Shell command to run this test:

	go test -v ./volcengine_clb_test.go -args \
	--CERTIMATE_SSLDEPLOYER_VOLCENGINECLB_INPUTCERTPATH="/path/to/your-input-cert.pem" \
	--CERTIMATE_SSLDEPLOYER_VOLCENGINECLB_INPUTKEYPATH="/path/to/your-input-key.pem" \
	--CERTIMATE_SSLDEPLOYER_VOLCENGINECLB_ACCESSKEYID="your-access-key-id" \
	--CERTIMATE_SSLDEPLOYER_VOLCENGINECLB_ACCESSKEYSECRET="your-access-key-secret" \
	--CERTIMATE_SSLDEPLOYER_VOLCENGINECLB_REGION="cn-beijing" \
	--CERTIMATE_SSLDEPLOYER_VOLCENGINECLB_LISTENERID="your-listener-id"
*/
func TestDeploy(t *testing.T) {
	flag.Parse()

	t.Run("Deploy", func(t *testing.T) {
		t.Log(strings.Join([]string{
			"args:",
			fmt.Sprintf("INPUTCERTPATH: %v", fInputCertPath),
			fmt.Sprintf("INPUTKEYPATH: %v", fInputKeyPath),
			fmt.Sprintf("ACCESSKEYID: %v", fAccessKeyId),
			fmt.Sprintf("ACCESSKEYSECRET: %v", fAccessKeySecret),
			fmt.Sprintf("REGION: %v", fRegion),
			fmt.Sprintf("LISTENERID: %v", fListenerId),
		}, "\n"))

		deployer, err := provider.NewSSLDeployerProvider(&provider.SSLDeployerProviderConfig{
			AccessKeyId:     fAccessKeyId,
			AccessKeySecret: fAccessKeySecret,
			Region:          fRegion,
			ResourceType:    provider.RESOURCE_TYPE_LISTENER,
			ListenerId:      fListenerId,
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
