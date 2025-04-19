package aliyunapigw_test

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"

	provider "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/aliyun-apigw"
)

var (
	fInputCertPath   string
	fInputKeyPath    string
	fAccessKeyId     string
	fAccessKeySecret string
	fRegion          string
	fServiceType     string
	fGatewayId       string
	fGroupId         string
	fDomain          string
)

func init() {
	argsPrefix := "CERTIMATE_DEPLOYER_ALIYUNAPIGW_"

	flag.StringVar(&fInputCertPath, argsPrefix+"INPUTCERTPATH", "", "")
	flag.StringVar(&fInputKeyPath, argsPrefix+"INPUTKEYPATH", "", "")
	flag.StringVar(&fAccessKeyId, argsPrefix+"ACCESSKEYID", "", "")
	flag.StringVar(&fAccessKeySecret, argsPrefix+"ACCESSKEYSECRET", "", "")
	flag.StringVar(&fRegion, argsPrefix+"REGION", "", "")
	flag.StringVar(&fGatewayId, argsPrefix+"GATEWARYID", "", "")
	flag.StringVar(&fGroupId, argsPrefix+"GROUPID", "", "")
	flag.StringVar(&fServiceType, argsPrefix+"SERVICETYPE", "", "")
	flag.StringVar(&fDomain, argsPrefix+"DOMAIN", "", "")
}

/*
Shell command to run this test:

	go test -v ./aliyun_apigw_test.go -args \
	--CERTIMATE_DEPLOYER_ALIYUNAPIGW_INPUTCERTPATH="/path/to/your-input-cert.pem" \
	--CERTIMATE_DEPLOYER_ALIYUNAPIGW_INPUTKEYPATH="/path/to/your-input-key.pem" \
	--CERTIMATE_DEPLOYER_ALIYUNAPIGW_ACCESSKEYID="your-access-key-id" \
	--CERTIMATE_DEPLOYER_ALIYUNAPIGW_ACCESSKEYSECRET="your-access-key-secret" \
	--CERTIMATE_DEPLOYER_ALIYUNAPIGW_REGION="cn-hangzhou" \
	--CERTIMATE_DEPLOYER_ALIYUNAPIGW_GATEWAYID="your-api-gateway-id" \
	--CERTIMATE_DEPLOYER_ALIYUNAPIGW_GROUPID="your-api-group-id" \
	--CERTIMATE_DEPLOYER_ALIYUNAPIGW_SERVICETYPE="cloudnative" \
	--CERTIMATE_DEPLOYER_ALIYUNAPIGW_DOMAIN="example.com"
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
			fmt.Sprintf("GATEWAYID: %v", fGatewayId),
			fmt.Sprintf("GROUPID: %v", fGroupId),
			fmt.Sprintf("SERVICETYPE: %v", fServiceType),
			fmt.Sprintf("DOMAIN: %v", fDomain),
		}, "\n"))

		deployer, err := provider.NewDeployer(&provider.DeployerConfig{
			AccessKeyId:     fAccessKeyId,
			AccessKeySecret: fAccessKeySecret,
			Region:          fRegion,
			ServiceType:     provider.ServiceType(fServiceType),
			GatewayId:       fGatewayId,
			GroupId:         fGroupId,
			Domain:          fDomain,
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
