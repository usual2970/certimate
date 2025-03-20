package aliyunwaf

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	aliyunOpen "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	aliyunWaf "github.com/alibabacloud-go/waf-openapi-20211001/v5/client"
	xerrors "github.com/pkg/errors"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
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
	// 服务版本。
	ServiceVersion string `json:"serviceVersion"`
	// WAF 实例 ID。
	InstanceId string `json:"instanceId"`
	// 接入域名（支持泛域名）。
	Domain string `json:"domain,omitempty"`
}

type DeployerProvider struct {
	config      *DeployerConfig
	logger      *slog.Logger
	sdkClient   *aliyunWaf.Client
	sslUploader uploader.Uploader
}

var _ deployer.Deployer = (*DeployerProvider)(nil)

func NewDeployer(config *DeployerConfig) (*DeployerProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	client, err := createSdkClient(config.AccessKeyId, config.AccessKeySecret, config.Region)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk client")
	}

	uploader, err := createSslUploader(config.AccessKeyId, config.AccessKeySecret, config.Region)
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
	if d.config.InstanceId == "" {
		return nil, errors.New("config `instanceId` is required")
	}

	switch d.config.ServiceVersion {
	case "3", "3.0":
		if err := d.deployToWAF3(ctx, certPem, privkeyPem); err != nil {
			return nil, err
		}

	default:
		return nil, xerrors.Errorf("unsupported service version: %s", d.config.ServiceVersion)
	}

	return &deployer.DeployResult{}, nil
}

func (d *DeployerProvider) deployToWAF3(ctx context.Context, certPem string, privkeyPem string) error {
	// 上传证书到 CAS
	upres, err := d.sslUploader.Upload(ctx, certPem, privkeyPem)
	if err != nil {
		return xerrors.Wrap(err, "failed to upload certificate file")
	} else {
		d.logger.Info("ssl certificate uploaded", slog.Any("result", upres))
	}

	if d.config.Domain == "" {
		// 未指定接入域名，只需替换默认证书即可

		// 查询默认 SSL/TLS 设置
		// REF: https://help.aliyun.com/zh/waf/web-application-firewall-3-0/developer-reference/api-waf-openapi-2021-10-01-describedefaulthttps
		describeDefaultHttpsReq := &aliyunWaf.DescribeDefaultHttpsRequest{
			InstanceId: tea.String(d.config.InstanceId),
			RegionId:   tea.String(d.config.Region),
		}
		describeDefaultHttpsResp, err := d.sdkClient.DescribeDefaultHttps(describeDefaultHttpsReq)
		d.logger.Debug("sdk request 'waf.DescribeDefaultHttps'", slog.Any("request", describeDefaultHttpsReq), slog.Any("response", describeDefaultHttpsResp))
		if err != nil {
			return xerrors.Wrap(err, "failed to execute sdk request 'waf.DescribeDefaultHttps'")
		}

		// 修改默认 SSL/TLS 设置
		// REF: https://help.aliyun.com/zh/waf/web-application-firewall-3-0/developer-reference/api-waf-openapi-2021-10-01-modifydefaulthttps
		modifyDefaultHttpsReq := &aliyunWaf.ModifyDefaultHttpsRequest{
			InstanceId:  tea.String(d.config.InstanceId),
			RegionId:    tea.String(d.config.Region),
			CertId:      tea.String(upres.CertId),
			TLSVersion:  tea.String("tlsv1"),
			EnableTLSv3: tea.Bool(false),
		}
		if describeDefaultHttpsResp.Body != nil && describeDefaultHttpsResp.Body.DefaultHttps != nil {
			modifyDefaultHttpsReq.TLSVersion = describeDefaultHttpsResp.Body.DefaultHttps.TLSVersion
			modifyDefaultHttpsReq.EnableTLSv3 = describeDefaultHttpsResp.Body.DefaultHttps.EnableTLSv3
		}
		modifyDefaultHttpsResp, err := d.sdkClient.ModifyDefaultHttps(modifyDefaultHttpsReq)
		d.logger.Debug("sdk request 'waf.ModifyDefaultHttps'", slog.Any("request", modifyDefaultHttpsReq), slog.Any("response", modifyDefaultHttpsResp))
		if err != nil {
			return xerrors.Wrap(err, "failed to execute sdk request 'waf.ModifyDefaultHttps'")
		}
	} else {
		// 指定接入域名

		// 查询 CNAME 接入详情
		// REF: https://help.aliyun.com/zh/waf/web-application-firewall-3-0/developer-reference/api-waf-openapi-2021-10-01-describedomaindetail
		describeDomainDetailReq := &aliyunWaf.DescribeDomainDetailRequest{
			InstanceId: tea.String(d.config.InstanceId),
			RegionId:   tea.String(d.config.Region),
			Domain:     tea.String(d.config.Domain),
		}
		describeDomainDetailResp, err := d.sdkClient.DescribeDomainDetail(describeDomainDetailReq)
		d.logger.Debug("sdk request 'waf.DescribeDomainDetail'", slog.Any("request", describeDomainDetailReq), slog.Any("response", describeDomainDetailResp))
		if err != nil {
			return xerrors.Wrap(err, "failed to execute sdk request 'waf.DescribeDomainDetail'")
		}

		// 修改 CNAME 接入资源
		// REF: https://help.aliyun.com/zh/waf/web-application-firewall-3-0/developer-reference/api-waf-openapi-2021-10-01-modifydomain
		modifyDomainReq := &aliyunWaf.ModifyDomainRequest{
			InstanceId: tea.String(d.config.InstanceId),
			RegionId:   tea.String(d.config.Region),
			Domain:     tea.String(d.config.Domain),
			Listen: &aliyunWaf.ModifyDomainRequestListen{
				CertId:      tea.String(upres.CertId),
				TLSVersion:  tea.String("tlsv1"),
				EnableTLSv3: tea.Bool(false),
			},
			Redirect: &aliyunWaf.ModifyDomainRequestRedirect{
				Loadbalance: tea.String("iphash"),
			},
		}
		if describeDomainDetailResp.Body != nil && describeDomainDetailResp.Body.Listen != nil {
			modifyDomainReq.Listen.TLSVersion = describeDomainDetailResp.Body.Listen.TLSVersion
			modifyDomainReq.Listen.EnableTLSv3 = describeDomainDetailResp.Body.Listen.EnableTLSv3
			modifyDomainReq.Listen.FocusHttps = describeDomainDetailResp.Body.Listen.FocusHttps
		}
		if describeDomainDetailResp.Body != nil && describeDomainDetailResp.Body.Redirect != nil {
			modifyDomainReq.Redirect.Loadbalance = describeDomainDetailResp.Body.Redirect.Loadbalance
			modifyDomainReq.Redirect.FocusHttpBackend = describeDomainDetailResp.Body.Redirect.FocusHttpBackend
			modifyDomainReq.Redirect.SniEnabled = describeDomainDetailResp.Body.Redirect.SniEnabled
			modifyDomainReq.Redirect.SniHost = describeDomainDetailResp.Body.Redirect.SniHost
		}
		modifyDomainResp, err := d.sdkClient.ModifyDomain(modifyDomainReq)
		d.logger.Debug("sdk request 'waf.ModifyDomain'", slog.Any("request", modifyDomainReq), slog.Any("response", modifyDomainResp))
		if err != nil {
			return xerrors.Wrap(err, "failed to execute sdk request 'waf.ModifyDomain'")
		}
	}

	return nil
}

func createSdkClient(accessKeyId, accessKeySecret, region string) (*aliyunWaf.Client, error) {
	// 接入点一览：https://api.aliyun.com/product/waf-openapi
	config := &aliyunOpen.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
		Endpoint:        tea.String(fmt.Sprintf("wafopenapi.%s.aliyuncs.com", region)),
	}

	client, err := aliyunWaf.NewClient(config)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func createSslUploader(accessKeyId, accessKeySecret, region string) (uploader.Uploader, error) {
	casRegion := region
	if casRegion != "" {
		// 阿里云 CAS 服务接入点是独立于 WAF 服务的
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
