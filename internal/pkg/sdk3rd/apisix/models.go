package apisix

type UpdateSSLRequest struct {
	ID     string    `json:"-"`
	Cert   *string   `json:"cert,omitempty"`
	Key    *string   `json:"key,omitempty"`
	SNIs   *[]string `json:"snis,omitempty"`
	Type   *string   `json:"type,omitempty"`
	Status *int32    `json:"status,omitempty"`
}

type UpdateSSLResponse struct{}
