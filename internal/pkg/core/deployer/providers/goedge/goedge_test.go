package goedge_test

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"

	provider "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/goedge"
)

var (
	fInputCertPath string
	fInputKeyPath  string
	fApiUrl        string
	fAccessKeyId   string
	fAccessKey     string
	fCertificateId int
)

func init() {
	argsPrefix := "CERTIMATE_DEPLOYER_GOEDGE_"

	flag.StringVar(&fInputCertPath, argsPrefix+"INPUTCERTPATH", "", "")
	flag.StringVar(&fInputKeyPath, argsPrefix+"INPUTKEYPATH", "", "")
	flag.StringVar(&fApiUrl, argsPrefix+"APIURL", "", "")
	flag.StringVar(&fAccessKeyId, argsPrefix+"ACCESSKEYID", "", "")
	flag.StringVar(&fAccessKey, argsPrefix+"ACCESSKEY", "", "")
	flag.IntVar(&fCertificateId, argsPrefix+"CERTIFICATEID", 0, "")
}

/*
Shell command to run this test:

	go test -v ./goedge_test.go -args \
	--CERTIMATE_DEPLOYER_GOEDGE_INPUTCERTPATH="/path/to/your-input-cert.pem" \
	--CERTIMATE_DEPLOYER_GOEDGE_INPUTKEYPATH="/path/to/your-input-key.pem" \
	--CERTIMATE_DEPLOYER_GOEDGE_APIURL="http://127.0.0.1:7788" \
	--CERTIMATE_DEPLOYER_GOEDGE_ACCESSKEYID="your-access-key-id" \
	--CERTIMATE_DEPLOYER_GOEDGE_ACCESSKEY="your-access-key" \
	--CERTIMATE_DEPLOYER_GOEDGE_CERTIFICATEID="your-cerficiate-id"
*/
func TestDeploy(t *testing.T) {
	flag.Parse()

	t.Run("Deploy", func(t *testing.T) {
		t.Log(strings.Join([]string{
			"args:",
			fmt.Sprintf("INPUTCERTPATH: %v", fInputCertPath),
			fmt.Sprintf("INPUTKEYPATH: %v", fInputKeyPath),
			fmt.Sprintf("APIURL: %v", fApiUrl),
			fmt.Sprintf("ACCESSKEYID: %v", fAccessKeyId),
			fmt.Sprintf("ACCESSKEY: %v", fAccessKey),
			fmt.Sprintf("CERTIFICATEID: %v", fCertificateId),
		}, "\n"))

		deployer, err := provider.NewDeployer(&provider.DeployerConfig{
			ApiUrl:                   fApiUrl,
			AccessKeyId:              fAccessKeyId,
			AccessKey:                fAccessKey,
			AllowInsecureConnections: true,
			ResourceType:             provider.RESOURCE_TYPE_CERTIFICATE,
			CertificateId:            int64(fCertificateId),
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
