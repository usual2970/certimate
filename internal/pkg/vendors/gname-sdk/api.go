package gnamesdk

type BaseResponse interface {
	GetCode() int
	GetMsg() string
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
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data int    `json:"data"`
}

func (r *AddDomainResolutionResponse) GetCode() int {
	return r.Code
}

func (r *AddDomainResolutionResponse) GetMsg() string {
	return r.Msg
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
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (r *ModifyDomainResolutionResponse) GetCode() int {
	return r.Code
}

func (r *ModifyDomainResolutionResponse) GetMsg() string {
	return r.Msg
}

type DeleteDomainResolutionRequest struct {
	ZoneName string `json:"ym"`
	RecordID string `json:"jxid"`
}

type DeleteDomainResolutionResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (r *DeleteDomainResolutionResponse) GetCode() int {
	return r.Code
}

func (r *DeleteDomainResolutionResponse) GetMsg() string {
	return r.Msg
}

type ListDomainResolutionRequest struct {
	ZoneName string `json:"ym"`
	Page     *int   `json:"page,omitempty"`
	PageSize *int   `json:"limit,omitempty"`
}

type ListDomainResolutionResponse struct {
	Code     int                 `json:"code"`
	Msg      string              `json:"msg"`
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

func (r *ListDomainResolutionResponse) GetCode() int {
	return r.Code
}

func (r *ListDomainResolutionResponse) GetMsg() string {
	return r.Msg
}
