package dtos

import "github.com/certimate-go/certimate/internal/domain"

type NotifyTestPushReq struct {
	Channel domain.NotifyChannelType `json:"channel"`
}
