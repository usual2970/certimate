package dtos

import "github.com/usual2970/certimate/internal/domain"

type NotifyTestPushReq struct {
	Channel domain.NotifyChannelType `json:"channel"`
}
