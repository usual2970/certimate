package applicant

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/challenge/dns01"
	"github.com/go-acme/lego/v4/lego"

	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/pkg/utils/slices"
	"github.com/usual2970/certimate/internal/repository"
)

type ApplyCertResult struct {
	CertificateFullChain string
	IssuerCertificate    string
	PrivateKey           string
	ACMECertUrl          string
	ACMECertStableUrl    string
	CSR                  string
}

type Applicant interface {
	Apply() (*ApplyCertResult, error)
}

type applicantOptions struct {
	Domains              []string
	ContactEmail         string
	Provider             domain.ApplyDNSProviderType
	ProviderAccessConfig map[string]any
	ProviderApplyConfig  map[string]any
	KeyAlgorithm         string
	Nameservers          []string
	PropagationTimeout   int32
	DisableFollowCNAME   bool
}

func NewWithApplyNode(node *domain.WorkflowNode) (Applicant, error) {
	if node.Type != domain.WorkflowNodeTypeApply {
		return nil, fmt.Errorf("node type is not apply")
	}

	accessRepo := repository.NewAccessRepository()
	accessId := node.GetConfigString("providerAccessId")
	access, err := accessRepo.GetById(context.Background(), accessId)
	if err != nil {
		return nil, fmt.Errorf("failed to get access #%s record: %w", accessId, err)
	}

	accessConfig, err := access.UnmarshalConfigToMap()
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal access config: %w", err)
	}

	options := &applicantOptions{
		Domains:              slices.Filter(strings.Split(node.GetConfigString("domains"), ";"), func(s string) bool { return s != "" }),
		ContactEmail:         node.GetConfigString("contactEmail"),
		Provider:             domain.ApplyDNSProviderType(node.GetConfigString("provider")),
		ProviderAccessConfig: accessConfig,
		ProviderApplyConfig:  node.GetConfigMap("providerConfig"),
		KeyAlgorithm:         node.GetConfigString("keyAlgorithm"),
		Nameservers:          slices.Filter(strings.Split(node.GetConfigString("nameservers"), ";"), func(s string) bool { return s != "" }),
		PropagationTimeout:   node.GetConfigInt32("propagationTimeout"),
		DisableFollowCNAME:   node.GetConfigBool("disableFollowCNAME"),
	}

	applicant, err := createApplicant(options)
	if err != nil {
		return nil, err
	}

	return &proxyApplicant{
		applicant: applicant,
		options:   options,
	}, nil
}

func apply(challengeProvider challenge.Provider, options *applicantOptions) (*ApplyCertResult, error) {
	settingsRepo := repository.NewSettingsRepository()
	settings, _ := settingsRepo.GetByName(context.Background(), "sslProvider")

	sslProviderConfig := &acmeSSLProviderConfig{
		Config:   acmeSSLProviderConfigContent{},
		Provider: defaultSSLProvider,
	}
	if settings != nil {
		if err := json.Unmarshal([]byte(settings.Content), sslProviderConfig); err != nil {
			return nil, err
		}
	}

	if sslProviderConfig.Provider == "" {
		sslProviderConfig.Provider = defaultSSLProvider
	}

	myUser, err := newAcmeUser(sslProviderConfig.Provider, options.ContactEmail)
	if err != nil {
		return nil, err
	}

	// Some unified lego environment variables are configured here.
	// link: https://github.com/go-acme/lego/issues/1867
	os.Setenv("LEGO_DISABLE_CNAME_SUPPORT", strconv.FormatBool(options.DisableFollowCNAME))

	config := lego.NewConfig(myUser)

	// This CA URL is configured for a local dev instance of Boulder running in Docker in a VM.
	config.CADirURL = sslProviderUrls[sslProviderConfig.Provider]
	config.Certificate.KeyType = parseKeyAlgorithm(options.KeyAlgorithm)

	// A client facilitates communication with the CA server.
	client, err := lego.NewClient(config)
	if err != nil {
		return nil, err
	}

	challengeOptions := make([]dns01.ChallengeOption, 0)
	if len(options.Nameservers) > 0 {
		challengeOptions = append(challengeOptions, dns01.AddRecursiveNameservers(dns01.ParseNameservers(options.Nameservers)))
		challengeOptions = append(challengeOptions, dns01.DisableAuthoritativeNssPropagationRequirement())
	}

	client.Challenge.SetDNS01Provider(challengeProvider, challengeOptions...)

	// New users will need to register
	if !myUser.hasRegistration() {
		reg, err := registerAcmeUser(client, sslProviderConfig, myUser)
		if err != nil {
			return nil, fmt.Errorf("failed to register: %w", err)
		}
		myUser.Registration = reg
	}

	certRequest := certificate.ObtainRequest{
		Domains: options.Domains,
		Bundle:  true,
	}
	certResource, err := client.Certificate.Obtain(certRequest)
	if err != nil {
		return nil, err
	}

	return &ApplyCertResult{
		CertificateFullChain: strings.TrimSpace(string(certResource.Certificate)),
		IssuerCertificate:    strings.TrimSpace(string(certResource.IssuerCertificate)),
		PrivateKey:           strings.TrimSpace(string(certResource.PrivateKey)),
		ACMECertUrl:          certResource.CertURL,
		ACMECertStableUrl:    certResource.CertStableURL,
		CSR:                  string(certResource.CSR),
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
	applicant challenge.Provider
	options   *applicantOptions
}

func (d *proxyApplicant) Apply() (*ApplyCertResult, error) {
	return apply(d.applicant, d.options)
}
