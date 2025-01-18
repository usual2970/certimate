package dtos

import "github.com/usual2970/certimate/internal/domain"

type WorkflowRunReq struct {
	WorkflowId string                     `json:"-"`
	Trigger    domain.WorkflowTriggerType `json:"trigger"`
}
