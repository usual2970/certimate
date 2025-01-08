package applicant

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"

	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"

	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/pkg/utils/x509"
	"github.com/usual2970/certimate/internal/repository"
)

type acmeUser struct {
	CA           string
	Email        string
	Registration *registration.Resource

	privkey string
}

func newAcmeUser(ca, email string) (*acmeUser, error) {
	repo := repository.NewAcmeAccountRepository()

	applyUser := &acmeUser{
		CA:    ca,
		Email: email,
	}

	acmeAccount, err := repo.GetByCAAndEmail(ca, email)
	if err != nil {
		key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			return nil, err
		}

		keyPEM, err := x509.ConvertECPrivateKeyToPEM(key)
		if err != nil {
			return nil, err
		}

		applyUser.privkey = keyPEM
		return applyUser, nil
	}

	applyUser.Registration = acmeAccount.Resource
	applyUser.privkey = acmeAccount.Key

	return applyUser, nil
}

func (u *acmeUser) GetEmail() string {
	return u.Email
}

func (u acmeUser) GetRegistration() *registration.Resource {
	return u.Registration
}

func (u *acmeUser) GetPrivateKey() crypto.PrivateKey {
	rs, _ := x509.ParseECPrivateKeyFromPEM(u.privkey)
	return rs
}

func (u *acmeUser) hasRegistration() bool {
	return u.Registration != nil
}

func (u *acmeUser) getPrivateKeyPEM() string {
	return u.privkey
}

type acmeAccountRepository interface {
	GetByCAAndEmail(ca, email string) (*domain.AcmeAccount, error)
	Save(ca, email, key string, resource *registration.Resource) error
}

func registerAcmeUser(client *lego.Client, sslProviderConfig *acmeSSLProviderConfig, user *acmeUser) (*registration.Resource, error) {
	// TODO: fix 潜在的并发问题

	var reg *registration.Resource
	var err error
	switch sslProviderConfig.Provider {
	case sslProviderZeroSSL:
		reg, err = client.Registration.RegisterWithExternalAccountBinding(registration.RegisterEABOptions{
			TermsOfServiceAgreed: true,
			Kid:                  sslProviderConfig.Config.ZeroSSL.EabKid,
			HmacEncoded:          sslProviderConfig.Config.ZeroSSL.EabHmacKey,
		})
	case sslProviderGoogleTrustServices:
		reg, err = client.Registration.RegisterWithExternalAccountBinding(registration.RegisterEABOptions{
			TermsOfServiceAgreed: true,
			Kid:                  sslProviderConfig.Config.GoogleTrustServices.EabKid,
			HmacEncoded:          sslProviderConfig.Config.GoogleTrustServices.EabHmacKey,
		})

	case sslProviderLetsEncrypt, sslProviderLetsEncryptStaging:
		reg, err = client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})

	default:
		err = fmt.Errorf("unsupported ssl provider: %s", sslProviderConfig.Provider)
	}

	if err != nil {
		return nil, err
	}

	repo := repository.NewAcmeAccountRepository()

	resp, err := repo.GetByCAAndEmail(sslProviderConfig.Provider, user.GetEmail())
	if err == nil {
		user.privkey = resp.Key
		return resp.Resource, nil
	}

	if err := repo.Save(sslProviderConfig.Provider, user.GetEmail(), user.getPrivateKeyPEM(), reg); err != nil {
		return nil, fmt.Errorf("failed to save registration: %w", err)
	}

	return reg, nil
}
