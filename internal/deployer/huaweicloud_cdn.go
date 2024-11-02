package deployer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/global"
	hcCdn "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cdn/v2"
	hcCdnModel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cdn/v2/model"
	hcCdnRegion "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cdn/v2/region"
	xerrors "github.com/pkg/errors"

	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	uploaderHcScm "github.com/usual2970/certimate/internal/pkg/core/uploader/providers/huaweicloud-scm"
	"github.com/usual2970/certimate/internal/pkg/utils/cast"
	hcCdnEx "github.com/usual2970/certimate/internal/pkg/vendors/huaweicloud-cdn-sdk"
)

type HuaweiCloudCDNDeployer struct {
	option *DeployerOption
	infos  []string

	sdkClient   *hcCdnEx.Client
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

	uploader, err := uploaderHcScm.New(&uploaderHcScm.HuaweiCloudSCMUploaderConfig{
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

func (d *HuaweiCloudCDNDeployer) GetInfos() []string {
	return d.infos
}

func (d *HuaweiCloudCDNDeployer) Deploy(ctx context.Context) error {
	// 上传证书到 SCM
	upres, err := d.sslUploader.Upload(ctx, d.option.Certificate.Certificate, d.option.Certificate.PrivateKey)
	if err != nil {
		return err
	}

	d.infos = append(d.infos, toStr("已上传证书", upres))

	// 查询加速域名配置
	// REF: https://support.huaweicloud.com/api-cdn/ShowDomainFullConfig.html
	showDomainFullConfigReq := &hcCdnModel.ShowDomainFullConfigRequest{
		DomainName: d.option.DeployConfig.GetConfigAsString("domain"),
	}
	showDomainFullConfigResp, err := d.sdkClient.ShowDomainFullConfig(showDomainFullConfigReq)
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'cdn.ShowDomainFullConfig'")
	}

	d.infos = append(d.infos, toStr("已查询到加速域名配置", showDomainFullConfigResp))

	// 更新加速域名配置
	// REF: https://support.huaweicloud.com/api-cdn/UpdateDomainMultiCertificates.html
	// REF: https://support.huaweicloud.com/usermanual-cdn/cdn_01_0306.html
	updateDomainMultiCertificatesReqBodyContent := &hcCdnEx.UpdateDomainMultiCertificatesExRequestBodyContent{}
	updateDomainMultiCertificatesReqBodyContent.DomainName = d.option.DeployConfig.GetConfigAsString("domain")
	updateDomainMultiCertificatesReqBodyContent.HttpsSwitch = 1
	updateDomainMultiCertificatesReqBodyContent.CertificateType = cast.Int32Ptr(2)
	updateDomainMultiCertificatesReqBodyContent.SCMCertificateId = cast.StringPtr(upres.CertId)
	updateDomainMultiCertificatesReqBodyContent.CertName = cast.StringPtr(upres.CertName)
	updateDomainMultiCertificatesReqBodyContent = updateDomainMultiCertificatesReqBodyContent.MergeConfig(showDomainFullConfigResp.Configs)
	updateDomainMultiCertificatesReq := &hcCdnEx.UpdateDomainMultiCertificatesExRequest{
		Body: &hcCdnEx.UpdateDomainMultiCertificatesExRequestBody{
			Https: updateDomainMultiCertificatesReqBodyContent,
		},
	}
	updateDomainMultiCertificatesResp, err := d.sdkClient.UploadDomainMultiCertificatesEx(updateDomainMultiCertificatesReq)
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'cdn.UploadDomainMultiCertificatesEx'")
	}

	d.infos = append(d.infos, toStr("已更新加速域名配置", updateDomainMultiCertificatesResp))

	return nil
}

func (d *HuaweiCloudCDNDeployer) createSdkClient(accessKeyId, secretAccessKey, region string) (*hcCdnEx.Client, error) {
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

	client := hcCdnEx.NewClient(hcClient)
	return client, nil
}
