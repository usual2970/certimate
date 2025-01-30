// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model

import (
    "gitlab.ecloud.com/ecloud/ecloudsdkcore/position"
)

type GetOperationLogBody struct {
    position.Body
	// 操作内容，支持模糊查询
	OperationLike string `json:"operationLike,omitempty"`

	// 操作人，支持模糊查询
	OperatorLike string `json:"operatorLike,omitempty"`

	// 操作时间
	OperationTime *GetOperationLogRequestOperationTime `json:"operationTime,omitempty"`
}
