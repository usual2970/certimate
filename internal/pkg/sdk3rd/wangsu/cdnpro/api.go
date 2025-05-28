package cdnpro

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-resty/resty/v2"
)

func (c *Client) CreateCertificate(req *CreateCertificateRequest) (*CreateCertificateResponse, error) {
	resp := &CreateCertificateResponse{}
	rres, err := c.client.SendRequestWithResult(http.MethodPost, "/cdn/certificates", req, resp, func(r *resty.Request) {
		r.SetHeader("X-CNC-Timestamp", fmt.Sprintf("%d", req.Timestamp))
	})
	if err != nil {
		return resp, err
	}

	resp.CertificateUrl = rres.Header().Get("Location")
	return resp, err
}

func (c *Client) UpdateCertificate(certificateId string, req *UpdateCertificateRequest) (*UpdateCertificateResponse, error) {
	if certificateId == "" {
		return nil, fmt.Errorf("wangsu api error: invalid parameter: certificateId")
	}

	resp := &UpdateCertificateResponse{}
	rres, err := c.client.SendRequestWithResult(http.MethodPatch, fmt.Sprintf("/cdn/certificates/%s", url.PathEscape(certificateId)), req, resp, func(r *resty.Request) {
		r.SetHeader("X-CNC-Timestamp", fmt.Sprintf("%d", req.Timestamp))
	})
	if err != nil {
		return resp, err
	}

	resp.CertificateUrl = rres.Header().Get("Location")
	return resp, err
}

func (c *Client) GetHostnameDetail(hostname string) (*GetHostnameDetailResponse, error) {
	if hostname == "" {
		return nil, fmt.Errorf("wangsu api error: invalid parameter: hostname")
	}

	resp := &GetHostnameDetailResponse{}
	_, err := c.client.SendRequestWithResult(http.MethodGet, fmt.Sprintf("/cdn/hostnames/%s", url.PathEscape(hostname)), nil, resp)
	return resp, err
}

func (c *Client) CreateDeploymentTask(req *CreateDeploymentTaskRequest) (*CreateDeploymentTaskResponse, error) {
	resp := &CreateDeploymentTaskResponse{}
	rres, err := c.client.SendRequestWithResult(http.MethodPost, "/cdn/deploymentTasks", req, resp)
	if err != nil {
		return resp, err
	}

	resp.DeploymentTaskUrl = rres.Header().Get("Location")
	return resp, err
}

func (c *Client) GetDeploymentTaskDetail(deploymentTaskId string) (*GetDeploymentTaskDetailResponse, error) {
	if deploymentTaskId == "" {
		return nil, fmt.Errorf("wangsu api error: invalid parameter: deploymentTaskId")
	}

	resp := &GetDeploymentTaskDetailResponse{}
	_, err := c.client.SendRequestWithResult(http.MethodGet, fmt.Sprintf("/cdn/deploymentTasks/%s", url.PathEscape(deploymentTaskId)), nil, resp)
	return resp, err
}
