package awscloudfront_test

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"

	provider "github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/aws-cloudfront"
)

var (
	fInputCertPath   string
	fInputKeyPath    string
	fAccessKeyId     string
	fSecretAccessKey string
	fRegion          string
	fDistribuitionId string
)

func init() {
	argsPrefix := "CERTIMATE_SSLDEPLOYER_AWSCLOUDFRONT_"

	flag.StringVar(&fInputCertPath, argsPrefix+"INPUTCERTPATH", "", "")
	flag.StringVar(&fInputKeyPath, argsPrefix+"INPUTKEYPATH", "", "")
	flag.StringVar(&fAccessKeyId, argsPrefix+"ACCESSKEYID", "", "")
	flag.StringVar(&fSecretAccessKey, argsPrefix+"SECRETACCESSKEY", "", "")
	flag.StringVar(&fRegion, argsPrefix+"REGION", "", "")
	flag.StringVar(&fDistribuitionId, argsPrefix+"DISTRIBUTIONID", "", "")
}

/*
Shell command to run this test:

	go test -v ./aws_cloudfront_test.go -args \
	--CERTIMATE_SSLDEPLOYER_AWSCLOUDFRONT_INPUTCERTPATH="/path/to/your-input-cert.pem" \
	--CERTIMATE_SSLDEPLOYER_AWSCLOUDFRONT_INPUTKEYPATH="/path/to/your-input-key.pem" \
	--CERTIMATE_SSLDEPLOYER_AWSCLOUDFRONT_ACCESSKEYID="your-access-key-id" \
	--CERTIMATE_SSLDEPLOYER_AWSCLOUDFRONT_SECRETACCESSKEY="your-secret-access-id" \
	--CERTIMATE_SSLDEPLOYER_AWSCLOUDFRONT_REGION="us-east-1" \
	--CERTIMATE_SSLDEPLOYER_AWSCLOUDFRONT_DISTRIBUTIONID="your-distribution-id"
*/
func TestDeploy(t *testing.T) {
	flag.Parse()

	t.Run("Deploy", func(t *testing.T) {
		t.Log(strings.Join([]string{
			"args:",
			fmt.Sprintf("INPUTCERTPATH: %v", fInputCertPath),
			fmt.Sprintf("INPUTKEYPATH: %v", fInputKeyPath),
			fmt.Sprintf("ACCESSKEYID: %v", fAccessKeyId),
			fmt.Sprintf("SECRETACCESSKEY: %v", fSecretAccessKey),
			fmt.Sprintf("REGION: %v", fRegion),
			fmt.Sprintf("DISTRIBUTIONID: %v", fDistribuitionId),
		}, "\n"))

		deployer, err := provider.NewSSLDeployerProvider(&provider.SSLDeployerProviderConfig{
			AccessKeyId:     fAccessKeyId,
			SecretAccessKey: fSecretAccessKey,
			Region:          fRegion,
			DistributionId:  fDistribuitionId,
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
