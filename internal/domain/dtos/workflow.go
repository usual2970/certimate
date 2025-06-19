package dtos

import "github.com/certimate-go/certimate/internal/domain"

type WorkflowStartRunReq struct {
	WorkflowId string                     `json:"-"`
	RunTrigger domain.WorkflowTriggerType `json:"trigger"`
}

type WorkflowCancelRunReq struct {
	WorkflowId string `json:"-"`
	RunId      string `json:"-"`
}
