// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model

import (
    "gitlab.ecloud.com/ecloud/ecloudsdkcore/position"
)
type ControlRecordOpenapiBodyOperateEnum string

// List of Operate
const (
    ControlRecordOpenapiBodyOperateEnumDisable ControlRecordOpenapiBodyOperateEnum = "DISABLE"
    ControlRecordOpenapiBodyOperateEnumEnable ControlRecordOpenapiBodyOperateEnum = "ENABLE"
)

type ControlRecordOpenapiBody struct {
    position.Body
	// 操作
	Operate ControlRecordOpenapiBodyOperateEnum `json:"operate"`

	// 域名
	DomainName string `json:"domainName"`

	// 解析记录ID列表
	RecordIdList []string `json:"recordIdList"`
}
