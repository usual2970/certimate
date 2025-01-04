package applicant

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/challenge/dns01"
	"github.com/go-acme/lego/v4/lego"

	"github.com/usual2970/certimate/internal/app"
	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/repository"
)

type applyConfig struct {
	Domains            string
	ContactEmail       string
	AccessConfig       string
	KeyAlgorithm       string
	Nameservers        string
	PropagationTimeout int32
	DisableFollowCNAME bool
}

type ApplyCertResult struct {
	Certificate       string
	PrivateKey        string
	IssuerCertificate string
	ACMECertUrl       string
	ACMECertStableUrl string
	CSR               string
}

type applicant interface {
	Apply() (*ApplyCertResult, error)
}

func NewWithApplyNode(node *domain.WorkflowNode) (applicant, error) {
	if node.Type != domain.WorkflowNodeTypeApply {
		return nil, fmt.Errorf("node type is not apply")
	}

	accessRepo := repository.NewAccessRepository()
	access, err := accessRepo.GetById(context.Background(), node.GetConfigString("providerAccessId"))
	if err != nil {
		return nil, fmt.Errorf("access record not found: %w", err)
	}

	applyConfig := &applyConfig{
		Domains:            node.GetConfigString("domains"),
		ContactEmail:       node.GetConfigString("contactEmail"),
		AccessConfig:       access.Config,
		KeyAlgorithm:       node.GetConfigString("keyAlgorithm"),
		Nameservers:        node.GetConfigString("nameservers"),
		PropagationTimeout: node.GetConfigInt32("propagationTimeout"),
		DisableFollowCNAME: node.GetConfigBool("disableFollowCNAME"),
	}

	challengeProvider, err := createChallengeProvider(domain.AccessProviderType(access.Provider), access.Config, applyConfig)
	if err != nil {
		return nil, err
	}

	return &proxyApplicant{
		applicant:   challengeProvider,
		applyConfig: applyConfig,
	}, nil
}

func apply(challengeProvider challenge.Provider, applyConfig *applyConfig) (*ApplyCertResult, error) {
	record, _ := app.GetApp().Dao().FindFirstRecordByFilter("settings", "name='sslProvider'")

	sslProvider := &acmeSSLProviderConfig{
		Config:   acmeSSLProviderConfigContent{},
		Provider: defaultSSLProvider,
	}
	if record != nil {
		if err := record.UnmarshalJSONField("content", sslProvider); err != nil {
			return nil, err
		}
	}

	// Some unified lego environment variables are configured here.
	// link: https://github.com/go-acme/lego/issues/1867
	os.Setenv("LEGO_DISABLE_CNAME_SUPPORT", strconv.FormatBool(applyConfig.DisableFollowCNAME))

	myUser, err := newAcmeUser(sslProvider.Provider, applyConfig.ContactEmail)
	if err != nil {
		return nil, err
	}

	config := lego.NewConfig(myUser)

	// This CA URL is configured for a local dev instance of Boulder running in Docker in a VM.
	config.CADirURL = sslProviderUrls[sslProvider.Provider]
	config.Certificate.KeyType = parseKeyAlgorithm(applyConfig.KeyAlgorithm)

	// A client facilitates communication with the CA server.
	client, err := lego.NewClient(config)
	if err != nil {
		return nil, err
	}

	challengeOptions := make([]dns01.ChallengeOption, 0)
	if len(applyConfig.Nameservers) > 0 {
		challengeOptions = append(challengeOptions, dns01.AddRecursiveNameservers(dns01.ParseNameservers(strings.Split(applyConfig.Nameservers, ";"))))
		challengeOptions = append(challengeOptions, dns01.DisableAuthoritativeNssPropagationRequirement())
	}

	client.Challenge.SetDNS01Provider(challengeProvider, challengeOptions...)

	// New users will need to register
	if !myUser.hasRegistration() {
		reg, err := registerAcmeUser(client, sslProvider, myUser)
		if err != nil {
			return nil, fmt.Errorf("failed to register: %w", err)
		}
		myUser.Registration = reg
	}

	request := certificate.ObtainRequest{
		Domains: strings.Split(applyConfig.Domains, ";"),
		Bundle:  true,
	}
	certificates, err := client.Certificate.Obtain(request)
	if err != nil {
		return nil, err
	}

	return &ApplyCertResult{
		PrivateKey:        string(certificates.PrivateKey),
		Certificate:       string(certificates.Certificate),
		IssuerCertificate: string(certificates.IssuerCertificate),
		CSR:               string(certificates.CSR),
		ACMECertUrl:       certificates.CertURL,
		ACMECertStableUrl: certificates.CertStableURL,
	}, nil
}

func parseKeyAlgorithm(algo string) certcrypto.KeyType {
	switch algo {
	case "RSA2048":
		return certcrypto.RSA2048
	case "RSA3072":
		return certcrypto.RSA3072
	case "RSA4096":
		return certcrypto.RSA4096
	case "RSA8192":
		return certcrypto.RSA8192
	case "EC256":
		return certcrypto.EC256
	case "EC384":
		return certcrypto.EC384
	default:
		return certcrypto.RSA2048
	}
}

// TODO: 暂时使用代理模式以兼容之前版本代码，后续重新实现此处逻辑
type proxyApplicant struct {
	applicant   challenge.Provider
	applyConfig *applyConfig
}

func (d *proxyApplicant) Apply() (*ApplyCertResult, error) {
	return apply(d.applicant, d.applyConfig)
}
