package huaweicloudwaf

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/global"
	hciam "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3"
	hciamModel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3/model"
	hciamregion "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3/region"
	hcwaf "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/waf/v1"
	hcwafmodel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/waf/v1/model"
	hcwafregion "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/waf/v1/region"
	xerrors "github.com/pkg/errors"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	uploadersp "github.com/usual2970/certimate/internal/pkg/core/uploader/providers/huaweicloud-waf"
	hwsdk "github.com/usual2970/certimate/internal/pkg/vendors/huaweicloud-sdk"
)

type DeployerConfig struct {
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

type DeployerProvider struct {
	config      *DeployerConfig
	logger      *slog.Logger
	sdkClient   *hcwaf.WafClient
	sslUploader uploader.Uploader
}

var _ deployer.Deployer = (*DeployerProvider)(nil)

func NewDeployer(config *DeployerConfig) (*DeployerProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	client, err := createSdkClient(config.AccessKeyId, config.SecretAccessKey, config.Region)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk client")
	}

	uploader, err := uploadersp.NewUploader(&uploadersp.UploaderConfig{
		AccessKeyId:     config.AccessKeyId,
		SecretAccessKey: config.SecretAccessKey,
		Region:          config.Region,
	})
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create ssl uploader")
	}

	return &DeployerProvider{
		config:      config,
		logger:      slog.Default(),
		sdkClient:   client,
		sslUploader: uploader,
	}, nil
}

func (d *DeployerProvider) WithLogger(logger *slog.Logger) deployer.Deployer {
	if logger == nil {
		d.logger = slog.Default()
	} else {
		d.logger = logger
	}
	d.sslUploader.WithLogger(logger)
	return d
}

func (d *DeployerProvider) Deploy(ctx context.Context, certPem string, privkeyPem string) (*deployer.DeployResult, error) {
	// 上传证书到 WAF
	upres, err := d.sslUploader.Upload(ctx, certPem, privkeyPem)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to upload certificate file")
	} else {
		d.logger.Info("ssl certificate uploaded", slog.Any("result", upres))
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

func (d *DeployerProvider) deployToCertificate(ctx context.Context, certPem string, privkeyPem string) error {
	if d.config.CertificateId == "" {
		return errors.New("config `certificateId` is required")
	}

	// 查询证书
	// REF: https://support.huaweicloud.com/api-waf/ShowCertificate.html
	showCertificateReq := &hcwafmodel.ShowCertificateRequest{
		CertificateId: d.config.CertificateId,
	}
	showCertificateResp, err := d.sdkClient.ShowCertificate(showCertificateReq)
	d.logger.Debug("sdk request 'waf.ShowCertificate'", slog.Any("request", showCertificateReq), slog.Any("response", showCertificateResp))
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'waf.ShowCertificate'")
	}

	// 更新证书
	// REF: https://support.huaweicloud.com/api-waf/UpdateCertificate.html
	updateCertificateReq := &hcwafmodel.UpdateCertificateRequest{
		CertificateId: d.config.CertificateId,
		Body: &hcwafmodel.UpdateCertificateRequestBody{
			Name:    *showCertificateResp.Name,
			Content: hwsdk.StringPtr(certPem),
			Key:     hwsdk.StringPtr(privkeyPem),
		},
	}
	updateCertificateResp, err := d.sdkClient.UpdateCertificate(updateCertificateReq)
	d.logger.Debug("sdk request 'waf.UpdateCertificate'", slog.Any("request", updateCertificateReq), slog.Any("response", updateCertificateResp))
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'waf.UpdateCertificate'")
	}

	return nil
}

func (d *DeployerProvider) deployToCloudServer(ctx context.Context, certPem string, privkeyPem string) error {
	if d.config.Domain == "" {
		return errors.New("config `domain` is required")
	}

	// 上传证书到 WAF
	upres, err := d.sslUploader.Upload(ctx, certPem, privkeyPem)
	if err != nil {
		return xerrors.Wrap(err, "failed to upload certificate file")
	} else {
		d.logger.Info("ssl certificate uploaded", slog.Any("result", upres))
	}

	// 遍历查询云模式防护域名列表，获取防护域名 ID
	// REF: https://support.huaweicloud.com/api-waf/ListHost.html
	hostId := ""
	listHostPage := int32(1)
	listHostPageSize := int32(100)
	for {
		listHostReq := &hcwafmodel.ListHostRequest{
			Hostname: hwsdk.StringPtr(strings.TrimPrefix(d.config.Domain, "*")),
			Page:     hwsdk.Int32Ptr(listHostPage),
			Pagesize: hwsdk.Int32Ptr(listHostPageSize),
		}
		listHostResp, err := d.sdkClient.ListHost(listHostReq)
		d.logger.Debug("sdk request 'waf.ListHost'", slog.Any("request", listHostReq), slog.Any("response", listHostResp))
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
	updateHostReq := &hcwafmodel.UpdateHostRequest{
		InstanceId: hostId,
		Body: &hcwafmodel.UpdateHostRequestBody{
			Certificateid:   hwsdk.StringPtr(upres.CertId),
			Certificatename: hwsdk.StringPtr(upres.CertName),
		},
	}
	updateHostResp, err := d.sdkClient.UpdateHost(updateHostReq)
	d.logger.Debug("sdk request 'waf.UpdateHost'", slog.Any("request", updateHostReq), slog.Any("response", updateHostResp))
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'waf.UpdateHost'")
	}

	return nil
}

func (d *DeployerProvider) deployToPremiumHost(ctx context.Context, certPem string, privkeyPem string) error {
	if d.config.Domain == "" {
		return errors.New("config `domain` is required")
	}

	// 上传证书到 WAF
	upres, err := d.sslUploader.Upload(ctx, certPem, privkeyPem)
	if err != nil {
		return xerrors.Wrap(err, "failed to upload certificate file")
	} else {
		d.logger.Info("ssl certificate uploaded", slog.Any("result", upres))
	}

	// 遍历查询独享模式域名列表，获取防护域名 ID
	// REF: https://support.huaweicloud.com/api-waf/ListPremiumHost.html
	hostId := ""
	listPremiumHostPage := int32(1)
	listPremiumHostPageSize := int32(100)
	for {
		listPremiumHostReq := &hcwafmodel.ListPremiumHostRequest{
			Hostname: hwsdk.StringPtr(strings.TrimPrefix(d.config.Domain, "*")),
			Page:     hwsdk.StringPtr(fmt.Sprintf("%d", listPremiumHostPage)),
			Pagesize: hwsdk.StringPtr(fmt.Sprintf("%d", listPremiumHostPageSize)),
		}
		listPremiumHostResp, err := d.sdkClient.ListPremiumHost(listPremiumHostReq)
		d.logger.Debug("sdk request 'waf.ListPremiumHost'", slog.Any("request", listPremiumHostReq), slog.Any("response", listPremiumHostResp))
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
	updatePremiumHostReq := &hcwafmodel.UpdatePremiumHostRequest{
		HostId: hostId,
		Body: &hcwafmodel.UpdatePremiumHostRequestBody{
			Certificateid:   hwsdk.StringPtr(upres.CertId),
			Certificatename: hwsdk.StringPtr(upres.CertName),
		},
	}
	updatePremiumHostResp, err := d.sdkClient.UpdatePremiumHost(updatePremiumHostReq)
	d.logger.Debug("sdk request 'waf.UpdatePremiumHost'", slog.Any("request", updatePremiumHostReq), slog.Any("response", updatePremiumHostResp))
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'waf.UpdatePremiumHost'")
	}

	return nil
}

func createSdkClient(accessKeyId, secretAccessKey, region string) (*hcwaf.WafClient, error) {
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

	hcRegion, err := hcwafregion.SafeValueOf(region)
	if err != nil {
		return nil, err
	}

	hcClient, err := hcwaf.WafClientBuilder().
		WithRegion(hcRegion).
		WithCredential(auth).
		SafeBuild()
	if err != nil {
		return nil, err
	}

	client := hcwaf.NewWafClient(hcClient)
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

	hcRegion, err := hciamregion.SafeValueOf(region)
	if err != nil {
		return "", err
	}

	hcClient, err := hciam.IamClientBuilder().
		WithRegion(hcRegion).
		WithCredential(auth).
		SafeBuild()
	if err != nil {
		return "", err
	}

	client := hciam.NewIamClient(hcClient)

	request := &hciamModel.KeystoneListProjectsRequest{
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
