package volcenginelive

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	velive "github.com/volcengine/volc-sdk-golang/service/live/v20230101"
	ve "github.com/volcengine/volcengine-go-sdk/volcengine"

	"github.com/certimate-go/certimate/pkg/core"
	sslmgrsp "github.com/certimate-go/certimate/pkg/core/ssl-manager/providers/volcengine-live"
)

type SSLDeployerProviderConfig struct {
	// 火山引擎 AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// 火山引擎 AccessKeySecret。
	AccessKeySecret string `json:"accessKeySecret"`
	// 直播流域名（支持泛域名）。
	Domain string `json:"domain"`
}

type SSLDeployerProvider struct {
	config     *SSLDeployerProviderConfig
	logger     *slog.Logger
	sdkClient  *velive.Live
	sslManager core.SSLManager
}

var _ core.SSLDeployer = (*SSLDeployerProvider)(nil)

func NewSSLDeployerProvider(config *SSLDeployerProviderConfig) (*SSLDeployerProvider, error) {
	if config == nil {
		return nil, errors.New("the configuration of the ssl deployer provider is nil")
	}

	client := velive.NewInstance()
	client.SetAccessKey(config.AccessKeyId)
	client.SetSecretKey(config.AccessKeySecret)

	sslmgr, err := sslmgrsp.NewSSLManagerProvider(&sslmgrsp.SSLManagerProviderConfig{
		AccessKeyId:     config.AccessKeyId,
		AccessKeySecret: config.AccessKeySecret,
	})
	if err != nil {
		return nil, fmt.Errorf("could not create ssl manager: %w", err)
	}

	return &SSLDeployerProvider{
		config:     config,
		logger:     slog.Default(),
		sdkClient:  client,
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

	domains := make([]string, 0)
	if strings.HasPrefix(d.config.Domain, "*.") {
		listDomainDetailPageNum := int32(1)
		listDomainDetailPageSize := int32(1000)
		listDomainDetailTotal := 0
		for {
			// 查询域名列表
			// REF: https://www.volcengine.com/docs/6469/1186277#%E6%9F%A5%E8%AF%A2%E5%9F%9F%E5%90%8D%E5%88%97%E8%A1%A8
			listDomainDetailReq := &velive.ListDomainDetailBody{
				PageNum:  listDomainDetailPageNum,
				PageSize: listDomainDetailPageSize,
			}
			listDomainDetailResp, err := d.sdkClient.ListDomainDetail(ctx, listDomainDetailReq)
			d.logger.Debug("sdk request 'live.ListDomainDetail'", slog.Any("request", listDomainDetailReq), slog.Any("response", listDomainDetailResp))
			if err != nil {
				return nil, fmt.Errorf("failed to execute sdk request 'live.ListDomainDetail': %w", err)
			}

			if listDomainDetailResp.Result.DomainList != nil {
				for _, item := range listDomainDetailResp.Result.DomainList {
					// 仅匹配泛域名的下一级子域名
					wildcardDomain := strings.TrimPrefix(d.config.Domain, "*")
					if strings.HasSuffix(item.Domain, wildcardDomain) && !strings.Contains(strings.TrimSuffix(item.Domain, wildcardDomain), ".") {
						domains = append(domains, item.Domain)
					}
				}
			}

			listDomainDetailLen := len(listDomainDetailResp.Result.DomainList)
			if listDomainDetailLen < int(listDomainDetailPageSize) || int(listDomainDetailResp.Result.Total) <= listDomainDetailTotal+listDomainDetailLen {
				break
			} else {
				listDomainDetailPageNum++
				listDomainDetailTotal += listDomainDetailLen
			}
		}

		if len(domains) == 0 {
			return nil, errors.New("domain not found")
		}
	} else {
		domains = append(domains, d.config.Domain)
	}

	if len(domains) > 0 {
		var errs []error

		for _, domain := range domains {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
				// 绑定证书
				// REF: https://www.volcengine.com/docs/6469/1186278#%E7%BB%91%E5%AE%9A%E8%AF%81%E4%B9%A6
				bindCertReq := &velive.BindCertBody{
					ChainID: upres.CertId,
					Domain:  domain,
					HTTPS:   ve.Bool(true),
				}
				bindCertResp, err := d.sdkClient.BindCert(ctx, bindCertReq)
				d.logger.Debug("sdk request 'live.BindCert'", slog.Any("request", bindCertReq), slog.Any("response", bindCertResp))
				if err != nil {
					errs = append(errs, err)
				}
			}
		}

		if len(errs) > 0 {
			return nil, errors.Join(errs...)
		}
	}

	return &core.SSLDeployResult{}, nil
}
