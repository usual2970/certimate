package aliyunnlb_test

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"

	provider "github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/aliyun-nlb"
)

var (
	fInputCertPath   string
	fInputKeyPath    string
	fAccessKeyId     string
	fAccessKeySecret string
	fRegion          string
	fLoadbalancerId  string
	fListenerId      string
)

func init() {
	argsPrefix := "CERTIMATE_SSLDEPLOYER_ALIYUNNLB_"

	flag.StringVar(&fInputCertPath, argsPrefix+"INPUTCERTPATH", "", "")
	flag.StringVar(&fInputKeyPath, argsPrefix+"INPUTKEYPATH", "", "")
	flag.StringVar(&fAccessKeyId, argsPrefix+"ACCESSKEYID", "", "")
	flag.StringVar(&fAccessKeySecret, argsPrefix+"ACCESSKEYSECRET", "", "")
	flag.StringVar(&fRegion, argsPrefix+"REGION", "", "")
	flag.StringVar(&fLoadbalancerId, argsPrefix+"LOADBALANCERID", "", "")
	flag.StringVar(&fListenerId, argsPrefix+"LISTENERID", "", "")
}

/*
Shell command to run this test:

	go test -v ./aliyun_nlb_test.go -args \
	--CERTIMATE_SSLDEPLOYER_ALIYUNNLB_INPUTCERTPATH="/path/to/your-input-cert.pem" \
	--CERTIMATE_SSLDEPLOYER_ALIYUNNLB_INPUTKEYPATH="/path/to/your-input-key.pem" \
	--CERTIMATE_SSLDEPLOYER_ALIYUNNLB_ACCESSKEYID="your-access-key-id" \
	--CERTIMATE_SSLDEPLOYER_ALIYUNNLB_ACCESSKEYSECRET="your-access-key-secret" \
	--CERTIMATE_SSLDEPLOYER_ALIYUNNLB_REGION="cn-hangzhou" \
	--CERTIMATE_SSLDEPLOYER_ALIYUNNLB_LOADBALANCERID="your-nlb-instance-id" \
	--CERTIMATE_SSLDEPLOYER_ALIYUNNLB_LISTENERID="your-nlb-listener-id"
*/
func TestDeploy(t *testing.T) {
	flag.Parse()

	t.Run("Deploy_ToLoadbalancer", func(t *testing.T) {
		t.Log(strings.Join([]string{
			"args:",
			fmt.Sprintf("INPUTCERTPATH: %v", fInputCertPath),
			fmt.Sprintf("INPUTKEYPATH: %v", fInputKeyPath),
			fmt.Sprintf("ACCESSKEYID: %v", fAccessKeyId),
			fmt.Sprintf("ACCESSKEYSECRET: %v", fAccessKeySecret),
			fmt.Sprintf("REGION: %v", fRegion),
			fmt.Sprintf("LOADBALANCERID: %v", fLoadbalancerId),
		}, "\n"))

		deployer, err := provider.NewSSLDeployerProvider(&provider.SSLDeployerProviderConfig{
			AccessKeyId:     fAccessKeyId,
			AccessKeySecret: fAccessKeySecret,
			Region:          fRegion,
			ResourceType:    provider.RESOURCE_TYPE_LOADBALANCER,
			LoadbalancerId:  fLoadbalancerId,
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

	t.Run("Deploy_ToListener", func(t *testing.T) {
		t.Log(strings.Join([]string{
			"args:",
			fmt.Sprintf("INPUTCERTPATH: %v", fInputCertPath),
			fmt.Sprintf("INPUTKEYPATH: %v", fInputKeyPath),
			fmt.Sprintf("ACCESSKEYID: %v", fAccessKeyId),
			fmt.Sprintf("ACCESSKEYSECRET: %v", fAccessKeySecret),
			fmt.Sprintf("REGION: %v", fRegion),
			fmt.Sprintf("LOADBALANCERID: %v", fLoadbalancerId),
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
