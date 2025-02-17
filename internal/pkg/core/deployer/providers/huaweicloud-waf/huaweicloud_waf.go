package huaweicloudwaf

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/global"
	hcIam "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3"
	hcIamModel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3/model"
	hcIamRegion "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3/region"
	hcWaf "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/waf/v1"
	hcWafModel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/waf/v1/model"
	hcWafRegion "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/waf/v1/region"
	xerrors "github.com/pkg/errors"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/logger"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	uploadersp "github.com/usual2970/certimate/internal/pkg/core/uploader/providers/huaweicloud-waf"
	hwsdk "github.com/usual2970/certimate/internal/pkg/vendors/huaweicloud-sdk"
)

type HuaweiCloudWAFDeployerConfig struct {
	// 华为云 AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// 华为云 SecretAccessKey。
	SecretAccessKey string `json:"secretAccessKey"`
	// 华为云区域。
	Region string `json:"region"`
	// 部署资源类型。
	ResourceType ResourceType `json:"resourceType"`
	// 证书 ID。
	// 部署资源类型为 [RESOURCE_TYPE_CERTIFICATE] 时必填。
	CertificateId string `json:"certificateId,omitempty"`
	// 防护域名（支持泛域名）。
	// 部署资源类型为 [RESOURCE_TYPE_CLOUDSERVER]、[RESOURCE_TYPE_PREMIUMHOST] 时必填。
	Domain string `json:"domain,omitempty"`
}

type HuaweiCloudWAFDeployer struct {
	config      *HuaweiCloudWAFDeployerConfig
	logger      logger.Logger
	sdkClient   *hcWaf.WafClient
	sslUploader uploader.Uploader
}

var _ deployer.Deployer = (*HuaweiCloudWAFDeployer)(nil)

func New(config *HuaweiCloudWAFDeployerConfig) (*HuaweiCloudWAFDeployer, error) {
	return NewWithLogger(config, logger.NewNilLogger())
}

func NewWithLogger(config *HuaweiCloudWAFDeployerConfig, logger logger.Logger) (*HuaweiCloudWAFDeployer, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}

	if logger == nil {
		return nil, errors.New("logger is nil")
	}

	client, err := createSdkClient(config.AccessKeyId, config.SecretAccessKey, config.Region)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk client")
	}

	uploader, err := uploadersp.New(&uploadersp.HuaweiCloudWAFUploaderConfig{
		AccessKeyId:     config.AccessKeyId,
		SecretAccessKey: config.SecretAccessKey,
		Region:          config.Region,
	})
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create ssl uploader")
	}

	return &HuaweiCloudWAFDeployer{
		logger:      logger,
		config:      config,
		sdkClient:   client,
		sslUploader: uploader,
	}, nil
}

func (d *HuaweiCloudWAFDeployer) Deploy(ctx context.Context, certPem string, privkeyPem string) (*deployer.DeployResult, error) {
	// 上传证书到 WAF
	upres, err := d.sslUploader.Upload(ctx, certPem, privkeyPem)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to upload certificate file")
	} else {
		d.logger.Logt("certificate file uploaded", upres)
	}

	// 根据部署资源类型决定部署方式
	switch d.config.ResourceType {
	case RESOURCE_TYPE_CERTIFICATE:
		if err := d.deployToCertificate(ctx, certPem, privkeyPem); err != nil {
			return nil, err
		}

	case RESOURCE_TYPE_CLOUDSERVER:
		if err := d.deployToCloudServer(ctx, certPem, privkeyPem); err != nil {
			return nil, err
		}

	case RESOURCE_TYPE_PREMIUMHOST:
		if err := d.deployToPremiumHost(ctx, certPem, privkeyPem); err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("unsupported resource type: %s", d.config.ResourceType)
	}

	return &deployer.DeployResult{}, nil
}

func (d *HuaweiCloudWAFDeployer) deployToCertificate(ctx context.Context, certPem string, privkeyPem string) error {
	if d.config.CertificateId == "" {
		return errors.New("config `certificateId` is required")
	}

	// 查询证书
	// REF: https://support.huaweicloud.com/api-waf/ShowCertificate.html
	showCertificateReq := &hcWafModel.ShowCertificateRequest{
		CertificateId: d.config.CertificateId,
	}
	showCertificateResp, err := d.sdkClient.ShowCertificate(showCertificateReq)
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'waf.ShowCertificate'")
	} else {
		d.logger.Logt("已获取 WAF 证书", showCertificateResp)
	}

	// 更新证书
	// REF: https://support.huaweicloud.com/api-waf/UpdateCertificate.html
	updateCertificateReq := &hcWafModel.UpdateCertificateRequest{
		CertificateId: d.config.CertificateId,
		Body: &hcWafModel.UpdateCertificateRequestBody{
			Name:    *showCertificateResp.Name,
			Content: hwsdk.StringPtr(certPem),
			Key:     hwsdk.StringPtr(privkeyPem),
		},
	}
	updateCertificateResp, err := d.sdkClient.UpdateCertificate(updateCertificateReq)
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'waf.UpdateCertificate'")
	} else {
		d.logger.Logt("已更新 WAF 证书", updateCertificateResp)
	}

	return nil
}

func (d *HuaweiCloudWAFDeployer) deployToCloudServer(ctx context.Context, certPem string, privkeyPem string) error {
	if d.config.Domain == "" {
		return errors.New("config `domain` is required")
	}

	// 上传证书到 WAF
	upres, err := d.sslUploader.Upload(ctx, certPem, privkeyPem)
	if err != nil {
		return xerrors.Wrap(err, "failed to upload certificate file")
	} else {
		d.logger.Logt("certificate file uploaded", upres)
	}

	// 遍历查询云模式防护域名列表，获取防护域名 ID
	// REF: https://support.huaweicloud.com/api-waf/ListHost.html
	hostId := ""
	listHostPage := int32(1)
	listHostPageSize := int32(100)
	for {
		listHostReq := &hcWafModel.ListHostRequest{
			Hostname: hwsdk.StringPtr(strings.TrimPrefix(d.config.Domain, "*")),
			Page:     hwsdk.Int32Ptr(listHostPage),
			Pagesize: hwsdk.Int32Ptr(listHostPageSize),
		}
		listHostResp, err := d.sdkClient.ListHost(listHostReq)
		if err != nil {
			return xerrors.Wrap(err, "failed to execute sdk request 'waf.ListHost'")
		}

		if listHostResp.Items != nil {
			for _, hostItem := range *listHostResp.Items {
				if strings.TrimPrefix(d.config.Domain, "*") == *hostItem.Hostname {
					hostId = *hostItem.Id
					break
				}
			}
		}

		if listHostResp.Items == nil || len(*listHostResp.Items) < int(listHostPageSize) {
			break
		} else {
			listHostPage++
		}
	}
	if hostId == "" {
		return errors.New("host not found")
	}

	// 更新云模式防护域名的配置
	// REF: https://support.huaweicloud.com/api-waf/UpdateHost.html
	updateHostReq := &hcWafModel.UpdateHostRequest{
		InstanceId: hostId,
		Body: &hcWafModel.UpdateHostRequestBody{
			Certificateid:   hwsdk.StringPtr(upres.CertId),
			Certificatename: hwsdk.StringPtr(upres.CertName),
		},
	}
	updateHostResp, err := d.sdkClient.UpdateHost(updateHostReq)
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'waf.UpdateHost'")
	} else {
		d.logger.Logt("已更新云模式防护域名的配置", updateHostResp)
	}

	return nil
}

func (d *HuaweiCloudWAFDeployer) deployToPremiumHost(ctx context.Context, certPem string, privkeyPem string) error {
	if d.config.Domain == "" {
		return errors.New("config `domain` is required")
	}

	// 上传证书到 WAF
	upres, err := d.sslUploader.Upload(ctx, certPem, privkeyPem)
	if err != nil {
		return xerrors.Wrap(err, "failed to upload certificate file")
	} else {
		d.logger.Logt("certificate file uploaded", upres)
	}

	// 遍历查询独享模式域名列表，获取防护域名 ID
	// REF: https://support.huaweicloud.com/api-waf/ListPremiumHost.html
	hostId := ""
	listPremiumHostPage := int32(1)
	listPremiumHostPageSize := int32(100)
	for {
		listPremiumHostReq := &hcWafModel.ListPremiumHostRequest{
			Hostname: hwsdk.StringPtr(strings.TrimPrefix(d.config.Domain, "*")),
			Page:     hwsdk.StringPtr(fmt.Sprintf("%d", listPremiumHostPage)),
			Pagesize: hwsdk.StringPtr(fmt.Sprintf("%d", listPremiumHostPageSize)),
		}
		listPremiumHostResp, err := d.sdkClient.ListPremiumHost(listPremiumHostReq)
		if err != nil {
			return xerrors.Wrap(err, "failed to execute sdk request 'waf.ListPremiumHost'")
		}

		if listPremiumHostResp.Items != nil {
			for _, hostItem := range *listPremiumHostResp.Items {
				if strings.TrimPrefix(d.config.Domain, "*") == *hostItem.Hostname {
					hostId = *hostItem.Id
					break
				}
			}
		}

		if listPremiumHostResp.Items == nil || len(*listPremiumHostResp.Items) < int(listPremiumHostPageSize) {
			break
		} else {
			listPremiumHostPage++
		}
	}
	if hostId == "" {
		return errors.New("host not found")
	}

	// 修改独享模式域名配置
	// REF: https://support.huaweicloud.com/api-waf/UpdatePremiumHost.html
	updatePremiumHostReq := &hcWafModel.UpdatePremiumHostRequest{
		HostId: hostId,
		Body: &hcWafModel.UpdatePremiumHostRequestBody{
			Certificateid:   hwsdk.StringPtr(upres.CertId),
			Certificatename: hwsdk.StringPtr(upres.CertName),
		},
	}
	updatePremiumHostResp, err := d.sdkClient.UpdatePremiumHost(updatePremiumHostReq)
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'waf.UpdatePremiumHost'")
	} else {
		d.logger.Logt("已修改独享模式域名配置", updatePremiumHostResp)
	}

	return nil
}

func createSdkClient(accessKeyId, secretAccessKey, region string) (*hcWaf.WafClient, error) {
	projectId, err := getSdkProjectId(accessKeyId, secretAccessKey, region)
	if err != nil {
		return nil, err
	}

	auth, err := basic.NewCredentialsBuilder().
		WithAk(accessKeyId).
		WithSk(secretAccessKey).
		WithProjectId(projectId).
		SafeBuild()
	if err != nil {
		return nil, err
	}

	hcRegion, err := hcWafRegion.SafeValueOf(region)
	if err != nil {
		return nil, err
	}

	hcClient, err := hcWaf.WafClientBuilder().
		WithRegion(hcRegion).
		WithCredential(auth).
		SafeBuild()
	if err != nil {
		return nil, err
	}

	client := hcWaf.NewWafClient(hcClient)
	return client, nil
}

func getSdkProjectId(accessKeyId, secretAccessKey, region string) (string, error) {
	auth, err := global.NewCredentialsBuilder().
		WithAk(accessKeyId).
		WithSk(secretAccessKey).
		SafeBuild()
	if err != nil {
		return "", err
	}

	hcRegion, err := hcIamRegion.SafeValueOf(region)
	if err != nil {
		return "", err
	}

	hcClient, err := hcIam.IamClientBuilder().
		WithRegion(hcRegion).
		WithCredential(auth).
		SafeBuild()
	if err != nil {
		return "", err
	}

	client := hcIam.NewIamClient(hcClient)

	request := &hcIamModel.KeystoneListProjectsRequest{
		Name: &region,
	}
	response, err := client.KeystoneListProjects(request)
	if err != nil {
		return "", err
	} else if response.Projects == nil || len(*response.Projects) == 0 {
		return "", errors.New("no project found")
	}

	return (*response.Projects)[0].Id, nil
}
