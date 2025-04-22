package dnslasdk

import (
	"fmt"
	"net/http"
	"net/url"
)

func (c *Client) ListDomains(req *ListDomainsRequest) (*ListDomainsResponse, error) {
	resp := &ListDomainsResponse{}
	err := c.sendRequestWithResult(http.MethodGet, "/domainList", req, resp)
	return resp, err
}

func (c *Client) ListRecords(req *ListRecordsRequest) (*ListRecordsResponse, error) {
	resp := &ListRecordsResponse{}
	err := c.sendRequestWithResult(http.MethodGet, "/recordList", req, resp)
	return resp, err
}

func (c *Client) CreateRecord(req *CreateRecordRequest) (*CreateRecordResponse, error) {
	resp := &CreateRecordResponse{}
	err := c.sendRequestWithResult(http.MethodPost, "/record", req, resp)
	return resp, err
}

func (c *Client) UpdateRecord(req *UpdateRecordRequest) (*UpdateRecordResponse, error) {
	resp := &UpdateRecordResponse{}
	err := c.sendRequestWithResult(http.MethodPut, "/record", req, resp)
	return resp, err
}

func (c *Client) DeleteRecord(req *DeleteRecordRequest) (*DeleteRecordResponse, error) {
	if req.Id == "" {
		return nil, fmt.Errorf("dnsla api error: invalid parameter: Id")
	}

	resp := &DeleteRecordResponse{}
	err := c.sendRequestWithResult(http.MethodDelete, fmt.Sprintf("/record?id=%s", url.QueryEscape(req.Id)), req, resp)
	return resp, err
}
