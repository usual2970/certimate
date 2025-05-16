package bunny

type AddCustomCertificateRequest struct {
	Hostname       string `json:"Hostname"`
	PullZoneId     string `json:"-"`
	Certificate    string `json:"Certificate"`
	CertificateKey string `json:"CertificateKey"`
}
