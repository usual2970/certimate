package aliyunwaf

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	aliopen "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	aliwaf "github.com/alibabacloud-go/waf-openapi-20211001/v5/client"
	xerrors "github.com/pkg/errors"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	uploadersp "github.com/usual2970/certimate/internal/pkg/core/uploader/providers/aliyun-cas"
	"github.com/usual2970/certimate/internal/pkg/utils/sliceutil"
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
	sdkClient   *aliwaf.Client
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
		describeDefaultHttpsReq := &aliwaf.DescribeDefaultHttpsRequest{
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
		modifyDefaultHttpsReq := &aliwaf.ModifyDefaultHttpsRequest{
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
		describeDomainDetailReq := &aliwaf.DescribeDomainDetailRequest{
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
		modifyDomainReq := &aliwaf.ModifyDomainRequest{
			InstanceId: tea.String(d.config.InstanceId),
			RegionId:   tea.String(d.config.Region),
			Domain:     tea.String(d.config.Domain),
			Listen:     &aliwaf.ModifyDomainRequestListen{CertId: tea.String(upres.CertId)},
			Redirect:   &aliwaf.ModifyDomainRequestRedirect{Loadbalance: tea.String("iphash")},
		}
		modifyDomainReq = assign(modifyDomainReq, describeDomainDetailResp.Body)
		modifyDomainResp, err := d.sdkClient.ModifyDomain(modifyDomainReq)
		d.logger.Debug("sdk request 'waf.ModifyDomain'", slog.Any("request", modifyDomainReq), slog.Any("response", modifyDomainResp))
		if err != nil {
			return xerrors.Wrap(err, "failed to execute sdk request 'waf.ModifyDomain'")
		}
	}

	return nil
}

func createSdkClient(accessKeyId, accessKeySecret, region string) (*aliwaf.Client, error) {
	// 接入点一览：https://api.aliyun.com/product/waf-openapi
	config := &aliopen.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
		Endpoint:        tea.String(fmt.Sprintf("wafopenapi.%s.aliyuncs.com", region)),
	}

	client, err := aliwaf.NewClient(config)
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

func assign(source *aliwaf.ModifyDomainRequest, target *aliwaf.DescribeDomainDetailResponseBody) *aliwaf.ModifyDomainRequest {
	// `ModifyDomain` 中不传的字段表示使用默认值、而非保留原值，
	// 因此这里需要把原配置中的参数重新赋值回去。

	if target == nil {
		return source
	}

	if target.Listen != nil {
		if source.Listen == nil {
			source.Listen = &aliwaf.ModifyDomainRequestListen{}
		}

		if target.Listen.CipherSuite != nil {
			source.Listen.CipherSuite = tea.Int32(int32(*target.Listen.CipherSuite))
		}

		if target.Listen.CustomCiphers != nil {
			source.Listen.CustomCiphers = target.Listen.CustomCiphers
		}

		if target.Listen.EnableTLSv3 != nil {
			source.Listen.EnableTLSv3 = target.Listen.EnableTLSv3
		}

		if target.Listen.ExclusiveIp != nil {
			source.Listen.ExclusiveIp = target.Listen.ExclusiveIp
		}

		if target.Listen.FocusHttps != nil {
			source.Listen.FocusHttps = target.Listen.FocusHttps
		}

		if target.Listen.Http2Enabled != nil {
			source.Listen.Http2Enabled = target.Listen.Http2Enabled
		}

		if target.Listen.HttpPorts != nil {
			source.Listen.HttpPorts = sliceutil.Map(target.Listen.HttpPorts, func(v *int64) *int32 {
				if v == nil {
					return nil
				}
				return tea.Int32(int32(*v))
			})
		}

		if target.Listen.HttpsPorts != nil {
			source.Listen.HttpsPorts = sliceutil.Map(target.Listen.HttpsPorts, func(v *int64) *int32 {
				if v == nil {
					return nil
				}
				return tea.Int32(int32(*v))
			})
		}

		if target.Listen.IPv6Enabled != nil {
			source.Listen.IPv6Enabled = target.Listen.IPv6Enabled
		}

		if target.Listen.ProtectionResource != nil {
			source.Listen.ProtectionResource = target.Listen.ProtectionResource
		}

		if target.Listen.TLSVersion != nil {
			source.Listen.TLSVersion = target.Listen.TLSVersion
		}

		if target.Listen.XffHeaderMode != nil {
			source.Listen.XffHeaderMode = tea.Int32(int32(*target.Listen.XffHeaderMode))
		}

		if target.Listen.XffHeaders != nil {
			source.Listen.XffHeaders = target.Listen.XffHeaders
		}
	}

	if target.Redirect != nil {
		if source.Redirect == nil {
			source.Redirect = &aliwaf.ModifyDomainRequestRedirect{}
		}

		if target.Redirect.Backends != nil {
			source.Redirect.Backends = sliceutil.Map(target.Redirect.Backends, func(v *aliwaf.DescribeDomainDetailResponseBodyRedirectBackends) *string {
				if v == nil {
					return nil
				}
				return v.Backend
			})
		}

		if target.Redirect.BackupBackends != nil {
			source.Redirect.BackupBackends = sliceutil.Map(target.Redirect.BackupBackends, func(v *aliwaf.DescribeDomainDetailResponseBodyRedirectBackupBackends) *string {
				if v == nil {
					return nil
				}
				return v.Backend
			})
		}

		if target.Redirect.ConnectTimeout != nil {
			source.Redirect.ConnectTimeout = target.Redirect.ConnectTimeout
		}

		if target.Redirect.FocusHttpBackend != nil {
			source.Redirect.FocusHttpBackend = target.Redirect.FocusHttpBackend
		}

		if target.Redirect.Keepalive != nil {
			source.Redirect.Keepalive = target.Redirect.Keepalive
		}

		if target.Redirect.KeepaliveRequests != nil {
			source.Redirect.KeepaliveRequests = target.Redirect.KeepaliveRequests
		}

		if target.Redirect.KeepaliveTimeout != nil {
			source.Redirect.KeepaliveTimeout = target.Redirect.KeepaliveTimeout
		}

		if target.Redirect.Loadbalance != nil {
			source.Redirect.Loadbalance = target.Redirect.Loadbalance
		}

		if target.Redirect.ReadTimeout != nil {
			source.Redirect.ReadTimeout = target.Redirect.ReadTimeout
		}

		if target.Redirect.RequestHeaders != nil {
			source.Redirect.RequestHeaders = sliceutil.Map(target.Redirect.RequestHeaders, func(v *aliwaf.DescribeDomainDetailResponseBodyRedirectRequestHeaders) *aliwaf.ModifyDomainRequestRedirectRequestHeaders {
				if v == nil {
					return nil
				}
				return &aliwaf.ModifyDomainRequestRedirectRequestHeaders{
					Key:   v.Key,
					Value: v.Value,
				}
			})
		}

		if target.Redirect.Retry != nil {
			source.Redirect.Retry = target.Redirect.Retry
		}

		if target.Redirect.SniEnabled != nil {
			source.Redirect.SniEnabled = target.Redirect.SniEnabled
		}

		if target.Redirect.SniHost != nil {
			source.Redirect.SniHost = target.Redirect.SniHost
		}

		if target.Redirect.WriteTimeout != nil {
			source.Redirect.WriteTimeout = target.Redirect.WriteTimeout
		}

		if target.Redirect.XffProto != nil {
			source.Redirect.XffProto = target.Redirect.XffProto
		}
	}

	return source
}
