package applicant

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/challenge/dns01"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"

	"github.com/usual2970/certimate/internal/app"
	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/pkg/utils/x509"
	"github.com/usual2970/certimate/internal/repository"
)

const defaultSSLProvider = "letsencrypt"
const (
	sslProviderLetsencrypt = "letsencrypt"
	sslProviderZeroSSL     = "zerossl"
	sslProviderGts         = "gts"
)

const (
	zerosslUrl     = "https://acme.zerossl.com/v2/DV90"
	letsencryptUrl = "https://acme-v02.api.letsencrypt.org/directory"
	gtsUrl         = "https://dv.acme-v02.api.pki.goog/directory"
)

var sslProviderUrls = map[string]string{
	sslProviderLetsencrypt: letsencryptUrl,
	sslProviderZeroSSL:     zerosslUrl,
	sslProviderGts:         gtsUrl,
}

type Certificate struct {
	CertUrl           string `json:"certUrl"`
	CertStableUrl     string `json:"certStableUrl"`
	PrivateKey        string `json:"privateKey"`
	Certificate       string `json:"certificate"`
	IssuerCertificate string `json:"issuerCertificate"`
	CSR               string `json:"csr"`
}

type applyConfig struct {
	Domains            string `json:"domains"`
	ContactEmail       string `json:"contactEmail"`
	AccessConfig       string `json:"accessConfig"`
	KeyAlgorithm       string `json:"keyAlgorithm"`
	Nameservers        string `json:"nameservers"`
	PropagationTimeout int32  `json:"propagationTimeout"`
	DisableFollowCNAME bool   `json:"disableFollowCNAME"`
}

type applyUser struct {
	CA           string
	Email        string
	Registration *registration.Resource

	privkey string
}

func newApplyUser(ca, email string) (*applyUser, error) {
	repo := getAcmeAccountRepository()

	applyUser := &applyUser{
		CA:    ca,
		Email: email,
	}

	acmeAccount, err := repo.GetByCAAndEmail(ca, email)
	if err != nil {
		key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			return nil, err
		}

		keyStr, err := x509.ConvertECPrivateKeyToPEM(key)
		if err != nil {
			return nil, err
		}

		applyUser.privkey = keyStr
		return applyUser, nil
	}

	applyUser.Registration = acmeAccount.Resource
	applyUser.privkey = acmeAccount.Key

	return applyUser, nil
}

func (u *applyUser) GetEmail() string {
	return u.Email
}

func (u applyUser) GetRegistration() *registration.Resource {
	return u.Registration
}

func (u *applyUser) GetPrivateKey() crypto.PrivateKey {
	rs, _ := x509.ParseECPrivateKeyFromPEM(u.privkey)
	return rs
}

func (u *applyUser) hasRegistration() bool {
	return u.Registration != nil
}

func (u *applyUser) getPrivateKeyString() string {
	return u.privkey
}

type Applicant interface {
	Apply() (*Certificate, error)
}

// TODO: 暂时使用代理模式以兼容之前版本代码，后续重新实现此处逻辑
type proxyApplicant struct {
	applyConfig *applyConfig
	applicant   challenge.Provider
}

func (d *proxyApplicant) Apply() (*Certificate, error) {
	return apply(d.applyConfig, d.applicant)
}

func GetWithApplyNode(node *domain.WorkflowNode) (Applicant, error) {
	// 获取授权配置
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
		applyConfig: applyConfig,
		applicant:   challengeProvider,
	}, nil
}

type SSLProviderConfig struct {
	Config   SSLProviderConfigContent `json:"config"`
	Provider string                   `json:"provider"`
}

type SSLProviderConfigContent struct {
	Zerossl SSLProviderEab `json:"zerossl"`
	Gts     SSLProviderEab `json:"gts"`
}

type SSLProviderEab struct {
	EabHmacKey string `json:"eabHmacKey"`
	EabKid     string `json:"eabKid"`
}

func apply(option *applyConfig, provider challenge.Provider) (*Certificate, error) {
	record, _ := app.GetApp().Dao().FindFirstRecordByFilter("settings", "name='sslProvider'")

	sslProvider := &SSLProviderConfig{
		Config:   SSLProviderConfigContent{},
		Provider: defaultSSLProvider,
	}
	if record != nil {
		if err := record.UnmarshalJSONField("content", sslProvider); err != nil {
			return nil, err
		}
	}

	// Some unified lego environment variables are configured here.
	// link: https://github.com/go-acme/lego/issues/1867
	os.Setenv("LEGO_DISABLE_CNAME_SUPPORT", strconv.FormatBool(option.DisableFollowCNAME))

	myUser, err := newApplyUser(sslProvider.Provider, option.ContactEmail)
	if err != nil {
		return nil, err
	}

	config := lego.NewConfig(myUser)

	// This CA URL is configured for a local dev instance of Boulder running in Docker in a VM.
	config.CADirURL = sslProviderUrls[sslProvider.Provider]
	config.Certificate.KeyType = parseKeyAlgorithm(option.KeyAlgorithm)

	// A client facilitates communication with the CA server.
	client, err := lego.NewClient(config)
	if err != nil {
		return nil, err
	}

	challengeOptions := make([]dns01.ChallengeOption, 0)
	if len(option.Nameservers) > 0 {
		challengeOptions = append(challengeOptions, dns01.AddRecursiveNameservers(dns01.ParseNameservers(strings.Split(option.Nameservers, ";"))))
		challengeOptions = append(challengeOptions, dns01.DisableAuthoritativeNssPropagationRequirement())
	}

	client.Challenge.SetDNS01Provider(provider, challengeOptions...)

	// New users will need to register
	if !myUser.hasRegistration() {
		reg, err := getReg(client, sslProvider, myUser)
		if err != nil {
			return nil, fmt.Errorf("failed to register: %w", err)
		}
		myUser.Registration = reg
	}

	domains := strings.Split(option.Domains, ";")
	request := certificate.ObtainRequest{
		Domains: domains,
		Bundle:  true,
	}
	certificates, err := client.Certificate.Obtain(request)
	if err != nil {
		return nil, err
	}

	return &Certificate{
		CertUrl:           certificates.CertURL,
		CertStableUrl:     certificates.CertStableURL,
		PrivateKey:        string(certificates.PrivateKey),
		Certificate:       string(certificates.Certificate),
		IssuerCertificate: string(certificates.IssuerCertificate),
		CSR:               string(certificates.CSR),
	}, nil
}

type AcmeAccountRepository interface {
	GetByCAAndEmail(ca, email string) (*domain.AcmeAccount, error)
	Save(ca, email, key string, resource *registration.Resource) error
}

func getAcmeAccountRepository() AcmeAccountRepository {
	return repository.NewAcmeAccountRepository()
}

func getReg(client *lego.Client, sslProvider *SSLProviderConfig, user *applyUser) (*registration.Resource, error) {
	// TODO: fix 潜在的并发问题

	var reg *registration.Resource
	var err error
	switch sslProvider.Provider {
	case sslProviderZeroSSL:
		reg, err = client.Registration.RegisterWithExternalAccountBinding(registration.RegisterEABOptions{
			TermsOfServiceAgreed: true,
			Kid:                  sslProvider.Config.Zerossl.EabKid,
			HmacEncoded:          sslProvider.Config.Zerossl.EabHmacKey,
		})
	case sslProviderGts:
		reg, err = client.Registration.RegisterWithExternalAccountBinding(registration.RegisterEABOptions{
			TermsOfServiceAgreed: true,
			Kid:                  sslProvider.Config.Gts.EabKid,
			HmacEncoded:          sslProvider.Config.Gts.EabHmacKey,
		})

	case sslProviderLetsencrypt:
		reg, err = client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})

	default:
		err = errors.New("unknown ssl provider")
	}

	if err != nil {
		return nil, err
	}

	repo := getAcmeAccountRepository()

	resp, err := repo.GetByCAAndEmail(sslProvider.Provider, user.GetEmail())
	if err == nil {
		user.privkey = resp.Key
		return resp.Resource, nil
	}

	if err := repo.Save(sslProvider.Provider, user.GetEmail(), user.getPrivateKeyString(), reg); err != nil {
		return nil, fmt.Errorf("failed to save registration: %w", err)
	}

	return reg, nil
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
