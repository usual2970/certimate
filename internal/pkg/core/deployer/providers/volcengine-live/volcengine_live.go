package volcenginelive

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	xerrors "github.com/pkg/errors"
	veLive "github.com/volcengine/volc-sdk-golang/service/live/v20230101"
	ve "github.com/volcengine/volcengine-go-sdk/volcengine"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	uploadersp "github.com/usual2970/certimate/internal/pkg/core/uploader/providers/volcengine-live"
)

type DeployerConfig struct {
	// 火山引擎 AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// 火山引擎 AccessKeySecret。
	AccessKeySecret string `json:"accessKeySecret"`
	// 直播流域名（支持泛域名）。
	Domain string `json:"domain"`
}

type DeployerProvider struct {
	config      *DeployerConfig
	logger      *slog.Logger
	sdkClient   *veLive.Live
	sslUploader uploader.Uploader
}

var _ deployer.Deployer = (*DeployerProvider)(nil)

func NewDeployer(config *DeployerConfig) (*DeployerProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	client := veLive.NewInstance()
	client.SetAccessKey(config.AccessKeyId)
	client.SetSecretKey(config.AccessKeySecret)

	uploader, err := uploadersp.NewUploader(&uploadersp.UploaderConfig{
		AccessKeyId:     config.AccessKeyId,
		AccessKeySecret: config.AccessKeySecret,
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
	// 上传证书到 Live
	upres, err := d.sslUploader.Upload(ctx, certPem, privkeyPem)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to upload certificate file")
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
			listDomainDetailReq := &veLive.ListDomainDetailBody{
				PageNum:  listDomainDetailPageNum,
				PageSize: listDomainDetailPageSize,
			}
			listDomainDetailResp, err := d.sdkClient.ListDomainDetail(ctx, listDomainDetailReq)
			d.logger.Debug("sdk request 'live.ListDomainDetail'", slog.Any("request", listDomainDetailReq), slog.Any("response", listDomainDetailResp))
			if err != nil {
				return nil, xerrors.Wrap(err, "failed to execute sdk request 'live.ListDomainDetail'")
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
			// 绑定证书
			// REF: https://www.volcengine.com/docs/6469/1186278#%E7%BB%91%E5%AE%9A%E8%AF%81%E4%B9%A6
			bindCertReq := &veLive.BindCertBody{
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

		if len(errs) > 0 {
			return nil, errors.Join(errs...)
		}
	}

	return &deployer.DeployResult{}, nil
}
