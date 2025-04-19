package cdn

import (
	"github.com/usual2970/certimate/internal/pkg/vendors/wangsu-sdk/openapi"
)

type baseResponse struct {
	RequestId *string `json:"-"`
	Code      *string `json:"code,omitempty"`
	Message   *string `json:"message,omitempty"`
}

var _ openapi.Result = (*baseResponse)(nil)

func (r *baseResponse) SetRequestId(requestId string) {
	r.RequestId = &requestId
}

type CertificateVersion struct {
	Comments           *string                               `json:"comments,omitempty"`
	PrivateKey         *string                               `json:"privateKey,omitempty"`
	Certificate        *string                               `json:"certificate,omitempty"`
	ChainCert          *string                               `json:"chainCert,omitempty"`
	IdentificationInfo *CertificateVersionIdentificationInfo `json:"identificationInfo,omitempty"`
}

type CertificateVersionIdentificationInfo struct {
	Country                 *string   `json:"country,omitempty"`
	State                   *string   `json:"state,omitempty"`
	City                    *string   `json:"city,omitempty"`
	Company                 *string   `json:"company,omitempty"`
	Department              *string   `json:"department,omitempty"`
	CommonName              *string   `json:"commonName,omitempty" required:"true"`
	Email                   *string   `json:"email,omitempty"`
	SubjectAlternativeNames *[]string `json:"subjectAlternativeNames,omitempty" required:"true"`
}

type CreateCertificateRequest struct {
	Timestamp   int64               `json:"-"`
	Name        *string             `json:"name,omitempty" required:"true"`
	Description *string             `json:"description,omitempty"`
	AutoRenew   *string             `json:"autoRenew,omitempty"`
	ForceRenew  *bool               `json:"forceRenew,omitempty"`
	NewVersion  *CertificateVersion `json:"newVersion,omitempty" required:"true"`
}

type CreateCertificateResponse struct {
	baseResponse
	CertificateUrl string `json:"-"`
}

type UpdateCertificateRequest struct {
	Timestamp   int64               `json:"-"`
	Name        *string             `json:"name,omitempty"`
	Description *string             `json:"description,omitempty"`
	AutoRenew   *string             `json:"autoRenew,omitempty"`
	ForceRenew  *bool               `json:"forceRenew,omitempty"`
	NewVersion  *CertificateVersion `json:"newVersion,omitempty" required:"true"`
}

type UpdateCertificateResponse struct {
	baseResponse
	CertificateUrl string `json:"-"`
}

type HostnameProperty struct {
	PropertyId    string  `json:"propertyId"`
	Version       int32   `json:"version"`
	CertificateId *string `json:"certificateId,omitempty"`
}

type GetHostnameDetailResponse struct {
	baseResponse
	Hostname             string            `json:"hostname"`
	PropertyInProduction *HostnameProperty `json:"propertyInProduction,omitempty"`
	PropertyInStaging    *HostnameProperty `json:"propertyInStaging,omitempty"`
}

type DeploymentTaskAction struct {
	Action        *string `json:"action,omitempty" required:"true"`
	PropertyId    *string `json:"propertyId,omitempty"`
	CertificateId *string `json:"certificateId,omitempty"`
	Version       *string `json:"version,omitempty"`
}

type CreateDeploymentTaskRequest struct {
	Name    *string                 `json:"name,omitempty"`
	Target  *string                 `json:"target,omitempty" required:"true"`
	Actions *[]DeploymentTaskAction `json:"actions,omitempty" required:"true"`
	Webhook *string                 `json:"webhook,omitempty"`
}

type CreateDeploymentTaskResponse struct {
	baseResponse
	DeploymentTaskUrl string `json:"-"`
}

type GetDeploymentTaskDetailResponse struct {
	baseResponse
	Id             string                 `json:"id"`
	Name           string                 `json:"name"`
	Target         string                 `json:"target"`
	Actions        []DeploymentTaskAction `json:"actions"`
	Status         string                 `json:"status"`
	StatusDetails  string                 `json:"statusDetails"`
	SubmissionTime string                 `json:"submissionTime"`
	FinishTime     string                 `json:"finishTime"`
	ApiRequestId   string                 `json:"apiRequestId"`
}
