// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model

import (
    "gitlab.ecloud.com/ecloud/ecloudsdkcore/position"
)

type RenewProductOpenapiBody struct {
    position.Body
	// 续订的时长
	Duration *int32 `json:"duration"`

	// 资源实例ID
	InstanceId string `json:"instanceId"`
}
