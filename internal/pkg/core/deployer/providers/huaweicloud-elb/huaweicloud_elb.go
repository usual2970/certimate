package huaweicloudelb

import (
	"context"
	"errors"
	"fmt"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/global"
	hcElb "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/elb/v3"
	hcElbModel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/elb/v3/model"
	hcElbRegion "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/elb/v3/region"
	hcIam "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3"
	hcIamModel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3/model"
	hcIamRegion "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3/region"
	xerrors "github.com/pkg/errors"
	"golang.org/x/exp/slices"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/logger"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	uploaderp "github.com/usual2970/certimate/internal/pkg/core/uploader/providers/huaweicloud-elb"
	hwsdk "github.com/usual2970/certimate/internal/pkg/vendors/huaweicloud-sdk"
)

type HuaweiCloudELBDeployerConfig struct {
	// 华为云 AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// 华为云 SecretAccessKey。
	SecretAccessKey string `json:"secretAccessKey"`
	// 华为云区域。
	Region string `json:"region"`
	// 部署资源类型。
	ResourceType DeployResourceType `json:"resourceType"`
	// 证书 ID。
	// 部署资源类型为 [DEPLOY_RESOURCE_CERTIFICATE] 时必填。
	CertificateId string `json:"certificateId,omitempty"`
	// 负载均衡器 ID。
	// 部署资源类型为 [DEPLOY_RESOURCE_LOADBALANCER] 时必填。
	LoadbalancerId string `json:"loadbalancerId,omitempty"`
	// 负载均衡监听 ID。
	// 部署资源类型为 [DEPLOY_RESOURCE_LISTENER] 时必填。
	ListenerId string `json:"listenerId,omitempty"`
}

type HuaweiCloudELBDeployer struct {
	config      *HuaweiCloudELBDeployerConfig
	logger      logger.Logger
	sdkClient   *hcElb.ElbClient
	sslUploader uploader.Uploader
}

var _ deployer.Deployer = (*HuaweiCloudELBDeployer)(nil)

func New(config *HuaweiCloudELBDeployerConfig) (*HuaweiCloudELBDeployer, error) {
	return NewWithLogger(config, logger.NewNilLogger())
}

func NewWithLogger(config *HuaweiCloudELBDeployerConfig, logger logger.Logger) (*HuaweiCloudELBDeployer, error) {
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

	uploader, err := uploaderp.New(&uploaderp.HuaweiCloudELBUploaderConfig{
		AccessKeyId:     config.AccessKeyId,
		SecretAccessKey: config.SecretAccessKey,
		Region:          config.Region,
	})
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create ssl uploader")
	}

	return &HuaweiCloudELBDeployer{
		logger:      logger,
		config:      config,
		sdkClient:   client,
		sslUploader: uploader,
	}, nil
}

func (d *HuaweiCloudELBDeployer) Deploy(ctx context.Context, certPem string, privkeyPem string) (*deployer.DeployResult, error) {
	// 上传证书到 SCM
	upres, err := d.sslUploader.Upload(ctx, certPem, privkeyPem)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to upload certificate file")
	}

	d.logger.Logt("certificate file uploaded", upres)

	// 根据部署资源类型决定部署方式
	switch d.config.ResourceType {
	case DEPLOY_RESOURCE_CERTIFICATE:
		if err := d.deployToCertificate(ctx, certPem, privkeyPem); err != nil {
			return nil, err
		}

	case DEPLOY_RESOURCE_LOADBALANCER:
		if err := d.deployToLoadbalancer(ctx, certPem, privkeyPem); err != nil {
			return nil, err
		}

	case DEPLOY_RESOURCE_LISTENER:
		if err := d.deployToListener(ctx, certPem, privkeyPem); err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("unsupported resource type: %s", d.config.ResourceType)
	}

	return &deployer.DeployResult{}, nil
}

func (d *HuaweiCloudELBDeployer) deployToCertificate(ctx context.Context, certPem string, privkeyPem string) error {
	if d.config.CertificateId == "" {
		return errors.New("config `certificateId` is required")
	}

	// 更新证书
	// REF: https://support.huaweicloud.com/api-elb/UpdateCertificate.html
	updateCertificateReq := &hcElbModel.UpdateCertificateRequest{
		CertificateId: d.config.CertificateId,
		Body: &hcElbModel.UpdateCertificateRequestBody{
			Certificate: &hcElbModel.UpdateCertificateOption{
				Certificate: hwsdk.StringPtr(certPem),
				PrivateKey:  hwsdk.StringPtr(privkeyPem),
			},
		},
	}
	updateCertificateResp, err := d.sdkClient.UpdateCertificate(updateCertificateReq)
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'elb.UpdateCertificate'")
	}

	d.logger.Logt("已更新 ELB 证书", updateCertificateResp)

	return nil
}

func (d *HuaweiCloudELBDeployer) deployToLoadbalancer(ctx context.Context, certPem string, privkeyPem string) error {
	if d.config.LoadbalancerId == "" {
		return errors.New("config `loadbalancerId` is required")
	}

	// 查询负载均衡器详情
	// REF: https://support.huaweicloud.com/api-elb/ShowLoadBalancer.html
	showLoadBalancerReq := &hcElbModel.ShowLoadBalancerRequest{
		LoadbalancerId: d.config.LoadbalancerId,
	}
	showLoadBalancerResp, err := d.sdkClient.ShowLoadBalancer(showLoadBalancerReq)
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'elb.ShowLoadBalancer'")
	}

	d.logger.Logt("已查询到 ELB 负载均衡器", showLoadBalancerResp)

	// 查询监听器列表
	// REF: https://support.huaweicloud.com/api-elb/ListListeners.html
	listenerIds := make([]string, 0)
	listListenersLimit := int32(2000)
	var listListenersMarker *string = nil
	for {
		listListenersReq := &hcElbModel.ListListenersRequest{
			Limit:          hwsdk.Int32Ptr(listListenersLimit),
			Marker:         listListenersMarker,
			Protocol:       &[]string{"HTTPS", "TERMINATED_HTTPS"},
			LoadbalancerId: &[]string{showLoadBalancerResp.Loadbalancer.Id},
		}
		listListenersResp, err := d.sdkClient.ListListeners(listListenersReq)
		if err != nil {
			return xerrors.Wrap(err, "failed to execute sdk request 'elb.ListListeners'")
		}

		if listListenersResp.Listeners != nil {
			for _, listener := range *listListenersResp.Listeners {
				listenerIds = append(listenerIds, listener.Id)
			}
		}

		if listListenersResp.Listeners == nil || len(*listListenersResp.Listeners) < int(listListenersLimit) {
			break
		} else {
			listListenersMarker = listListenersResp.PageInfo.NextMarker
		}
	}

	d.logger.Logt("已查询到 ELB 负载均衡器下的监听器", listenerIds)

	// 上传证书到 SCM
	upres, err := d.sslUploader.Upload(ctx, certPem, privkeyPem)
	if err != nil {
		return xerrors.Wrap(err, "failed to upload certificate file")
	}

	d.logger.Logt("certificate file uploaded", upres)

	// 遍历更新监听器证书
	if len(listenerIds) == 0 {
		return errors.New("listener not found")
	} else {
		var errs []error

		for _, listenerId := range listenerIds {
			if err := d.modifyListenerCertificate(ctx, listenerId, upres.CertId); err != nil {
				errs = append(errs, err)
			}
		}

		if len(errs) > 0 {
			return errors.Join(errs...)
		}
	}

	return nil
}

func (d *HuaweiCloudELBDeployer) deployToListener(ctx context.Context, certPem string, privkeyPem string) error {
	if d.config.ListenerId == "" {
		return errors.New("config `listenerId` is required")
	}

	// 上传证书到 SCM
	upres, err := d.sslUploader.Upload(ctx, certPem, privkeyPem)
	if err != nil {
		return xerrors.Wrap(err, "failed to upload certificate file")
	}

	d.logger.Logt("certificate file uploaded", upres)

	// 更新监听器证书
	if err := d.modifyListenerCertificate(ctx, d.config.ListenerId, upres.CertId); err != nil {
		return err
	}

	return nil
}

func (d *HuaweiCloudELBDeployer) modifyListenerCertificate(ctx context.Context, cloudListenerId string, cloudCertId string) error {
	// 查询监听器详情
	// REF: https://support.huaweicloud.com/api-elb/ShowListener.html
	showListenerReq := &hcElbModel.ShowListenerRequest{
		ListenerId: cloudListenerId,
	}
	showListenerResp, err := d.sdkClient.ShowListener(showListenerReq)
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'elb.ShowListener'")
	}

	d.logger.Logt("已查询到 ELB 监听器", showListenerResp)

	// 更新监听器
	// REF: https://support.huaweicloud.com/api-elb/UpdateListener.html
	updateListenerReq := &hcElbModel.UpdateListenerRequest{
		ListenerId: cloudListenerId,
		Body: &hcElbModel.UpdateListenerRequestBody{
			Listener: &hcElbModel.UpdateListenerOption{
				DefaultTlsContainerRef: hwsdk.StringPtr(cloudCertId),
			},
		},
	}
	if showListenerResp.Listener.SniContainerRefs != nil {
		if len(showListenerResp.Listener.SniContainerRefs) > 0 {
			// 如果开启 SNI，需替换同 SAN 的证书
			sniCertIds := make([]string, 0)
			sniCertIds = append(sniCertIds, cloudCertId)

			listOldCertificateReq := &hcElbModel.ListCertificatesRequest{
				Id: &showListenerResp.Listener.SniContainerRefs,
			}
			listOldCertificateResp, err := d.sdkClient.ListCertificates(listOldCertificateReq)
			if err != nil {
				return xerrors.Wrap(err, "failed to execute sdk request 'elb.ListCertificates'")
			}

			showNewCertificateReq := &hcElbModel.ShowCertificateRequest{
				CertificateId: cloudCertId,
			}
			showNewCertificateResp, err := d.sdkClient.ShowCertificate(showNewCertificateReq)
			if err != nil {
				return xerrors.Wrap(err, "failed to execute sdk request 'elb.ShowCertificate'")
			}

			for _, certificate := range *listOldCertificateResp.Certificates {
				oldCertificate := certificate
				newCertificate := showNewCertificateResp.Certificate

				if oldCertificate.SubjectAlternativeNames != nil && newCertificate.SubjectAlternativeNames != nil {
					if slices.Equal(*oldCertificate.SubjectAlternativeNames, *newCertificate.SubjectAlternativeNames) {
						continue
					}
				} else {
					if oldCertificate.Domain == newCertificate.Domain {
						continue
					}
				}

				sniCertIds = append(sniCertIds, certificate.Id)
			}

			updateListenerReq.Body.Listener.SniContainerRefs = &sniCertIds
		}

		if showListenerResp.Listener.SniMatchAlgo != "" {
			updateListenerReq.Body.Listener.SniMatchAlgo = hwsdk.StringPtr(showListenerResp.Listener.SniMatchAlgo)
		}
	}
	updateListenerResp, err := d.sdkClient.UpdateListener(updateListenerReq)
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'elb.UpdateListener'")
	}

	d.logger.Logt("已更新 ELB 监听器", updateListenerResp)

	return nil
}

func createSdkClient(accessKeyId, secretAccessKey, region string) (*hcElb.ElbClient, error) {
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

	hcRegion, err := hcElbRegion.SafeValueOf(region)
	if err != nil {
		return nil, err
	}

	hcClient, err := hcElb.ElbClientBuilder().
		WithRegion(hcRegion).
		WithCredential(auth).
		SafeBuild()
	if err != nil {
		return nil, err
	}

	client := hcElb.NewElbClient(hcClient)
	return client, nil
}

func getSdkProjectId(accessKeyId, secretAccessKey, region string) (string, error) {
	if region == "" {
		region = "cn-north-4" // IAM 服务默认区域：华北四北京
	}

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
