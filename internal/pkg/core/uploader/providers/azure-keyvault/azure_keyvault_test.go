package azurekeyvault_test

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"

	provider "github.com/usual2970/certimate/internal/pkg/core/uploader/providers/azure-keyvault"
)

var (
	fInputCertPath string
	fInputKeyPath  string
	fTenantId      string
	fClientId      string
	fClientSecret  string
	fCloudName     string
	fKeyVaultName  string
)

func init() {
	argsPrefix := "CERTIMATE_UPLOADER_AZUREKEYVAULT_"

	flag.StringVar(&fInputCertPath, argsPrefix+"INPUTCERTPATH", "", "")
	flag.StringVar(&fInputKeyPath, argsPrefix+"INPUTKEYPATH", "", "")
	flag.StringVar(&fTenantId, argsPrefix+"TENANTID", "", "")
	flag.StringVar(&fClientId, argsPrefix+"CLIENTID", "", "")
	flag.StringVar(&fClientSecret, argsPrefix+"CLIENTSECRET", "", "")
	flag.StringVar(&fCloudName, argsPrefix+"CLOUDNAME", "", "")
	flag.StringVar(&fKeyVaultName, argsPrefix+"KEYVAULTNAME", "", "")
}

/*
Shell command to run this test:

	go test -v ./azure_keyvault_test.go -args \
	--CERTIMATE_UPLOADER_AZUREKEYVAULT_INPUTCERTPATH="/path/to/your-input-cert.pem" \
	--CERTIMATE_UPLOADER_AZUREKEYVAULT_INPUTKEYPATH="/path/to/your-input-key.pem" \
	--CERTIMATE_UPLOADER_AZUREKEYVAULT_TENANTID="your-tenant-id" \
	--CERTIMATE_UPLOADER_AZUREKEYVAULT_CLIENTID="your-app-registration-client-id" \
	--CERTIMATE_UPLOADER_AZUREKEYVAULT_CLIENTSECRET="your-app-registration-client-secret" \
	--CERTIMATE_UPLOADER_AZUREKEYVAULT_CLOUDNAME="china" \
	--CERTIMATE_UPLOADER_AZUREKEYVAULT_KEYVAULTNAME="your-keyvault-name"
*/
func TestDeploy(t *testing.T) {
	flag.Parse()

	t.Run("Deploy", func(t *testing.T) {
		t.Log(strings.Join([]string{
			"args:",
			fmt.Sprintf("INPUTCERTPATH: %v", fInputCertPath),
			fmt.Sprintf("INPUTKEYPATH: %v", fInputKeyPath),
			fmt.Sprintf("TENANTID: %v", fTenantId),
			fmt.Sprintf("CLIENTID: %v", fClientId),
			fmt.Sprintf("CLIENTSECRET: %v", fClientSecret),
			fmt.Sprintf("CLOUDNAME: %v", fCloudName),
			fmt.Sprintf("KEYVAULTNAME: %v", fKeyVaultName),
		}, "\n"))

		uploader, err := provider.NewUploader(&provider.UploaderConfig{
			TenantId:     fTenantId,
			ClientId:     fClientId,
			ClientSecret: fClientSecret,
			CloudName:    fCloudName,
			KeyVaultName: fKeyVaultName,
		})
		if err != nil {
			t.Errorf("err: %+v", err)
			return
		}

		fInputCertData, _ := os.ReadFile(fInputCertPath)
		fInputKeyData, _ := os.ReadFile(fInputKeyPath)
		res, err := uploader.Upload(context.Background(), string(fInputCertData), string(fInputKeyData))
		if err != nil {
			t.Errorf("err: %+v", err)
			return
		}

		sres, _ := json.Marshal(res)
		t.Logf("ok: %s", string(sres))
	})
}
