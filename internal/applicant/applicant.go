package applicant

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"strings"

	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/challenge/dns01"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
	"github.com/pocketbase/pocketbase/models"
)

const (
	configTypeTencent    = "tencent"
	configTypeAliyun     = "aliyun"
	configTypeCloudflare = "cloudflare"
	configTypeNamesilo   = "namesilo"
	configTypeGodaddy    = "godaddy"
)

const defaultEmail = "536464346@qq.com"

type Certificate struct {
	CertUrl           string `json:"certUrl"`
	CertStableUrl     string `json:"certStableUrl"`
	PrivateKey        string `json:"privateKey"`
	Certificate       string `json:"certificate"`
	IssuerCertificate string `json:"issuerCertificate"`
	Csr               string `json:"csr"`
}

type ApplyOption struct {
	Email       string `json:"email"`
	Domain      string `json:"domain"`
	Access      string `json:"access"`
	Nameservers string `json:"nameservers"`
}

type MyUser struct {
	Email        string
	Registration *registration.Resource
	key          crypto.PrivateKey
}

func (u *MyUser) GetEmail() string {
	return u.Email
}
func (u MyUser) GetRegistration() *registration.Resource {
	return u.Registration
}
func (u *MyUser) GetPrivateKey() crypto.PrivateKey {
	return u.key
}

type Applicant interface {
	Apply() (*Certificate, error)
}

func Get(record *models.Record) (Applicant, error) {
	access := record.ExpandedOne("access")
	email := record.GetString("email")
	if email == "" {
		email = defaultEmail
	}
	option := &ApplyOption{
		Email:       email,
		Domain:      record.GetString("domain"),
		Access:      access.GetString("config"),
		Nameservers: record.GetString("nameservers"),
	}
	switch access.GetString("configType") {
	case configTypeTencent:
		return NewTencent(option), nil
	case configTypeAliyun:
		return NewAliyun(option), nil
	case configTypeCloudflare:
		return NewCloudflare(option), nil
	case configTypeNamesilo:
		return NewNamesilo(option), nil
	case configTypeGodaddy:
		return NewGodaddy(option), nil
	default:
		return nil, errors.New("unknown config type")
	}

}

func apply(option *ApplyOption, provider challenge.Provider) (*Certificate, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}

	myUser := MyUser{
		Email: option.Email,
		key:   privateKey,
	}

	config := lego.NewConfig(&myUser)

	// This CA URL is configured for a local dev instance of Boulder running in Docker in a VM.
	config.CADirURL = "https://acme-v02.api.letsencrypt.org/directory"
	config.Certificate.KeyType = certcrypto.RSA2048

	// A client facilitates communication with the CA server.
	client, err := lego.NewClient(config)
	if err != nil {
		return nil, err
	}

	challengeOptions := make([]dns01.ChallengeOption, 0)
	nameservers := ParseNameservers(option.Nameservers)
	if len(nameservers) > 0 {
		challengeOptions = append(challengeOptions, dns01.AddRecursiveNameservers(nameservers))
	}

	client.Challenge.SetDNS01Provider(provider, challengeOptions...)

	// New users will need to register
	reg, err := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	if err != nil {
		return nil, err
	}
	myUser.Registration = reg

	domains := []string{option.Domain}

	// 如果是通配置符域名，把根域名也加入
	if strings.HasPrefix(option.Domain, "*.") && len(strings.Split(option.Domain, ".")) == 3 {
		rootDomain := strings.TrimPrefix(option.Domain, "*.")
		domains = append(domains, rootDomain)
	}

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

func ParseNameservers(ns string) []string {
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
