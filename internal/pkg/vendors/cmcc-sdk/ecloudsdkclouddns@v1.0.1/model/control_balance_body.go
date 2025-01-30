// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model

import (
    "gitlab.ecloud.com/ecloud/ecloudsdkcore/position"
)

type ControlBalanceBody struct {
    position.Body
	// 操作列表
	OperateList *[]ControlBalanceRequestOperateList `json:"operateList"`

	// 域名
	DomainName string `json:"domainName"`
}
