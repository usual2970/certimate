package applicant

import "github.com/usual2970/certimate/internal/domain"

const (
	sslProviderLetsEncrypt         = string(domain.ApplyCAProviderTypeLetsEncrypt)
	sslProviderLetsEncryptStaging  = string(domain.ApplyCAProviderTypeLetsEncryptStaging)
	sslProviderGoogleTrustServices = string(domain.ApplyCAProviderTypeGoogleTrustServices)
	sslProviderZeroSSL             = string(domain.ApplyCAProviderTypeZeroSSL)

	sslProviderDefault = sslProviderLetsEncrypt
)

var sslProviderUrls = map[string]string{
	sslProviderLetsEncrypt:         "https://acme-v02.api.letsencrypt.org/directory",
	sslProviderLetsEncryptStaging:  "https://acme-staging-v02.api.letsencrypt.org/directory",
	sslProviderGoogleTrustServices: "https://dv.acme-v02.api.pki.goog/directory",
	sslProviderZeroSSL:             "https://acme.zerossl.com/v2/DV90",
}

type acmeSSLProviderConfig struct {
	Config   map[domain.ApplyCAProviderType]map[string]any `json:"config"`
	Provider string                                        `json:"provider"`
}
