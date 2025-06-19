package udnr

import (
	"github.com/ucloud/ucloud-sdk-go/ucloud/request"
	"github.com/ucloud/ucloud-sdk-go/ucloud/response"
)

type QueryDomainDNSRequest struct {
	request.CommonBase

	Dn *string `required:"true"`
}

type QueryDomainDNSResponse struct {
	response.CommonBase

	Data []DomainDNSRecord
}

func (c *UDNRClient) NewQueryDomainDNSRequest() *QueryDomainDNSRequest {
	req := &QueryDomainDNSRequest{}

	c.Client.SetupRequest(req)

	req.SetRetryable(false)
	return req
}

func (c *UDNRClient) QueryDomainDNS(req *QueryDomainDNSRequest) (*QueryDomainDNSResponse, error) {
	var err error
	var res QueryDomainDNSResponse

	reqCopier := *req

	err = c.Client.InvokeAction("UdnrDomainDNSQuery", &reqCopier, &res)
	if err != nil {
		return &res, err
	}

	return &res, nil
}
