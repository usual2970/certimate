// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model


type UpgradeOrderOpenapiResponseStateEnum string

// List of State
const (
    UpgradeOrderOpenapiResponseStateEnumError UpgradeOrderOpenapiResponseStateEnum = "ERROR"
    UpgradeOrderOpenapiResponseStateEnumException UpgradeOrderOpenapiResponseStateEnum = "EXCEPTION"
    UpgradeOrderOpenapiResponseStateEnumForbidden UpgradeOrderOpenapiResponseStateEnum = "FORBIDDEN"
    UpgradeOrderOpenapiResponseStateEnumOk UpgradeOrderOpenapiResponseStateEnum = "OK"
)

type UpgradeOrderOpenapiResponse struct {

	RequestId string `json:"requestId,omitempty"`

	ErrorMessage string `json:"errorMessage,omitempty"`

	ErrorCode string `json:"errorCode,omitempty"`

	State UpgradeOrderOpenapiResponseStateEnum `json:"state,omitempty"`

	Body *UpgradeOrderOpenapiResponseBody `json:"body,omitempty"`
}
