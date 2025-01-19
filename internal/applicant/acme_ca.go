package applicant

const (
	sslProviderLetsEncrypt         = "letsencrypt"
	sslProviderLetsEncryptStaging  = "letsencrypt_staging"
	sslProviderZeroSSL             = "zerossl"
	sslProviderGoogleTrustServices = "gts"
)
const defaultSSLProvider = sslProviderLetsEncrypt

const (
	letsencryptUrl        = "https://acme-v02.api.letsencrypt.org/directory"
	letsencryptStagingUrl = "https://acme-staging-v02.api.letsencrypt.org/directory"
	zerosslUrl            = "https://acme.zerossl.com/v2/DV90"
	gtsUrl                = "https://dv.acme-v02.api.pki.goog/directory"
)

var sslProviderUrls = map[string]string{
	sslProviderLetsEncrypt:         letsencryptUrl,
	sslProviderLetsEncryptStaging:  letsencryptStagingUrl,
	sslProviderZeroSSL:             zerosslUrl,
	sslProviderGoogleTrustServices: gtsUrl,
}

type acmeSSLProviderConfig struct {
	Config   acmeSSLProviderConfigContent `json:"config"`
	Provider string                       `json:"provider"`
}

type acmeSSLProviderConfigContent struct {
	ZeroSSL             acmeSSLProviderEabConfig `json:"zerossl"`
	GoogleTrustServices acmeSSLProviderEabConfig `json:"gts"`
}

type acmeSSLProviderEabConfig struct {
	EabHmacKey string `json:"eabHmacKey"`
	EabKid     string `json:"eabKid"`
}
