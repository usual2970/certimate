package k8ssecret_test

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"

	provider "github.com/certimate-go/certimate/pkg/core/ssl-deployer/providers/k8s-secret"
)

var (
	fInputCertPath       string
	fInputKeyPath        string
	fNamespace           string
	fSecretName          string
	fSecretDataKeyForCrt string
	fSecretDataKeyForKey string
)

func init() {
	argsPrefix := "CERTIMATE_SSLDEPLOYER_K8SSECRET_"

	flag.StringVar(&fInputCertPath, argsPrefix+"INPUTCERTPATH", "", "")
	flag.StringVar(&fInputKeyPath, argsPrefix+"INPUTKEYPATH", "", "")
	flag.StringVar(&fNamespace, argsPrefix+"NAMESPACE", "default", "")
	flag.StringVar(&fSecretName, argsPrefix+"SECRETNAME", "", "")
	flag.StringVar(&fSecretDataKeyForCrt, argsPrefix+"SECRETDATAKEYFORCRT", "tls.crt", "")
	flag.StringVar(&fSecretDataKeyForKey, argsPrefix+"SECRETDATAKEYFORKEY", "tls.key", "")
}

/*
Shell command to run this test:

	go test -v ./k8s_secret_test.go -args \
	--CERTIMATE_SSLDEPLOYER_K8SSECRET_INPUTCERTPATH="/path/to/your-input-cert.pem" \
	--CERTIMATE_SSLDEPLOYER_K8SSECRET_INPUTKEYPATH="/path/to/your-input-key.pem" \
	--CERTIMATE_SSLDEPLOYER_K8SSECRET_NAMESPACE="default" \
	--CERTIMATE_SSLDEPLOYER_K8SSECRET_SECRETNAME="secret" \
	--CERTIMATE_SSLDEPLOYER_K8SSECRET_SECRETDATAKEYFORCRT="tls.crt" \
	--CERTIMATE_SSLDEPLOYER_K8SSECRET_SECRETDATAKEYFORKEY="tls.key"
*/
func TestDeploy(t *testing.T) {
	flag.Parse()

	t.Run("Deploy", func(t *testing.T) {
		t.Log(strings.Join([]string{
			"args:",
			fmt.Sprintf("INPUTCERTPATH: %v", fInputCertPath),
			fmt.Sprintf("INPUTKEYPATH: %v", fInputKeyPath),
			fmt.Sprintf("NAMESPACE: %v", fNamespace),
			fmt.Sprintf("SECRETNAME: %v", fSecretName),
			fmt.Sprintf("SECRETDATAKEYFORCRT: %v", fSecretDataKeyForCrt),
			fmt.Sprintf("SECRETDATAKEYFORKEY: %v", fSecretDataKeyForKey),
		}, "\n"))

		deployer, err := provider.NewSSLDeployerProvider(&provider.SSLDeployerProviderConfig{
			Namespace:           fNamespace,
			SecretName:          fSecretName,
			SecretDataKeyForCrt: fSecretDataKeyForCrt,
			SecretDataKeyForKey: fSecretDataKeyForKey,
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
