package cdnpro

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type GetDeploymentTaskDetailResponse struct {
	apiResponseBase

	Name           string                     `json:"name"`
	Target         string                     `json:"target"`
	Actions        []DeploymentTaskActionInfo `json:"actions"`
	Status         string                     `json:"status"`
	StatusDetails  string                     `json:"statusDetails"`
	SubmissionTime string                     `json:"submissionTime"`
	FinishTime     string                     `json:"finishTime"`
	ApiRequestId   string                     `json:"apiRequestId"`
}

func (c *Client) GetDeploymentTaskDetail(deploymentTaskId string) (*GetDeploymentTaskDetailResponse, error) {
	return c.GetDeploymentTaskDetailWithContext(context.Background(), deploymentTaskId)
}

func (c *Client) GetDeploymentTaskDetailWithContext(ctx context.Context, deploymentTaskId string) (*GetDeploymentTaskDetailResponse, error) {
	if deploymentTaskId == "" {
		return nil, fmt.Errorf("sdkerr: unset deploymentTaskId")
	}

	httpreq, err := c.newRequest(http.MethodGet, fmt.Sprintf("/cdn/deploymentTasks/%s", url.PathEscape(deploymentTaskId)))
	if err != nil {
		return nil, err
	} else {
		httpreq.SetContext(ctx)
	}

	result := &GetDeploymentTaskDetailResponse{}
	if _, err := c.doRequestWithResult(httpreq, result); err != nil {
		return result, err
	}

	return result, nil
}
