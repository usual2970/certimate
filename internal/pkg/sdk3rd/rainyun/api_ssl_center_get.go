package rainyun

import (
	"context"
	"fmt"
	"net/http"
)

type SslCenterGetResponse struct {
	apiResponseBase

	Data *SslDetail `json:"data,omitempty"`
}

func (c *Client) SslCenterGet(sslId int32) (*SslCenterGetResponse, error) {
	return c.SslCenterGetWithContext(context.Background(), sslId)
}

func (c *Client) SslCenterGetWithContext(ctx context.Context, sslId int32) (*SslCenterGetResponse, error) {
	if sslId == 0 {
		return nil, fmt.Errorf("sdkerr: unset sslId")
	}

	httpreq, err := c.newRequest(http.MethodGet, fmt.Sprintf("/product/sslcenter/%d", sslId))
	if err != nil {
		return nil, err
	} else {
		httpreq.SetContext(ctx)
	}

	result := &SslCenterGetResponse{}
	if _, err := c.doRequestWithResult(httpreq, result); err != nil {
		return result, err
	}

	return result, nil
}
