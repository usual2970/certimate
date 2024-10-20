package deployer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/global"
	cdn "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cdn/v2"
	cdnModel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cdn/v2/model"
	cdnRegion "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cdn/v2/region"

	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/utils/rand"
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

	client, err := d.createClient(access)
	if err != nil {
		return err
	}

	d.infos = append(d.infos, toStr("HuaweiCloudCdnClient 创建成功", nil))

	// 查询加速域名配置
	showDomainFullConfigReq := &cdnModel.ShowDomainFullConfigRequest{
		DomainName: getDeployString(d.option.DeployConfig, "domain"),
	}
	showDomainFullConfigResp, err := client.ShowDomainFullConfig(showDomainFullConfigReq)
	if err != nil {
		return err
	}

	d.infos = append(d.infos, toStr("已查询到加速域名配置", showDomainFullConfigResp))

	// 更新加速域名配置
	certName := fmt.Sprintf("%s-%s", d.option.DomainId, rand.RandStr(12))
	updateDomainMultiCertificatesReq := &cdnModel.UpdateDomainMultiCertificatesRequest{
		Body: &cdnModel.UpdateDomainMultiCertificatesRequestBody{
			Https: mergeHuaweiCloudCDNConfig(showDomainFullConfigResp.Configs, &cdnModel.UpdateDomainMultiCertificatesRequestBodyContent{
				DomainName:  getDeployString(d.option.DeployConfig, "domain"),
				HttpsSwitch: 1,
				CertName:    &certName,
				Certificate: &d.option.Certificate.Certificate,
				PrivateKey:  &d.option.Certificate.PrivateKey,
			}),
		},
	}
	updateDomainMultiCertificatesResp, err := client.UpdateDomainMultiCertificates(updateDomainMultiCertificatesReq)
	if err != nil {
		return err
	}

	d.infos = append(d.infos, toStr("已更新加速域名配置", updateDomainMultiCertificatesResp))

	return nil
}

func (d *HuaweiCloudCDNDeployer) createClient(access *domain.HuaweiCloudAccess) (*cdn.CdnClient, error) {
	auth, err := global.NewCredentialsBuilder().
		WithAk(access.AccessKeyId).
		WithSk(access.SecretAccessKey).
		SafeBuild()
	if err != nil {
		return nil, err
	}

	region, err := cdnRegion.SafeValueOf(access.Region)
	if err != nil {
		return nil, err
	}

	hcClient, err := cdn.CdnClientBuilder().
		WithRegion(region).
		WithCredential(auth).
		SafeBuild()
	if err != nil {
		return nil, err
	}

	client := cdn.NewCdnClient(hcClient)
	return client, nil
}

func mergeHuaweiCloudCDNConfig(src *cdnModel.ConfigsGetBody, dest *cdnModel.UpdateDomainMultiCertificatesRequestBodyContent) *cdnModel.UpdateDomainMultiCertificatesRequestBodyContent {
	if src == nil {
		return dest
	}

	// 华为云 API 中不传的字段表示使用默认值、而非保留原值，因此这里需要把原配置中的参数重新赋值回去
	// 而且蛋疼的是查询接口返回的数据结构和更新接口传入的参数结构不一致，需要做很多转化
	// REF: https://support.huaweicloud.com/api-cdn/ShowDomainFullConfig.html
	// REF: https://support.huaweicloud.com/api-cdn/UpdateDomainMultiCertificates.html

	if *src.OriginProtocol == "follow" {
		accessOriginWay := int32(1)
		dest.AccessOriginWay = &accessOriginWay
	} else if *src.OriginProtocol == "http" {
		accessOriginWay := int32(2)
		dest.AccessOriginWay = &accessOriginWay
	} else if *src.OriginProtocol == "https" {
		accessOriginWay := int32(3)
		dest.AccessOriginWay = &accessOriginWay
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
			http2 := int32(1)
			dest.Http2 = &http2
		}
	}

	return dest
}
