package cdnpro

type apiResponse interface {
	GetCode() string
	GetMessage() string
}

type apiResponseBase struct {
	Code    *string `json:"code,omitempty"`
	Message *string `json:"message,omitempty"`
}

var _ apiResponse = (*apiResponseBase)(nil)

func (r *apiResponseBase) GetCode() string {
	if r.Code == nil {
		return ""
	}

	return *r.Code
}

func (r *apiResponseBase) GetMessage() string {
	if r.Message == nil {
		return ""
	}

	return *r.Message
}

type CertificateVersionInfo struct {
	Comments           *string                               `json:"comments,omitempty"`
	PrivateKey         *string                               `json:"privateKey,omitempty"`
	Certificate        *string                               `json:"certificate,omitempty"`
	ChainCert          *string                               `json:"chainCert,omitempty"`
	IdentificationInfo *CertificateVersionIdentificationInfo `json:"identificationInfo,omitempty"`
}

type CertificateVersionIdentificationInfo struct {
	Country                 *string   `json:"country,omitempty"`
	State                   *string   `json:"state,omitempty"`
	City                    *string   `json:"city,omitempty"`
	Company                 *string   `json:"company,omitempty"`
	Department              *string   `json:"department,omitempty"`
	CommonName              *string   `json:"commonName,omitempty"`
	Email                   *string   `json:"email,omitempty"`
	SubjectAlternativeNames *[]string `json:"subjectAlternativeNames,omitempty"`
}

type HostnamePropertyInfo struct {
	PropertyId    string  `json:"propertyId"`
	Version       int32   `json:"version"`
	CertificateId *string `json:"certificateId,omitempty"`
}

type DeploymentTaskActionInfo struct {
	Action        *string `json:"action,omitempty"`
	PropertyId    *string `json:"propertyId,omitempty"`
	CertificateId *string `json:"certificateId,omitempty"`
	Version       *int32  `json:"version,omitempty"`
}
