package domain

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"strings"
	"time"

	certutil "github.com/usual2970/certimate/internal/pkg/utils/cert"
)

const CollectionNameCertificate = "certificate"

type Certificate struct {
	Meta
	Source            CertificateSourceType       `json:"source" db:"source"`
	SubjectAltNames   string                      `json:"subjectAltNames" db:"subjectAltNames"`
	SerialNumber      string                      `json:"serialNumber" db:"serialNumber"`
	Certificate       string                      `json:"certificate" db:"certificate"`
	PrivateKey        string                      `json:"privateKey" db:"privateKey"`
	IssuerOrg         string                      `json:"issuerOrg" db:"issuerOrg"`
	IssuerCertificate string                      `json:"issuerCertificate" db:"issuerCertificate"`
	KeyAlgorithm      CertificateKeyAlgorithmType `json:"keyAlgorithm" db:"keyAlgorithm"`
	EffectAt          time.Time                   `json:"effectAt" db:"effectAt"`
	ExpireAt          time.Time                   `json:"expireAt" db:"expireAt"`
	ACMEAccountUrl    string                      `json:"acmeAccountUrl" db:"acmeAccountUrl"`
	ACMECertUrl       string                      `json:"acmeCertUrl" db:"acmeCertUrl"`
	ACMECertStableUrl string                      `json:"acmeCertStableUrl" db:"acmeCertStableUrl"`
	WorkflowId        string                      `json:"workflowId" db:"workflowId"`
	WorkflowNodeId    string                      `json:"workflowNodeId" db:"workflowNodeId"`
	WorkflowRunId     string                      `json:"workflowRunId" db:"workflowRunId"`
	WorkflowOutputId  string                      `json:"workflowOutputId" db:"workflowOutputId"`
	DeletedAt         *time.Time                  `json:"deleted" db:"deleted"`
}

func (c *Certificate) PopulateFromX509(certX509 *x509.Certificate) *Certificate {
	c.SubjectAltNames = strings.Join(certX509.DNSNames, ";")
	c.SerialNumber = strings.ToUpper(certX509.SerialNumber.Text(16))
	c.IssuerOrg = strings.Join(certX509.Issuer.Organization, ";")
	c.EffectAt = certX509.NotBefore
	c.ExpireAt = certX509.NotAfter

	switch certX509.PublicKeyAlgorithm {
	case x509.RSA:
		{
			len := 0
			if pubkey, ok := certX509.PublicKey.(*rsa.PublicKey); ok {
				len = pubkey.N.BitLen()
			}

			switch len {
			case 0:
				c.KeyAlgorithm = CertificateKeyAlgorithmType("RSA")
			case 2048:
				c.KeyAlgorithm = CertificateKeyAlgorithmTypeRSA2048
			case 3072:
				c.KeyAlgorithm = CertificateKeyAlgorithmTypeRSA3072
			case 4096:
				c.KeyAlgorithm = CertificateKeyAlgorithmTypeRSA4096
			case 8192:
				c.KeyAlgorithm = CertificateKeyAlgorithmTypeRSA8192
			default:
				c.KeyAlgorithm = CertificateKeyAlgorithmType(fmt.Sprintf("RSA%d", len))
			}
		}

	case x509.ECDSA:
		{
			len := 0
			if pubkey, ok := certX509.PublicKey.(*ecdsa.PublicKey); ok {
				if pubkey.Curve != nil && pubkey.Curve.Params() != nil {
					len = pubkey.Curve.Params().BitSize
				}
			}

			switch len {
			case 0:
				c.KeyAlgorithm = CertificateKeyAlgorithmType("EC")
			case 256:
				c.KeyAlgorithm = CertificateKeyAlgorithmTypeEC256
			case 384:
				c.KeyAlgorithm = CertificateKeyAlgorithmTypeEC384
			case 521:
				c.KeyAlgorithm = CertificateKeyAlgorithmTypeEC512
			default:
				c.KeyAlgorithm = CertificateKeyAlgorithmType(fmt.Sprintf("EC%d", len))
			}
		}

	case x509.Ed25519:
		{
			c.KeyAlgorithm = CertificateKeyAlgorithmType("ED25519")
		}

	default:
		c.KeyAlgorithm = CertificateKeyAlgorithmType("")
	}

	return c
}

func (c *Certificate) PopulateFromPEM(certPEM, privkeyPEM string) *Certificate {
	c.Certificate = certPEM
	c.PrivateKey = privkeyPEM

	_, issuerCertPEM, _ := certutil.ExtractCertificatesFromPEM(certPEM)
	c.IssuerCertificate = issuerCertPEM

	certX509, _ := certutil.ParseCertificateFromPEM(certPEM)
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
