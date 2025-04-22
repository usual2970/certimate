package rainyunsslcenter_test

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"

	provider "github.com/usual2970/certimate/internal/pkg/core/uploader/providers/rainyun-sslcenter"
)

var (
	fInputCertPath string
	fInputKeyPath  string
	fApiKey        string
)

func init() {
	argsPrefix := "CERTIMATE_UPLOADER_RAINYUNSSLCENTER_"

	flag.StringVar(&fInputCertPath, argsPrefix+"INPUTCERTPATH", "", "")
	flag.StringVar(&fInputKeyPath, argsPrefix+"INPUTKEYPATH", "", "")
	flag.StringVar(&fApiKey, argsPrefix+"APIKEY", "", "")
}

/*
Shell command to run this test:

	go test -v ./rainyun_sslcenter_test.go -args \
	--CERTIMATE_UPLOADER_RAINYUNSSLCENTER_INPUTCERTPATH="/path/to/your-input-cert.pem" \
	--CERTIMATE_UPLOADER_RAINYUNSSLCENTER_INPUTKEYPATH="/path/to/your-input-key.pem" \
	--CERTIMATE_UPLOADER_RAINYUNSSLCENTER_APIKEY="your-api-key"
*/
func TestDeploy(t *testing.T) {
	flag.Parse()

	t.Run("Deploy", func(t *testing.T) {
		t.Log(strings.Join([]string{
			"args:",
			fmt.Sprintf("INPUTCERTPATH: %v", fInputCertPath),
			fmt.Sprintf("INPUTKEYPATH: %v", fInputKeyPath),
			fmt.Sprintf("APIKEY: %v", fApiKey),
		}, "\n"))

		uploader, err := provider.NewUploader(&provider.UploaderConfig{
			ApiKey: fApiKey,
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
