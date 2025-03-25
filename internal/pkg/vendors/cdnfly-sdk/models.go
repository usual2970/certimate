package cdnflysdk

import "encoding/json"

type BaseResponse interface {
	GetCode() string
	GetMessage() string
}

type baseResponse struct {
	Code    json.Number `json:"code"`
	Message string      `json:"msg"`
}

func (r *baseResponse) GetCode() string {
	return r.Code.String()
}

func (r *baseResponse) GetMessage() string {
	return r.Message
}

type GetSiteRequest struct {
	Id string `json:"-"`
}

type GetSiteResponse struct {
	baseResponse
	Data *struct {
		Id          int64  `json:"id"`
		Name        string `json:"name"`
		Domain      string `json:"domain"`
		HttpsListen string `json:"https_listen"`
	} `json:"data,omitempty"`
}

type UpdateSiteRequest struct {
	Id          string  `json:"-"`
	HttpsListen *string `json:"https_listen,omitempty"`
	Enable      *bool   `json:"enable,omitempty"`
}

type UpdateSiteResponse struct {
	baseResponse
}

type CreateCertificateRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"des,omitempty"`
	Type        string  `json:"type"`
	Cert        string  `json:"cert"`
	Key         string  `json:"key"`
}

type CreateCertificateResponse struct {
	baseResponse
	Data string `json:"data"`
}

type UpdateCertificateRequest struct {
	Id          string  `json:"-"`
	Name        *string `json:"name,omitempty"`
	Description *string `json:"des,omitempty"`
	Type        *string `json:"type,omitempty"`
	Cert        *string `json:"cert,omitempty"`
	Key         *string `json:"key,omitempty"`
	Enable      *bool   `json:"enable,omitempty"`
}

type UpdateCertificateResponse struct {
	baseResponse
}
