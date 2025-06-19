package tencentcloudecdn

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	tccdn "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn/v20180606"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tcssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"

	"github.com/certimate-go/certimate/pkg/core"
	sslmgrsp "github.com/certimate-go/certimate/pkg/core/ssl-manager/providers/tencentcloud-ssl"
)

type SSLDeployerProviderConfig struct {
	// 腾讯云 SecretId。
	SecretId string `json:"secretId"`
	// 腾讯云 SecretKey。
	SecretKey string `json:"secretKey"`
	// 加速域名（支持泛域名）。
	Domain string `json:"domain"`
}

type SSLDeployerProvider struct {
	config     *SSLDeployerProviderConfig
	logger     *slog.Logger
	sdkClients *wSDKClients
	sslManager core.SSLManager
}

var _ core.SSLDeployer = (*SSLDeployerProvider)(nil)

type wSDKClients struct {
	SSL *tcssl.Client
	CDN *tccdn.Client
}

func NewSSLDeployerProvider(config *SSLDeployerProviderConfig) (*SSLDeployerProvider, error) {
	if config == nil {
		return nil, errors.New("the configuration of the ssl deployer provider is nil")
	}

	clients, err := createSDKClients(config.SecretId, config.SecretKey)
	if err != nil {
		return nil, fmt.Errorf("could not create sdk client: %w", err)
	}

	sslmgr, err := sslmgrsp.NewSSLManagerProvider(&sslmgrsp.SSLManagerProviderConfig{
		SecretId:  config.SecretId,
		SecretKey: config.SecretKey,
	})
	if err != nil {
		return nil, fmt.Errorf("could not create ssl manager: %w", err)
	}

	return &SSLDeployerProvider{
		config:     config,
		logger:     slog.Default(),
		sdkClients: clients,
		sslManager: sslmgr,
	}, nil
}

func (d *SSLDeployerProvider) SetLogger(logger *slog.Logger) {
	if logger == nil {
		d.logger = slog.New(slog.DiscardHandler)
	} else {
		d.logger = logger
	}

	d.sslManager.SetLogger(logger)
}

func (d *SSLDeployerProvider) Deploy(ctx context.Context, certPEM string, privkeyPEM string) (*core.SSLDeployResult, error) {
	if d.config.Domain == "" {
		return nil, errors.New("config `domain` is required")
	}

	// 上传证书
	upres, err := d.sslManager.Upload(ctx, certPEM, privkeyPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to upload certificate file: %w", err)
	} else {
		d.logger.Info("ssl certificate uploaded", slog.Any("result", upres))
	}

	// 获取待部署的 CDN 实例
	// 如果是泛域名，根据证书匹配 CDN 实例
	instanceIds := make([]string, 0)
	if strings.HasPrefix(d.config.Domain, "*.") {
		domains, err := d.getDomainsByCertificateId(upres.CertId)
		if err != nil {
			return nil, err
		}

		instanceIds = domains
	} else {
		instanceIds = append(instanceIds, d.config.Domain)
	}

	if len(instanceIds) == 0 {
		d.logger.Info("no ecdn instances to deploy")
	} else {
		d.logger.Info("found ecdn instances to deploy", slog.Any("instanceIds", instanceIds))

		// 证书部署到 CDN 实例
		// REF: https://cloud.tencent.com/document/product/400/91667
		deployCertificateInstanceReq := tcssl.NewDeployCertificateInstanceRequest()
		deployCertificateInstanceReq.CertificateId = common.StringPtr(upres.CertId)
		deployCertificateInstanceReq.ResourceType = common.StringPtr("cdn")
		deployCertificateInstanceReq.Status = common.Int64Ptr(1)
		deployCertificateInstanceReq.InstanceIdList = common.StringPtrs(instanceIds)
		deployCertificateInstanceResp, err := d.sdkClients.SSL.DeployCertificateInstance(deployCertificateInstanceReq)
		d.logger.Debug("sdk request 'ssl.DeployCertificateInstance'", slog.Any("request", deployCertificateInstanceReq), slog.Any("response", deployCertificateInstanceResp))
		if err != nil {
			return nil, fmt.Errorf("failed to execute sdk request 'ssl.DeployCertificateInstance': %w", err)
		}

		// 循环获取部署任务详情，等待任务状态变更
		// REF: https://cloud.tencent.com/document/api/400/91658
		for {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
			}

			describeHostDeployRecordDetailReq := tcssl.NewDescribeHostDeployRecordDetailRequest()
			describeHostDeployRecordDetailReq.DeployRecordId = common.StringPtr(fmt.Sprintf("%d", *deployCertificateInstanceResp.Response.DeployRecordId))
			describeHostDeployRecordDetailResp, err := d.sdkClients.SSL.DescribeHostDeployRecordDetail(describeHostDeployRecordDetailReq)
			d.logger.Debug("sdk request 'ssl.DescribeHostDeployRecordDetail'", slog.Any("request", describeHostDeployRecordDetailReq), slog.Any("response", describeHostDeployRecordDetailResp))
			if err != nil {
				return nil, fmt.Errorf("failed to execute sdk request 'ssl.DescribeHostDeployRecordDetail': %w", err)
			}

			var runningCount, succeededCount, failedCount, totalCount int64
			if describeHostDeployRecordDetailResp.Response.TotalCount == nil {
				return nil, errors.New("unexpected deployment job status")
			} else {
				if describeHostDeployRecordDetailResp.Response.RunningTotalCount != nil {
					runningCount = *describeHostDeployRecordDetailResp.Response.RunningTotalCount
				}
				if describeHostDeployRecordDetailResp.Response.SuccessTotalCount != nil {
					succeededCount = *describeHostDeployRecordDetailResp.Response.SuccessTotalCount
				}
				if describeHostDeployRecordDetailResp.Response.FailedTotalCount != nil {
					failedCount = *describeHostDeployRecordDetailResp.Response.FailedTotalCount
				}
				if describeHostDeployRecordDetailResp.Response.TotalCount != nil {
					totalCount = *describeHostDeployRecordDetailResp.Response.TotalCount
				}

				if succeededCount+failedCount == totalCount {
					break
				}
			}

			d.logger.Info(fmt.Sprintf("waiting for deployment job completion (running: %d, succeeded: %d, failed: %d, total: %d) ...", runningCount, succeededCount, failedCount, totalCount))
			time.Sleep(time.Second * 5)
		}
	}

	return &core.SSLDeployResult{}, nil
}

func (d *SSLDeployerProvider) getDomainsByCertificateId(cloudCertId string) ([]string, error) {
	// 获取证书中的可用域名
	// REF: https://cloud.tencent.com/document/product/228/42491
	describeCertDomainsReq := tccdn.NewDescribeCertDomainsRequest()
	describeCertDomainsReq.CertId = common.StringPtr(cloudCertId)
	describeCertDomainsReq.Product = common.StringPtr("ecdn")
	describeCertDomainsResp, err := d.sdkClients.CDN.DescribeCertDomains(describeCertDomainsReq)
	d.logger.Debug("sdk request 'cdn.DescribeCertDomains'", slog.Any("request", describeCertDomainsReq), slog.Any("response", describeCertDomainsResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'cdn.DescribeCertDomains': %w", err)
	}

	domains := make([]string, 0)
	if describeCertDomainsResp.Response.Domains != nil {
		for _, domain := range describeCertDomainsResp.Response.Domains {
			domains = append(domains, *domain)
		}
	}

	return domains, nil
}

func createSDKClients(secretId, secretKey string) (*wSDKClients, error) {
	credential := common.NewCredential(secretId, secretKey)

	sslClient, err := tcssl.NewClient(credential, "", profile.NewClientProfile())
	if err != nil {
		return nil, err
	}

	cdnClient, err := tccdn.NewClient(credential, "", profile.NewClientProfile())
	if err != nil {
		return nil, err
	}

	return &wSDKClients{
		SSL: sslClient,
		CDN: cdnClient,
	}, nil
}
