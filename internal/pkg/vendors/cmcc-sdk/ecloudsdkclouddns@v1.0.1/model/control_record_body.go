// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model

import (
    "gitlab.ecloud.com/ecloud/ecloudsdkcore/position"
)
type ControlRecordBodyOperateEnum string

// List of Operate
const (
    ControlRecordBodyOperateEnumDisable ControlRecordBodyOperateEnum = "DISABLE"
    ControlRecordBodyOperateEnumEnable ControlRecordBodyOperateEnum = "ENABLE"
)

type ControlRecordBody struct {
    position.Body
	// 操作
	Operate ControlRecordBodyOperateEnum `json:"operate"`

	// 域名
	DomainName string `json:"domainName"`

	// 解析记录ID列表
	RecordIdList []string `json:"recordIdList"`
}
