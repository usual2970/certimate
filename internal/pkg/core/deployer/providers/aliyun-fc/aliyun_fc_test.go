package aliyunfc_test

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"

	provider "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/aliyun-fc"
)

var (
	fInputCertPath   string
	fInputKeyPath    string
	fAccessKeyId     string
	fAccessKeySecret string
	fRegion          string
	fSiteId          int64
)

func init() {
	argsPrefix := "CERTIMATE_DEPLOYER_ALIYUNFC_"

	flag.StringVar(&fInputCertPath, argsPrefix+"INPUTCERTPATH", "", "")
	flag.StringVar(&fInputKeyPath, argsPrefix+"INPUTKEYPATH", "", "")
	flag.StringVar(&fAccessKeyId, argsPrefix+"ACCESSKEYID", "", "")
	flag.StringVar(&fAccessKeySecret, argsPrefix+"ACCESSKEYSECRET", "", "")
	flag.StringVar(&fRegion, argsPrefix+"REGION", "", "")
	flag.Int64Var(&fSiteId, argsPrefix+"SITEID", "", "")
}

/*
Shell command to run this test:

	go test -v ./aliyun_fc_test.go -args \
	--CERTIMATE_DEPLOYER_ALIYUNFC_INPUTCERTPATH="/path/to/your-input-cert.pem" \
	--CERTIMATE_DEPLOYER_ALIYUNFC_INPUTKEYPATH="/path/to/your-input-key.pem" \
	--CERTIMATE_DEPLOYER_ALIYUNFC_ACCESSKEYID="your-access-key-id" \
	--CERTIMATE_DEPLOYER_ALIYUNFC_ACCESSKEYSECRET="your-access-key-secret" \
	--CERTIMATE_DEPLOYER_ALIYUNFC_REGION="cn-hangzhou" \
	--CERTIMATE_DEPLOYER_ALIYUNFC_SITEID="your-fc-site-id"
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
			fmt.Sprintf("SITEID: %v", fSiteId),
		}, "\n"))

		deployer, err := provider.NewDeployer(&provider.DeployerConfig{
			AccessKeyId:     fAccessKeyId,
			AccessKeySecret: fAccessKeySecret,
			Region:          fRegion,
			SiteId:          fSiteId,
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
