package onepanelsdk

import (
	"fmt"
	"net/http"
)

func (c *Client) UpdateSystemSSL(req *UpdateSystemSSLRequest) (*UpdateSystemSSLResponse, error) {
	resp := &UpdateSystemSSLResponse{}
	err := c.sendRequestWithResult(http.MethodPost, "/settings/ssl/update", req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) SearchWebsiteSSL(req *SearchWebsiteSSLRequest) (*SearchWebsiteSSLResponse, error) {
	resp := &SearchWebsiteSSLResponse{}
	err := c.sendRequestWithResult(http.MethodPost, "/websites/ssl/search", req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) UploadWebsiteSSL(req *UploadWebsiteSSLRequest) (*UploadWebsiteSSLResponse, error) {
	resp := &UploadWebsiteSSLResponse{}
	err := c.sendRequestWithResult(http.MethodPost, "/websites/ssl/upload", req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) GetHttpsConf(req *GetHttpsConfRequest) (*GetHttpsConfResponse, error) {
	resp := &GetHttpsConfResponse{}
	err := c.sendRequestWithResult(http.MethodGet, fmt.Sprintf("/websites/%d/https", req.WebsiteID), req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) UpdateHttpsConf(req *UpdateHttpsConfRequest) (*UpdateHttpsConfResponse, error) {
	resp := &UpdateHttpsConfResponse{}
	err := c.sendRequestWithResult(http.MethodPost, fmt.Sprintf("/websites/%d/https", req.WebsiteID), req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
