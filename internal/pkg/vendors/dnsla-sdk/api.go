package dnslasdk

import (
	"fmt"
	"net/http"
	"net/url"
)

func (c *Client) ListDomains(req *ListDomainsRequest) (*ListDomainsResponse, error) {
	resp := ListDomainsResponse{}
	err := c.sendRequestWithResult(http.MethodGet, "/domainList", req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) ListRecords(req *ListRecordsRequest) (*ListRecordsResponse, error) {
	resp := ListRecordsResponse{}
	err := c.sendRequestWithResult(http.MethodGet, "/recordList", req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) CreateRecord(req *CreateRecordRequest) (*CreateRecordResponse, error) {
	resp := CreateRecordResponse{}
	err := c.sendRequestWithResult(http.MethodPost, "/record", req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) UpdateRecord(req *UpdateRecordRequest) (*UpdateRecordResponse, error) {
	resp := UpdateRecordResponse{}
	err := c.sendRequestWithResult(http.MethodPut, "/record", req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) DeleteRecord(req *DeleteRecordRequest) (*DeleteRecordResponse, error) {
	resp := DeleteRecordResponse{}
	err := c.sendRequestWithResult(http.MethodDelete, fmt.Sprintf("/record?id=%s", url.QueryEscape(req.Id)), req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
