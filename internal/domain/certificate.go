package domain

import "time"

type CertificateSourceType string

const (
	CertificateSourceTypeWorkflow = CertificateSourceType("workflow")
	CertificateSourceTypeUpload   = CertificateSourceType("upload")
)

type Certificate struct {
	Meta
	Source            CertificateSourceType `json:"source" db:"source"`
	SubjectAltNames   string                `json:"subjectAltNames" db:"subjectAltNames"`
	Certificate       string                `json:"certificate" db:"certificate"`
	PrivateKey        string                `json:"privateKey" db:"privateKey"`
	IssuerCertificate string                `json:"issuerCertificate" db:"issuerCertificate"`
	EffectAt          time.Time             `json:"effectAt" db:"effectAt"`
	ExpireAt          time.Time             `json:"expireAt" db:"expireAt"`
	ACMECertUrl       string                `json:"acmeCertUrl" db:"acmeCertUrl"`
	ACMECertStableUrl string                `json:"acmeCertStableUrl" db:"acmeCertStableUrl"`
	WorkflowId        string                `json:"workflowId" db:"workflowId"`
	WorkflowNodeId    string                `json:"workflowNodeId" db:"workflowNodeId"`
	WorkflowOutputId  string                `json:"workflowOutputId" db:"workflowOutputId"`
}
