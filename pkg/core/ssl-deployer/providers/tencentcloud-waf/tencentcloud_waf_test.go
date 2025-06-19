package tencentcloudwaf_test

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"

	provider "github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/tencentcloud-waf"
)

var (
	fInputCertPath string
	fInputKeyPath  string
	fSecretId      string
	fSecretKey     string
	fRegion        string
	fDomain        string
	fDomainId      string
	fInstanceId    string
)

func init() {
	argsPrefix := "CERTIMATE_SSLDEPLOYER_TENCENTCLOUDWAF_"

	flag.StringVar(&fInputCertPath, argsPrefix+"INPUTCERTPATH", "", "")
	flag.StringVar(&fInputKeyPath, argsPrefix+"INPUTKEYPATH", "", "")
	flag.StringVar(&fSecretId, argsPrefix+"SECRETID", "", "")
	flag.StringVar(&fSecretKey, argsPrefix+"SECRETKEY", "", "")
	flag.StringVar(&fRegion, argsPrefix+"REGION", "", "")
	flag.StringVar(&fDomain, argsPrefix+"DOMAIN", "", "")
	flag.StringVar(&fDomainId, argsPrefix+"DOMAINID", "", "")
	flag.StringVar(&fInstanceId, argsPrefix+"INSTANCEID", "", "")
}

/*
Shell command to run this test:

	go test -v ./tencentcloud_waf_test.go -args \
	--CERTIMATE_SSLDEPLOYER_TENCENTCLOUDWAF_INPUTCERTPATH="/path/to/your-input-cert.pem" \
	--CERTIMATE_SSLDEPLOYER_TENCENTCLOUDWAF_INPUTKEYPATH="/path/to/your-input-key.pem" \
	--CERTIMATE_SSLDEPLOYER_TENCENTCLOUDWAF_SECRETID="your-secret-id" \
	--CERTIMATE_SSLDEPLOYER_TENCENTCLOUDWAF_SECRETKEY="your-secret-key" \
	--CERTIMATE_SSLDEPLOYER_TENCENTCLOUDWAF_REGION="ap-guangzhou" \
	--CERTIMATE_SSLDEPLOYER_TENCENTCLOUDWAF_DOMAIN="example.com" \
	--CERTIMATE_SSLDEPLOYER_TENCENTCLOUDWAF_DOMAINID="your-domain-id" \
	--CERTIMATE_SSLDEPLOYER_TENCENTCLOUDWAF_INSTANCEID="your-instance-id"
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
			fmt.Sprintf("REGION: %v", fRegion),
			fmt.Sprintf("DOMAIN: %v", fDomain),
			fmt.Sprintf("INSTANCEID: %v", fInstanceId),
		}, "\n"))

		deployer, err := provider.NewSSLDeployerProvider(&provider.SSLDeployerProviderConfig{
			SecretId:   fSecretId,
			SecretKey:  fSecretKey,
			Region:     fRegion,
			Domain:     fDomain,
			DomainId:   fDomainId,
			InstanceId: fInstanceId,
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
