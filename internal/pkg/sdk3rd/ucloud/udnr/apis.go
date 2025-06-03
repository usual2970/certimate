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

type AddDomainDNSRequest struct {
	request.CommonBase

	Dn         *string `required:"true"`
	DnsType    *string `required:"true"`
	RecordName *string `required:"true"`
	Content    *string `required:"true"`
	TTL        *int    `required:"true"`
	Prio       *int    `required:"false"`
}

type AddDomainDNSResponse struct {
	response.CommonBase
}

func (c *UDNRClient) NewAddDomainDNSRequest() *AddDomainDNSRequest {
	req := &AddDomainDNSRequest{}

	c.Client.SetupRequest(req)

	req.SetRetryable(false)
	return req
}

func (c *UDNRClient) AddDomainDNS(req *AddDomainDNSRequest) (*AddDomainDNSResponse, error) {
	var err error
	var res AddDomainDNSResponse

	reqCopier := *req

	err = c.Client.InvokeAction("UdnrDomainDNSAdd", &reqCopier, &res)
	if err != nil {
		return &res, err
	}

	return &res, nil
}

type DeleteDomainDNSRequest struct {
	request.CommonBase

	Dn         *string `required:"true"`
	DnsType    *string `required:"true"`
	RecordName *string `required:"true"`
	Content    *string `required:"true"`
}

type DeleteDomainDNSResponse struct {
	response.CommonBase
}

func (c *UDNRClient) NewDeleteDomainDNSRequest() *DeleteDomainDNSRequest {
	req := &DeleteDomainDNSRequest{}

	c.Client.SetupRequest(req)

	req.SetRetryable(false)
	return req
}

func (c *UDNRClient) DeleteDomainDNS(req *DeleteDomainDNSRequest) (*DeleteDomainDNSResponse, error) {
	var err error
	var res DeleteDomainDNSResponse

	reqCopier := *req

	err = c.Client.InvokeAction("UdnrDeleteDnsRecord", &reqCopier, &res)
	if err != nil {
		return &res, err
	}

	return &res, nil
}
