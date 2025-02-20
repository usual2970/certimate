package jdcloudalb

import (
	"context"
	"errors"
	"fmt"
	"strings"

	jdCore "github.com/jdcloud-api/jdcloud-sdk-go/core"
	jdCommon "github.com/jdcloud-api/jdcloud-sdk-go/services/common/models"
	jdLbApi "github.com/jdcloud-api/jdcloud-sdk-go/services/lb/apis"
	jdLbClient "github.com/jdcloud-api/jdcloud-sdk-go/services/lb/client"
	jdLbModel "github.com/jdcloud-api/jdcloud-sdk-go/services/lb/models"
	xerrors "github.com/pkg/errors"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/logger"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	uploadersp "github.com/usual2970/certimate/internal/pkg/core/uploader/providers/jdcloud-ssl"
	"github.com/usual2970/certimate/internal/pkg/utils/slices"
)

type DeployerConfig struct {
	// 京东云 AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// 京东云 AccessKeySecret。
	AccessKeySecret string `json:"accessKeySecret"`
	// 京东云地域 ID。
	RegionId string `json:"regionId"`
	// 部署资源类型。
	ResourceType ResourceType `json:"resourceType"`
	// 负载均衡器 ID。
	// 部署资源类型为 [RESOURCE_TYPE_LOADBALANCER] 时必填。
	LoadbalancerId string `json:"loadbalancerId,omitempty"`
	// 监听器 ID。
	// 部署资源类型为 [RESOURCE_TYPE_LISTENER] 时必填。
	ListenerId string `json:"listenerId,omitempty"`
	// SNI 域名（支持泛域名）。
	// 部署资源类型为 [RESOURCE_TYPE_LOADBALANCER]、[RESOURCE_TYPE_LISTENER] 时选填。
	Domain string `json:"domain,omitempty"`
}

type DeployerProvider struct {
	config      *DeployerConfig
	logger      logger.Logger
	sdkClient   *jdLbClient.LbClient
	sslUploader uploader.Uploader
}

var _ deployer.Deployer = (*DeployerProvider)(nil)

func NewDeployer(config *DeployerConfig) (*DeployerProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	client, err := createSdkClient(config.AccessKeyId, config.AccessKeySecret)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk client")
	}

	uploader, err := uploadersp.NewUploader(&uploadersp.UploaderConfig{
		AccessKeyId:     config.AccessKeyId,
		AccessKeySecret: config.AccessKeySecret,
	})
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create ssl uploader")
	}

	return &DeployerProvider{
		config:      config,
		logger:      logger.NewNilLogger(),
		sdkClient:   client,
		sslUploader: uploader,
	}, nil
}

func (d *DeployerProvider) WithLogger(logger logger.Logger) *DeployerProvider {
	d.logger = logger
	return d
}

func (d *DeployerProvider) Deploy(ctx context.Context, certPem string, privkeyPem string) (*deployer.DeployResult, error) {
	// 上传证书到 SSL
	upres, err := d.sslUploader.Upload(ctx, certPem, privkeyPem)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to upload certificate file")
	} else {
		d.logger.Logt("certificate file uploaded", upres)
	}

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

	// 查询负载均衡器详情
	// REF: https://docs.jdcloud.com/cn/load-balancer/api/describeloadbalancer
	describeLoadBalancerReq := jdLbApi.NewDescribeLoadBalancerRequest(d.config.RegionId, d.config.LoadbalancerId)
	describeLoadBalancerResp, err := d.sdkClient.DescribeLoadBalancer(describeLoadBalancerReq)
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'lb.DescribeLoadBalancer'")
	} else {
		d.logger.Logt("已查询到负载均衡器详情", describeLoadBalancerResp)
	}

	// 查询监听器列表
	// REF: https://docs.jdcloud.com/cn/load-balancer/api/describelisteners
	listenerIds := make([]string, 0)
	describeListenersPageNumber := 1
	describeListenersPageSize := 100
	for {
		describeListenersReq := jdLbApi.NewDescribeListenersRequest(d.config.RegionId)
		describeListenersReq.SetFilters([]jdCommon.Filter{{Name: "loadBalancerId", Values: []string{d.config.LoadbalancerId}}})
		describeListenersReq.SetPageSize(describeListenersPageNumber)
		describeListenersReq.SetPageSize(describeListenersPageSize)
		describeListenersResp, err := d.sdkClient.DescribeListeners(describeListenersReq)
		if err != nil {
			return xerrors.Wrap(err, "failed to execute sdk request 'lb.DescribeListeners'")
		}

		for _, listener := range describeListenersResp.Result.Listeners {
			if strings.EqualFold(listener.Protocol, "https") || strings.EqualFold(listener.Protocol, "tls") {
				listenerIds = append(listenerIds, listener.ListenerId)
			}
		}

		if len(describeListenersResp.Result.Listeners) < int(describeListenersPageSize) {
			break
		} else {
			describeListenersPageNumber++
		}
	}

	// 遍历更新监听器证书
	if len(listenerIds) == 0 {
		return errors.New("listener not found")
	} else {
		d.logger.Logt("已查询到负载均衡器下的全部 HTTPS/TLS 监听器", listenerIds)

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

	// 更新监听器证书
	if err := d.updateListenerCertificate(ctx, d.config.ListenerId, cloudCertId); err != nil {
		return err
	}

	return nil
}

func (d *DeployerProvider) updateListenerCertificate(ctx context.Context, cloudListenerId string, cloudCertId string) error {
	// 查询监听器详情
	// REF: https://docs.jdcloud.com/cn/load-balancer/api/describelistener
	describeListenerReq := jdLbApi.NewDescribeListenerRequest(d.config.RegionId, cloudListenerId)
	describeListenerResp, err := d.sdkClient.DescribeListener(describeListenerReq)
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'lb.DescribeListener'")
	} else {
		d.logger.Logt("已查询到监听器详情", describeListenerResp)
	}

	if d.config.Domain == "" {
		// 未指定 SNI，只需部署到监听器

		// 修改监听器信息
		// REF: https://docs.jdcloud.com/cn/load-balancer/api/updatelistener
		updateListenerReq := jdLbApi.NewUpdateListenerRequest(d.config.RegionId, cloudListenerId)
		updateListenerReq.SetCertificateSpecs([]jdLbModel.CertificateSpec{{CertificateId: cloudCertId}})
		updateListenerResp, err := d.sdkClient.UpdateListener(updateListenerReq)
		if err != nil {
			return xerrors.Wrap(err, "failed to execute sdk request 'lb.UpdateListener'")
		} else {
			d.logger.Logt("已修改监听器信息", updateListenerResp)
		}
	} else {
		// 指定 SNI，需部署到扩展证书

		extCertSpecs := slices.Filter(describeListenerResp.Result.Listener.ExtensionCertificateSpecs, func(extCertSpec jdLbModel.ExtensionCertificateSpec) bool {
			return extCertSpec.Domain == d.config.Domain
		})
		if len(extCertSpecs) == 0 {
			return errors.New("extension certificate spec not found")
		}

		// 批量修改扩展证书
		// REF: https://docs.jdcloud.com/cn/load-balancer/api/updatelistenercertificates
		updateListenerCertificatesReq := jdLbApi.NewUpdateListenerCertificatesRequest(
			d.config.RegionId,
			cloudListenerId,
			slices.Map(extCertSpecs, func(extCertSpec jdLbModel.ExtensionCertificateSpec) jdLbModel.ExtCertificateUpdateSpec {
				return jdLbModel.ExtCertificateUpdateSpec{
					CertificateBindId: extCertSpec.CertificateBindId,
					CertificateId:     &cloudCertId,
					Domain:            &extCertSpec.Domain,
				}
			}),
		)
		updateListenerCertificatesResp, err := d.sdkClient.UpdateListenerCertificates(updateListenerCertificatesReq)
		if err != nil {
			return xerrors.Wrap(err, "failed to execute sdk request 'lb.UpdateListenerCertificates'")
		} else {
			d.logger.Logt("已批量修改扩展证书", updateListenerCertificatesResp)
		}
	}

	return nil
}

func createSdkClient(accessKeyId, accessKeySecret string) (*jdLbClient.LbClient, error) {
	clientCredentials := jdCore.NewCredentials(accessKeyId, accessKeySecret)
	client := jdLbClient.NewLbClient(clientCredentials)
	return client, nil
}
