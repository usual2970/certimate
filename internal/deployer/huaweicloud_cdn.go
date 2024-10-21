package deployer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/global"
	cdn "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cdn/v2"
	cdnModel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cdn/v2/model"
	cdnRegion "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cdn/v2/region"

	"github.com/usual2970/certimate/internal/domain"
	uploader "github.com/usual2970/certimate/internal/pkg/core/uploader"
	"github.com/usual2970/certimate/internal/pkg/utils/cast"
)

type HuaweiCloudCDNDeployer struct {
	option *DeployerOption
	infos  []string
}

func NewHuaweiCloudCDNDeployer(option *DeployerOption) (Deployer, error) {
	return &HuaweiCloudCDNDeployer{
		option: option,
		infos:  make([]string, 0),
	}, nil
}

func (d *HuaweiCloudCDNDeployer) GetID() string {
	return fmt.Sprintf("%s-%s", d.option.AccessRecord.GetString("name"), d.option.AccessRecord.Id)
}

func (d *HuaweiCloudCDNDeployer) GetInfo() []string {
	return d.infos
}

func (d *HuaweiCloudCDNDeployer) Deploy(ctx context.Context) error {
	access := &domain.HuaweiCloudAccess{}
	if err := json.Unmarshal([]byte(d.option.Access), access); err != nil {
		return err
	}

	// TODO: CDN 服务与 DNS 服务所支持的区域可能不一致，这里暂时不传而是使用默认值，仅支持华为云国内版
	client, err := d.createClient("", access.AccessKeyId, access.SecretAccessKey)
	if err != nil {
		return err
	}

	d.infos = append(d.infos, toStr("SDK 客户端创建成功", nil))

	// 查询加速域名配置
	// REF: https://support.huaweicloud.com/api-cdn/ShowDomainFullConfig.html
	showDomainFullConfigReq := &cdnModel.ShowDomainFullConfigRequest{
		DomainName: d.option.DeployConfig.GetConfigAsString("domain"),
	}
	showDomainFullConfigResp, err := client.ShowDomainFullConfig(showDomainFullConfigReq)
	if err != nil {
		return err
	}

	d.infos = append(d.infos, toStr("已查询到加速域名配置", showDomainFullConfigResp))

	// 更新加速域名配置
	// REF: https://support.huaweicloud.com/api-cdn/UpdateDomainMultiCertificates.html
	// REF: https://support.huaweicloud.com/usermanual-cdn/cdn_01_0306.html
	updateDomainMultiCertificatesReqBodyContent := &huaweicloudCDNUpdateDomainMultiCertificatesRequestBodyContent{}
	updateDomainMultiCertificatesReqBodyContent.DomainName = d.option.DeployConfig.GetConfigAsString("domain")
	updateDomainMultiCertificatesReqBodyContent.HttpsSwitch = 1
	var updateDomainMultiCertificatesResp *cdnModel.UpdateDomainMultiCertificatesResponse
	if d.option.DeployConfig.GetConfigAsBool("useSCM") {
		uploader, err := uploader.NewHuaweiCloudSCMUploader(&uploader.HuaweiCloudSCMUploaderConfig{
			Region:          "", // TODO: SCM 服务与 DNS 服务所支持的区域可能不一致，这里暂时不传而是使用默认值，仅支持华为云国内版
			AccessKeyId:     access.AccessKeyId,
			SecretAccessKey: access.SecretAccessKey,
		})
		if err != nil {
			return err
		}

		// 上传证书到 SCM
		uploadResult, err := uploader.Upload(ctx, d.option.Certificate.Certificate, d.option.Certificate.PrivateKey)
		if err != nil {
			return err
		}

		d.infos = append(d.infos, toStr("已上传证书", uploadResult))

		updateDomainMultiCertificatesReqBodyContent.CertificateType = cast.Int32Ptr(2)
		updateDomainMultiCertificatesReqBodyContent.SCMCertificateId = cast.StringPtr(uploadResult.CertId)
		updateDomainMultiCertificatesReqBodyContent.CertName = cast.StringPtr(uploadResult.CertName)
	} else {
		updateDomainMultiCertificatesReqBodyContent.CertificateType = cast.Int32Ptr(0)
		updateDomainMultiCertificatesReqBodyContent.CertName = cast.StringPtr(fmt.Sprintf("certimate-%d", time.Now().UnixMilli()))
		updateDomainMultiCertificatesReqBodyContent.Certificate = cast.StringPtr(d.option.Certificate.Certificate)
		updateDomainMultiCertificatesReqBodyContent.PrivateKey = cast.StringPtr(d.option.Certificate.PrivateKey)
	}
	updateDomainMultiCertificatesReqBodyContent = mergeHuaweiCloudCDNConfig(showDomainFullConfigResp.Configs, updateDomainMultiCertificatesReqBodyContent)
	updateDomainMultiCertificatesReq := &huaweicloudCDNUpdateDomainMultiCertificatesRequest{
		Body: &huaweicloudCDNUpdateDomainMultiCertificatesRequestBody{
			Https: updateDomainMultiCertificatesReqBodyContent,
		},
	}
	updateDomainMultiCertificatesResp, err = executeHuaweiCloudCDNUploadDomainMultiCertificates(client, updateDomainMultiCertificatesReq)
	if err != nil {
		return err
	}

	d.infos = append(d.infos, toStr("已更新加速域名配置", updateDomainMultiCertificatesResp))

	return nil
}

func (d *HuaweiCloudCDNDeployer) createClient(region, accessKeyId, secretAccessKey string) (*cdn.CdnClient, error) {
	auth, err := global.NewCredentialsBuilder().
		WithAk(accessKeyId).
		WithSk(secretAccessKey).
		SafeBuild()
	if err != nil {
		return nil, err
	}

	if region == "" {
		region = "cn-north-1" // CDN 服务默认区域：华北北京一
	}

	hcRegion, err := cdnRegion.SafeValueOf(region)
	if err != nil {
		return nil, err
	}

	hcClient, err := cdn.CdnClientBuilder().
		WithRegion(hcRegion).
		WithCredential(auth).
		SafeBuild()
	if err != nil {
		return nil, err
	}

	client := cdn.NewCdnClient(hcClient)
	return client, nil
}

type huaweicloudCDNUpdateDomainMultiCertificatesRequestBodyContent struct {
	cdnModel.UpdateDomainMultiCertificatesRequestBodyContent `json:",inline"`

	SCMCertificateId *string `json:"scm_certificate_id,omitempty"`
}

type huaweicloudCDNUpdateDomainMultiCertificatesRequestBody struct {
	Https *huaweicloudCDNUpdateDomainMultiCertificatesRequestBodyContent `json:"https,omitempty"`
}

type huaweicloudCDNUpdateDomainMultiCertificatesRequest struct {
	Body *huaweicloudCDNUpdateDomainMultiCertificatesRequestBody `json:"body,omitempty"`
}

func executeHuaweiCloudCDNUploadDomainMultiCertificates(client *cdn.CdnClient, request *huaweicloudCDNUpdateDomainMultiCertificatesRequest) (*cdnModel.UpdateDomainMultiCertificatesResponse, error) {
	// 华为云官方 SDK 中目前提供的字段缺失，这里暂时先需自定义请求
	// 可能需要等之后 SDK 更新

	requestDef := cdn.GenReqDefForUpdateDomainMultiCertificates()

	if resp, err := client.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*cdnModel.UpdateDomainMultiCertificatesResponse), nil
	}
}

func mergeHuaweiCloudCDNConfig(src *cdnModel.ConfigsGetBody, dest *huaweicloudCDNUpdateDomainMultiCertificatesRequestBodyContent) *huaweicloudCDNUpdateDomainMultiCertificatesRequestBodyContent {
	if src == nil {
		return dest
	}

	// 华为云 API 中不传的字段表示使用默认值、而非保留原值，因此这里需要把原配置中的参数重新赋值回去
	// 而且蛋疼的是查询接口返回的数据结构和更新接口传入的参数结构不一致，需要做很多转化

	if *src.OriginProtocol == "follow" {
		dest.AccessOriginWay = cast.Int32Ptr(1)
	} else if *src.OriginProtocol == "http" {
		dest.AccessOriginWay = cast.Int32Ptr(2)
	} else if *src.OriginProtocol == "https" {
		dest.AccessOriginWay = cast.Int32Ptr(3)
	}

	if src.ForceRedirect != nil {
		dest.ForceRedirectConfig = &cdnModel.ForceRedirect{}

		if src.ForceRedirect.Status == "on" {
			dest.ForceRedirectConfig.Switch = 1
			dest.ForceRedirectConfig.RedirectType = src.ForceRedirect.Type
		} else {
			dest.ForceRedirectConfig.Switch = 0
		}
	}

	if src.Https != nil {
		if *src.Https.Http2Status == "on" {
			dest.Http2 = cast.Int32Ptr(1)
		}
	}

	return dest
}
