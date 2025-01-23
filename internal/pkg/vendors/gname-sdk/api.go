package gnamesdk

type BaseResponse interface {
	GetCode() int
	GetMsg() string
}

type AddDNSRecordRequest struct {
	ZoneName    string `json:"ym"`
	RecordType  string `json:"lx"`
	RecordName  string `json:"zj"`
	RecordValue string `json:"jlz"`
	MX          int    `json:"mx"`
	TTL         int    `json:"ttl"`
}

type AddDNSRecordResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data int    `json:"data"`
}

func (r *AddDNSRecordResponse) GetCode() int {
	return r.Code
}

func (r *AddDNSRecordResponse) GetMsg() string {
	return r.Msg
}

type EditDNSRecordRequest struct {
	ID          string `json:"jxid"`
	ZoneName    string `json:"ym"`
	RecordType  string `json:"lx"`
	RecordName  string `json:"zj"`
	RecordValue string `json:"jlz"`
	MX          int    `json:"mx"`
	TTL         int    `json:"ttl"`
}

type EditDNSRecordResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (r *EditDNSRecordResponse) GetCode() int {
	return r.Code
}

func (r *EditDNSRecordResponse) GetMsg() string {
	return r.Msg
}

type DeleteDNSRecordRequest struct {
	ZoneName string `json:"ym"`
	RecordId int    `json:"jxid"`
}

type DeleteDNSRecordResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (r *DeleteDNSRecordResponse) GetCode() int {
	return r.Code
}

func (r *DeleteDNSRecordResponse) GetMsg() string {
	return r.Msg
}

type ListDNSRecordRequest struct {
	ZoneName string `json:"ym"`
	Page     *int   `json:"page,omitempty"`
	PageSize *int   `json:"limit,omitempty"`
}

type ListDNSRecordResponse struct {
	Code     int          `json:"code"`
	Msg      string       `json:"msg"`
	Count    int          `json:"count"`
	Data     []*DNSRecord `json:"data"`
	Page     int          `json:"page"`
	PageSize int          `json:"pagesize"`
}

type DNSRecord struct {
	ID          string `json:"id"`
	ZoneName    string `json:"ym"`
	RecordType  string `json:"lx"`
	RecordName  string `json:"zjt"`
	RecordValue string `json:"jxz"`
	MX          int    `json:"mx"`
}

func (r *ListDNSRecordResponse) GetCode() int {
	return r.Code
}

func (r *ListDNSRecordResponse) GetMsg() string {
	return r.Msg
}
