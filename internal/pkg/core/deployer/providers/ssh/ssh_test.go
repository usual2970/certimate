package ssh_test

import (
	"context"
	"os"
	"strconv"
	"testing"

	dSsh "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/ssh"
)

/*
Shell command to run this test:

	CERTIMATE_DEPLOYER_SSH_INPUTCERTPATH="/path/to/your-input-cert.pem" \
	CERTIMATE_DEPLOYER_SSH_INPUTKEYPATH="/path/to/your-input-key.pem" \
	CERTIMATE_DEPLOYER_SSH_SSHHOST="localhost" \
	CERTIMATE_DEPLOYER_SSH_SSHPORT=22 \
	CERTIMATE_DEPLOYER_SSH_SSHUSERNAME="root" \
	CERTIMATE_DEPLOYER_SSH_SSHPASSWORD="password" \
	CERTIMATE_DEPLOYER_SSH_OUTPUTCERTPATH="/path/to/your-output-cert.pem" \
	CERTIMATE_DEPLOYER_SSH_OUTPUTKEYPATH="/path/to/your-output-key.pem" \
	go test -v -run TestDeploy ssh_test.go
*/
func TestDeploy(t *testing.T) {
	envPrefix := "CERTIMATE_DEPLOYER_LOCAL_"
	tInputCertData, _ := os.ReadFile(os.Getenv(envPrefix + "INPUTCERTPATH"))
	tInputKeyData, _ := os.ReadFile(os.Getenv(envPrefix + "INPUTKEYPATH"))
	tSshHost := os.Getenv(envPrefix + "SSHHOST")
	tSshPort, _ := strconv.ParseInt(os.Getenv(envPrefix+"SSHPORT"), 10, 32)
	tSshUsername := os.Getenv(envPrefix + "SSHUSERNAME")
	tSshPassword := os.Getenv(envPrefix + "SSHPASSWORD")
	tOutputCertPath := os.Getenv(envPrefix + "OUTPUTCERTPATH")
	tOutputKeyPath := os.Getenv(envPrefix + "OUTPUTKEYPATH")

	deployer, err := dSsh.New(&dSsh.SshDeployerConfig{
		SshHost:        tSshHost,
		SshPort:        int32(tSshPort),
		SshUsername:    tSshUsername,
		SshPassword:    tSshPassword,
		OutputCertPath: tOutputCertPath,
		OutputKeyPath:  tOutputKeyPath,
	})
	if err != nil {
		t.Errorf("err: %+v", err)
		panic(err)
	}

	res, err := deployer.Deploy(context.Background(), string(tInputCertData), string(tInputKeyData))
	if err != nil {
		t.Errorf("err: %+v", err)
		panic(err)
	}

	t.Logf("ok: %v", res)
}
