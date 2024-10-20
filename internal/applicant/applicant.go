package applicant

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"fmt"
	"strings"

	"certimate/internal/domain"
	"certimate/internal/utils/app"

	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/challenge/dns01"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
	"github.com/pocketbase/pocketbase/models"
)

const (
	configTypeAliyun      = "aliyun"
	configTypeTencent     = "tencent"
	configTypeHuaweiCloud = "huaweicloud"
	configTypeAws         = "aws"
	configTypeCloudflare  = "cloudflare"
	configTypeNamesilo    = "namesilo"
	configTypeGodaddy     = "godaddy"
	configTypePdns        = "pdns"
	configTypeHttpreq     = "httpreq"
)

const defaultSSLProvider = "letsencrypt"
const (
	sslProviderLetsencrypt = "letsencrypt"
	sslProviderZeroSSL     = "zerossl"
)

const (
	zerosslUrl     = "https://acme.zerossl.com/v2/DV90"
	letsencryptUrl = "https://acme-v02.api.letsencrypt.org/directory"
)

var sslProviderUrls = map[string]string{
	sslProviderLetsencrypt: letsencryptUrl,
	sslProviderZeroSSL:     zerosslUrl,
}

const defaultEmail = "536464346@qq.com"

const defaultTimeout = 60

type Certificate struct {
	CertUrl           string `json:"certUrl"`
	CertStableUrl     string `json:"certStableUrl"`
	PrivateKey        string `json:"privateKey"`
	Certificate       string `json:"certificate"`
	IssuerCertificate string `json:"issuerCertificate"`
	Csr               string `json:"csr"`
}

type ApplyOption struct {
	Email        string `json:"email"`
	Domain       string `json:"domain"`
	Access       string `json:"access"`
	KeyAlgorithm string `json:"keyAlgorithm"`
	Nameservers  string `json:"nameservers"`
	Timeout      int64  `json:"timeout"`
}

type ApplyUser struct {
	Email        string
	Registration *registration.Resource
	key          crypto.PrivateKey
}

func (u *ApplyUser) GetEmail() string {
	return u.Email
}

func (u ApplyUser) GetRegistration() *registration.Resource {
	return u.Registration
}

func (u *ApplyUser) GetPrivateKey() crypto.PrivateKey {
	return u.key
}

type Applicant interface {
	Apply() (*Certificate, error)
}

func Get(record *models.Record) (Applicant, error) {
	if record.GetString("applyConfig") == "" {
		return nil, errors.New("applyConfig is empty")
	}

	applyConfig := &domain.ApplyConfig{}
	record.UnmarshalJSONField("applyConfig", applyConfig)

	access, err := app.GetApp().Dao().FindRecordById("access", applyConfig.Access)
	if err != nil {
		return nil, fmt.Errorf("access record not found: %w", err)
	}

	if applyConfig.Email == "" {
		applyConfig.Email = defaultEmail
	}

	if applyConfig.Timeout == 0 {
		applyConfig.Timeout = defaultTimeout
	}

	option := &ApplyOption{
		Email:        applyConfig.Email,
		Domain:       record.GetString("domain"),
		Access:       access.GetString("config"),
		KeyAlgorithm: applyConfig.KeyAlgorithm,
		Nameservers:  applyConfig.Nameservers,
		Timeout:      applyConfig.Timeout,
	}

	switch access.GetString("configType") {
	case configTypeAliyun:
		return NewAliyun(option), nil
	case configTypeTencent:
		return NewTencent(option), nil
	case configTypeHuaweiCloud:
		return NewHuaweiCloud(option), nil
	case configTypeAws:
		return NewAws(option), nil
	case configTypeCloudflare:
		return NewCloudflare(option), nil
	case configTypeNamesilo:
		return NewNamesilo(option), nil
	case configTypeGodaddy:
		return NewGodaddy(option), nil
	case configTypePdns:
		return NewPdns(option), nil
	case configTypeHttpreq:
		return NewHttpreq(option), nil
	default:
		return nil, errors.New("unknown config type")
	}
}

type SSLProviderConfig struct {
	Config   SSLProviderConfigContent `json:"config"`
	Provider string                   `json:"provider"`
}

type SSLProviderConfigContent struct {
	Zerossl struct {
		EabHmacKey string `json:"eabHmacKey"`
		EabKid     string `json:"eabKid"`
	}
}

func apply(option *ApplyOption, provider challenge.Provider) (*Certificate, error) {
	record, _ := app.GetApp().Dao().FindFirstRecordByFilter("settings", "name='ssl-provider'")

	sslProvider := &SSLProviderConfig{
		Config:   SSLProviderConfigContent{},
		Provider: defaultSSLProvider,
	}
	if record != nil {
		if err := record.UnmarshalJSONField("content", sslProvider); err != nil {
			return nil, err
		}
	}

	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}

	myUser := ApplyUser{
		Email: option.Email,
		key:   privateKey,
	}

	config := lego.NewConfig(&myUser)

	// This CA URL is configured for a local dev instance of Boulder running in Docker in a VM.
	config.CADirURL = sslProviderUrls[sslProvider.Provider]
	config.Certificate.KeyType = parseKeyAlgorithm(option.KeyAlgorithm)

	// A client facilitates communication with the CA server.
	client, err := lego.NewClient(config)
	if err != nil {
		return nil, err
	}

	challengeOptions := make([]dns01.ChallengeOption, 0)
	nameservers := parseNameservers(option.Nameservers)
	if len(nameservers) > 0 {
		challengeOptions = append(challengeOptions, dns01.AddRecursiveNameservers(nameservers))
	}

	client.Challenge.SetDNS01Provider(provider, challengeOptions...)

	// New users will need to register
	reg, err := getReg(client, sslProvider)
	if err != nil {
		return nil, fmt.Errorf("failed to register: %w", err)
	}
	myUser.Registration = reg

	domains := strings.Split(option.Domain, ";")
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
		Csr:               string(certificates.CSR),
	}, nil
}

func getReg(client *lego.Client, sslProvider *SSLProviderConfig) (*registration.Resource, error) {
	var reg *registration.Resource
	var err error
	switch sslProvider.Provider {
	case sslProviderZeroSSL:
		reg, err = client.Registration.RegisterWithExternalAccountBinding(registration.RegisterEABOptions{
			TermsOfServiceAgreed: true,
			Kid:                  sslProvider.Config.Zerossl.EabKid,
			HmacEncoded:          sslProvider.Config.Zerossl.EabHmacKey,
		})

	case sslProviderLetsencrypt:
		reg, err = client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})

	default:
		err = errors.New("unknown ssl provider")
	}

	if err != nil {
		return nil, err
	}

	return reg, nil
}

func parseNameservers(ns string) []string {
	nameservers := make([]string, 0)

	lines := strings.Split(ns, ";")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		nameservers = append(nameservers, line)
	}

	return nameservers
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
