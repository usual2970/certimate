package applicant

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"strings"

	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
	"golang.org/x/sync/singleflight"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/internal/repository"
	xcert "github.com/certimate-go/certimate/pkg/utils/cert"
	xmaps "github.com/certimate-go/certimate/pkg/utils/maps"
)

type acmeUser struct {
	// 证书颁发机构标识。
	// 通常等同于 [CAProviderType] 的值。
	// 对于自定义 ACME CA，值为 "custom#{access_id}"。
	CA string
	// 邮箱。
	Email string
	// 注册信息。
	Registration *registration.Resource

	// CSR 私钥。
	privkey string
}

func newAcmeUser(ca, caAccessId, email string) (*acmeUser, error) {
	repo := repository.NewAcmeAccountRepository()

	applyUser := &acmeUser{
		CA:    ca,
		Email: email,
	}
	if ca == caCustom {
		applyUser.CA = fmt.Sprintf("%s#%s", ca, caAccessId)
	}

	acmeAccount, err := repo.GetByCAAndEmail(applyUser.CA, applyUser.Email)
	if err != nil {
		key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			return nil, err
		}

		keyPEM, err := xcert.ConvertECPrivateKeyToPEM(key)
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
	rs, _ := xcert.ParseECPrivateKeyFromPEM(u.privkey)
	return rs
}

func (u *acmeUser) hasRegistration() bool {
	return u.Registration != nil
}

func (u *acmeUser) getCAProvider() string {
	return strings.Split(u.CA, "#")[0]
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
	switch user.getCAProvider() {
	case caLetsEncrypt, caLetsEncryptStaging:
		reg, err = client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})

	case caBuypass:
		{
			reg, err = client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
		}

	case caGoogleTrustServices:
		{
			access := domain.AccessConfigForGoogleTrustServices{}
			if err := xmaps.Populate(userRegisterOptions, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			reg, err = client.Registration.RegisterWithExternalAccountBinding(registration.RegisterEABOptions{
				TermsOfServiceAgreed: true,
				Kid:                  access.EabKid,
				HmacEncoded:          access.EabHmacKey,
			})
		}

	case caSSLCom:
		{
			access := domain.AccessConfigForSSLCom{}
			if err := xmaps.Populate(userRegisterOptions, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			reg, err = client.Registration.RegisterWithExternalAccountBinding(registration.RegisterEABOptions{
				TermsOfServiceAgreed: true,
				Kid:                  access.EabKid,
				HmacEncoded:          access.EabHmacKey,
			})
		}

	case caZeroSSL:
		{
			access := domain.AccessConfigForZeroSSL{}
			if err := xmaps.Populate(userRegisterOptions, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			reg, err = client.Registration.RegisterWithExternalAccountBinding(registration.RegisterEABOptions{
				TermsOfServiceAgreed: true,
				Kid:                  access.EabKid,
				HmacEncoded:          access.EabHmacKey,
			})
		}

	case caCustom:
		{
			access := domain.AccessConfigForACMECA{}
			if err := xmaps.Populate(userRegisterOptions, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			if access.EabKid == "" && access.EabHmacKey == "" {
				reg, err = client.Registration.Register(registration.RegisterOptions{
					TermsOfServiceAgreed: true,
				})
			} else {
				reg, err = client.Registration.RegisterWithExternalAccountBinding(registration.RegisterEABOptions{
					TermsOfServiceAgreed: true,
					Kid:                  access.EabKid,
					HmacEncoded:          access.EabHmacKey,
				})
			}
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
