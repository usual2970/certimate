package certificate

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

type CertificateRecord struct {
	CertificateId string `json:"certificate-id"`
	Name          string `json:"name"`
	Comment       string `json:"comment"`
	ValidityFrom  string `json:"certificate-validity-from"`
	ValidityTo    string `json:"certificate-validity-to"`
	Serial        string `json:"certificate-serial"`
}
