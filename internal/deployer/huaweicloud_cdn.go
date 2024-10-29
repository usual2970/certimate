package deployer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/global"
	hcCdn "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cdn/v2"
	hcCdnModel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cdn/v2/model"
	hcCdnRegion "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cdn/v2/region"
	xerrors "github.com/pkg/errors"

	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	"github.com/usual2970/certimate/internal/pkg/utils/cast"
)

type HuaweiCloudCDNDeployer struct {
	option *DeployerOption
	infos  []string

	sdkClient   *hcCdn.CdnClient
	sslUploader uploader.Uploader
}

func NewHuaweiCloudCDNDeployer(option *DeployerOption) (Deployer, error) {
	access := &domain.HuaweiCloudAccess{}
	if err := json.Unmarshal([]byte(option.Access), access); err != nil {
		return nil, xerrors.Wrap(err, "failed to get access")
	}

	client, err := (&HuaweiCloudCDNDeployer{}).createSdkClient(
		access.AccessKeyId,
		access.SecretAccessKey,
		option.DeployConfig.GetConfigAsString("region"),
	)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk client")
	}

	// TODO: SCM 服务与 DNS 服务所支持的区域可能不一致，这里暂时不传而是使用默认值，仅支持华为云国内版
	uploader, err := uploader.NewHuaweiCloudSCMUploader(&uploader.HuaweiCloudSCMUploaderConfig{
		AccessKeyId:     access.AccessKeyId,
		SecretAccessKey: access.SecretAccessKey,
		Region:          "",
	})
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create ssl uploader")
	}

	return &HuaweiCloudCDNDeployer{
		option:      option,
		infos:       make([]string, 0),
		sdkClient:   client,
		sslUploader: uploader,
	}, nil
}

func (d *HuaweiCloudCDNDeployer) GetID() string {
	return fmt.Sprintf("%s-%s", d.option.AccessRecord.GetString("name"), d.option.AccessRecord.Id)
}

func (d *HuaweiCloudCDNDeployer) GetInfo() []string {
	return d.infos
}

func (d *HuaweiCloudCDNDeployer) Deploy(ctx context.Context) error {
	// 查询加速域名配置
	// REF: https://support.huaweicloud.com/api-cdn/ShowDomainFullConfig.html
	showDomainFullConfigReq := &hcCdnModel.ShowDomainFullConfigRequest{
		DomainName: d.option.DeployConfig.GetConfigAsString("domain"),
	}
	showDomainFullConfigResp, err := d.sdkClient.ShowDomainFullConfig(showDomainFullConfigReq)
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
	var updateDomainMultiCertificatesResp *hcCdnModel.UpdateDomainMultiCertificatesResponse
	if d.option.DeployConfig.GetConfigAsBool("useSCM") {
		// 上传证书到 SCM
		upres, err := d.sslUploader.Upload(ctx, d.option.Certificate.Certificate, d.option.Certificate.PrivateKey)
		if err != nil {
			return err
		}

		d.infos = append(d.infos, toStr("已上传证书", upres))

		updateDomainMultiCertificatesReqBodyContent.CertificateType = cast.Int32Ptr(2)
		updateDomainMultiCertificatesReqBodyContent.SCMCertificateId = cast.StringPtr(upres.CertId)
		updateDomainMultiCertificatesReqBodyContent.CertName = cast.StringPtr(upres.CertName)
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
	updateDomainMultiCertificatesResp, err = executeHuaweiCloudCDNUploadDomainMultiCertificates(d.sdkClient, updateDomainMultiCertificatesReq)
	if err != nil {
		return err
	}

	d.infos = append(d.infos, toStr("已更新加速域名配置", updateDomainMultiCertificatesResp))

	return nil
}

func (d *HuaweiCloudCDNDeployer) createSdkClient(accessKeyId, secretAccessKey, region string) (*hcCdn.CdnClient, error) {
	if region == "" {
		region = "cn-north-1" // CDN 服务默认区域：华北一北京
	}

	auth, err := global.NewCredentialsBuilder().
		WithAk(accessKeyId).
		WithSk(secretAccessKey).
		SafeBuild()
	if err != nil {
		return nil, err
	}

	hcRegion, err := hcCdnRegion.SafeValueOf(region)
	if err != nil {
		return nil, err
	}

	hcClient, err := hcCdn.CdnClientBuilder().
		WithRegion(hcRegion).
		WithCredential(auth).
		SafeBuild()
	if err != nil {
		return nil, err
	}

	client := hcCdn.NewCdnClient(hcClient)
	return client, nil
}

type huaweicloudCDNUpdateDomainMultiCertificatesRequestBodyContent struct {
	hcCdnModel.UpdateDomainMultiCertificatesRequestBodyContent `json:",inline"`

	SCMCertificateId *string `json:"scm_certificate_id,omitempty"`
}

type huaweicloudCDNUpdateDomainMultiCertificatesRequestBody struct {
	Https *huaweicloudCDNUpdateDomainMultiCertificatesRequestBodyContent `json:"https,omitempty"`
}

type huaweicloudCDNUpdateDomainMultiCertificatesRequest struct {
	Body *huaweicloudCDNUpdateDomainMultiCertificatesRequestBody `json:"body,omitempty"`
}

func executeHuaweiCloudCDNUploadDomainMultiCertificates(client *hcCdn.CdnClient, request *huaweicloudCDNUpdateDomainMultiCertificatesRequest) (*hcCdnModel.UpdateDomainMultiCertificatesResponse, error) {
	// 华为云官方 SDK 中目前提供的字段缺失，这里暂时先需自定义请求
	// 可能需要等之后 SDK 更新

	requestDef := hcCdn.GenReqDefForUpdateDomainMultiCertificates()

	if resp, err := client.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*hcCdnModel.UpdateDomainMultiCertificatesResponse), nil
	}
}

func mergeHuaweiCloudCDNConfig(src *hcCdnModel.ConfigsGetBody, dest *huaweicloudCDNUpdateDomainMultiCertificatesRequestBodyContent) *huaweicloudCDNUpdateDomainMultiCertificatesRequestBodyContent {
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
		dest.ForceRedirectConfig = &hcCdnModel.ForceRedirect{}

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
