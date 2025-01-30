// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model

import (
    "gitlab.ecloud.com/ecloud/ecloudsdkcore/position"
)

type CancelDnsOrderBody struct {
    position.Body
	// 退订的实例ID
	InstanceId string `json:"instanceId"`
}
