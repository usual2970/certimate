package applicant

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/challenge/dns01"
	"github.com/go-acme/lego/v4/lego"
	"golang.org/x/exp/slices"
	"golang.org/x/time/rate"

	"github.com/usual2970/certimate/internal/domain"
	uslices "github.com/usual2970/certimate/internal/pkg/utils/slices"
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
	Domains               []string
	ContactEmail          string
	Provider              domain.ApplyDNSProviderType
	ProviderAccessConfig  map[string]any
	ProviderApplyConfig   map[string]any
	KeyAlgorithm          string
	Nameservers           []string
	DnsPropagationTimeout int32
	DnsTTL                int32
	DisableFollowCNAME    bool
	DisableARI            bool
	SkipBeforeExpiryDays  int32
	ReplacedARICertId     string
}

func NewWithApplyNode(node *domain.WorkflowNode) (Applicant, error) {
	if node.Type != domain.WorkflowNodeTypeApply {
		return nil, fmt.Errorf("node type is not apply")
	}

	nodeConfig := node.GetConfigForApply()
	options := &applicantOptions{
		Domains:               uslices.Filter(strings.Split(nodeConfig.Domains, ";"), func(s string) bool { return s != "" }),
		ContactEmail:          nodeConfig.ContactEmail,
		Provider:              domain.ApplyDNSProviderType(nodeConfig.Provider),
		ProviderApplyConfig:   nodeConfig.ProviderConfig,
		KeyAlgorithm:          nodeConfig.KeyAlgorithm,
		Nameservers:           uslices.Filter(strings.Split(nodeConfig.Nameservers, ";"), func(s string) bool { return s != "" }),
		DnsPropagationTimeout: nodeConfig.DnsPropagationTimeout,
		DnsTTL:                nodeConfig.DnsTTL,
		DisableFollowCNAME:    nodeConfig.DisableFollowCNAME,
		DisableARI:            nodeConfig.DisableARI,
		SkipBeforeExpiryDays:  nodeConfig.SkipBeforeExpiryDays,
	}

	accessRepo := repository.NewAccessRepository()
	if access, err := accessRepo.GetById(context.Background(), nodeConfig.ProviderAccessId); err != nil {
		return nil, fmt.Errorf("failed to get access #%s record: %w", nodeConfig.ProviderAccessId, err)
	} else {
		accessConfig, err := access.UnmarshalConfigToMap()
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal access config: %w", err)
		}

		options.ProviderAccessConfig = accessConfig
	}

	certRepo := repository.NewCertificateRepository()
	lastCertificate, _ := certRepo.GetByWorkflowNodeId(context.Background(), node.Id)
	if lastCertificate != nil {
		newCertSan := slices.Clone(options.Domains)
		oldCertSan := strings.Split(lastCertificate.SubjectAltNames, ";")
		slices.Sort(newCertSan)
		slices.Sort(oldCertSan)

		if slices.Equal(newCertSan, oldCertSan) {
			lastCertX509, _ := certcrypto.ParsePEMCertificate([]byte(lastCertificate.Certificate))
			if lastCertX509 != nil {
				replacedARICertId, _ := certificate.MakeARICertID(lastCertX509)
				options.ReplacedARICertId = replacedARICertId
			}
		}
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

	acmeUser, err := newAcmeUser(sslProviderConfig.Provider, options.ContactEmail)
	if err != nil {
		return nil, err
	}

	// Some unified lego environment variables are configured here.
	// link: https://github.com/go-acme/lego/issues/1867
	os.Setenv("LEGO_DISABLE_CNAME_SUPPORT", strconv.FormatBool(options.DisableFollowCNAME))

	// Create an ACME client config
	config := lego.NewConfig(acmeUser)
	config.CADirURL = sslProviderUrls[sslProviderConfig.Provider]
	config.Certificate.KeyType = parseKeyAlgorithm(options.KeyAlgorithm)

	// Create an ACME client
	client, err := lego.NewClient(config)
	if err != nil {
		return nil, err
	}

	// Set the DNS01 challenge provider
	challengeOptions := make([]dns01.ChallengeOption, 0)
	if len(options.Nameservers) > 0 {
		challengeOptions = append(challengeOptions, dns01.AddRecursiveNameservers(dns01.ParseNameservers(options.Nameservers)))
		challengeOptions = append(challengeOptions, dns01.DisableAuthoritativeNssPropagationRequirement())
	}
	client.Challenge.SetDNS01Provider(challengeProvider, challengeOptions...)

	// New users need to register first
	if !acmeUser.hasRegistration() {
		reg, err := registerAcmeUserWithSingleFlight(client, sslProviderConfig, acmeUser)
		if err != nil {
			return nil, fmt.Errorf("failed to register: %w", err)
		}
		acmeUser.Registration = reg
	}

	// Obtain a certificate
	certRequest := certificate.ObtainRequest{
		Domains: options.Domains,
		Bundle:  true,
	}
	if !options.DisableARI {
		certRequest.ReplacesCertID = options.ReplacedARICertId
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
		CSR:                  strings.TrimSpace(string(certResource.CSR)),
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

var limiters sync.Map

const (
	limitBurst         = 300
	limitRate  float64 = float64(1) / float64(36)
)

func getLimiter(key string) *rate.Limiter {
	limiter, _ := limiters.LoadOrStore(key, rate.NewLimiter(rate.Limit(limitRate), 300))
	return limiter.(*rate.Limiter)
}

func (d *proxyApplicant) Apply() (*ApplyCertResult, error) {
	limiter := getLimiter(fmt.Sprintf("apply_%s", d.options.ContactEmail))
	limiter.Wait(context.Background())
	return apply(d.applicant, d.options)
}
