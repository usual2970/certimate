package cdn

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-resty/resty/v2"
)

func (c *Client) CreateCertificate(req *CreateCertificateRequest) (*CreateCertificateResponse, error) {
	resp := &CreateCertificateResponse{}
	r, err := c.client.SendRequestWithResult(http.MethodPost, "/cdn/certificates", req, resp, func(r *resty.Request) {
		r.SetHeader("x-cnc-timestamp", fmt.Sprintf("%d", req.Timestamp))
	})
	if err != nil {
		return resp, err
	}

	resp.CertificateUrl = r.Header().Get("Location")
	return resp, err
}

func (c *Client) UpdateCertificate(certificateId string, req *UpdateCertificateRequest) (*UpdateCertificateResponse, error) {
	resp := &UpdateCertificateResponse{}
	r, err := c.client.SendRequestWithResult(http.MethodPatch, fmt.Sprintf("/cdn/certificates/%s", url.PathEscape(certificateId)), req, resp, func(r *resty.Request) {
		r.SetHeader("x-cnc-timestamp", fmt.Sprintf("%d", req.Timestamp))
	})
	if err != nil {
		return resp, err
	}

	resp.CertificateUrl = r.Header().Get("Location")
	return resp, err
}

func (c *Client) GetHostnameDetail(hostname string) (*GetHostnameDetailResponse, error) {
	resp := &GetHostnameDetailResponse{}
	_, err := c.client.SendRequestWithResult(http.MethodGet, fmt.Sprintf("/cdn/hostnames/%s", url.PathEscape(hostname)), nil, resp)
	return resp, err
}

func (c *Client) CreateDeploymentTask(req *CreateDeploymentTaskRequest) (*CreateDeploymentTaskResponse, error) {
	resp := &CreateDeploymentTaskResponse{}
	r, err := c.client.SendRequestWithResult(http.MethodPost, "/cdn/deploymentTasks", req, resp)
	if err != nil {
		return resp, err
	}

	resp.DeploymentTaskUrl = r.Header().Get("Location")
	return resp, err
}

func (c *Client) GetDeploymentTaskDetail(deploymentTaskId string) (*GetDeploymentTaskDetailResponse, error) {
	resp := &GetDeploymentTaskDetailResponse{}
	_, err := c.client.SendRequestWithResult(http.MethodGet, fmt.Sprintf("/cdn/deploymentTasks/%s", deploymentTaskId), nil, resp)
	return resp, err
}
