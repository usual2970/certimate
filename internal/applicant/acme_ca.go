package applicant

import "github.com/certimate-go/certimate/internal/domain"

const (
	caLetsEncrypt         = string(domain.CAProviderTypeLetsEncrypt)
	caLetsEncryptStaging  = string(domain.CAProviderTypeLetsEncryptStaging)
	caBuypass             = string(domain.CAProviderTypeBuypass)
	caGoogleTrustServices = string(domain.CAProviderTypeGoogleTrustServices)
	caSSLCom              = string(domain.CAProviderTypeSSLCom)
	caZeroSSL             = string(domain.CAProviderTypeZeroSSL)
	caCustom              = string(domain.CAProviderTypeACMECA)

	caDefault = caLetsEncrypt
)

var caDirUrls = map[string]string{
	caLetsEncrypt:         "https://acme-v02.api.letsencrypt.org/directory",
	caLetsEncryptStaging:  "https://acme-staging-v02.api.letsencrypt.org/directory",
	caBuypass:             "https://api.buypass.com/acme/directory",
	caGoogleTrustServices: "https://dv.acme-v02.api.pki.goog/directory",
	caSSLCom:              "https://acme.ssl.com/sslcom-dv-rsa",
	caSSLCom + "RSA":      "https://acme.ssl.com/sslcom-dv-rsa",
	caSSLCom + "ECC":      "https://acme.ssl.com/sslcom-dv-ecc",
	caZeroSSL:             "https://acme.zerossl.com/v2/DV90",
}

type acmeSSLProviderConfig struct {
	Config   map[domain.CAProviderType]map[string]any `json:"config"`
	Provider string                                   `json:"provider"`
}
