package ssh_test

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"

	provider "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/ssh"
)

var (
	fInputCertPath  string
	fInputKeyPath   string
	fSshHost        string
	fSshPort        int
	fSshUsername    string
	fSshPassword    string
	fOutputCertPath string
	fOutputKeyPath  string
)

func init() {
	argsPrefix := "CERTIMATE_DEPLOYER_SSH_"

	flag.StringVar(&fInputCertPath, argsPrefix+"INPUTCERTPATH", "", "")
	flag.StringVar(&fInputKeyPath, argsPrefix+"INPUTKEYPATH", "", "")
	flag.StringVar(&fSshHost, argsPrefix+"SSHHOST", "", "")
	flag.IntVar(&fSshPort, argsPrefix+"SSHPORT", 0, "")
	flag.StringVar(&fSshUsername, argsPrefix+"SSHUSERNAME", "", "")
	flag.StringVar(&fSshPassword, argsPrefix+"SSHPASSWORD", "", "")
	flag.StringVar(&fOutputCertPath, argsPrefix+"OUTPUTCERTPATH", "", "")
	flag.StringVar(&fOutputKeyPath, argsPrefix+"OUTPUTKEYPATH", "", "")
}

/*
Shell command to run this test:

	go test -v ./ssh_test.go -args \
	--CERTIMATE_DEPLOYER_SSH_INPUTCERTPATH="/path/to/your-input-cert.pem" \
	--CERTIMATE_DEPLOYER_SSH_INPUTKEYPATH="/path/to/your-input-key.pem" \
	--CERTIMATE_DEPLOYER_SSH_SSHHOST="localhost" \
	--CERTIMATE_DEPLOYER_SSH_SSHPORT=22 \
	--CERTIMATE_DEPLOYER_SSH_SSHUSERNAME="root" \
	--CERTIMATE_DEPLOYER_SSH_SSHPASSWORD="password" \
	--CERTIMATE_DEPLOYER_SSH_OUTPUTCERTPATH="/path/to/your-output-cert.pem" \
	--CERTIMATE_DEPLOYER_SSH_OUTPUTKEYPATH="/path/to/your-output-key.pem"
*/
func TestDeploy(t *testing.T) {
	flag.Parse()

	t.Run("Deploy", func(t *testing.T) {
		t.Log(strings.Join([]string{
			"args:",
			fmt.Sprintf("INPUTCERTPATH: %v", fInputCertPath),
			fmt.Sprintf("INPUTKEYPATH: %v", fInputKeyPath),
			fmt.Sprintf("SSHHOST: %v", fSshHost),
			fmt.Sprintf("SSHPORT: %v", fSshPort),
			fmt.Sprintf("SSHUSERNAME: %v", fSshUsername),
			fmt.Sprintf("SSHPASSWORD: %v", fSshPassword),
			fmt.Sprintf("OUTPUTCERTPATH: %v", fOutputCertPath),
			fmt.Sprintf("OUTPUTKEYPATH: %v", fOutputKeyPath),
		}, "\n"))

		deployer, err := provider.NewDeployer(&provider.DeployerConfig{
			SshHost:        fSshHost,
			SshPort:        int32(fSshPort),
			SshUsername:    fSshUsername,
			SshPassword:    fSshPassword,
			OutputFormat:   provider.OUTPUT_FORMAT_PEM,
			OutputCertPath: fOutputCertPath,
			OutputKeyPath:  fOutputKeyPath,
		})
		if err != nil {
			t.Errorf("err: %+v", err)
		}

		fInputCertData, _ := os.ReadFile(fInputCertPath)
		fInputKeyData, _ := os.ReadFile(fInputKeyPath)
		res, err := deployer.Deploy(context.Background(), string(fInputCertData), string(fInputKeyData))
		if err != nil {
			t.Errorf("err: %+v", err)
		}

		t.Logf("ok: %v", res)
	})
}
