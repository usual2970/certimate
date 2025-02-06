package domain

import (
	"crypto/x509"
	"strings"
	"time"

	"github.com/usual2970/certimate/internal/pkg/utils/certs"
)

const CollectionNameCertificate = "certificate"

type Certificate struct {
	Meta
	Source            CertificateSourceType       `json:"source" db:"source"`
	SubjectAltNames   string                      `json:"subjectAltNames" db:"subjectAltNames"`
	SerialNumber      string                      `json:"serialNumber" db:"serialNumber"`
	Certificate       string                      `json:"certificate" db:"certificate"`
	PrivateKey        string                      `json:"privateKey" db:"privateKey"`
	Issuer            string                      `json:"issuer" db:"issuer"`
	IssuerCertificate string                      `json:"issuerCertificate" db:"issuerCertificate"`
	KeyAlgorithm      CertificateKeyAlgorithmType `json:"keyAlgorithm" db:"keyAlgorithm"`
	EffectAt          time.Time                   `json:"effectAt" db:"effectAt"`
	ExpireAt          time.Time                   `json:"expireAt" db:"expireAt"`
	ACMEAccountUrl    string                      `json:"acmeAccountUrl" db:"acmeAccountUrl"`
	ACMECertUrl       string                      `json:"acmeCertUrl" db:"acmeCertUrl"`
	ACMECertStableUrl string                      `json:"acmeCertStableUrl" db:"acmeCertStableUrl"`
	WorkflowId        string                      `json:"workflowId" db:"workflowId"`
	WorkflowNodeId    string                      `json:"workflowNodeId" db:"workflowNodeId"`
	WorkflowOutputId  string                      `json:"workflowOutputId" db:"workflowOutputId"`
	DeletedAt         *time.Time                  `json:"deleted" db:"deleted"`
}

func (c *Certificate) PopulateFromX509(certX509 *x509.Certificate) *Certificate {
	c.SubjectAltNames = strings.Join(certX509.DNSNames, ";")
	c.SerialNumber = strings.ToUpper(certX509.SerialNumber.Text(16))
	c.Issuer = strings.Join(certX509.Issuer.Organization, ";")
	c.EffectAt = certX509.NotBefore
	c.ExpireAt = certX509.NotAfter

	switch certX509.SignatureAlgorithm {
	case x509.SHA256WithRSA, x509.SHA256WithRSAPSS:
		c.KeyAlgorithm = CertificateKeyAlgorithmTypeRSA2048
	case x509.SHA384WithRSA, x509.SHA384WithRSAPSS:
		c.KeyAlgorithm = CertificateKeyAlgorithmTypeRSA3072
	case x509.SHA512WithRSA, x509.SHA512WithRSAPSS:
		c.KeyAlgorithm = CertificateKeyAlgorithmTypeRSA4096
	case x509.ECDSAWithSHA256:
		c.KeyAlgorithm = CertificateKeyAlgorithmTypeEC256
	case x509.ECDSAWithSHA384:
		c.KeyAlgorithm = CertificateKeyAlgorithmTypeEC384
	case x509.ECDSAWithSHA512:
		c.KeyAlgorithm = CertificateKeyAlgorithmTypeEC512
	default:
		c.KeyAlgorithm = CertificateKeyAlgorithmType("")
	}

	return c
}

func (c *Certificate) PopulateFromPEM(certPEM, privkeyPEM string) *Certificate {
	c.Certificate = certPEM
	c.PrivateKey = privkeyPEM

	_, issuerCertPEM, _ := certs.ExtractCertificatesFromPEM(certPEM)
	c.IssuerCertificate = issuerCertPEM

	certX509, _ := certs.ParseCertificateFromPEM(certPEM)
	if certX509 != nil {
		c.PopulateFromX509(certX509)
	}

	return c
}

type CertificateSourceType string

const (
	CertificateSourceTypeWorkflow = CertificateSourceType("workflow")
	CertificateSourceTypeUpload   = CertificateSourceType("upload")
)

type CertificateKeyAlgorithmType string

const (
	CertificateKeyAlgorithmTypeRSA2048 = CertificateKeyAlgorithmType("RSA2048")
	CertificateKeyAlgorithmTypeRSA3072 = CertificateKeyAlgorithmType("RSA3072")
	CertificateKeyAlgorithmTypeRSA4096 = CertificateKeyAlgorithmType("RSA4096")
	CertificateKeyAlgorithmTypeRSA8192 = CertificateKeyAlgorithmType("RSA8192")
	CertificateKeyAlgorithmTypeEC256   = CertificateKeyAlgorithmType("EC256")
	CertificateKeyAlgorithmTypeEC384   = CertificateKeyAlgorithmType("EC384")
	CertificateKeyAlgorithmTypeEC512   = CertificateKeyAlgorithmType("EC512")
)
