package dnslasdk

type BaseResponse interface {
	GetCode() int32
	GetMessage() string
}

type baseResponse struct {
	Code    *int32  `json:"code,omitempty"`
	Message *string `json:"message,omitempty"`
}

func (r *baseResponse) GetCode() int32 {
	if r.Code != nil {
		return *r.Code
	}
	return 0
}

func (r *baseResponse) GetMessage() string {
	if r.Message != nil {
		return *r.Message
	}
	return ""
}

type DomainInfo struct {
	Id            string `json:"id"`
	GroupId       string `json:"groupId"`
	GroupName     string `json:"groupName"`
	Domain        string `json:"domain"`
	DisplayDomain string `json:"displayDomain"`
	CreatedAt     int64  `json:"createdAt"`
	UpdatedAt     int64  `json:"updatedAt"`
}

type RecordInfo struct {
	Id          string `json:"id"`
	DomainId    string `json:"domainId"`
	GroupId     string `json:"groupId"`
	GroupName   string `json:"groupName"`
	LineId      string `json:"lineId"`
	LineCode    string `json:"lineCode"`
	LineName    string `json:"lineName"`
	Type        int32  `json:"type"`
	Host        string `json:"host"`
	DisplayHost string `json:"displayHost"`
	Data        string `json:"data"`
	DisplayData string `json:"displayData"`
	Ttl         int32  `json:"ttl"`
	Weight      int32  `json:"weight"`
	Preference  int32  `json:"preference"`
	CreatedAt   int64  `json:"createdAt"`
	UpdatedAt   int64  `json:"updatedAt"`
}

type ListDomainsRequest struct {
	PageIndex int32   `json:"pageIndex"`
	PageSize  int32   `json:"pageSize"`
	GroupId   *string `json:"groupId,omitempty"`
}

type ListDomainsResponse struct {
	baseResponse
	Data *struct {
		Total   int32         `json:"total"`
		Results []*DomainInfo `json:"results"`
	} `json:"data,omitempty"`
}

type ListRecordsRequest struct {
	PageIndex int32   `json:"pageIndex"`
	PageSize  int32   `json:"pageSize"`
	DomainId  string  `json:"domainId"`
	GroupId   *string `json:"groupId,omitempty"`
	LineId    *string `json:"lineId,omitempty"`
	Type      *int32  `json:"type,omitempty"`
	Host      *string `json:"host,omitempty"`
	Data      *string `json:"data,omitempty"`
}

type ListRecordsResponse struct {
	baseResponse
	Data *struct {
		Total   int32         `json:"total"`
		Results []*RecordInfo `json:"results"`
	} `json:"data,omitempty"`
}

type CreateRecordRequest struct {
	DomainId   string  `json:"domainId"`
	GroupId    *string `json:"groupId,omitempty"`
	LineId     *string `json:"lineId,omitempty"`
	Type       int32   `json:"type"`
	Host       string  `json:"host"`
	Data       string  `json:"data"`
	Ttl        int32   `json:"ttl"`
	Weight     *int32  `json:"weight,omitempty"`
	Preference *int32  `json:"preference,omitempty"`
}

type CreateRecordResponse struct {
	baseResponse
	Data *struct {
		Id string `json:"id"`
	} `json:"data,omitempty"`
}

type UpdateRecordRequest struct {
	Id         string  `json:"id"`
	GroupId    *string `json:"groupId,omitempty"`
	LineId     *string `json:"lineId,omitempty"`
	Type       *int32  `json:"type,omitempty"`
	Host       *string `json:"host,omitempty"`
	Data       *string `json:"data,omitempty"`
	Ttl        *int32  `json:"ttl,omitempty"`
	Weight     *int32  `json:"weight,omitempty"`
	Preference *int32  `json:"preference,omitempty"`
}

type UpdateRecordResponse struct {
	baseResponse
}

type DeleteRecordRequest struct {
	Id string `json:"-"`
}

type DeleteRecordResponse struct {
	baseResponse
}
