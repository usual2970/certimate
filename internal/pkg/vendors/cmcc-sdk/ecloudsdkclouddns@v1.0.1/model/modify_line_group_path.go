// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model

import (
    "gitlab.ecloud.com/ecloud/ecloudsdkcore/position"
)

type ModifyLineGroupPath struct {
    position.Path
	// 线路分组ID
	GroupId string `json:"groupId"`
}
