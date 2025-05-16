package netlify

import (
	"fmt"
	"net/http"
	"net/url"
)

func (c *Client) ProvisionSiteTLSCertificate(siteId string, params *ProvisionSiteTLSCertificateParams) (*ProvisionSiteTLSCertificateResponse, error) {
	if siteId == "" {
		return nil, fmt.Errorf("netlify api error: invalid parameter: SiteId")
	}

	resp := &ProvisionSiteTLSCertificateResponse{}
	err := c.sendRequestWithResult(http.MethodPost, fmt.Sprintf("/sites/%s/ssl", url.PathEscape(siteId)), params, nil, resp)
	return resp, err
}
