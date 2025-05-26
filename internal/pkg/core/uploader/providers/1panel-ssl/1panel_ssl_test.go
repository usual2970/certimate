package onepanelssl_test

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"

	provider "github.com/usual2970/certimate/internal/pkg/core/uploader/providers/1panel-ssl"
)

var (
	fInputCertPath string
	fInputKeyPath  string
	fServerUrl     string
	fApiVersion    string
	fApiKey        string
)

func init() {
	argsPrefix := "CERTIMATE_UPLOADER_1PANELSSL_"

	flag.StringVar(&fInputCertPath, argsPrefix+"INPUTCERTPATH", "", "")
	flag.StringVar(&fInputKeyPath, argsPrefix+"INPUTKEYPATH", "", "")
	flag.StringVar(&fServerUrl, argsPrefix+"SERVERURL", "", "")
	flag.StringVar(&fApiVersion, argsPrefix+"APIVERSION", "v1", "")
	flag.StringVar(&fApiKey, argsPrefix+"APIKEY", "", "")
}

/*
Shell command to run this test:

	go test -v ./1panel_ssl_test.go -args \
	--CERTIMATE_UPLOADER_1PANELSSL_INPUTCERTPATH="/path/to/your-input-cert.pem" \
	--CERTIMATE_UPLOADER_1PANELSSL_INPUTKEYPATH="/path/to/your-input-key.pem" \
	--CERTIMATE_UPLOADER_1PANELSSL_SERVERURL="http://127.0.0.1:20410" \
	--CERTIMATE_UPLOADER_1PANELSSL_APIVERSION="v1" \
	--CERTIMATE_UPLOADER_1PANELSSL_APIKEY="your-api-key"
*/
func TestDeploy(t *testing.T) {
	flag.Parse()

	t.Run("Deploy", func(t *testing.T) {
		t.Log(strings.Join([]string{
			"args:",
			fmt.Sprintf("INPUTCERTPATH: %v", fInputCertPath),
			fmt.Sprintf("INPUTKEYPATH: %v", fInputKeyPath),
			fmt.Sprintf("SERVERURL: %v", fServerUrl),
			fmt.Sprintf("APIVERSION: %v", fApiVersion),
			fmt.Sprintf("APIKEY: %v", fApiKey),
		}, "\n"))

		uploader, err := provider.NewUploader(&provider.UploaderConfig{
			ServerUrl:  fServerUrl,
			ApiVersion: fApiVersion,
			ApiKey:     fApiKey,
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
