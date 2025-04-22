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
	certutil "github.com/usual2970/certimate/internal/pkg/utils/cert"
	maputil "github.com/usual2970/certimate/internal/pkg/utils/map"
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

var registerGroup singleflight.Group

func registerAcmeUserWithSingleFlight(client *lego.Client, user *acmeUser, userRegisterOptions map[string]any) (*registration.Resource, error) {
	resp, err, _ := registerGroup.Do(fmt.Sprintf("register_acme_user_%s_%s", user.CA, user.Email), func() (interface{}, error) {
		return registerAcmeUser(client, user, userRegisterOptions)
	})

	if err != nil {
		return nil, err
	}

	return resp.(*registration.Resource), nil
}

func registerAcmeUser(client *lego.Client, user *acmeUser, userRegisterOptions map[string]any) (*registration.Resource, error) {
	var reg *registration.Resource
	var err error
	switch user.CA {
	case sslProviderLetsEncrypt, sslProviderLetsEncryptStaging:
		reg, err = client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})

	case sslProviderBuypass:
		{
			reg, err = client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
		}

	case sslProviderGoogleTrustServices:
		{
			access := domain.AccessConfigForGoogleTrustServices{}
			if err := maputil.Populate(userRegisterOptions, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			reg, err = client.Registration.RegisterWithExternalAccountBinding(registration.RegisterEABOptions{
				TermsOfServiceAgreed: true,
				Kid:                  access.EabKid,
				HmacEncoded:          access.EabHmacKey,
			})
		}

	case sslProviderSSLCom:
		{
			access := domain.AccessConfigForSSLCom{}
			if err := maputil.Populate(userRegisterOptions, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			reg, err = client.Registration.RegisterWithExternalAccountBinding(registration.RegisterEABOptions{
				TermsOfServiceAgreed: true,
				Kid:                  access.EabKid,
				HmacEncoded:          access.EabHmacKey,
			})
		}

	case sslProviderZeroSSL:
		{
			access := domain.AccessConfigForZeroSSL{}
			if err := maputil.Populate(userRegisterOptions, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			reg, err = client.Registration.RegisterWithExternalAccountBinding(registration.RegisterEABOptions{
				TermsOfServiceAgreed: true,
				Kid:                  access.EabKid,
				HmacEncoded:          access.EabHmacKey,
			})
		}

	default:
		err = fmt.Errorf("unsupported ca provider '%s'", user.CA)
	}
	if err != nil {
		return nil, err
	}

	repo := repository.NewAcmeAccountRepository()
	resp, err := repo.GetByCAAndEmail(user.CA, user.Email)
	if err == nil {
		user.privkey = resp.Key
		return resp.Resource, nil
	}

	if _, err := repo.Save(context.Background(), &domain.AcmeAccount{
		CA:       user.CA,
		Email:    user.Email,
		Key:      user.getPrivateKeyPEM(),
		Resource: reg,
	}); err != nil {
		return nil, fmt.Errorf("failed to save acme account registration: %w", err)
	}

	return reg, nil
}
