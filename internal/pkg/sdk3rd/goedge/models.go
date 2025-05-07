package goedge

type BaseResponse interface {
	GetCode() int32
	GetMessage() string
}

type baseResponse struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

func (r *baseResponse) GetCode() int32 {
	return r.Code
}

func (r *baseResponse) GetMessage() string {
	return r.Message
}

type getAPIAccessTokenRequest struct {
	Type        string `json:"type"`
	AccessKeyId string `json:"accessKeyId"`
	AccessKey   string `json:"accessKey"`
}

type getAPIAccessTokenResponse struct {
	baseResponse
	Data *struct {
		Token     string `json:"token"`
		ExpiresAt int64  `json:"expiresAt"`
	} `json:"data,omitempty"`
}

type UpdateSSLCertRequest struct {
	SSLCertId   int64    `json:"sslCertId"`
	IsOn        bool     `json:"isOn"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	ServerName  string   `json:"serverName"`
	IsCA        bool     `json:"isCA"`
	CertData    string   `json:"certData"`
	KeyData     string   `json:"keyData"`
	TimeBeginAt int64    `json:"timeBeginAt"`
	TimeEndAt   int64    `json:"timeEndAt"`
	DNSNames    []string `json:"dnsNames"`
	CommonNames []string `json:"commonNames"`
}

type UpdateSSLCertResponse struct {
	baseResponse
}
