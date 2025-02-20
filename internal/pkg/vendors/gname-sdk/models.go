package gnamesdk

type BaseResponse interface {
	GetCode() int
	GetMessage() string
}

type baseResponse struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
}

func (r *baseResponse) GetCode() int {
	return r.Code
}

func (r *baseResponse) GetMessage() string {
	return r.Message
}

type AddDomainResolutionRequest struct {
	ZoneName    string `json:"ym"`
	RecordType  string `json:"lx"`
	RecordName  string `json:"zj"`
	RecordValue string `json:"jlz"`
	MX          int    `json:"mx"`
	TTL         int    `json:"ttl"`
}

type AddDomainResolutionResponse struct {
	baseResponse
	Data string `json:"data"`
}

type ModifyDomainResolutionRequest struct {
	ID          string `json:"jxid"`
	ZoneName    string `json:"ym"`
	RecordType  string `json:"lx"`
	RecordName  string `json:"zj"`
	RecordValue string `json:"jlz"`
	MX          int    `json:"mx"`
	TTL         int    `json:"ttl"`
}

type ModifyDomainResolutionResponse struct {
	baseResponse
}

type DeleteDomainResolutionRequest struct {
	ZoneName string `json:"ym"`
	RecordID string `json:"jxid"`
}

type DeleteDomainResolutionResponse struct {
	baseResponse
}

type ListDomainResolutionRequest struct {
	ZoneName string `json:"ym"`
	Page     *int   `json:"page,omitempty"`
	PageSize *int   `json:"limit,omitempty"`
}

type ListDomainResolutionResponse struct {
	baseResponse
	Count    int                 `json:"count"`
	Data     []*ResolutionRecord `json:"data"`
	Page     int                 `json:"page"`
	PageSize int                 `json:"pagesize"`
}

type ResolutionRecord struct {
	ID          string `json:"id"`
	ZoneName    string `json:"ym"`
	RecordType  string `json:"lx"`
	RecordName  string `json:"zjt"`
	RecordValue string `json:"jxz"`
	MX          int    `json:"mx"`
}
