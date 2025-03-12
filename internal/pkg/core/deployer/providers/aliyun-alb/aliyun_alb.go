package aliyunalb

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	aliyunAlb "github.com/alibabacloud-go/alb-20200616/v2/client"
	aliyunCas "github.com/alibabacloud-go/cas-20200407/v3/client"
	aliyunOpen "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	xerrors "github.com/pkg/errors"
	"golang.org/x/exp/slices"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/logger"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	uploadersp "github.com/usual2970/certimate/internal/pkg/core/uploader/providers/aliyun-cas"
)

type DeployerConfig struct {
	// 阿里云 AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// 阿里云 AccessKeySecret。
	AccessKeySecret string `json:"accessKeySecret"`
	// 阿里云地域。
	Region string `json:"region"`
	// 部署资源类型。
	ResourceType ResourceType `json:"resourceType"`
	// 负载均衡实例 ID。
	// 部署资源类型为 [RESOURCE_TYPE_LOADBALANCER] 时必填。
	LoadbalancerId string `json:"loadbalancerId,omitempty"`
	// 负载均衡监听 ID。
	// 部署资源类型为 [RESOURCE_TYPE_LISTENER] 时必填。
	ListenerId string `json:"listenerId,omitempty"`
	// SNI 域名（支持泛域名）。
	// 部署资源类型为 [RESOURCE_TYPE_LOADBALANCER]、[RESOURCE_TYPE_LISTENER] 时选填。
	Domain string `json:"domain,omitempty"`
}

type DeployerProvider struct {
	config      *DeployerConfig
	logger      logger.Logger
	sdkClients  *wSdkClients
	sslUploader uploader.Uploader
}

var _ deployer.Deployer = (*DeployerProvider)(nil)

type wSdkClients struct {
	alb *aliyunAlb.Client
	cas *aliyunCas.Client
}

func NewDeployer(config *DeployerConfig) (*DeployerProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	clients, err := createSdkClients(config.AccessKeyId, config.AccessKeySecret, config.Region)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk clients")
	}

	uploader, err := createSslUploader(config.AccessKeyId, config.AccessKeySecret, config.Region)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create ssl uploader")
	}

	return &DeployerProvider{
		config:      config,
		logger:      logger.NewNilLogger(),
		sdkClients:  clients,
		sslUploader: uploader,
	}, nil
}

func (d *DeployerProvider) WithLogger(logger logger.Logger) *DeployerProvider {
	d.logger = logger
	return d
}

func (d *DeployerProvider) Deploy(ctx context.Context, certPem string, privkeyPem string) (*deployer.DeployResult, error) {
	// 上传证书到 CAS
	upres, err := d.sslUploader.Upload(ctx, certPem, privkeyPem)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to upload certificate file")
	}

	d.logger.Logt("certificate file uploaded", upres)

	// 根据部署资源类型决定部署方式
	switch d.config.ResourceType {
	case RESOURCE_TYPE_LOADBALANCER:
		if err := d.deployToLoadbalancer(ctx, upres.CertId); err != nil {
			return nil, err
		}

	case RESOURCE_TYPE_LISTENER:
		if err := d.deployToListener(ctx, upres.CertId); err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("unsupported resource type: %s", d.config.ResourceType)
	}

	return &deployer.DeployResult{}, nil
}

func (d *DeployerProvider) deployToLoadbalancer(ctx context.Context, cloudCertId string) error {
	if d.config.LoadbalancerId == "" {
		return errors.New("config `loadbalancerId` is required")
	}

	// 查询负载均衡实例的详细信息
	// REF: https://help.aliyun.com/zh/slb/application-load-balancer/developer-reference/api-alb-2020-06-16-getloadbalancerattribute
	getLoadBalancerAttributeReq := &aliyunAlb.GetLoadBalancerAttributeRequest{
		LoadBalancerId: tea.String(d.config.LoadbalancerId),
	}
	getLoadBalancerAttributeResp, err := d.sdkClients.alb.GetLoadBalancerAttribute(getLoadBalancerAttributeReq)
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'alb.GetLoadBalancerAttribute'")
	}

	d.logger.Logt("已查询到 ALB 负载均衡实例", getLoadBalancerAttributeResp)

	// 查询 HTTPS 监听列表
	// REF: https://help.aliyun.com/zh/slb/application-load-balancer/developer-reference/api-alb-2020-06-16-listlisteners
	listenerIds := make([]string, 0)
	listListenersLimit := int32(100)
	var listListenersToken *string = nil
	for {
		listListenersReq := &aliyunAlb.ListListenersRequest{
			MaxResults:       tea.Int32(listListenersLimit),
			NextToken:        listListenersToken,
			LoadBalancerIds:  []*string{tea.String(d.config.LoadbalancerId)},
			ListenerProtocol: tea.String("HTTPS"),
		}
		listListenersResp, err := d.sdkClients.alb.ListListeners(listListenersReq)
		if err != nil {
			return xerrors.Wrap(err, "failed to execute sdk request 'alb.ListListeners'")
		}

		if listListenersResp.Body.Listeners != nil {
			for _, listener := range listListenersResp.Body.Listeners {
				listenerIds = append(listenerIds, *listener.ListenerId)
			}
		}

		if len(listListenersResp.Body.Listeners) == 0 || listListenersResp.Body.NextToken == nil {
			break
		} else {
			listListenersToken = listListenersResp.Body.NextToken
		}
	}

	d.logger.Logt("已查询到 ALB 负载均衡实例下的全部 HTTPS 监听", listenerIds)

	// 查询 QUIC 监听列表
	// REF: https://help.aliyun.com/zh/slb/application-load-balancer/developer-reference/api-alb-2020-06-16-listlisteners
	listListenersToken = nil
	for {
		listListenersReq := &aliyunAlb.ListListenersRequest{
			MaxResults:       tea.Int32(listListenersLimit),
			NextToken:        listListenersToken,
			LoadBalancerIds:  []*string{tea.String(d.config.LoadbalancerId)},
			ListenerProtocol: tea.String("QUIC"),
		}
		listListenersResp, err := d.sdkClients.alb.ListListeners(listListenersReq)
		if err != nil {
			return xerrors.Wrap(err, "failed to execute sdk request 'alb.ListListeners'")
		}

		if listListenersResp.Body.Listeners != nil {
			for _, listener := range listListenersResp.Body.Listeners {
				listenerIds = append(listenerIds, *listener.ListenerId)
			}
		}

		if len(listListenersResp.Body.Listeners) == 0 || listListenersResp.Body.NextToken == nil {
			break
		} else {
			listListenersToken = listListenersResp.Body.NextToken
		}
	}

	d.logger.Logt("已查询到 ALB 负载均衡实例下的全部 QUIC 监听", listenerIds)

	// 遍历更新监听证书
	if len(listenerIds) == 0 {
		return errors.New("listener not found")
	} else {
		var errs []error

		for _, listenerId := range listenerIds {
			if err := d.updateListenerCertificate(ctx, listenerId, cloudCertId); err != nil {
				errs = append(errs, err)
			}
		}

		if len(errs) > 0 {
			return errors.Join(errs...)
		}
	}

	return nil
}

func (d *DeployerProvider) deployToListener(ctx context.Context, cloudCertId string) error {
	if d.config.ListenerId == "" {
		return errors.New("config `listenerId` is required")
	}

	// 更新监听
	if err := d.updateListenerCertificate(ctx, d.config.ListenerId, cloudCertId); err != nil {
		return err
	}

	return nil
}

func (d *DeployerProvider) updateListenerCertificate(ctx context.Context, cloudListenerId string, cloudCertId string) error {
	// 查询监听的属性
	// REF: https://help.aliyun.com/zh/slb/application-load-balancer/developer-reference/api-alb-2020-06-16-getlistenerattribute
	getListenerAttributeReq := &aliyunAlb.GetListenerAttributeRequest{
		ListenerId: tea.String(cloudListenerId),
	}
	getListenerAttributeResp, err := d.sdkClients.alb.GetListenerAttribute(getListenerAttributeReq)
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'alb.GetListenerAttribute'")
	}

	d.logger.Logt("已查询到 ALB 监听配置", getListenerAttributeResp)

	if d.config.Domain == "" {
		// 未指定 SNI，只需部署到监听器

		// 修改监听的属性
		// REF: https://help.aliyun.com/zh/slb/application-load-balancer/developer-reference/api-alb-2020-06-16-updatelistenerattribute
		updateListenerAttributeReq := &aliyunAlb.UpdateListenerAttributeRequest{
			ListenerId: tea.String(cloudListenerId),
			Certificates: []*aliyunAlb.UpdateListenerAttributeRequestCertificates{{
				CertificateId: tea.String(cloudCertId),
			}},
		}
		updateListenerAttributeResp, err := d.sdkClients.alb.UpdateListenerAttribute(updateListenerAttributeReq)
		if err != nil {
			return xerrors.Wrap(err, "failed to execute sdk request 'alb.UpdateListenerAttribute'")
		}

		d.logger.Logt("已更新 ALB 监听配置", updateListenerAttributeResp)
	} else {
		// 指定 SNI，需部署到扩展域名

		// 查询监听证书列表
		// REF: https://help.aliyun.com/zh/slb/application-load-balancer/developer-reference/api-alb-2020-06-16-listlistenercertificates
		listenerCertificates := make([]aliyunAlb.ListListenerCertificatesResponseBodyCertificates, 0)
		listListenerCertificatesLimit := int32(100)
		var listListenerCertificatesToken *string = nil
		for {
			listListenerCertificatesReq := &aliyunAlb.ListListenerCertificatesRequest{
				NextToken:       listListenerCertificatesToken,
				MaxResults:      tea.Int32(listListenerCertificatesLimit),
				ListenerId:      tea.String(cloudListenerId),
				CertificateType: tea.String("Server"),
			}
			listListenerCertificatesResp, err := d.sdkClients.alb.ListListenerCertificates(listListenerCertificatesReq)
			if err != nil {
				return xerrors.Wrap(err, "failed to execute sdk request 'alb.ListListenerCertificates'")
			}

			if listListenerCertificatesResp.Body.Certificates != nil {
				for _, listenerCertificate := range listListenerCertificatesResp.Body.Certificates {
					listenerCertificates = append(listenerCertificates, *listenerCertificate)
				}
			}

			if len(listListenerCertificatesResp.Body.Certificates) == 0 || listListenerCertificatesResp.Body.NextToken == nil {
				break
			} else {
				listListenerCertificatesToken = listListenerCertificatesResp.Body.NextToken
			}
		}

		d.logger.Logt("已查询到 ALB 监听下全部证书", listenerCertificates)

		// 遍历查询监听证书，并找出需要解除关联的证书
		// REF: https://help.aliyun.com/zh/slb/application-load-balancer/developer-reference/api-alb-2020-06-16-listlistenercertificates
		// REF: https://help.aliyun.com/zh/ssl-certificate/developer-reference/api-cas-2020-04-07-getusercertificatedetail
		certificateIsAssociated := false
		certificateIdsExpired := make([]string, 0)
		if len(listenerCertificates) > 0 {
			var errs []error

			for _, listenerCertificate := range listenerCertificates {
				// 监听证书 ID 格式：${证书 ID}-${地域}
				certificateId := strings.Split(*listenerCertificate.CertificateId, "-")[0]
				if certificateId == cloudCertId {
					certificateIsAssociated = true
					continue
				}

				if *listenerCertificate.IsDefault || !strings.EqualFold(*listenerCertificate.Status, "Associated") {
					continue
				}

				certificateIdAsInt64, err := strconv.ParseInt(certificateId, 10, 64)
				if err != nil {
					errs = append(errs, err)
					continue
				}

				getUserCertificateDetailReq := &aliyunCas.GetUserCertificateDetailRequest{
					CertId: tea.Int64(certificateIdAsInt64),
				}
				getUserCertificateDetailResp, err := d.sdkClients.cas.GetUserCertificateDetail(getUserCertificateDetailReq)
				if err != nil {
					errs = append(errs, xerrors.Wrap(err, "failed to execute sdk request 'cas.GetUserCertificateDetail'"))
					continue
				}

				certCnMatched := getUserCertificateDetailResp.Body.Common != nil && *getUserCertificateDetailResp.Body.Common == d.config.Domain
				certSanMatched := getUserCertificateDetailResp.Body.Sans != nil && slices.Contains(strings.Split(*getUserCertificateDetailResp.Body.Sans, ","), d.config.Domain)
				if !certCnMatched && !certSanMatched {
					continue
				}

				certEndDate, _ := time.Parse("2006-01-02", *getUserCertificateDetailResp.Body.EndDate)
				if time.Now().Before(certEndDate) {
					continue
				}

				certificateIdsExpired = append(certificateIdsExpired, certificateId)
			}

			if len(errs) > 0 {
				return errors.Join(errs...)
			}
		}

		// 关联监听和扩展证书
		// REF: https://help.aliyun.com/zh/slb/application-load-balancer/developer-reference/api-alb-2020-06-16-associateadditionalcertificateswithlistener
		if !certificateIsAssociated {
			associateAdditionalCertificatesFromListenerReq := &aliyunAlb.AssociateAdditionalCertificatesWithListenerRequest{
				ListenerId: tea.String(cloudListenerId),
				Certificates: []*aliyunAlb.AssociateAdditionalCertificatesWithListenerRequestCertificates{
					{
						CertificateId: tea.String(cloudCertId),
					},
				},
			}
			associateAdditionalCertificatesFromListenerResp, err := d.sdkClients.alb.AssociateAdditionalCertificatesWithListener(associateAdditionalCertificatesFromListenerReq)
			if err != nil {
				return xerrors.Wrap(err, "failed to execute sdk request 'alb.AssociateAdditionalCertificatesWithListener'")
			}

			d.logger.Logt("已关联 ALB 监听和扩展证书", associateAdditionalCertificatesFromListenerResp)
		}

		// 解除关联监听和扩展证书
		// REF: https://help.aliyun.com/zh/slb/application-load-balancer/developer-reference/api-alb-2020-06-16-dissociateadditionalcertificatesfromlistener
		if len(certificateIdsExpired) > 0 {
			dissociateAdditionalCertificates := make([]*aliyunAlb.DissociateAdditionalCertificatesFromListenerRequestCertificates, 0)
			for _, certificateId := range certificateIdsExpired {
				dissociateAdditionalCertificates = append(dissociateAdditionalCertificates, &aliyunAlb.DissociateAdditionalCertificatesFromListenerRequestCertificates{
					CertificateId: tea.String(certificateId),
				})
			}

			dissociateAdditionalCertificatesFromListenerReq := &aliyunAlb.DissociateAdditionalCertificatesFromListenerRequest{
				ListenerId:   tea.String(cloudListenerId),
				Certificates: dissociateAdditionalCertificates,
			}
			dissociateAdditionalCertificatesFromListenerResp, err := d.sdkClients.alb.DissociateAdditionalCertificatesFromListener(dissociateAdditionalCertificatesFromListenerReq)
			if err != nil {
				return xerrors.Wrap(err, "failed to execute sdk request 'alb.DissociateAdditionalCertificatesFromListener'")
			}

			d.logger.Logt("已解除关联 ALB 监听和扩展证书", dissociateAdditionalCertificatesFromListenerResp)
		}
	}

	return nil
}

func createSdkClients(accessKeyId, accessKeySecret, region string) (*wSdkClients, error) {
	// 接入点一览 https://api.aliyun.com/product/Alb
	var albEndpoint string
	switch region {
	case "cn-hangzhou-finance":
		albEndpoint = "alb.cn-hangzhou.aliyuncs.com"
	default:
		albEndpoint = fmt.Sprintf("alb.%s.aliyuncs.com", region)
	}

	albConfig := &aliyunOpen.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
		Endpoint:        tea.String(albEndpoint),
	}
	albClient, err := aliyunAlb.NewClient(albConfig)
	if err != nil {
		return nil, err
	}

	// 接入点一览 https://api.aliyun.com/product/cas
	var casEndpoint string
	if !strings.HasPrefix(region, "cn-") {
		casEndpoint = "cas.ap-southeast-1.aliyuncs.com"
	} else {
		casEndpoint = "cas.aliyuncs.com"
	}

	casConfig := &aliyunOpen.Config{
		Endpoint:        tea.String(casEndpoint),
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
	}
	casClient, err := aliyunCas.NewClient(casConfig)
	if err != nil {
		return nil, err
	}

	return &wSdkClients{
		alb: albClient,
		cas: casClient,
	}, nil
}

func createSslUploader(accessKeyId, accessKeySecret, region string) (uploader.Uploader, error) {
	casRegion := region
	if casRegion != "" {
		// 阿里云 CAS 服务接入点是独立于 ALB 服务的
		// 国内版固定接入点：华东一杭州
		// 国际版固定接入点：亚太东南一新加坡
		if casRegion != "" && !strings.HasPrefix(casRegion, "cn-") {
			casRegion = "ap-southeast-1"
		} else {
			casRegion = "cn-hangzhou"
		}
	}

	uploader, err := uploadersp.NewUploader(&uploadersp.UploaderConfig{
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
		Region:          casRegion,
	})
	return uploader, err
}
