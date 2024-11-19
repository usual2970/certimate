package local_test

import (
	"context"
	"os"
	"testing"

	dLocal "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/local"
)

/*
Shell command to run this test:

	CERTIMATE_DEPLOYER_LOCAL_INPUTCERTPATH="/path/to/your-input-cert.pem" \
	CERTIMATE_DEPLOYER_LOCAL_INPUTKEYPATH="/path/to/your-input-key.pem" \
	CERTIMATE_DEPLOYER_LOCAL_OUTPUTCERTPATH="/path/to/your-output-cert.pem" \
	CERTIMATE_DEPLOYER_LOCAL_OUTPUTKEYPATH="/path/to/your-output-key.pem" \
	go test -v -run TestDeploy local_test.go
*/
func TestDeploy(t *testing.T) {
	envPrefix := "CERTIMATE_DEPLOYER_LOCAL_"
	tInputCertData, _ := os.ReadFile(os.Getenv(envPrefix + "INPUTCERTPATH"))
	tInputKeyData, _ := os.ReadFile(os.Getenv(envPrefix + "INPUTKEYPATH"))
	tOutputCertPath := os.Getenv(envPrefix + "OUTPUTCERTPATH")
	tOutputKeyPath := os.Getenv(envPrefix + "OUTPUTKEYPATH")

	deployer, err := dLocal.New(&dLocal.LocalDeployerConfig{
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

/*
Shell command to run this test:

	CERTIMATE_DEPLOYER_LOCAL_INPUTCERTPATH="/path/to/your-input-cert.pem" \
	CERTIMATE_DEPLOYER_LOCAL_INPUTKEYPATH="/path/to/your-input-key.pem" \
	CERTIMATE_DEPLOYER_LOCAL_OUTPUTCERTPATH="/path/to/your-output-cert.pem" \
	CERTIMATE_DEPLOYER_LOCAL_PFXPASSWORD="your-pfx-password" \
	go test -v -run TestDeploy_PFX local_test.go
*/
func TestDeploy_PFX(t *testing.T) {
	envPrefix := "CERTIMATE_DEPLOYER_LOCAL_"
	tInputCertData, _ := os.ReadFile(os.Getenv(envPrefix + "INPUTCERTPATH"))
	tInputKeyData, _ := os.ReadFile(os.Getenv(envPrefix + "INPUTKEYPATH"))
	tOutputCertPath := os.Getenv(envPrefix + "OUTPUTCERTPATH")
	tPfxPassword := os.Getenv(envPrefix + "PFXPASSWORD")

	deployer, err := dLocal.New(&dLocal.LocalDeployerConfig{
		OutputFormat:   dLocal.OUTPUT_FORMAT_PFX,
		OutputCertPath: tOutputCertPath,
		PfxPassword:    tPfxPassword,
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

/*
Shell command to run this test:

		CERTIMATE_DEPLOYER_LOCAL_INPUTCERTPATH="/path/to/your-input-cert.pem" \
		CERTIMATE_DEPLOYER_LOCAL_INPUTKEYPATH="/path/to/your-input-key.pem" \
		CERTIMATE_DEPLOYER_LOCAL_OUTPUTCERTPATH="/path/to/your-output-cert.pem" \
		CERTIMATE_DEPLOYER_LOCAL_JKSALIAS="your-jks-alias" \
	  CERTIMATE_DEPLOYER_LOCAL_JKSKEYPASS="your-jks-keypass" \
	  CERTIMATE_DEPLOYER_LOCAL_JKSSTOREPASS="your-jks-storepass" \
		go test -v -run TestDeploy_JKS local_test.go
*/
func TestDeploy_JKS(t *testing.T) {
	envPrefix := "CERTIMATE_DEPLOYER_LOCAL_"
	tInputCertData, _ := os.ReadFile(os.Getenv(envPrefix + "INPUTCERTPATH"))
	tInputKeyData, _ := os.ReadFile(os.Getenv(envPrefix + "INPUTKEYPATH"))
	tOutputCertPath := os.Getenv(envPrefix + "OUTPUTCERTPATH")
	tJksAlias := os.Getenv(envPrefix + "JKSALIAS")
	tJksKeypass := os.Getenv(envPrefix + "JKSKEYPASS")
	tJksStorepass := os.Getenv(envPrefix + "JKSSTOREPASS")

	deployer, err := dLocal.New(&dLocal.LocalDeployerConfig{
		OutputFormat:   dLocal.OUTPUT_FORMAT_JKS,
		OutputCertPath: tOutputCertPath,
		JksAlias:       tJksAlias,
		JksKeypass:     tJksKeypass,
		JksStorepass:   tJksStorepass,
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
