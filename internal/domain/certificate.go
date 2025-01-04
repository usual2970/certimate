package domain

import "time"

var ValidityDuration = time.Hour * 24 * 10

type Certificate struct {
	Meta
	SubjectAltNames   string    `json:"san" db:"san"`
	Certificate       string    `json:"certificate" db:"certificate"`
	PrivateKey        string    `json:"privateKey" db:"privateKey"`
	IssuerCertificate string    `json:"issuerCertificate" db:"issuerCertificate"`
	CertUrl           string    `json:"certUrl" db:"certUrl"`
	CertStableUrl     string    `json:"certStableUrl" db:"certStableUrl"`
	WorkflowId        string    `json:"workflow" db:"workflow"`
	WorkflowNodeId    string    `json:"nodeId" db:"nodeId"`
	WorkflowOutputId  string    `json:"output" db:"output"`
	ExpireAt          time.Time `json:"expireAt" db:"expireAt"`
}

type CertificateMeta struct {
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
