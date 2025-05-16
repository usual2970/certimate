package applicant

import "github.com/usual2970/certimate/internal/domain"

const (
	sslProviderLetsEncrypt         = string(domain.CAProviderTypeLetsEncrypt)
	sslProviderLetsEncryptStaging  = string(domain.CAProviderTypeLetsEncryptStaging)
	sslProviderBuypass             = string(domain.CAProviderTypeBuypass)
	sslProviderGoogleTrustServices = string(domain.CAProviderTypeGoogleTrustServices)
	sslProviderSSLCom              = string(domain.CAProviderTypeSSLCom)
	sslProviderZeroSSL             = string(domain.CAProviderTypeZeroSSL)

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

func getCAProviderURL(providerType string, config map[string]any) string {
	baseURL, ok := sslProviderUrls[providerType]
	if !ok {
		return sslProviderUrls[sslProviderDefault]
	}

	if config != nil {
		// 如果配置了代理域名，则替换URL
		if proxyDomain, ok := config["proxyDomain"].(string); ok && proxyDomain != "" {
			return proxyDomain
		}
	}
	return baseURL
}

type acmeSSLProviderConfig struct {
	Config   map[domain.CAProviderType]map[string]any `json:"config"`
	Provider string                                   `json:"provider"`
}
