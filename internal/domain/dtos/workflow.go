package dtos

import "github.com/usual2970/certimate/internal/domain"

type WorkflowStartRunReq struct {
	WorkflowId string                     `json:"-"`
	Trigger    domain.WorkflowTriggerType `json:"trigger"`
}

type WorkflowCancelRunReq struct {
	WorkflowId string `json:"-"`
	RunId      string `json:"-"`
}
