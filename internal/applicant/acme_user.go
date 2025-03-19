package applicant

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"

	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
	"golang.org/x/sync/singleflight"

	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/pkg/utils/certutil"
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

		keyPEM, err := certutil.ConvertECPrivateKeyToPEM(key)
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
	rs, _ := certutil.ParseECPrivateKeyFromPEM(u.privkey)
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

var registerGroup singleflight.Group

func registerAcmeUserWithSingleFlight(client *lego.Client, sslProviderConfig *acmeSSLProviderConfig, user *acmeUser) (*registration.Resource, error) {
	resp, err, _ := registerGroup.Do(fmt.Sprintf("register_acme_user_%s_%s", sslProviderConfig.Provider, user.GetEmail()), func() (interface{}, error) {
		return registerAcmeUser(client, sslProviderConfig, user)
	})

	if err != nil {
		return nil, err
	}

	return resp.(*registration.Resource), nil
}

func registerAcmeUser(client *lego.Client, sslProviderConfig *acmeSSLProviderConfig, user *acmeUser) (*registration.Resource, error) {
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

	if _, err := repo.Save(context.Background(), &domain.AcmeAccount{
		CA:       sslProviderConfig.Provider,
		Email:    user.GetEmail(),
		Key:      user.getPrivateKeyPEM(),
		Resource: reg,
	}); err != nil {
		return nil, fmt.Errorf("failed to save registration: %w", err)
	}

	return reg, nil
}
