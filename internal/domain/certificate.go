package domain

import "time"

type Certificate struct {
	Meta
	Source            string    `json:"source" db:"source"`
	SubjectAltNames   string    `json:"subjectAltNames" db:"subjectAltNames"`
	Certificate       string    `json:"certificate" db:"certificate"`
	PrivateKey        string    `json:"privateKey" db:"privateKey"`
	IssuerCertificate string    `json:"issuerCertificate" db:"issuerCertificate"`
	EffectAt          time.Time `json:"effectAt" db:"effectAt"`
	ExpireAt          time.Time `json:"expireAt" db:"expireAt"`
	AcmeCertUrl       string    `json:"acmeCertUrl" db:"acmeCertUrl"`
	AcmeCertStableUrl string    `json:"acmeCertStableUrl" db:"acmeCertStableUrl"`
	WorkflowId        string    `json:"workflowId" db:"workflowId"`
	WorkflowNodeId    string    `json:"workflowNodeId" db:"workflowNodeId"`
	WorkflowOutputId  string    `json:"workflowOutputId" db:"workflowOutputId"`
}

type CertificateSourceType string

const (
	CERTIFICATE_SOURCE_WORKFLOW = CertificateSourceType("workflow")
	CERTIFICATE_SOURCE_UPLOAD   = CertificateSourceType("upload")
)
