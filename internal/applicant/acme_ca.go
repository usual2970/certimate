package applicant

import "github.com/usual2970/certimate/internal/domain"

const (
	sslProviderLetsEncrypt         = string(domain.ApplyCAProviderTypeLetsEncrypt)
	sslProviderLetsEncryptStaging  = string(domain.ApplyCAProviderTypeLetsEncryptStaging)
	sslProviderBuypass             = string(domain.ApplyCAProviderTypeBuypass)
	sslProviderGoogleTrustServices = string(domain.ApplyCAProviderTypeGoogleTrustServices)
	sslProviderSSLCom              = string(domain.ApplyCAProviderTypeSSLCom)
	sslProviderZeroSSL             = string(domain.ApplyCAProviderTypeZeroSSL)

	sslProviderDefault = sslProviderLetsEncrypt
)

var sslProviderUrls = map[string]string{
	sslProviderLetsEncrypt:         "https://acme-v02.api.letsencrypt.org/directory",
	sslProviderLetsEncryptStaging:  "https://acme-staging-v02.api.letsencrypt.org/directory",
	sslProviderBuypass:             "https://api.buypass.com/acme/directory",
	sslProviderGoogleTrustServices: "https://dv.acme-v02.api.pki.goog/directory",
	sslProviderSSLCom:              "https://acme.ssl.com/sslcom-dv-rsa",
	sslProviderSSLCom + "RSA":      "https://acme.ssl.com/sslcom-dv-rsa",
	sslProviderSSLCom + "ECC":      "https://acme.ssl.com/sslcom-dv-ecc",
	sslProviderZeroSSL:             "https://acme.zerossl.com/v2/DV90",
}

type acmeSSLProviderConfig struct {
	Config   map[domain.ApplyCAProviderType]map[string]any `json:"config"`
	Provider string                                        `json:"provider"`
}
