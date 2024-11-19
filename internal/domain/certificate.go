package domain

import "time"

type Certificate struct {
	Meta
	SAN               string    `json:"san"`
	Certificate       string    `json:"certificate"`
	PrivateKey        string    `json:"privateKey"`
	IssuerCertificate string    `json:"issuerCertificate"`
	CertUrl           string    `json:"certUrl"`
	CertStableUrl     string    `json:"certStableUrl"`
	Output            string    `json:"output"`
	Workflow          string    `json:"workflow"`
	ExpireAt          time.Time `json:"ExpireAt"`
	NodeId            string    `json:"nodeId"`
}

type MetaData struct {
	Version            string              `json:"version"`
	SerialNumber       string              `json:"serialNumber"`
	Validity           CertificateValidity `json:"validity"`
	SignatureAlgorithm string              `json:"signatureAlgorithm"`
	Issuer             CertificateIssuer   `json:"issuer"`
	Subject            CertificateSubject  `json:"subject"`
}

type CertificateIssuer struct {
	Country      string `json:"country"`
	Organization string `json:"organization"`
	CommonName   string `json:"commonName"`
}

type CertificateSubject struct {
	CN string `json:"CN"`
}

type CertificateValidity struct {
	NotBefore string `json:"notBefore"`
	NotAfter  string `json:"notAfter"`
}
