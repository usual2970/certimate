package applicant

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/challenge/dns01"
	"github.com/go-acme/lego/v4/lego"
	"golang.org/x/exp/slices"
	"golang.org/x/time/rate"

	"github.com/usual2970/certimate/internal/domain"
	sliceutil "github.com/usual2970/certimate/internal/pkg/utils/slice"
	"github.com/usual2970/certimate/internal/repository"
)

type ApplyResult struct {
	CertificateFullChain string
	IssuerCertificate    string
	PrivateKey           string
	ACMEAccountUrl       string
	ACMECertUrl          string
	ACMECertStableUrl    string
	CSR                  string
}

type Applicant interface {
	Apply(ctx context.Context) (*ApplyResult, error)
}

type ApplicantWithWorkflowNodeConfig struct {
	Node   *domain.WorkflowNode
	Logger *slog.Logger
}

func NewWithWorkflowNode(config ApplicantWithWorkflowNodeConfig) (Applicant, error) {
	if config.Node == nil {
		return nil, fmt.Errorf("node is nil")
	}
	if config.Node.Type != domain.WorkflowNodeTypeApply {
		return nil, fmt.Errorf("node type is not '%s'", string(domain.WorkflowNodeTypeApply))
	}

	nodeConfig := config.Node.GetConfigForApply()
	options := &applicantProviderOptions{
		Domains:                  sliceutil.Filter(strings.Split(nodeConfig.Domains, ";"), func(s string) bool { return s != "" }),
		ContactEmail:             nodeConfig.ContactEmail,
		Provider:                 domain.ACMEDns01ProviderType(nodeConfig.Provider),
		ProviderAccessConfig:     make(map[string]any),
		ProviderExtendedConfig:   nodeConfig.ProviderConfig,
		CAProvider:               domain.CAProviderType(nodeConfig.CAProvider),
		CAProviderAccessConfig:   make(map[string]any),
		CAProviderExtendedConfig: nodeConfig.CAProviderConfig,
		KeyAlgorithm:             nodeConfig.KeyAlgorithm,
		Nameservers:              sliceutil.Filter(strings.Split(nodeConfig.Nameservers, ";"), func(s string) bool { return s != "" }),
		DnsPropagationWait:       nodeConfig.DnsPropagationWait,
		DnsPropagationTimeout:    nodeConfig.DnsPropagationTimeout,
		DnsTTL:                   nodeConfig.DnsTTL,
		DisableFollowCNAME:       nodeConfig.DisableFollowCNAME,
	}

	accessRepo := repository.NewAccessRepository()
	if nodeConfig.ProviderAccessId != "" {
		if access, err := accessRepo.GetById(context.Background(), nodeConfig.ProviderAccessId); err != nil {
			return nil, fmt.Errorf("failed to get access #%s record: %w", nodeConfig.ProviderAccessId, err)
		} else {
			options.ProviderAccessConfig = access.Config
		}
	}
	if nodeConfig.CAProviderAccessId != "" {
		if access, err := accessRepo.GetById(context.Background(), nodeConfig.CAProviderAccessId); err != nil {
			return nil, fmt.Errorf("failed to get access #%s record: %w", nodeConfig.CAProviderAccessId, err)
		} else {
			options.CAProviderAccessConfig = access.Config
		}
	}

	settingsRepo := repository.NewSettingsRepository()
	if string(options.CAProvider) == "" {
		settings, _ := settingsRepo.GetByName(context.Background(), "sslProvider")

		sslProviderConfig := &acmeSSLProviderConfig{
			Config:   make(map[domain.CAProviderType]map[string]any),
			Provider: sslProviderDefault,
		}
		if settings != nil {
			if err := json.Unmarshal([]byte(settings.Content), sslProviderConfig); err != nil {
				return nil, err
			} else if sslProviderConfig.Provider == "" {
				sslProviderConfig.Provider = sslProviderDefault
			}
		}

		options.CAProvider = domain.CAProviderType(sslProviderConfig.Provider)
		options.CAProviderAccessConfig = sslProviderConfig.Config[options.CAProvider]
	}

	certRepo := repository.NewCertificateRepository()
	lastCertificate, _ := certRepo.GetByWorkflowNodeId(context.Background(), config.Node.Id)
	if lastCertificate != nil {
		newCertSan := slices.Clone(options.Domains)
		oldCertSan := strings.Split(lastCertificate.SubjectAltNames, ";")
		slices.Sort(newCertSan)
		slices.Sort(oldCertSan)

		if slices.Equal(newCertSan, oldCertSan) {
			lastCertX509, _ := certcrypto.ParsePEMCertificate([]byte(lastCertificate.Certificate))
			if lastCertX509 != nil {
				replacedARICertId, _ := certificate.MakeARICertID(lastCertX509)
				options.ReplacedARIAcct = lastCertificate.ACMEAccountUrl
				options.ReplacedARICert = replacedARICertId
			}
		}
	}

	applicant, err := createApplicantProvider(options)
	if err != nil {
		return nil, err
	}

	return &applicantImpl{
		applicant: applicant,
		options:   options,
	}, nil
}

type applicantImpl struct {
	applicant challenge.Provider
	options   *applicantProviderOptions
}

var _ Applicant = (*applicantImpl)(nil)

func (d *applicantImpl) Apply(ctx context.Context) (*ApplyResult, error) {
	limiter := getLimiter(fmt.Sprintf("apply_%s", d.options.ContactEmail))
	if err := limiter.Wait(ctx); err != nil {
		return nil, err
	}

	return applyUseLego(d.applicant, d.options)
}

const (
	limitBurst         = 300
	limitRate  float64 = float64(1) / float64(36)
)

var limiters sync.Map

func getLimiter(key string) *rate.Limiter {
	limiter, _ := limiters.LoadOrStore(key, rate.NewLimiter(rate.Limit(limitRate), 300))
	return limiter.(*rate.Limiter)
}

func applyUseLego(legoProvider challenge.Provider, options *applicantProviderOptions) (*ApplyResult, error) {
	user, err := newAcmeUser(string(options.CAProvider), options.ContactEmail)
	if err != nil {
		return nil, err
	}

	// Some unified lego environment variables are configured here.
	// link: https://github.com/go-acme/lego/issues/1867
	os.Setenv("LEGO_DISABLE_CNAME_SUPPORT", strconv.FormatBool(options.DisableFollowCNAME))

	// Create an ACME client config
	config := lego.NewConfig(user)
	config.Certificate.KeyType = parseLegoKeyAlgorithm(domain.CertificateKeyAlgorithmType(options.KeyAlgorithm))
	config.CADirURL = sslProviderUrls[user.CA]
	if user.CA == sslProviderSSLCom {
		if strings.HasPrefix(options.KeyAlgorithm, "RSA") {
			config.CADirURL = sslProviderUrls[sslProviderSSLCom+"RSA"]
		} else if strings.HasPrefix(options.KeyAlgorithm, "EC") {
			config.CADirURL = sslProviderUrls[sslProviderSSLCom+"ECC"]
		}
	}

	// Create an ACME client
	client, err := lego.NewClient(config)
	if err != nil {
		return nil, err
	}

	// Set the DNS01 challenge provider
	client.Challenge.SetDNS01Provider(legoProvider,
		dns01.CondOption(
			len(options.Nameservers) > 0,
			dns01.AddRecursiveNameservers(dns01.ParseNameservers(options.Nameservers)),
		),
		dns01.CondOption(
			options.DnsPropagationWait > 0,
			dns01.PropagationWait(time.Duration(options.DnsPropagationWait)*time.Second, true),
		),
		dns01.CondOption(
			len(options.Nameservers) > 0 || options.DnsPropagationWait > 0,
			dns01.DisableAuthoritativeNssPropagationRequirement(),
		),
	)

	// New users need to register first
	if !user.hasRegistration() {
		reg, err := registerAcmeUserWithSingleFlight(client, user, options.CAProviderAccessConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to register acme user: %w", err)
		}
		user.Registration = reg
	}

	// Obtain a certificate
	certRequest := certificate.ObtainRequest{
		Domains: options.Domains,
		Bundle:  true,
	}
	if options.ReplacedARIAcct == user.Registration.URI {
		certRequest.ReplacesCertID = options.ReplacedARICert
	}
	certResource, err := client.Certificate.Obtain(certRequest)
	if err != nil {
		return nil, err
	}

	return &ApplyResult{
		CertificateFullChain: strings.TrimSpace(string(certResource.Certificate)),
		IssuerCertificate:    strings.TrimSpace(string(certResource.IssuerCertificate)),
		PrivateKey:           strings.TrimSpace(string(certResource.PrivateKey)),
		ACMEAccountUrl:       user.Registration.URI,
		ACMECertUrl:          certResource.CertURL,
		ACMECertStableUrl:    certResource.CertStableURL,
		CSR:                  strings.TrimSpace(string(certResource.CSR)),
	}, nil
}

func parseLegoKeyAlgorithm(algo domain.CertificateKeyAlgorithmType) certcrypto.KeyType {
	alogMap := map[domain.CertificateKeyAlgorithmType]certcrypto.KeyType{
		domain.CertificateKeyAlgorithmTypeRSA2048: certcrypto.RSA2048,
		domain.CertificateKeyAlgorithmTypeRSA3072: certcrypto.RSA3072,
		domain.CertificateKeyAlgorithmTypeRSA4096: certcrypto.RSA4096,
		domain.CertificateKeyAlgorithmTypeRSA8192: certcrypto.RSA8192,
		domain.CertificateKeyAlgorithmTypeEC256:   certcrypto.EC256,
		domain.CertificateKeyAlgorithmTypeEC384:   certcrypto.EC384,
		domain.CertificateKeyAlgorithmTypeEC512:   certcrypto.KeyType("P512"),
	}

	if keyType, ok := alogMap[algo]; ok {
		return keyType
	}

	return certcrypto.RSA2048
}
