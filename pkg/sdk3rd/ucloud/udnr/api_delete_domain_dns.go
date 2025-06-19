package udnr

import (
	"github.com/ucloud/ucloud-sdk-go/ucloud/request"
	"github.com/ucloud/ucloud-sdk-go/ucloud/response"
)

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
