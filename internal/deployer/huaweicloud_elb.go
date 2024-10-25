package deployer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/global"
	hcElb "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/elb/v3"
	hcElbModel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/elb/v3/model"
	hcElbRegion "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/elb/v3/region"
	hcIam "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3"
	hcIamModel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3/model"
	hcIamRegion "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3/region"

	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	"github.com/usual2970/certimate/internal/pkg/utils/cast"
)

type HuaweiCloudELBDeployer struct {
	option *DeployerOption
	infos  []string

	sdkClient   *hcElb.ElbClient
	sslUploader uploader.Uploader
}

func NewHuaweiCloudELBDeployer(option *DeployerOption) (Deployer, error) {
	access := &domain.HuaweiCloudAccess{}
	if err := json.Unmarshal([]byte(option.Access), access); err != nil {
		return nil, err
	}

	client, err := (&HuaweiCloudELBDeployer{}).createSdkClient(
		access.AccessKeyId,
		access.SecretAccessKey,
		option.DeployConfig.GetConfigAsString("region"),
	)
	if err != nil {
		return nil, err
	}

	uploader, err := uploader.NewHuaweiCloudELBUploader(&uploader.HuaweiCloudELBUploaderConfig{
		AccessKeyId:     access.AccessKeyId,
		SecretAccessKey: access.SecretAccessKey,
		Region:          option.DeployConfig.GetConfigAsString("region"),
	})
	if err != nil {
		return nil, err
	}

	return &HuaweiCloudELBDeployer{
		option:      option,
		infos:       make([]string, 0),
		sdkClient:   client,
		sslUploader: uploader,
	}, nil
}

func (d *HuaweiCloudELBDeployer) GetID() string {
	return fmt.Sprintf("%s-%s", d.option.AccessRecord.GetString("name"), d.option.AccessRecord.Id)
}

func (d *HuaweiCloudELBDeployer) GetInfo() []string {
	return d.infos
}

func (d *HuaweiCloudELBDeployer) Deploy(ctx context.Context) error {
	switch d.option.DeployConfig.GetConfigAsString("resourceType") {
	case "certificate":
		if err := d.deployToCertificate(ctx); err != nil {
			return err
		}
	case "loadbalancer":
		if err := d.deployToLoadbalancer(ctx); err != nil {
			return err
		}
	case "listener":
		if err := d.deployToListener(ctx); err != nil {
			return err
		}
	default:
		return errors.New("unsupported resource type")
	}

	return nil
}

func (d *HuaweiCloudELBDeployer) createSdkClient(accessKeyId, secretAccessKey, region string) (*hcElb.ElbClient, error) {
	if region == "" {
		region = "cn-north-4" // ELB 服务默认区域：华北四北京
	}

	projectId, err := (&HuaweiCloudELBDeployer{}).getSdkProjectId(
		accessKeyId,
		secretAccessKey,
		region,
	)
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

func (u *HuaweiCloudELBDeployer) getSdkProjectId(accessKeyId, secretAccessKey, region string) (string, error) {
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
	if err != nil {
		return "", err
	}

	request := &hcIamModel.KeystoneListProjectsRequest{
		Name: &region,
	}
	response, err := client.KeystoneListProjects(request)
	if err != nil {
		return "", err
	} else if response.Projects == nil || len(*response.Projects) == 0 {
		return "", fmt.Errorf("no project found")
	}

	return (*response.Projects)[0].Id, nil
}

func (d *HuaweiCloudELBDeployer) deployToCertificate(ctx context.Context) error {
	hcCertId := d.option.DeployConfig.GetConfigAsString("certificateId")
	if hcCertId == "" {
		return errors.New("`certificateId` is required")
	}

	// 更新证书
	// REF: https://support.huaweicloud.com/api-elb/UpdateCertificate.html
	updateCertificateReq := &hcElbModel.UpdateCertificateRequest{
		CertificateId: hcCertId,
		Body: &hcElbModel.UpdateCertificateRequestBody{
			Certificate: &hcElbModel.UpdateCertificateOption{
				Certificate: cast.StringPtr(d.option.Certificate.Certificate),
				PrivateKey:  cast.StringPtr(d.option.Certificate.PrivateKey),
			},
		},
	}
	updateCertificateResp, err := d.sdkClient.UpdateCertificate(updateCertificateReq)
	if err != nil {
		return fmt.Errorf("failed to execute sdk request 'elb.UpdateCertificate': %w", err)
	}

	d.infos = append(d.infos, toStr("已更新 ELB 证书", updateCertificateResp))

	return nil
}

func (d *HuaweiCloudELBDeployer) deployToLoadbalancer(ctx context.Context) error {
	hcLoadbalancerId := d.option.DeployConfig.GetConfigAsString("loadbalancerId")
	if hcLoadbalancerId == "" {
		return errors.New("`loadbalancerId` is required")
	}

	// 查询负载均衡器详情
	// REF: https://support.huaweicloud.com/api-elb/ShowLoadBalancer.html
	showLoadBalancerReq := &hcElbModel.ShowLoadBalancerRequest{
		LoadbalancerId: hcLoadbalancerId,
	}
	showLoadBalancerResp, err := d.sdkClient.ShowLoadBalancer(showLoadBalancerReq)
	if err != nil {
		return fmt.Errorf("failed to execute sdk request 'elb.ShowLoadBalancer': %w", err)
	}

	d.infos = append(d.infos, toStr("已查询到 ELB 负载均衡器", showLoadBalancerResp))

	// 查询监听器列表
	// REF: https://support.huaweicloud.com/api-elb/ListListeners.html
	hcListenerIds := make([]string, 0)
	listListenersLimit := int32(2000)
	var listListenersMarker *string = nil
	for {
		listListenersReq := &hcElbModel.ListListenersRequest{
			Limit:          cast.Int32Ptr(listListenersLimit),
			Marker:         listListenersMarker,
			Protocol:       &[]string{"HTTPS", "TERMINATED_HTTPS"},
			LoadbalancerId: &[]string{showLoadBalancerResp.Loadbalancer.Id},
		}
		listListenersResp, err := d.sdkClient.ListListeners(listListenersReq)
		if err != nil {
			return fmt.Errorf("failed to execute sdk request 'elb.ListListeners': %w", err)
		}

		if listListenersResp.Listeners != nil {
			for _, listener := range *listListenersResp.Listeners {
				hcListenerIds = append(hcListenerIds, listener.Id)
			}
		}

		if listListenersResp.Listeners == nil || len(*listListenersResp.Listeners) < int(listListenersLimit) {
			break
		} else {
			listListenersMarker = listListenersResp.PageInfo.NextMarker
		}
	}

	d.infos = append(d.infos, toStr("已查询到 ELB 负载均衡器下的监听器", hcListenerIds))

	// 上传证书到 SCM
	uploadResult, err := d.sslUploader.Upload(ctx, d.option.Certificate.Certificate, d.option.Certificate.PrivateKey)
	if err != nil {
		return err
	}

	d.infos = append(d.infos, toStr("已上传证书", uploadResult))

	// 批量更新监听器证书
	var errs []error
	for _, hcListenerId := range hcListenerIds {
		if err := d.updateListenerCertificate(ctx, hcListenerId, uploadResult.CertId); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func (d *HuaweiCloudELBDeployer) deployToListener(ctx context.Context) error {
	hcListenerId := d.option.DeployConfig.GetConfigAsString("listenerId")
	if hcListenerId == "" {
		return errors.New("`listenerId` is required")
	}

	// 上传证书到 SCM
	uploadResult, err := d.sslUploader.Upload(ctx, d.option.Certificate.Certificate, d.option.Certificate.PrivateKey)
	if err != nil {
		return err
	}

	d.infos = append(d.infos, toStr("已上传证书", uploadResult))

	// 更新监听器证书
	if err := d.updateListenerCertificate(ctx, hcListenerId, uploadResult.CertId); err != nil {
		return err
	}

	return nil
}

func (d *HuaweiCloudELBDeployer) updateListenerCertificate(ctx context.Context, hcListenerId string, hcCertId string) error {
	// 查询监听器详情
	// REF: https://support.huaweicloud.com/api-elb/ShowListener.html
	showListenerReq := &hcElbModel.ShowListenerRequest{
		ListenerId: hcListenerId,
	}
	showListenerResp, err := d.sdkClient.ShowListener(showListenerReq)
	if err != nil {
		return fmt.Errorf("failed to execute sdk request 'elb.ShowListener': %w", err)
	}

	d.infos = append(d.infos, toStr("已查询到 ELB 监听器", showListenerResp))

	// 更新监听器
	// REF: https://support.huaweicloud.com/api-elb/UpdateListener.html
	updateListenerReq := &hcElbModel.UpdateListenerRequest{
		ListenerId: hcListenerId,
		Body: &hcElbModel.UpdateListenerRequestBody{
			Listener: &hcElbModel.UpdateListenerOption{
				DefaultTlsContainerRef: cast.StringPtr(hcCertId),
			},
		},
	}
	if showListenerResp.Listener.SniContainerRefs != nil {
		if len(showListenerResp.Listener.SniContainerRefs) > 0 {
			// 如果开启 SNI，需替换同 SAN 的证书
			sniCertIds := make([]string, 0)
			sniCertIds = append(sniCertIds, hcCertId)

			listOldCertificateReq := &hcElbModel.ListCertificatesRequest{
				Id: &showListenerResp.Listener.SniContainerRefs,
			}
			listOldCertificateResp, err := d.sdkClient.ListCertificates(listOldCertificateReq)
			if err != nil {
				return fmt.Errorf("failed to execute sdk request 'elb.ListCertificates': %w", err)
			}

			showNewCertificateReq := &hcElbModel.ShowCertificateRequest{
				CertificateId: hcCertId,
			}
			showNewCertificateResp, err := d.sdkClient.ShowCertificate(showNewCertificateReq)
			if err != nil {
				return fmt.Errorf("failed to execute sdk request 'elb.ShowCertificate': %w", err)
			}

			for _, certificate := range *listOldCertificateResp.Certificates {
				oldCertificate := certificate
				newCertificate := showNewCertificateResp.Certificate

				if oldCertificate.SubjectAlternativeNames != nil && newCertificate.SubjectAlternativeNames != nil {
					oldCertificateSans := oldCertificate.SubjectAlternativeNames
					newCertificateSans := newCertificate.SubjectAlternativeNames
					sort.Strings(*oldCertificateSans)
					sort.Strings(*newCertificateSans)
					if strings.Join(*oldCertificateSans, ";") == strings.Join(*newCertificateSans, ";") {
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
			updateListenerReq.Body.Listener.SniMatchAlgo = cast.StringPtr(showListenerResp.Listener.SniMatchAlgo)
		}
	}
	updateListenerResp, err := d.sdkClient.UpdateListener(updateListenerReq)
	if err != nil {
		return fmt.Errorf("failed to execute sdk request 'elb.UpdateListener': %w", err)
	}

	d.infos = append(d.infos, toStr("已更新 ELB 监听器", updateListenerResp))

	return nil
}
