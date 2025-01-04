package applicant

const defaultSSLProvider = "letsencrypt"
const (
	sslProviderLetsencrypt = "letsencrypt"
	sslProviderZeroSSL     = "zerossl"
	sslProviderGts         = "gts"
)

const (
	zerosslUrl     = "https://acme.zerossl.com/v2/DV90"
	letsencryptUrl = "https://acme-v02.api.letsencrypt.org/directory"
	gtsUrl         = "https://dv.acme-v02.api.pki.goog/directory"
)

var sslProviderUrls = map[string]string{
	sslProviderLetsencrypt: letsencryptUrl,
	sslProviderZeroSSL:     zerosslUrl,
	sslProviderGts:         gtsUrl,
}

type acmeSSLProviderConfig struct {
	Config   acmeSSLProviderConfigContent `json:"config"`
	Provider string                       `json:"provider"`
}

type acmeSSLProviderConfigContent struct {
	Zerossl acmeSSLProviderEab `json:"zerossl"`
	Gts     acmeSSLProviderEab `json:"gts"`
}

type acmeSSLProviderEab struct {
	EabHmacKey string `json:"eabHmacKey"`
	EabKid     string `json:"eabKid"`
}
