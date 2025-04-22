package rainyunsdk

import (
	"fmt"
	"net/http"
)

func (c *Client) SslCenterList(req *SslCenterListRequest) (*SslCenterListResponse, error) {
	resp := &SslCenterListResponse{}
	err := c.sendRequestWithResult(http.MethodGet, "/product/sslcenter", req, resp)
	return resp, err
}

func (c *Client) SslCenterGet(id int32) (*SslCenterGetResponse, error) {
	if id == 0 {
		return nil, fmt.Errorf("rainyun api error: invalid parameter: id")
	}

	resp := &SslCenterGetResponse{}
	err := c.sendRequestWithResult(http.MethodGet, fmt.Sprintf("/product/sslcenter/%d", id), nil, resp)
	return resp, err
}

func (c *Client) SslCenterCreate(req *SslCenterCreateRequest) (*SslCenterCreateResponse, error) {
	resp := &SslCenterCreateResponse{}
	err := c.sendRequestWithResult(http.MethodPost, "/product/sslcenter/", req, resp)
	return resp, err
}

func (c *Client) RcdnInstanceSslBind(id int32, req *RcdnInstanceSslBindRequest) (*RcdnInstanceSslBindResponse, error) {
	if id == 0 {
		return nil, fmt.Errorf("rainyun api error: invalid parameter: id")
	}

	resp := &RcdnInstanceSslBindResponse{}
	err := c.sendRequestWithResult(http.MethodPost, fmt.Sprintf("/product/rcdn/instance/%d/ssl_bind", id), req, resp)
	return resp, err
}
