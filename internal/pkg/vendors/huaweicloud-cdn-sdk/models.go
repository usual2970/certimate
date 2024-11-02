package huaweicloudcdnsdk

import (
	hcCdnModel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cdn/v2/model"

	"github.com/usual2970/certimate/internal/pkg/utils/cast"
)

type UpdateDomainMultiCertificatesExRequestBodyContent struct {
	hcCdnModel.UpdateDomainMultiCertificatesRequestBodyContent `json:",inline"`

	// 华为云官方 SDK 中目前提供的字段缺失，这里暂时先需自定义请求，可能需要等之后 SDK 更新。
	SCMCertificateId *string `json:"scm_certificate_id,omitempty"`
}

type UpdateDomainMultiCertificatesExRequestBody struct {
	Https *UpdateDomainMultiCertificatesExRequestBodyContent `json:"https,omitempty"`
}

type UpdateDomainMultiCertificatesExRequest struct {
	Body *UpdateDomainMultiCertificatesExRequestBody `json:"body,omitempty"`
}

type UpdateDomainMultiCertificatesExResponse struct {
	hcCdnModel.UpdateDomainMultiCertificatesResponse
}

func (m *UpdateDomainMultiCertificatesExRequestBodyContent) MergeConfig(src *hcCdnModel.ConfigsGetBody) *UpdateDomainMultiCertificatesExRequestBodyContent {
	if src == nil {
		return m
	}

	// 华为云 API 中不传的字段表示使用默认值、而非保留原值，因此这里需要把原配置中的参数重新赋值回去。
	// 而且蛋疼的是查询接口返回的数据结构和更新接口传入的参数结构不一致，需要做很多转化。

	if *src.OriginProtocol == "follow" {
		m.AccessOriginWay = cast.Int32Ptr(1)
	} else if *src.OriginProtocol == "http" {
		m.AccessOriginWay = cast.Int32Ptr(2)
	} else if *src.OriginProtocol == "https" {
		m.AccessOriginWay = cast.Int32Ptr(3)
	}

	if src.ForceRedirect != nil {
		m.ForceRedirectConfig = &hcCdnModel.ForceRedirect{}

		if src.ForceRedirect.Status == "on" {
			m.ForceRedirectConfig.Switch = 1
			m.ForceRedirectConfig.RedirectType = src.ForceRedirect.Type
		} else {
			m.ForceRedirectConfig.Switch = 0
		}
	}

	if src.Https != nil {
		if *src.Https.Http2Status == "on" {
			m.Http2 = cast.Int32Ptr(1)
		}
	}

	return m
}
